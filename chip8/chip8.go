package chip8

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/arjenvanderende/chip8/io"
)

const (
	// clockRate represents the number of operations that the CPU can process per second
	clockRate int = 540
	// programOffset represents the offset in memory where the program is loaded
	programOffset int = 0x200
)

var (
	rnd    = rand.New(rand.NewSource(time.Now().UnixNano()))
	digits = []byte{
		0xf0, 0x90, 0x90, 0x90, 0xf0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xf0, 0x10, 0xf0, 0x80, 0xf0, // 2,
		0xf0, 0x10, 0xf0, 0x10, 0xf0, // 3
		0x90, 0x90, 0xf0, 0x10, 0x10, // 4
		0xf0, 0x80, 0xf0, 0x10, 0xf0, // 5
		0xf0, 0x80, 0xf0, 0x90, 0xf0, // 6
		0xf0, 0x10, 0x20, 0x40, 0x40, // 7
		0xf0, 0x90, 0xf0, 0x90, 0xf0, // 8
		0xf0, 0x90, 0xf0, 0x10, 0xf0, // 9
		0xf0, 0x90, 0xf0, 0x90, 0x90, // A
		0xe0, 0x90, 0xe0, 0x90, 0xe0, // B
		0xf0, 0x80, 0x80, 0x80, 0xf0, // C
		0xe0, 0x90, 0x90, 0x90, 0xe0, // D
		0xf0, 0x80, 0xf0, 0x80, 0xf0, // E
		0xf0, 0x80, 0xf0, 0x80, 0x80, // F
	}
)

// Memory represents the memory address space of the Chip-8
type Memory [0x1000]byte

// CPU represents the Chip8 CPU
type CPU struct {
	pc     int // program counter
	memory Memory
	i      uint16   // 16-bit register
	v      [16]byte // 8-bit general purpose registers
	sp     uint8    // stack pointer
	stack  [16]int

	programSize int
}

// Load reads the program stored in the file into memory
func Load(filename string) (*CPU, error) {
	// read ROM file
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Unable to load Chip8 file %s: %v", filename, err)
	}

	// copy ROM into memory at program address
	cpu := CPU{
		pc:          programOffset,
		programSize: len(bytes),
		i:           0,
		v:           [16]byte{},
		sp:          0,
	}
	// copy digits for op: Fx29
	for i, b := range digits {
		cpu.memory[i] = b
	}
	// copy program
	for i, b := range bytes {
		cpu.memory[programOffset+i] = b
	}
	return &cpu, nil
}

// Run starts running the program
func (cpu *CPU) Run(display io.Display, keyboard io.Keyboard) error {
	clock := time.NewTicker(time.Second / time.Duration(clockRate))
	defer clock.Stop()

	frame := time.NewTicker(time.Second / time.Duration(60))
	defer frame.Stop()

	for {
		select {
		// // possible alternative
		// case <-keyboard.Quit():
		// 	return nil
		case <-clock.C:
			// run the next tick of the program
			err := cpu.interpret(display, keyboard)
			if err != nil {
				return fmt.Errorf("Could not interpret op: %v", err)
			}

			// check if the user tried to quit the program
			if keyboard.IsPressed(io.KeyEsc) {
				return nil
			}
			keyboard.Tick()
		case <-frame.C:
			display.Flush()
		}
	}
}

func (cpu *CPU) printState(pc int, op string) {
	log.Printf("op=%-40s pc=%03x next pc=%03x i=%03x v=%v\n", op, pc, cpu.pc, cpu.i, cpu.v)
}

