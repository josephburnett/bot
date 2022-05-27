package main

import (
	"image/color"
	"time"

	"github.com/josephburnett/bot/pkg/express"
)

func main() {
	c := newClock()
	for {
		c.readButtons()
		c.displayTime()
		time.Sleep(time.Millisecond)
	}
}

const (
	milliseconds = iota
	seconds
	minutes
	hours
)

type clock struct {
	mode  int
	board *express.Board
}

func newClock() *clock {
	return &clock{
		mode:  seconds,
		board: express.NewBoard(),
	}
}

func (c *clock) readButtons() {
	if _, push := c.board.HandleButtonA(); push {
		if c.mode < hours {
			c.mode++
		}
	}
	if _, push := c.board.HandleButtonB(); push {
		if c.mode > milliseconds {
			c.mode--
		}
	}
}

func (c *clock) displayTime() {
	t := time.Now()
	var displayNumber int
	var displayColor color.RGBA

	switch c.mode {
	case milliseconds:
		displayNumber = t.Nanosecond()/1000/1000 + 1
		displayColor = color.RGBA{R: 0x01, G: 0x00, B: 0x00}
	case seconds:
		displayNumber = t.Second() + 1
		displayColor = color.RGBA{R: 0x01, G: 0x01, B: 0x00}
	case minutes:
		displayNumber = t.Minute() + 1
		displayColor = color.RGBA{R: 0x00, G: 0x01, B: 0x00}
	case hours:
		displayColor = color.RGBA{R: 0x00, G: 0x00, B: 0x01}
		displayNumber = t.Hour() + 1
	}

	var lights [10]color.RGBA
	for i := range lights {
		mask := 1 << i
		if on := displayNumber & mask; on > 0 {
			lights[i] = displayColor
		} else {
			lights[i] = color.RGBA{R: 0x00, G: 0x00, B: 0x00}
		}
	}
	c.board.SetLights(lights)
}
