package main

type Message struct {
	ClientID     int     `json:"client_id"`
	PlayerX      float64 `json:"pos_x"`
	PlayerY      float64 `json:"pos_y"`
	PlayerHeight int     `json:"height"`
	PlayerWidth  int     `json:"width"`
}
