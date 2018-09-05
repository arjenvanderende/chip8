package io

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

const (
	width  = 64
	height = 32
)

// Graphics does something
type Graphics interface {
	Clear()
	Close()
	Draw(x, y int, sprite []byte) bool
	Flush()
}

type tb struct {
	pixels [width][height]bool
}

// NewTermbox does something
func NewTermbox() (Graphics, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}

	return &tb{}, nil
}

func (t *tb) Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func (t *tb) Flush() {
	termbox.Flush()
}

func (t *tb) Draw(x, y int, sprite []byte) bool {
	collision := false
	for dy, line := range sprite {
		fmt.Printf("DRAW X=%d, Y=%d: %08b\n", x, y, line)
		for dx := 0; dx < 8; dx++ {
			// determine if pixel is on or off
			a := line&(1<<uint(7-dx)) > 0
			b := t.pixels[x+dx][y+dy]
			on := a != b

			// collision detection
			if a == b && a == true {
				collision = true
			}

			// draw pixel
			if on {
				termbox.SetCell(x+dx, y+dy, '█', termbox.ColorGreen, termbox.ColorDefault)
			} else {
				termbox.SetCell(x+dx, y+dy, ' ', termbox.ColorDefault, termbox.ColorDefault)
			}
		}
	}
	return collision
}

func (t *tb) Close() {
	termbox.Close()
}
