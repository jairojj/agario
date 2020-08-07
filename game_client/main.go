package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
}

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

var playerCircle *PlayerCircle
var playerMoves chan PlayerCircle

func (g *Game) Update(screen *ebiten.Image) error {
	handleInput()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	circle, _ := ebiten.NewImage(playerCircle.Width, playerCircle.Height, ebiten.FilterDefault)
	circle.Fill(color.White)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(playerCircle.PosX, playerCircle.PosY)

	screen.DrawImage(circle, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	playerMoves = make(chan PlayerCircle)

	go startWsClient(playerMoves)

	playerCircle = &PlayerCircle{
		PosX:   0,
		PosY:   0,
		Width:  20,
		Height: 20,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("AGARIO")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func handleInput() {
	anyKeyPressed := false

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		playerCircle.PosY--
		anyKeyPressed = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		playerCircle.PosX++
		anyKeyPressed = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		playerCircle.PosY++
		anyKeyPressed = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		playerCircle.PosX--
		anyKeyPressed = true
	}

	if anyKeyPressed {
		playerMoves <- *playerCircle
	}
}
