package main

import (
	"image/color"
	"math/rand"
)

var Colors map[string]color.Color = map[string]color.Color{
	"white":  color.White,
	"red":    color.RGBA{255, 0, 0, 255},
	"yellow": color.RGBA{255, 255, 0, 255},
	"blue":   color.RGBA{0, 0, 255, 255},
	"green":  color.RGBA{0, 255, 0, 255}}

func getRandomColor() string {
	colors := []string{}

	for color := range Colors {
		colors = append(colors, color)
	}

	return colors[rand.Intn(len(colors))]
}
