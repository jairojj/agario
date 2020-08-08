package main

type ConsumableSquare struct {
	PosX   float64 `json:"pos_x,omitempty"`
	PosY   float64 `json:"pos_y,omitempty"`
	Height int     `json:"height,omitempty"`
	Width  int     `json:"width,omitempty"`
	Color  string  `json:"color,omitempty"`
}
