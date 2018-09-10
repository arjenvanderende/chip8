package termbox

import (
	"log"
	"sync"

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
	mutex  sync.Mutex
}

func newKeyboard() *keyboard {
	return &keyboard{
		// Create an event buffer with a large enough size to prevent keys
		// from being dropped when too many events are being sent.
		events: make(chan (io.Key), 50),
	}
}

func (k *keyboard) poll() {
	for {
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey:
			if event.Key == termbox.KeyEsc {
				k.mutex.Lock()
				select {
				case k.events <- io.KeyEsc:
				default:
					// Fallback for when channel is (being) closed, but key presses are still received/buffered before
					// shutdown is complete. Otherwise the goroutine might hang trying to send on a full channel and/or
					// panic on a closed channel.
					log.Printf("Keyboard buffer full, dropped key: %v", event.Key)
				}
				k.mutex.Unlock()
			}
		case termbox.EventInterrupt:
			return
		}
	}
}

func (k *keyboard) close() {
	k.mutex.Lock()
	close(k.events)
	k.events = nil
	k.mutex.Unlock()

	// send interrupt to unblock poll()
	termbox.Interrupt()
}

func (k *keyboard) Events() <-chan (io.Key) {
	return k.events
}
