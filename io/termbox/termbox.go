package termbox

import (
	"github.com/arjenvanderende/chip8/io"
	"github.com/nsf/termbox-go"
)

// Closer disposes the display and keyboard that the termbox library initialises
type Closer func()

// New initialises a display and keyboard device via the termbox library
func New() (io.Display, io.Keyboard, Closer, error) {
	err := termbox.Init()
	if err != nil {
		return nil, nil, nil, err
	}

	termbox.SetInputMode(termbox.InputEsc)
	keyboard := &keyboard{
		events: make(chan (io.Key)),
	}
	go keyboard.poll()

	return &display{}, keyboard, func() {
		// release all resources
		keyboard.close()
		termbox.Close()
	}, nil
}
