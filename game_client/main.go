package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type PlayerCircle struct {
	PosX   float64 `json:"pos_x"`
	PosY   float64 `json:"pos_y"`
	Height int     `json:"height"`
	Width  int     `json:"width"`
	Color  string  `json:"color"`
}

func (p PlayerCircle) String() string {
	jsonPlayerCircle, _ := json.Marshal(p)
	return fmt.Sprint(string(jsonPlayerCircle))
}

type Game struct {
	OtherPlayers map[int]PlayerCircle
}

var playerCircle *PlayerCircle
var playerMoves chan PlayerCircle

func (g *Game) Update(screen *ebiten.Image) error {
	handleInput()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw current player
	circle, _ := ebiten.NewImage(playerCircle.Width, playerCircle.Height, ebiten.FilterDefault)
	circle.Fill(Colors[playerCircle.Color])

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(playerCircle.PosX, playerCircle.PosY)

	screen.DrawImage(circle, op)

	// Draw other players
	for _, otherPlayerCircle := range g.OtherPlayers {
		circle, _ := ebiten.NewImage(otherPlayerCircle.Width, otherPlayerCircle.Height, ebiten.FilterDefault)
		circle.Fill(Colors[otherPlayerCircle.Color])

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(otherPlayerCircle.PosX, otherPlayerCircle.PosY)

		screen.DrawImage(circle, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	playerMoves = make(chan PlayerCircle)

	game := &Game{OtherPlayers: map[int]PlayerCircle{}}

	randomColor := getRandomColor()

	go startWsClient(playerMoves, game, randomColor)

	playerCircle = &PlayerCircle{
		PosX:   0,
		PosY:   0,
		Width:  20,
		Height: 20,
		Color:  randomColor,
	}

	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("AGARIO")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func handleInput() {
	hasPlayerMoved := false

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		if playerCircle.PosY <= 0 {
			return
		}

		playerCircle.PosY--
		hasPlayerMoved = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		if playerCircle.PosX+float64(playerCircle.Width) > screenWidth {
			return
		}

		playerCircle.PosX++
		hasPlayerMoved = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		if playerCircle.PosY+float64(playerCircle.Height) > screenHeight {
			return
		}

		playerCircle.PosY++
		hasPlayerMoved = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		if playerCircle.PosX <= 0 {
			return
		}

		playerCircle.PosX--
		hasPlayerMoved = true
	}

	if hasPlayerMoved {
		playerMoves <- *playerCircle
	}
}
