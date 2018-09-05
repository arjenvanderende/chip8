package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/arjenvanderende/chip8/chip8"
	"github.com/arjenvanderende/chip8/io/termbox"
)

func main() {
	filename := flag.String("romfile", "roms/fishie.ch8", "The ROM file to load")
	decompile := flag.Bool("decompile", false, "Print opcodes of the loaded ROM")
	flag.Parse()

	// load the ROM file
	cpu, err := chip8.Load(*filename)
	if err != nil {
		log.Fatal(err)
	}

	// disassemble opcodes
	if *decompile {
		printOpcodes(cpu)
	} else {
		err = run(cpu)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func printOpcodes(cpu *chip8.CPU) {
	for {
		op := cpu.DisassembleOp()
		fmt.Printf("%s\n", op)

		if !cpu.NextOp() {
			break
		}
	}
}

func run(cpu *chip8.CPU) error {
	// initialise I/O devices
	display, keyboard, closer, err := termbox.New()
	if err != nil {
		return fmt.Errorf("Unable to initialise graphics: %v", err)
	}
	defer closer()

	// run the program
	err = cpu.Run(display, keyboard)
	if err != nil {
		return fmt.Errorf("Program failed to run: %v", err)
	}
	return nil
}