func (cpu *CPU) interpret(display io.Display, keyboard io.Keyboard) error {
	op := cpu.DisassembleOp()
	defer cpu.printState(cpu.pc, op)

	nib1 := cpu.memory[cpu.pc] >> 4
	vx := cpu.memory[cpu.pc] & 0x0f
	vy := cpu.memory[cpu.pc+1] >> 4
	n := cpu.memory[cpu.pc+1] & 0x0f
	nn := cpu.memory[cpu.pc+1]
	nnn := uint16(cpu.memory[cpu.pc]&0x0f)<<8 + uint16(cpu.memory[cpu.pc+1])

	switch nib1 {
	case 0x0:
		switch cpu.memory[cpu.pc+1] {
		case 0xe0:
			display.Clear()
		case 0xee:
			cpu.sp--
			cpu.pc = cpu.stack[cpu.sp]
		default:
			return fmt.Errorf("Unknown 0")
		}
	case 0x1:
		cpu.pc = int(nnn)
		return nil
	case 0x2:
		cpu.stack[cpu.sp] = cpu.pc
		cpu.sp++
		cpu.pc = int(nnn)
		return nil
	case 0x3:
		if cpu.v[vx] == nn {
			cpu.pc += 2
		}
	case 0x4:
		if cpu.v[vx] != nn {
			cpu.pc += 2
		}
	case 0x6:
		cpu.v[vx] = nn
	case 0x7:
		cpu.v[vx] = cpu.v[vx] + nn
	case 0x8:
		lastNib := cpu.memory[cpu.pc+1] & 0x0f
		switch lastNib {
		case 0x0:
			cpu.v[vx] = cpu.v[vy]
		case 0x1:
			cpu.v[vx] = cpu.v[vx] | cpu.v[vy]
		case 0x2:
			cpu.v[vx] = cpu.v[vx] & cpu.v[vy]
		case 0x3:
			cpu.v[vx] = cpu.v[vx] ^ cpu.v[vy]
		case 0x4:
			// set carry flag
			acc := int16(cpu.v[vx]) + int16(cpu.v[vy])
			if acc > 255 {
				cpu.v[0xf] = 1
			} else {
				cpu.v[0xf] = 0
			}
			cpu.v[vx] = byte(acc)
		case 0x5:
			// set borrow flag
			if cpu.v[vx] > cpu.v[vy] {
				cpu.v[0xf] = 1
			} else {
				cpu.v[0xf] = 0
			}
			cpu.v[vx] = cpu.v[vx] - cpu.v[vy]
		case 0x6:
			if cpu.v[vx]&0x1 > 0 {
				cpu.v[0xf] = 1
			} else {
				cpu.v[0xf] = 0
			}
			cpu.v[vx] = cpu.v[vx] / 2
		case 0x7:
			// set borrow flag
			if cpu.v[vy] > cpu.v[vx] {
				cpu.v[0xf] = 1
			} else {
				cpu.v[0xf] = 0
			}
			cpu.v[vx] = cpu.v[vy] - cpu.v[vx]
		case 0xe:
			if cpu.v[vx]&0x80 > 0 {
				cpu.v[0xf] = 1
			} else {
				cpu.v[0xf] = 0
			}
			cpu.v[vx] = cpu.v[vx] * 2
		default:
			return fmt.Errorf("Unknown 8: %1x", lastNib)
		}
	case 0xa:
		cpu.i = nnn
	case 0xc:
		cpu.v[vx] = byte(rnd.Intn(256)) & nn
	case 0xd:
		x := int(cpu.v[vx])
		y := int(cpu.v[vy])
		sprite := cpu.memory[cpu.i : cpu.i+uint16(n)]
		collision := display.Draw(x, y, sprite)
		if collision {
			cpu.v[0xf] = 0x1
		} else {
			cpu.v[0xf] = 0x0
		}
	case 0xe:
		switch cpu.memory[cpu.pc+1] {
		case 0x9e:
			if keyboard.IsPressed(io.Key(cpu.v[vx])) {
				cpu.pc += 2
			}
		case 0xa1:
			if !keyboard.IsPressed(io.Key(cpu.v[vx])) {
				cpu.pc += 2
			}
		default:
			return fmt.Errorf("Unknown E: %2x", cpu.memory[cpu.pc+1])
		}
	case 0xf:
		switch cpu.memory[cpu.pc+1] {
		case 0x0a:
			key := keyboard.PressedButton()
			if key == nil || io.IsOperationalKey(*key) {
				// Skip processing the op at the current PC, allow the emulator
				// to process the operational key and let it loop to the same
				// op to start waiting for a key press again
				return nil
			}
			cpu.v[vx] = byte(*key)
		case 0x1e:
			cpu.i += uint16(cpu.v[vx])
		case 0x29:
			cpu.i = uint16(cpu.v[vx]) * 5
		case 0x33:
			v := uint16(cpu.v[vx])
			cpu.memory[cpu.i+0] = byte((v / 100) % 10)
			cpu.memory[cpu.i+1] = byte((v / 10) % 10)
			cpu.memory[cpu.i+2] = byte(v % 10)
		case 0x55:
			for i := uint16(0); i <= uint16(vx); i++ {
				cpu.memory[cpu.i+i] = cpu.v[i]
			}
		case 0x65:
			for i := uint16(0); i <= uint16(vx); i++ {
				cpu.v[i] = cpu.memory[cpu.i+i]
			}
		default:
			return fmt.Errorf("Unknown F: %2x", cpu.memory[cpu.pc+1])
		}
	default:
		return fmt.Errorf("Unknown nib: %d", nib1)
	}

	cpu.pc += 2
	return nil
}

// NextOp increments the PC to the next operation
// Returns false when there are no more operations to read
func (cpu *CPU) NextOp() bool {
	cpu.pc += 2
	return cpu.pc <= programOffset+cpu.programSize
}

