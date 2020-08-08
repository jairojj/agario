package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	OtherPlayers      map[int]PlayerCircle
	ConsumableSquares []ConsumableSquare
	MessageQueue      chan Message
	CurrentPlayerID   int
	PlayerCircle      *PlayerCircle
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.HandleInput()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw current player
	circle, _ := ebiten.NewImage(g.PlayerCircle.Width, g.PlayerCircle.Height, ebiten.FilterDefault)
	circle.Fill(Colors[g.PlayerCircle.Color])

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.PlayerCircle.PosX, g.PlayerCircle.PosY)

	screen.DrawImage(circle, op)

	// Draw other players
	for _, otherPlayerCircle := range g.OtherPlayers {
		circle, _ := ebiten.NewImage(otherPlayerCircle.Width, otherPlayerCircle.Height, ebiten.FilterDefault)
		circle.Fill(Colors[otherPlayerCircle.Color])

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(otherPlayerCircle.PosX, otherPlayerCircle.PosY)

		screen.DrawImage(circle, op)
	}

	// Draw consumable squares
	for _, consumableSquare := range g.ConsumableSquares {
		square, _ := ebiten.NewImage(consumableSquare.Width, consumableSquare.Height, ebiten.FilterDefault)
		square.Fill(Colors[consumableSquare.Color])

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(consumableSquare.PosX, consumableSquare.PosY)

		screen.DrawImage(square, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		OtherPlayers:      map[int]PlayerCircle{},
		ConsumableSquares: []ConsumableSquare{},
		MessageQueue:      make(chan Message),
		CurrentPlayerID:   rand.Intn(100),
		PlayerCircle: &PlayerCircle{
			PosX:   0,
			PosY:   0,
			Width:  20,
			Height: 20,
			Color:  getRandomColor(),
			Points: 0,
		},
	}

	go game.startWsClient()

	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("AGARIO")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) HandleInput() {
	hasPlayerMoved := false

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		if g.PlayerCircle.PosY <= 0 {
			return
		}

		g.PlayerCircle.PosY--
		hasPlayerMoved = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		if g.PlayerCircle.PosX+float64(g.PlayerCircle.Width) > screenWidth {
			return
		}

		g.PlayerCircle.PosX++
		hasPlayerMoved = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		if g.PlayerCircle.PosY+float64(g.PlayerCircle.Height) > screenHeight {
			return
		}

		g.PlayerCircle.PosY++
		hasPlayerMoved = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		if g.PlayerCircle.PosX <= 0 {
			return
		}

		g.PlayerCircle.PosX--
		hasPlayerMoved = true
	}

	if hasPlayerMoved {
		message := Message{ClientID: g.CurrentPlayerID, Event: PlayerMoved, PlayerCircle: *g.PlayerCircle}
		g.MessageQueue <- message

		g.HandleCollision()
	}
}

func (g *Game) HandleCollision() {
	for i, consumableSquare := range g.ConsumableSquares {
		if (consumableSquare.PosX >= g.PlayerCircle.PosX &&
			consumableSquare.PosX <= g.PlayerCircle.PosX+float64(g.PlayerCircle.Width)) &&
			(consumableSquare.PosY >= g.PlayerCircle.PosY &&
				consumableSquare.PosY <= g.PlayerCircle.PosY+float64(g.PlayerCircle.Height)) {
			log.Println("Square consumabled: ", consumableSquare)

			g.ConsumableSquares = append(g.ConsumableSquares[0:i], g.ConsumableSquares[i+1:]...)
			g.PlayerCircle.Height += consumableSquare.Height
			g.PlayerCircle.Width += consumableSquare.Width
			g.PlayerCircle.Points++

			message := Message{ClientID: g.CurrentPlayerID, Event: PlayerMoved, PlayerCircle: *g.PlayerCircle}
			g.MessageQueue <- message

			message = Message{ClientID: g.CurrentPlayerID, Event: ConsumableSquareChanged, ConsumableSquares: g.ConsumableSquares}
			g.MessageQueue <- message
		}
	}
}
