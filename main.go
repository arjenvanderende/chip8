package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/arjenvanderende/chip8/chip8"
)

func main() {
	filename := flag.String("romfile", "roms/fishie.ch8", "The ROM file to load")
	flag.Parse()

	// load the ROM file
	cpu, err := chip8.Load(*filename)
	if err != nil {
		log.Fatal(err)
	}

	// disassemble opcodes
	for {
		cpu.DisassembleOp()
		fmt.Printf("\n")

		if !cpu.NextOp() {
			break
		}
	}
}
