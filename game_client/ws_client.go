package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func startWsClient(playerMoves chan PlayerCircle, game *Game, randomColor string) {
	rand.Seed(time.Now().UnixNano())
	clientID := rand.Intn(100)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:3000", Path: "/ws", RawQuery: fmt.Sprintf("id=%d&color=%s", clientID, randomColor)}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go readMessages(done, c, game)
	writeMessages(done, c, clientID, interrupt)
}

func readMessages(done chan struct{}, c *websocket.Conn, game *Game) {
	defer close(done)
	for {
		message := Message{}
		err := c.ReadJSON(&message)
		if err != nil {
			log.Println("read:", err)
			return
		}

		switch message.Event {
		case PlayerMoved:
			game.OtherPlayers[message.ClientID] = message.PlayerCircle
		case PlayerDisconnected:
			delete(game.OtherPlayers, message.ClientID)
		case ConsumableSquareChanged:
			game.ConsumableSquares = message.ConsumableSquares
		}
	}
}

func writeMessages(done chan struct{}, c *websocket.Conn, clientID int, interrupt chan os.Signal) {
	for {
		select {
		case <-done:
			return
		case playerCircle := <-playerMoves:
			message := Message{PlayerCircle: playerCircle, ClientID: clientID, Event: PlayerMoved}

			err := c.WriteJSON(message)
			if err != nil {
				log.Println("write:", err)
				return
			}
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
