package main

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	ClientID     int          `json:"client_id"`
	PlayerCircle PlayerCircle `json:"player_circle"`
}

func (m Message) String() string {
	jsonMessage, _ := json.Marshal(m)
	return fmt.Sprint(string(jsonMessage))
}
