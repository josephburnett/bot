package main

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/josephburnett/bot/pkg/express"
)

func main() {
	s := newSecret()
	for {
		s.prompt()
		if _, push := s.board.HandleButtonA(); push {
			if s.isNext(true) {
				s.correct()
			} else {
				s.incorrect()
			}
		}
		if _, push := s.board.HandleButtonB(); push {
			if s.isNext(false) {
				s.correct()
			} else {
				s.incorrect()
			}
		}
		if s.isWin() {
			s.win()
			s.newCode()
		}
		time.Sleep(10 * time.Millisecond)
	}
}

type secret struct {
	board *express.Board
	code  [10]bool
	count int
}

func newSecret() *secret {
	s := &secret{
		board: express.NewBoard(),
	}
	s.newCode()
	return s
}

func (s *secret) isNext(guess bool) (correct bool) {
	if s.count == 10 {
		return true
	}
	if s.code[s.count] == guess {
		return true
	}
	return false
}

func (s *secret) isWin() bool {
	return s.count == 10
}

func (s *secret) newCode() {
	rand.Seed(time.Now().UnixNano())
	for i := range s.code {
		r := rand.Intn(2)
		if r == 0 {
			s.code[i] = true
		} else {
			s.code[i] = false
		}
	}
	s.count = 0
}

func (s *secret) getColors() [10]color.RGBA {
	colors := [10]color.RGBA{}
	for i := range s.code {
		if i >= s.count {
			// Unknown
			colors[i] = color.RGBA{R: 0x00, G: 0x00, B: 0x00}
			continue
		}
		// Correct guess
		if s.code[i] {
			colors[i] = color.RGBA{R: 0x00, G: 0x01, B: 0x00}
		} else {
			colors[i] = color.RGBA{R: 0x00, G: 0x00, B: 0x01}
		}
	}
	return colors
}

func (s *secret) correct() {
	if s.count < 10 {
		s.count++
	}
	colors := s.getColors()
	s.board.SetLights(colors)
	time.Sleep(200 * time.Millisecond)
}

func (s *secret) incorrect() {
	for i := s.count - 1; i >= 0; i-- {
		colors := s.getColors()
		colors[i] = color.RGBA{R: 0x01, G: 0x00, B: 0x00}
		s.board.SetLights(colors)
		s.count = i
		time.Sleep(50 * time.Millisecond)
	}
	time.Sleep(500 * time.Millisecond)
}

func (s *secret) win() {
	colors := s.getColors()
	for j := 0; j < 5*len(colors); j++ {
		temp := colors[0]
		for i := 1; i < len(colors); i++ {
			colors[i-1] = colors[i]
		}
		colors[len(colors)-1] = temp
		s.board.SetLights(colors)
		time.Sleep(75 * time.Millisecond)
	}
}

const blinkRate = 500

func (s *secret) prompt() {
	colors := s.getColors()
	t := time.Now().UnixMilli() % blinkRate
	if t > blinkRate/2 {
		t = blinkRate - t
	}
	t = t / (blinkRate / 10)
	colors[s.count] = color.RGBA{R: 0x00, G: uint8(t), B: uint8(t)}
	s.board.SetLights(colors)
}
