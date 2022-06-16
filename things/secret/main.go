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
		s.showPrompt()
		if _, push := s.board.HandleButtonA(); push {
			if s.guessNext(true) {
				s.showCorrect()
			} else {
				s.showIncorrect()
			}
		}
		if _, push := s.board.HandleButtonB(); push {
			if s.guessNext(false) {
				s.showCorrect()
			} else {
				s.showIncorrect()
			}
		}
		if s.isWin() {
			s.showWin()
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

func (s *secret) guessNext(guess bool) (correct bool) {
	if s.count == 10 {
		return true
	}
	if s.code[s.count] == guess {
		s.count++
		return true
	}
	s.count = 0
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

func (s *secret) showCorrect() {
	colors := s.getColors()
	s.board.SetLights(colors)
	time.Sleep(200 * time.Millisecond)
}

func (s *secret) showIncorrect() {
	colors := s.getColors()
	if s.count < 10 {
		colors[s.count] = color.RGBA{R: 0x01, G: 0x00, B: 0x00}
	}
	s.board.SetLights(colors)
	time.Sleep(200 * time.Millisecond)
}

func (s *secret) showWin() {
	colors := s.getColors()
	s.board.SetLights(colors)
}

func (s *secret) showPrompt() {
	colors := s.getColors()
	s.board.SetLights(colors)
}
