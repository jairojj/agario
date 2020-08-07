package main

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	ClientID     int `json:"client_id"`
	PlayerCircle struct {
		PosX   float64 `json:"pos_x"`
		PosY   float64 `json:"pos_y"`
		Height int     `json:"height"`
		Width  int     `json:"width"`
	} `json:"player_circle"`
}

func (m Message) String() string {
	jsonMessage, _ := json.Marshal(m)
	return fmt.Sprint(string(jsonMessage))
}
