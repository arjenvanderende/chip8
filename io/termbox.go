package io

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

// Graphics does something
type Graphics interface {
	Draw(number int)
	Close()
}

type tb struct {
	width  int
	height int
}

// NewTermbox does something
func NewTermbox() (Graphics, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}

	return &tb{
		width:  64,
		height: 32,
	}, nil
}

func (t *tb) Draw(number int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// draw the border
	for y := 0; y < t.height; y++ {
		for x := 0; x < t.width; x++ {
			rune := ' '
			if y == 0 || y == t.height-1 || x == 0 || x == t.width-1 {
				rune = '█'
			}
			termbox.SetCell(x, y, rune, termbox.ColorDefault, termbox.ColorDefault)
		}
	}

	// draw the text
	text := fmt.Sprintf("%d", number)
	for pos, char := range text {
		termbox.SetCell(2+pos, 2, char, termbox.ColorDefault, termbox.ColorDefault)
	}

	termbox.Flush()
}

func (t *tb) Close() {
	termbox.Close()
}
