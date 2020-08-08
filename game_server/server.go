package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan Message

	ID int

	PlayerCircle PlayerCircle
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()

		message := Message{ClientID: c.ID, Event: PlayerDisconnected}
		c.hub.broadcast <- message
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		message := Message{}
		err := c.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.PlayerCircle = message.PlayerCircle

		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				log.Println("Send close message to client: ", c.ID)
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if message.ClientID == c.ID {
				continue
			}

			c.conn.WriteJSON(message)

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request, consumableSquares *[]ConsumableSquare) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	queryParams := r.URL.Query()
	clientID, _ := strconv.Atoi(queryParams.Get("id"))

	initialPlayerCircle := PlayerCircle{
		PosX:   0,
		PosY:   0,
		Width:  20,
		Height: 20,
		Color:  queryParams.Get("color"),
	}

	client := &Client{hub: hub, conn: conn, send: make(chan Message, 256), ID: clientID, PlayerCircle: initialPlayerCircle}

	// Send current player to all other players
	message := Message{ClientID: client.ID, PlayerCircle: initialPlayerCircle, Event: PlayerMoved}
	hub.broadcast <- message

	for otherClient := range hub.clients {
		// Send all other players to current player
		message = Message{ClientID: otherClient.ID, PlayerCircle: otherClient.PlayerCircle, Event: PlayerMoved}
		client.conn.WriteJSON(message)
	}

	//Send current consumable squares when client connects
	message = Message{ClientID: client.ID, Event: ConsumableSquareChanged, ConsumableSquares: *consumableSquares}
	client.conn.WriteJSON(message)

	client.hub.register <- client

	log.Println("New client: ", client)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