// DisassembleOp output the assembly for the operation at the PC.
func (cpu *CPU) DisassembleOp() string {
	nib1 := cpu.memory[cpu.pc] >> 4

	vx := cpu.memory[cpu.pc] & 0x0f
	vy := cpu.memory[cpu.pc+1] >> 4
	n := cpu.memory[cpu.pc+1] & 0x0f
	nn := cpu.memory[cpu.pc+1]
	nnn := int16(cpu.memory[cpu.pc]&0x0f)<<8 + int16(cpu.memory[cpu.pc+1])

	op := "not implemented"
	switch nib1 {
	case 0x0:
		switch cpu.memory[cpu.pc+1] {
		case 0xe0:
			op = fmt.Sprintf("%-10s", "CLS")
		case 0xee:
			op = fmt.Sprintf("%-10s", "RET")
		default:
			op = fmt.Sprintf("%-10s %03x", "SYS", nnn)
		}
	case 0x1:
		op = fmt.Sprintf("%-10s %03x", "JP", nnn)
	case 0x2:
		op = fmt.Sprintf("%-10s %03x", "CALL", nnn)
	case 0x3:
		op = fmt.Sprintf("%-10s V%01x, %02x", "SE", vx, nn)
	case 0x4:
		op = fmt.Sprintf("%-10s V%01x, %02x", "SNE", vx, nn)
	case 0x5:
		op = fmt.Sprintf("%-10s V%01x, V%01x", "SE", vx, vy)
	case 0x6:
		op = fmt.Sprintf("%-10s V%01x, %02x", "LD", vx, nn)
	case 0x7:
		op = fmt.Sprintf("%-10s V%01x, %02x", "ADD", vx, nn)
	case 0x8:
		lastNib := cpu.memory[cpu.pc+1] & 0x0f
		switch lastNib {
		case 0x0:
			op = fmt.Sprintf("%-10s V%01x, V%01x", "LD", vx, vy)
		case 0x1:
			op = fmt.Sprintf("%-10s V%01x, V%01x", "OR", vx, vy)
		case 0x2:
			op = fmt.Sprintf("%-10s V%01x, V%01x", "AND", vx, vy)
		case 0x3:
			op = fmt.Sprintf("%-10s V%01x, V%01x", "XOR", vx, vy)
		case 0x4:
			op = fmt.Sprintf("%-10s V%01x, V%01x", "ADD", vx, vy)
		case 0x5:
			op = fmt.Sprintf("%-10s V%01x, V%01x, V%01x", "SUB", vx, vx, vy)
		case 0x6:
			op = fmt.Sprintf("%-10s V%01x, V%01x", "SHR", vx, vy)
		case 0x7:
			op = fmt.Sprintf("%-10s V%01x, V%01x, V%01x", "SUBN", vx, vy, vy)
		case 0xe:
			op = fmt.Sprintf("%-10s V%01x, V%01x", "SHL", vx, vy)
		default:
			op = fmt.Sprintf("UNKNOWN 8")
		}
	case 0x9:
		op = fmt.Sprintf("%-10s V%01x, V%01x", "SNE", vx, vy)
	case 0xa:
		op = fmt.Sprintf("%-10s I,%03x", "LD", nnn)
	case 0xb:
		op = fmt.Sprintf("%-10s V0,%03x", "JP", nnn)
	case 0xc:
		op = fmt.Sprintf("%-10s V%01x, %02x", "RND", vx, nn)
	case 0xd:
		op = fmt.Sprintf("%-10s V%01x, V%01x, %01x", "DRW", vx, vy, n)
	case 0xe:
		switch cpu.memory[cpu.pc+1] {
		case 0x9e:
			op = fmt.Sprintf("%-10s V%01x", "SKP", vx)
		case 0xa1:
			op = fmt.Sprintf("%-10s V%01x", "SKNP", vx)
		default:
			op = fmt.Sprintf("UNKNOWN E")
		}
	case 0xf:
		switch cpu.memory[cpu.pc+1] {
		case 0x07:
			op = fmt.Sprintf("%-10s V%01x, DELAY", "LD", vx)
		case 0x0a:
			op = fmt.Sprintf("%-10s V%01x, KEY", "LD", vx)
		case 0x15:
			op = fmt.Sprintf("%-10s DELAY, V%01x", "LD", vx)
		case 0x18:
			op = fmt.Sprintf("%-10s SOUND, V%01x", "LD", vx)
		case 0x1e:
			op = fmt.Sprintf("%-10s I, V%01x", "ADD", vx)
		case 0x29:
			op = fmt.Sprintf("%-10s F, V%01x", "LD", vx)
		case 0x33:
			op = fmt.Sprintf("%-10s B, V%01x", "LD", vx)
		case 0x55:
			op = fmt.Sprintf("%-10s [I], V%01x", "LD", vx)
		case 0x65:
			op = fmt.Sprintf("%-10s V%01x,[I]", "LD", vx)
		default:
			op = fmt.Sprintf("UNKNOWN F")
		}
	}

	return fmt.Sprintf("%04x %02x %02x %s", cpu.pc, cpu.memory[cpu.pc], cpu.memory[cpu.pc+1], op)
}
