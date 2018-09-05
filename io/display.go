package io

const (
	// DisplayWidth represents the width of the displat in pixels
	DisplayWidth = 64
	// DisplayHeight represents the height of the display in pixels
	DisplayHeight = 32
)

// Display can draw pixels onto the display
type Display interface {
	Clear()
	Draw(x, y int, sprite []byte) bool
	Flush()
}
