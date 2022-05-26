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
		c.tick()
		time.Sleep(time.Second)
	}
}

type clock struct {
	lights []color.RGBA
	device ws2812.Device
}

func newClock() clock {
	neo := machine.NEOPIXELS
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	device := ws2812.New(neo)
	return clock{
		lights: make([]color.RGBA, 10),
		device: device,
	}
}

func (c clock) tick() {
	t := time.Now()
	s := t.Second()
	for i := range c.lights {
		mask := 1 << i
		if on := s & mask; on > 0 {
			c.lights[i] = color.RGBA{R: 0x01, G: 0x01, B: 0x01}
		} else {
			c.lights[i] = color.RGBA{R: 0x00, G: 0x00, B: 0x00}
		}
	}
	c.device.WriteColors(c.lights)
}
