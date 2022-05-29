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
	colorOff         = color.RGBA{R: 0x00, G: 0x00, B: 0x00}
	colorMillisHigh  = color.RGBA{R: 0x50, G: 0x00, B: 0x00}
	colorMillisLow   = color.RGBA{R: 0x01, G: 0x00, B: 0x00}
	colorSecondsHigh = color.RGBA{R: 0x50, G: 0x50, B: 0x00}
	colorSecondsLow  = color.RGBA{R: 0x01, G: 0x01, B: 0x00}
	colorMinutesHigh = color.RGBA{R: 0x00, G: 0x50, B: 0x00}
	colorMinutesLow  = color.RGBA{R: 0x00, G: 0x01, B: 0x00}
	colorHoursHigh   = color.RGBA{R: 0x00, G: 0x00, B: 0x50}
	colorHoursLow    = color.RGBA{R: 0x00, G: 0x00, B: 0x01}
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

	set := func(number, ringOffset int, bitRange [2]int, high, low color.RGBA) {
		for i := range lights {
			if i < bitRange[0] || i >= bitRange[1] {
				// Display a subset of the bits.
				continue
			}
			// Read up to 10 bits from the number.
			mask := 1 << i
			// Display the bits at offset along ring.
			index := (i + ringOffset) % len(lights)
			if on := number & mask; on > 0 {
				// Bit is a high (1).
				lights[index] = high
			} else {
				// Bit is a low (0).
				lights[index] = low
			}
		}
	}

	t := time.Now()
	switch c.display {
	case displayMillisAndSeconds:
		millis := t.Nanosecond()/1000/1000 + 1
		seconds := t.Second() + 1
		set(millis, 0, [2]int{6, 10}, colorMillisHigh, colorMillisLow)
		set(seconds, 0, [2]int{0, 6}, colorSecondsHigh, colorSecondsLow)
	case displaySecondsAndMinutes:
		seconds := t.Second() + 1
		minutes := t.Minute() + 1
		set(seconds, 4, [2]int{2, 6}, colorSecondsHigh, colorSecondsLow)
		set(minutes, 0, [2]int{0, 6}, colorMinutesHigh, colorMinutesLow)
	case displayMinutesAndHours:
		minutes := t.Minute() + 1
		hours := t.Hour()%12 + 1
		set(minutes, 4, [2]int{0, 6}, colorMinutesHigh, colorMinutesLow)
		set(hours, 0, [2]int{0, 4}, colorHoursHigh, colorHoursLow)
	}
	c.board.SetLights(lights)
}
