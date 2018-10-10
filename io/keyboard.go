package io

type Keyboard interface {
	Events() <-chan (Key)
}

type Key int

const (
	KeyEsc Key = iota
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
)

// KeyValue returns the hexadecimal value of the key
func KeyValue(k Key) byte {
	switch k {
	case Key0:
		return 0x0
	case Key1:
		return 0x1
	case Key2:
		return 0x2
	case Key3:
		return 0x3
	case Key4:
		return 0x4
	case Key5:
		return 0x5
	case Key6:
		return 0x6
	case Key7:
		return 0x7
	case Key8:
		return 0x8
	case Key9:
		return 0x9
	case KeyA:
		return 0xA
	case KeyB:
		return 0xB
	case KeyC:
		return 0xC
	case KeyD:
		return 0xD
	case KeyE:
		return 0xE
	case KeyF:
		return 0xF
	default:
		return 0x0
	}
}
