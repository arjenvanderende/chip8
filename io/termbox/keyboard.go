package termbox

import (
	"sync"
	"unicode"

	"github.com/arjenvanderende/chip8/io"
	"github.com/nsf/termbox-go"
)

const keyPressDuration uint8 = 1

type keyboard struct {
	pressedKeys  map[io.Key]uint8
	waitForPress chan io.Key
	mutex        sync.RWMutex
}

func newKeyboard() *keyboard {
	return &keyboard{
		pressedKeys:  make(map[io.Key]uint8),
		waitForPress: nil,
	}
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

func (k *keyboard) poll() {
	for {
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey:
			if key, ok := keyMapping[unicode.ToLower(event.Ch)]; ok {
				k.registerKeyPress(key)
			} else if event.Key == termbox.KeyEsc {
				k.registerKeyPress(io.KeyEsc)
			}
		case termbox.EventInterrupt:
			return
		}
	}
}

func (k *keyboard) registerKeyPress(key io.Key) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	k.pressedKeys[key] = keyPressDuration

	// emit event if WaitForKeyPress was invoked
	if k.waitForPress != nil {
		k.waitForPress <- key
		close(k.waitForPress)
		k.waitForPress = nil
	}
}

func (k *keyboard) close() {
	// send interrupt to unblock poll()
	termbox.Interrupt()
}

func (k *keyboard) Tick() {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	for key := range k.pressedKeys {
		if k.pressedKeys[key] <= 0 {
			delete(k.pressedKeys, key)
		} else {
			k.pressedKeys[key] = k.pressedKeys[key] - 1
		}
	}
}

func (k *keyboard) PressedButton() *io.Key {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	for key := range k.pressedKeys {
		return &key
	}
	return nil
}

func (k *keyboard) IsPressed(key io.Key) bool {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	_, ok := k.pressedKeys[key]
	return ok
}
