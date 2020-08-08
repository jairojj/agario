package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func (g *Game) startWsClient() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:3000", Path: "/ws", RawQuery: fmt.Sprintf("id=%d&color=%s", g.CurrentPlayerID, g.PlayerCircle.Color)}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go g.ReadMessages(done, c)
	g.WriteMessages(done, c, interrupt)
}

func (g *Game) ReadMessages(done chan struct{}, c *websocket.Conn) {
	defer close(done)
	for {
		message := Message{}
		err := c.ReadJSON(&message)
		if err != nil {
			log.Println("read:", err)
			return
		}

		log.Println("Received: ", message)

		switch message.Event {
		case PlayerMoved:
			g.OtherPlayers[message.ClientID] = message.PlayerCircle
		case PlayerDisconnected:
			delete(g.OtherPlayers, message.ClientID)
		case ConsumableSquareChanged:
			g.ConsumableSquares = message.ConsumableSquares
		}
	}
}

func (g *Game) WriteMessages(done chan struct{}, c *websocket.Conn, interrupt chan os.Signal) {
	for {
		select {
		case <-done:
			return
		case message := <-g.MessageQueue:
			err := c.WriteJSON(message)
			if err != nil {
				log.Println("write:", err)
				return
			}

			log.Println("Sending: ", message)
		case <-interrupt:
			closeConnection(c, done)
		}
	}
}

func closeConnection(c *websocket.Conn, done chan struct{}) {
	log.Println("interrupt")

	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
	select {
	case <-done:
	case <-time.After(time.Second):
	}
	return
}
