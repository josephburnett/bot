package express

import (
	"image/color"
	"machine"
	"sync"

	"tinygo.org/x/drivers/ws2812"
)

type Board struct {
	mux         sync.Mutex
	buttonA     machine.Pin
	buttonALast bool
	buttonB     machine.Pin
	buttonBLast bool
	lights      ws2812.Device
	colors      []color.RGBA
}

func NewBoard() *Board {

	buttonA := machine.BUTTONA
	buttonA.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	buttonB := machine.BUTTONB
	buttonB.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	neo := machine.NEOPIXELS
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	lights := ws2812.New(neo)

	return &Board{
		buttonA: buttonA,
		buttonB: buttonB,
		lights:  lights,
		colors:  make([]color.RGBA, 10),
	}
}

func (b *Board) SetLights(colors [10]color.RGBA) {
	b.mux.Lock()
	defer b.mux.Unlock()
	for i := range colors {
		b.colors[i] = colors[i]
	}
	b.lights.WriteColors(b.colors)
}

func (b *Board) HandleButtonA() (down, push bool) {
	down = b.buttonA.Get()
	if down && !b.buttonALast {
		// First time
		push = true
	}
	b.buttonALast = down
	return
}

func (b *Board) HandleButtonB() (down, push bool) {
	down = b.buttonB.Get()
	if down && !b.buttonBLast {
		// First time
		push = true
	}
	b.buttonBLast = down
	return
}
