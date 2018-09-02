package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/arjenvanderende/chip8/chip8"
	"github.com/arjenvanderende/chip8/io"
)

func main() {
	filename := flag.String("romfile", "roms/fishie.ch8", "The ROM file to load")
	decompile := flag.Bool("decompile", false, "Print opcodes of the loaded ROM")
	flag.Parse()

	// initialise the graphics
	graphics, err := io.NewTermbox()
	if err != nil {
		log.Fatal(err)
	}
	defer graphics.Close()

	// load the ROM file
	cpu, err := chip8.Load(*filename)
	if err != nil {
		log.Fatal(err)
	}

	// disassemble opcodes
	if *decompile {
		printOpcodes(cpu)
	} else {
		cpu.Run(graphics)
	}
}

func printOpcodes(cpu *chip8.CPU) {
	for {
		cpu.DisassembleOp()
		fmt.Printf("\n")

		if !cpu.NextOp() {
			break
		}
	}
}
