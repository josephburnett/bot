package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ws2812"
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
	mode         int
	click        bool
	topButton    machine.Pin
	bottomButton machine.Pin
	device       ws2812.Device
	lights       []color.RGBA
}

func newClock() *clock {

	neo := machine.NEOPIXELS
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	device := ws2812.New(neo)

	top := machine.BUTTONA
	top.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	bottom := machine.BUTTONB
	bottom.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	return &clock{
		mode:         seconds,
		topButton:    top,
		bottomButton: bottom,
		device:       device,
		lights:       make([]color.RGBA, 10),
	}
}

func (c *clock) readButtons() {
	if c.click {
		if c.bottomButton.Get() || c.topButton.Get() {
			return
		}
	}
	c.click = false

	if c.topButton.Get() {
		c.click = true
		if c.mode < hours {
			c.mode++
		}
	}
	if c.bottomButton.Get() {
		c.click = true
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

	for i := range c.lights {
		mask := 1 << i
		if on := displayNumber & mask; on > 0 {
			c.lights[i] = displayColor
		} else {
			c.lights[i] = color.RGBA{R: 0x00, G: 0x00, B: 0x00}
		}
	}
	c.device.WriteColors(c.lights)
}
