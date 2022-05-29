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
	displayMillisAndSeconds = iota
	displaySecondsAndMinutes
	displayMinutesAndHours
)

var (
	colorOff     = color.RGBA{R: 0x00, G: 0x00, B: 0x00}
	colorMillis  = color.RGBA{R: 0x01, G: 0x00, B: 0x00}
	colorSeconds = color.RGBA{R: 0x01, G: 0x01, B: 0x00}
	colorMinutes = color.RGBA{R: 0x00, G: 0x01, B: 0x00}
	colorHours   = color.RGBA{R: 0x00, G: 0x00, B: 0x01}
)

type clock struct {
	display int
	board   *express.Board
}

func newClock() *clock {
	return &clock{
		display: displayMillisAndSeconds,
		board:   express.NewBoard(),
	}
}

func (c *clock) readButtons() {
	if _, push := c.board.HandleButtonA(); push {
		if c.display < displayMinutesAndHours {
			c.display++
		}
	}
	if _, push := c.board.HandleButtonB(); push {
		if c.display > displayMillisAndSeconds {
			c.display--
		}
	}
}

func (c *clock) displayTime() {

	var lights [10]color.RGBA
	for i := range lights {
		lights[i] = colorOff
	}

	set := func(number, offset int, lightsRange [2]int, color color.RGBA) {
		for i := range lights {
			if i < lightsRange[0] || i >= lightsRange[1] {
				continue
			}
			mask := 1 << i
			index := (i + offset) % len(lights)
			if on := number & mask; on > 0 {
				lights[index] = color
			} else {
				lights[index] = colorOff
			}
		}
	}

	t := time.Now()
	switch c.display {
	case displayMillisAndSeconds:
		millis := t.Nanosecond() / 1000 / 1000
		seconds := t.Second()
		set(millis, 0, [2]int{6, 10}, colorMillis)
		set(seconds, 0, [2]int{0, 5}, colorSeconds)
	}
	c.board.SetLights(lights)
}
