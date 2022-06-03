package main

import (
	"image/color"
	"time"

	"github.com/josephburnett/bot/pkg/express"
)

func main() {
	g := newGame()
	for {
		g.tick()
		time.Sleep(10 * time.Millisecond)
	}
}

const (
	decelerationRate = 0.005 // Rate of ball's deceleration per tick
	minSpeed         = 0.003 // Minimum speed of the ball
	maxSpeed         = 0.05  // Maximum speed of the ball
	smackPower       = 0.01  // Max speed added by a smack
	smackPointA      = 0.25  // Player A smacking clockwise
	smackPointB      = 0.75  // Player B smacking anti-clockwise
	smackRange       = 0.05  // Range of the players' rackets +/-
)

type game struct {
	board *express.Board
	point float64 // 0-1 around the board
	speed float64 // Current speed of the ball. Positive is anti-clockwise
}

func newGame() *game {
	return &game{
		board: express.NewBoard(),
		point: 0,
		speed: 0.05,
	}
}

func (g *game) tick() {
	g.slowDownBall()
	g.advanceBall()
	g.smackBall()
	g.updateDisplay()
}

func (g *game) slowDownBall() {
	g.speed = g.speed * (1 - decelerationRate)
	if g.speed > 0 {
		if g.speed < minSpeed {
			g.speed = minSpeed
		}
		if g.speed > maxSpeed {
			g.speed = maxSpeed
		}
	}
	if g.speed < 0 {
		if g.speed > -minSpeed {
			g.speed = -minSpeed
		}
		if g.speed < -maxSpeed {
			g.speed = -maxSpeed
		}
	}
}

func (g *game) advanceBall() {
	g.point += g.speed
	if g.point > 1 {
		g.point -= 1
	}
	if g.point < 0 {
		g.point += 1
	}
}

func (g *game) smackBall() {
	if _, push := g.board.HandleButtonA(); push {
		// Player A swings their racket
		from := smackPointA - smackRange*2
		to := smackPointA + smackRange
		if g.point > from && g.point < to {
			// Player A smacks the ball
			if g.speed > 0 {
				// Ball changes direction
				g.speed *= -1
			}
			// Ball picks up speed clockwise
			g.speed -= smackPower
		}
	}
	if _, push := g.board.HandleButtonB(); push {
		// Player B swings their racket
		from := smackPointB - smackRange
		to := smackPointB + smackRange*2
		if g.point > from && g.point < to {
			// Play B smacks the ball
			if g.speed < 0 {
				// Ball changes direction
				g.speed *= -1
			}
			// Ball picks up speed clockwise
			g.speed += smackPower
		}
	}
}

var segmentEnds = [12]float64{}

func init() {
	// Board is dipvided into 12 segments between 0 and 1
	segmentSize := 1.0 / 12
	var s float64
	for i := range segmentEnds {
		s += segmentSize
		segmentEnds[i] = s
	}
}

func (g *game) updateDisplay() {
	var light int
	for i, s := range segmentEnds {
		if g.point < s {
			light = i
			break
		}
	}
	lights := [10]color.RGBA{}
	// Player A smack point
	lights[2] = color.RGBA{B: 0x01}
	// Player B smack point
	lights[7] = color.RGBA{G: 0x01}
	// USB is where the first light should be
	if light > 0 && light < 6 {
		lights[light-1] = color.RGBA{R: 0x01}
	}
	// Power is where the sixth light should be
	if light > 6 {
		lights[light-2] = color.RGBA{R: 0x01}
	}
	g.board.SetLights(lights)
}
