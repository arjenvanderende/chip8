package chip8

import (
	"fmt"
	"io/ioutil"
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
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// Memory represents the memory address space of the Chip-8
type Memory [0x1000]byte

// CPU represents the Chip8 CPU
type CPU struct {
	pc     int
	memory Memory
	i      uint16
	v      [16]byte

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
	}
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
		case <-clock.C:
			err := cpu.interpret(display)
			if err != nil {
				return fmt.Errorf("Could not interpret op: %v", err)
			}
		case <-frame.C:
			display.Flush()
		case key := <-keyboard.Events():
			if key == io.KeyEsc {
				return nil
			}
		}
	}
}

func (cpu *CPU) printState(pc int, op string) {
	fmt.Printf("op=%-40s pc=%03x next pc=%03x i=%03x v=%v\n", op, pc, cpu.pc, cpu.i, cpu.v)
}

func (cpu *CPU) interpret(display io.Display) error {
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
		default:
			return fmt.Errorf("Unknown 0")
		}
	case 0x1:
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
	case 0xf:
		switch cpu.memory[cpu.pc+1] {
		case 0x1e:
			cpu.i += uint16(cpu.v[vx])
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
