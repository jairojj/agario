package main

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	ClientID          int                `json:"client_id,omitempty"`
	PlayerCircle      PlayerCircle       `json:"player_circle,omitempty"`
	Event             Event              `json:"event,omitempty"`
	ConsumableSquares []ConsumableSquare `json:"consumable_squares,omitempty"`
}

func (m Message) String() string {
	jsonMessage, _ := json.Marshal(m)
	return fmt.Sprint(string(jsonMessage))
}

type Event int

const (
	PlayerMoved Event = iota + 1
	PlayerDisconnected
	ConsumableSquareChanged
)
