package chip8

import (
	"fmt"
	"io/ioutil"
	"time"
)

const (
	// clockRate represents the number of operations that the CPU can process per second
	clockRate int = 540
	// programOffset represents the offset in memory where the program is loaded
	programOffset int = 0x200
)

// Memory represents the memory address space of the Chip-8
type Memory [0x1000]byte

// CPU represents the Chip8 CPU
type CPU struct {
	pc          int
	memory      Memory
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
	}
	for i, b := range bytes {
		cpu.memory[programOffset+i] = b
	}
	return &cpu, nil
}

// Run starts running the program
func (cpu *CPU) Run() {
	clock := time.NewTicker(time.Second / time.Duration(clockRate))
	defer clock.Stop()

	frame := time.NewTicker(time.Second / time.Duration(60))
	defer frame.Stop()

	for {
		select {
		case <-clock.C:
			fmt.Printf(".")
		case <-frame.C:
			fmt.Printf("F")
		}
	}
}

// NextOp increments the PC to the next operation
// Returns false when there are no more operations to read
func (cpu *CPU) NextOp() bool {
	cpu.pc += 2
	return cpu.pc <= programOffset+cpu.programSize
}

// DisassembleOp output the assembly for the operation at the PC.
func (cpu *CPU) DisassembleOp() {
	nib1 := cpu.memory[cpu.pc] >> 4

	vx := cpu.memory[cpu.pc] & 0x0f
	vy := cpu.memory[cpu.pc+1] >> 4
	n := cpu.memory[cpu.pc+1] & 0x0f
	nn := cpu.memory[cpu.pc+1]
	nnn := int16(cpu.memory[cpu.pc]&0x0f)<<8 + int16(cpu.memory[cpu.pc+1])

	fmt.Printf("%04x %02x %02x ", cpu.pc, cpu.memory[cpu.pc], cpu.memory[cpu.pc+1])
	switch nib1 {
	case 0x0:
		switch cpu.memory[cpu.pc+1] {
		case 0xe0:
			fmt.Printf("%-10s", "CLS")
		case 0xee:
			fmt.Printf("%-10s", "RET")
		default:
			fmt.Printf("%-10s %03x", "SYS", nnn)
		}
	case 0x1:
		fmt.Printf("%-10s %03x", "JP", nnn)
	case 0x2:
		fmt.Printf("%-10s %03x", "CALL", nnn)
	case 0x3:
		fmt.Printf("%-10s V%01x, %02x", "SE", vx, nn)
	case 0x4:
		fmt.Printf("%-10s V%01x, %02x", "SNE", vx, nn)
	case 0x5:
		fmt.Printf("%-10s V%01x, V%01x", "SE", vx, vy)
	case 0x6:
		fmt.Printf("%-10s V%01x, %02x", "LD", vx, nn)
	case 0x7:
		fmt.Printf("%-10s V%01x, %02x", "ADD", vx, nn)
	case 0x8:
		lastNib := cpu.memory[cpu.pc+1] & 0x0f
		switch lastNib {
		case 0x0:
			fmt.Printf("%-10s V%01x, V%01x", "LD", vx, vy)
		case 0x1:
			fmt.Printf("%-10s V%01x, V%01x", "OR", vx, vy)
		case 0x2:
			fmt.Printf("%-10s V%01x, V%01x", "AND", vx, vy)
		case 0x3:
			fmt.Printf("%-10s V%01x, V%01x", "XOR", vx, vy)
		case 0x4:
			fmt.Printf("%-10s V%01x, V%01x", "ADD", vx, vy)
		case 0x5:
			fmt.Printf("%-10s V%01x, V%01x, V%01x", "SUB", vx, vx, vy)
		case 0x6:
			fmt.Printf("%-10s V%01x, V%01x", "SHR", vx, vy)
		case 0x7:
			fmt.Printf("%-10s V%01x, V%01x, V%01x", "SUBN", vx, vy, vy)
		case 0xe:
			fmt.Printf("%-10s V%01x, V%01x", "SHL", vx, vy)
		default:
			fmt.Printf("UNKNOWN 8")
		}
	case 0x9:
		fmt.Printf("%-10s V%01x, V%01x", "SNE", vx, vy)
	case 0xa:
		fmt.Printf("%-10s I,%03x", "LD", nnn)
	case 0xb:
		fmt.Printf("%-10s V0,%03x", "JP", nnn)
	case 0xc:
		fmt.Printf("%-10s V%01x, %02x", "RND", vx, nn)
	case 0xd:
		fmt.Printf("%-10s V%01x, V%01x, %01x", "DRW", vx, vy, n)
	case 0xe:
		switch cpu.memory[cpu.pc+1] {
		case 0x9e:
			fmt.Printf("%-10s V%01x", "SKP", vx)
		case 0xa1:
			fmt.Printf("%-10s V%01x", "SKNP", vx)
		default:
			fmt.Printf("UNKNOWN E")
		}
	case 0xf:
		switch cpu.memory[cpu.pc+1] {
		case 0x07:
			fmt.Printf("%-10s V%01x, DELAY", "LD", vx)
		case 0x0a:
			fmt.Printf("%-10s V%01x, KEY", "LD", vx)
		case 0x15:
			fmt.Printf("%-10s DELAY, V%01x", "LD", vx)
		case 0x18:
			fmt.Printf("%-10s SOUND, V%01x", "LD", vx)
		case 0x1e:
			fmt.Printf("%-10s I, V%01x", "ADD", vx)
		case 0x29:
			fmt.Printf("%-10s F, V%01x", "LD", vx)
		case 0x33:
			fmt.Printf("%-10s B, V%01x", "LD", vx)
		case 0x55:
			fmt.Printf("%-10s [I], V%01x", "LD", vx)
		case 0x65:
			fmt.Printf("%-10s V%01x,[I]", "LD", vx)
		default:
			fmt.Printf("UNKNOWN F")
		}
	default:
		fmt.Print("not implemented")
	}
}
