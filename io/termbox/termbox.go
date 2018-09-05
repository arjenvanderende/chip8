package termbox

import (
	"github.com/arjenvanderende/chip8/io"
	"github.com/nsf/termbox-go"
)

// Termbox represents the I/O devices implemented via termbox-go
type Termbox struct {
	Display  io.Display
	Keyboard io.Keyboard
	keyboard *keyboard
}

// Close finalizes usage of the termbox library
func (tb *Termbox) Close() {
	tb.keyboard.close()
	termbox.Close()
}

// New initialises a display and keyboard device via the termbox library
func New() (*Termbox, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}

	termbox.SetInputMode(termbox.InputEsc)
	keyboard := &keyboard{
		events: make(chan (io.Key)),
	}
	go keyboard.poll()

	return &Termbox{
		Display:  &display{},
		Keyboard: keyboard,
		keyboard: keyboard,
	}, nil
}
