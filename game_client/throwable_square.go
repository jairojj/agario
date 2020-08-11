package main

import (
	"encoding/json"
	"fmt"
)

type ThrowableSquare struct {
	CurrentX int    `json:"current_x,omitempty"`
	CurrentY int    `json:"current_y,omitempty"`
	DestX    int    `json:"dest_x,omitempty"`
	DestY    int    `json:"dest_y,omitempty"`
	Reached  bool   `json:"reached,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
	Color    string `json:"color,omitempty"`
}

func (t ThrowableSquare) String() string {
	jsonThrowableSquare, _ := json.Marshal(t)
	return fmt.Sprint(string(jsonThrowableSquare))
}
