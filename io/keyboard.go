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
