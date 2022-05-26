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
	top    [5]color.RGBA
	bottom [5]color.RGBA
	device ws2812.Device
}

func newClock() clock {
	neo := machine.NEOPIXELS
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	device := ws2812.New(neo)
	return clock{
		device: device,
	}
}

func (c clock) tick() {
	t := time.Now()
	s := t.Second()
	for i := range c.top {
		mask := 1 << i
		if on := s & mask; on > 0 {
			c.top[i] = color.RGBA{R: 0x01, G: 0x01, B: 0x01}
		} else {
			c.top[i] = color.RGBA{R: 0x00, G: 0x00, B: 0x00}
		}
	}
	var all [10]color.RGBA
	for i := range c.top {
		all[i] = c.top[i]
	}
	c.device.WriteColors(all[0:len(all)])
}
