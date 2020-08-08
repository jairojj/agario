package main

import (
	"encoding/json"
	"fmt"
)

type PlayerCircle struct {
	PosX   float64 `json:"pos_x,omitempty"`
	PosY   float64 `json:"pos_y,omitempty"`
	Height int     `json:"height,omitempty"`
	Width  int     `json:"width,omitempty"`
	Color  string  `json:"color,omitempty"`
	Points int     `json:"points,omitempty"`
}

func (p PlayerCircle) String() string {
	jsonPlayerCircle, _ := json.Marshal(p)
	return fmt.Sprint(string(jsonPlayerCircle))
}
