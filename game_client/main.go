package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
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
	ThrowableSquares  []*ThrowableSquare
	MouseClickedAt    time.Time
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.HandleKeyInput()
	g.HandleMouseInput()

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
	i := 0
	for otherPlayerID, otherPlayerCircle := range g.OtherPlayers {
		i++
		circle, _ := ebiten.NewImage(otherPlayerCircle.Width, otherPlayerCircle.Height, ebiten.FilterDefault)
		circle.Fill(Colors[otherPlayerCircle.Color])

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(otherPlayerCircle.PosX, otherPlayerCircle.PosY)

		screen.DrawImage(circle, op)

		//Shor other players points
		otherPlayerPoint := fmt.Sprintf("Player %d - %d", otherPlayerID, otherPlayerCircle.Points)
		ebitenutil.DebugPrintAt(screen, otherPlayerPoint, screenWidth-100, (i+1)*10)
	}

	// Draw consumable squares
	for _, consumableSquare := range g.ConsumableSquares {
		square, _ := ebiten.NewImage(consumableSquare.Width, consumableSquare.Height, ebiten.FilterDefault)
		square.Fill(Colors[consumableSquare.Color])

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(consumableSquare.PosX, consumableSquare.PosY)

		screen.DrawImage(square, op)
	}

	// Show current player points
	playerPoints := fmt.Sprintf("Your points %d", g.PlayerCircle.Points)
	ebitenutil.DebugPrintAt(screen, playerPoints, screenWidth-100, 0)

	// Draw throwable squares
	log.Println(g.ThrowableSquares)

	for _, throwableSquare := range g.ThrowableSquares {
		if throwableSquare.CurrentX == throwableSquare.DestX && throwableSquare.CurrentY == throwableSquare.DestY {
			throwableSquare.Reached = true
			g.ThrowableSquares = append(g.ThrowableSquares[:i], g.ThrowableSquares[i+1:]...)
			continue
		}

		square, _ := ebiten.NewImage(throwableSquare.Width, throwableSquare.Height, ebiten.FilterDefault)
		square.Fill(Colors[throwableSquare.Color])

		op := &ebiten.DrawImageOptions{}

		if throwableSquare.CurrentX < throwableSquare.DestX {
			throwableSquare.CurrentX++
		} else if throwableSquare.CurrentX > throwableSquare.DestX {
			throwableSquare.CurrentX--
		}

		if throwableSquare.CurrentY < throwableSquare.DestY {
			throwableSquare.CurrentY++
		} else if throwableSquare.CurrentY > throwableSquare.DestY {
			throwableSquare.CurrentY--
		}

		op.GeoM.Translate(float64(throwableSquare.CurrentX), float64(throwableSquare.CurrentY))
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
		ThrowableSquares: []*ThrowableSquare{},
	}

	go game.startWsClient()

	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("AGARIO")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) HandleKeyInput() {
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

func (g *Game) HandleMouseInput() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if g.MouseClickedAt.IsZero() {
			g.MouseClickedAt = time.Now()
		}

		diff := time.Now().Sub(g.MouseClickedAt).Seconds()

		if diff < 1.0 {
			return
		}

		log.Println("Mouse button pressed")

		g.MouseClickedAt = time.Now()
		destX, destY := ebiten.CursorPosition()
		g.ThrowSquare(int(g.PlayerCircle.PosX), int(g.PlayerCircle.PosY), destX, destY)
	}
}

func (g *Game) ThrowSquare(sourceX, sourceY, destX, destY int) {
	throwableSquare := ThrowableSquare{
		CurrentX: sourceX,
		CurrentY: sourceY,
		DestX:    destX,
		DestY:    destY,
		Reached:  false,
		Width:    5,
		Height:   5,
		Color:    getRandomColor(),
	}

	g.ThrowableSquares = append(g.ThrowableSquares, &throwableSquare)
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

			break
		}
	}
}
