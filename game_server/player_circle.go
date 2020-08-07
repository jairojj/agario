package main

import (
	"encoding/json"
	"fmt"
)

type PlayerCircle struct {
	PosX   float64 `json:"pos_x"`
	PosY   float64 `json:"pos_y"`
	Height int     `json:"height"`
	Width  int     `json:"width"`
}

func (p PlayerCircle) String() string {
	jsonPlayerCircle, _ := json.Marshal(p)
	return fmt.Sprint(string(jsonPlayerCircle))
}
