package termbox

import (
	"github.com/arjenvanderende/chip8/io"
	"github.com/nsf/termbox-go"
)

// Maps keyboard to the following layout
// 1 2 3 C
// 4 5 6 D
// 7 8 9 E
// A 0 B F
type keyboard struct {
	events chan (io.Key)
}

func (k *keyboard) poll() {
	for {
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey:
			if event.Key == termbox.KeyEsc {
				k.events <- io.KeyEsc
			}
		case termbox.EventInterrupt:
			return
		}
	}
}

func (k *keyboard) close() {
	termbox.Interrupt()
}

func (k *keyboard) Events() <-chan (io.Key) {
	return k.events
}
