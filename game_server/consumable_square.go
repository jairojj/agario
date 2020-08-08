package main

import "math/rand"

type ConsumableSquare struct {
	PosX   float64 `json:"pos_x,omitempty"`
	PosY   float64 `json:"pos_y,omitempty"`
	Height int     `json:"height,omitempty"`
	Width  int     `json:"width,omitempty"`
	Color  string  `json:"color,omitempty"`
}

func GenerateConsumableSquares(count int) []ConsumableSquare {
	consumableSquares := []ConsumableSquare{}

	for i := 0; i < count; i++ {
		consumableSquares = append(consumableSquares, ConsumableSquare{
			PosX:   float64(rand.Intn(screenWidth - 10)),
			PosY:   float64(rand.Intn(screenHeight - 10)),
			Height: 5,
			Width:  5,
			Color:  getRandomColor(),
		})
	}

	return consumableSquares
}
