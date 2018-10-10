package termbox

import (
	"log"
	"sync"
	"unicode"

	"github.com/arjenvanderende/chip8/io"
	"github.com/nsf/termbox-go"
)

type keyboard struct {
	events chan (io.Key)
	mutex  sync.Mutex
}

// Maps keyboard to the following layout:
// 1 2 3 C
// 4 5 6 D
// 7 8 9 E
// A 0 B F
var keyMapping = map[rune]io.Key{
	rune('1'): io.Key1, rune('2'): io.Key2, rune('3'): io.Key3, rune('4'): io.KeyC,
	rune('q'): io.Key4, rune('w'): io.Key5, rune('e'): io.Key6, rune('r'): io.KeyD,
	rune('a'): io.Key7, rune('s'): io.Key8, rune('d'): io.Key9, rune('f'): io.KeyE,
	rune('z'): io.KeyA, rune('x'): io.Key0, rune('c'): io.KeyB, rune('v'): io.KeyF,
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
			if key, ok := keyMapping[unicode.ToLower(event.Ch)]; ok {
				k.sendKeyPress(key)
			} else if event.Key == termbox.KeyEsc {
				k.sendKeyPress(io.KeyEsc)
			}
		case termbox.EventInterrupt:
			return
		}
	}
}

func (k *keyboard) sendKeyPress(key io.Key) {
	k.mutex.Lock()
	select {
	case k.events <- key:
	default:
		// Fallback for when channel is (being) closed, but key presses are still received/buffered before
		// shutdown is complete. Otherwise the goroutine might hang trying to send on a full channel and/or
		// panic on a closed channel.
		log.Printf("Keyboard buffer full, dropped key: %v", key)
	}
	k.mutex.Unlock()
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
