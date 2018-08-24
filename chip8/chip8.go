package chip8

import (
	"fmt"
	"io/ioutil"
)

// ProgramOffset represents the offset in memory where the program is loaded
const programOffset int = 0x200

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
			fmt.Printf("%-10s", "RTS")
		default:
			fmt.Printf("UNKNOWN 0")
		}
	case 0x1:
		fmt.Printf("%-10s $%03x", "JUMP", nnn)
	case 0x2:
		fmt.Printf("%-10s $%03x", "CALL", nnn)
	case 0x3:
		fmt.Printf("%-10s V%01X,#$%02x", "SKIP.EQ", vx, nn)
	case 0x4:
		fmt.Printf("%-10s V%01X,#$%02x", "SKIP.NE", vx, nn)
	case 0x5:
		fmt.Printf("%-10s V%01X,V%01x", "SKIP.EQ", vx, vy)
	case 0x6:
		fmt.Printf("%-10s V%01X,#$%02x", "MVI", vx, nn)
	case 0x7:
		fmt.Printf("%-10s V%01X,#$%02x", "ADI", vx, nn)
	case 0x8:
		lastNib := cpu.memory[cpu.pc+1] & 0x0f
		switch lastNib {
		case 0x0:
			fmt.Printf("%-10s V%01X,V%01X", "MOV.", vx, vy)
		case 0x1:
			fmt.Printf("%-10s V%01X,V%01X", "OR.", vx, vy)
		case 0x2:
			fmt.Printf("%-10s V%01X,V%01X", "AND.", vx, vy)
		case 0x3:
			fmt.Printf("%-10s V%01X,V%01X", "XOR.", vx, vy)
		case 0x4:
			fmt.Printf("%-10s V%01X,V%01X", "ADD.", vx, vy)
		case 0x5:
			fmt.Printf("%-10s V%01X,V%01X,V%01X", "SUB.", vx, vx, vy)
		case 0x6:
			fmt.Printf("%-10s V%01X,V%01X", "SHIFTR.", vx, vy)
		case 0x7:
			fmt.Printf("%-10s V%01X,V%01X,V%01X", "SUB.", vx, vy, vy)
		case 0xe:
			fmt.Printf("%-10s V%01X,V%01X", "SHIFTL.", vx, vy)
		default:
			fmt.Printf("UNKNOWN 8")
		}
	case 0x9:
		fmt.Printf("%-10s V%01X,V%01X", "SKIP.NE", vx, vy)
	case 0xa:
		fmt.Printf("%-10s I,#$%03x", "MVI", nnn)
	case 0xb:
		fmt.Printf("%-10s $%03x(V0)", "JUMP", nnn)
	case 0xc:
		fmt.Printf("%-10s V%01x,#$%02x", "RNDMSK", vx, nn)
	case 0xd:
		fmt.Printf("%-10s V%01X,V%01X,#$%01x", "SPRITE", vx, vy, n)
	case 0xe:
		switch cpu.memory[cpu.pc+1] {
		case 0x9e:
			fmt.Printf("%-10s V%01X", "SKIP.KEYDOWN", vx)
		case 0xa1:
			fmt.Printf("%-10s V%01X", "SKIP.KEYUP", vx)
		default:
			fmt.Printf("UNKNOWN E")
		}
	case 0xf:
		switch cpu.memory[cpu.pc+1] {
		case 0x07:
			fmt.Printf("%-10s V%01X,DELAY", "MOV", vx)
		case 0x0a:
			fmt.Printf("%-10s V%01X", "KEY", vx)
		case 0x15:
			fmt.Printf("%-10s DELAY,V%01X", "MOV", vx)
		case 0x18:
			fmt.Printf("%-10s SOUND,V%01X", "MOV", vx)
		case 0x1e:
			fmt.Printf("%-10s I,V%01X", "ADI", vx)
		case 0x29:
			fmt.Printf("%-10s I,V%01X", "SPRITECHAR", vx)
		case 0x33:
			fmt.Printf("%-10s (I),V%01X", "MOVBCD", vx)
		case 0x55:
			fmt.Printf("%-10s (I),V0-V%01X", "MOVM", vx)
		case 0x65:
			fmt.Printf("%-10s V0-V%01X,(I)", "MOVM", vx)
		default:
			fmt.Printf("UNKNOWN F")
		}
	default:
		fmt.Print("not implemented")
	}
}
