package main

import (
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
	posX   float64 `json:"pos_x"`
	posY   float64 `json:"pos_y"`
	height int     `json:"height"`
	width  int     `json:"width"`
}

func (p PlayerCircle) String() string {
	return fmt.Sprint(p.posX, p.posY)
}

var playerCircle *PlayerCircle
var playerMoves chan PlayerCircle

func (g *Game) Update(screen *ebiten.Image) error {
	handleInput()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	circle, _ := ebiten.NewImage(playerCircle.width, playerCircle.height, ebiten.FilterDefault)
	circle.Fill(color.White)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(playerCircle.posX, playerCircle.posY)

	screen.DrawImage(circle, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	playerMoves = make(chan PlayerCircle)

	go startWsClient(playerMoves)

	playerCircle = &PlayerCircle{
		posX:   0,
		posY:   0,
		width:  20,
		height: 20,
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
		playerCircle.posY--
		anyKeyPressed = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		playerCircle.posX++
		anyKeyPressed = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		playerCircle.posY++
		anyKeyPressed = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		playerCircle.posX--
		anyKeyPressed = true
	}

	if anyKeyPressed {
		playerMoves <- *playerCircle
	}
}
