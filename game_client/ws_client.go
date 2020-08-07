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

func startWsClient(playerMoves chan PlayerCircle) {
	rand.Seed(time.Now().UnixNano())
	clientID := rand.Intn(100)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:3000", Path: "/ws", RawQuery: fmt.Sprint("id=", clientID)}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			message := Message{}
			err := c.ReadJSON(&message)
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Print("recv:", message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case playerCircle := <-playerMoves:
			message := Message{PlayerCircle: playerCircle, ClientID: clientID}
			log.Println("Sending: ", message)

			err = c.WriteJSON(message)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
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
	}
}
