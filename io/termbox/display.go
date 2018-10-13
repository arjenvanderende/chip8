package termbox

import (
	"log"

	"github.com/arjenvanderende/chip8/io"
	"github.com/nsf/termbox-go"
)

type display struct {
	pixels [io.DisplayWidth * io.DisplayHeight]bool
}

func (s *display) Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	s.pixels = [io.DisplayWidth * io.DisplayHeight]bool{}
}

func (s *display) Flush() {
	termbox.Flush()
}

func (s *display) Draw(x, y int, sprite []byte) bool {
	collision := false
	for dy, line := range sprite {
		log.Printf("DRAW X=%d, Y=%d: %08b\n", x, y, line)
		for dx := 0; dx < 8; dx++ {
			// determine if pixel is on or off
			p := (((y + dy) * io.DisplayWidth) + x + dx) % (io.DisplayWidth * io.DisplayHeight)
			a := line&(1<<uint(7-dx)) > 0
			b := s.pixels[p]
			on := a != b

			// collision detection
			if a == b && a == true {
				collision = true
			}

			// draw pixel
			rx := p % io.DisplayWidth
			ry := p / io.DisplayWidth
			if on {
				termbox.SetCell(rx, ry, '█', termbox.ColorGreen, termbox.ColorDefault)
			} else {
				termbox.SetCell(rx, ry, ' ', termbox.ColorDefault, termbox.ColorDefault)
			}

			// remember the state
			s.pixels[p] = on
		}
	}
	return collision
}
