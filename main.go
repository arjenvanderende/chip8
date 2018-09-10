package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/arjenvanderende/chip8/chip8"
	"github.com/arjenvanderende/chip8/io/termbox"
)

func main() {
	decompile := flag.Bool("decompile", false, "Print opcodes of the loaded ROM")
	filename := flag.String("romfile", "roms/fishie.ch8", "The ROM file to load")
	logfile := flag.String("logfile", "", "The file to log to")
	flag.Parse()

	// setup logging
	if *logfile != "" {
		f, err := os.Create(*logfile)
		if err != nil {
			log.Fatal(fmt.Errorf("Unable to create logfile: %v", err))
		}
		log.SetOutput(f)
	}

	// load the ROM file
	cpu, err := chip8.Load(*filename)
	if err != nil {
		log.Fatal(err)
	}

	// disassemble opcodes
	if *decompile {
		printOpcodes(os.Stdout, cpu)
	} else {
		err = run(cpu)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func printOpcodes(w io.Writer, cpu *chip8.CPU) {
	for {
		op := cpu.DisassembleOp()
		fmt.Fprintf(w, "%s\n", op)

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
