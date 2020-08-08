package main

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	ClientID     int          `json:"client_id"`
	PlayerCircle PlayerCircle `json:"player_circle"`
	Event        Event        `json:"event"`
}

func (m Message) String() string {
	jsonMessage, _ := json.Marshal(m)
	return fmt.Sprint(string(jsonMessage))
}

type Event int

const (
	PlayerMoved Event = iota + 1
	PlayerDisconnected
)
