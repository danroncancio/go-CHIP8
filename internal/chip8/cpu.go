package chip8

import (
	"fmt"
)

type Nibbles struct {
	first  uint16
	second uint16
	third  uint16
	n      uint16
	nn     uint16
	nnn    uint16
}

type CPU struct {
	registers [16]uint8
	I         uint16
	delay     uint8
	sound     uint8
	pc        uint16
	stack     [16]uint16
	opcode    uint16
	nibbles   Nibbles

	mem           *Memory
	displayBuffer []byte
	screenWidth   int
	screenHeight  int
}

func newCPU(mem *Memory, displayBuffer []byte, screenWidth int, screenHeight int) (*CPU, error) {
	cpu := &CPU{}

	cpu.screenWidth = screenWidth
	cpu.screenHeight = screenHeight

	cpu.displayBuffer = displayBuffer
	cpu.mem = mem
	cpu.pc = 0x200 // Start of the rom data

	return cpu, nil
}

func (cpu *CPU) fetch() {
	// Instructions are two bytes (Big endian)
	// First we grab the first byte and shifted 8 bits to the left (hight byte of opcode)
	// Second we get the next byte, no need to shift here

	cpu.opcode = 0
	cpu.opcode |= uint16(cpu.mem.memory[cpu.pc]) << 8
	cpu.opcode |= uint16(cpu.mem.memory[cpu.pc+1])
	cpu.pc += 2
}

func (cpu *CPU) decode() {
	cpu.nibbles.first = (cpu.opcode & 0xF000) >> 12
	cpu.nibbles.second = (cpu.opcode & 0x0F00) >> 8
	cpu.nibbles.third = (cpu.opcode & 0x00F0) >> 4

	cpu.nibbles.n = cpu.opcode & 0x000F
	cpu.nibbles.nn = cpu.opcode & 0x00FF
	cpu.nibbles.nnn = cpu.opcode & 0x0FFF
}

func (cpu *CPU) execute() {
	fmt.Printf("Nibbles: First: %X - Second: %X - Third: %X\n", cpu.nibbles.first, cpu.nibbles.second, cpu.nibbles.third)

	switch cpu.nibbles.first {
	case 0x0:
		switch cpu.nibbles.nn {
		case 0xE0: // (CLS) Clear the display
			clear(cpu.displayBuffer)
		}
	case 0x1: // (JP) Jump to adress NNN
		cpu.pc = cpu.nibbles.nnn
	case 0x6: // (LD) Set register VX to NN
		cpu.registers[cpu.nibbles.second] = uint8(cpu.nibbles.nn)
	case 0x7: // (ADD) Add value NN to VX
		cpu.registers[cpu.nibbles.second] += uint8(cpu.nibbles.nn)
	case 0xA: // (LD I) Set I index to NNN
		cpu.I = cpu.nibbles.nnn
	case 0xD: // (DRW) Draw
		x := cpu.registers[cpu.nibbles.second] % uint8(cpu.screenWidth)
		y := cpu.registers[cpu.nibbles.third] % uint8(cpu.screenHeight)

		cpu.registers[0xF] = 0

		bytesToRead := cpu.nibbles.n

		for row := 0; uint16(row) < bytesToRead; row++ {
			spriteData := cpu.mem.memory[cpu.I+uint16(row)]

			for col := 0; col < 8; col++ {
				pixelBit := (spriteData >> (7 - col)) & 1

				idx := (int(y) * cpu.screenWidth) + (int(x) + col)

				if pixelBit == 1 {
					if cpu.displayBuffer[idx] == 1 {
						cpu.registers[0xF] = 1 // Set VF for collision
					}
					cpu.displayBuffer[idx] ^= 1
				}
			}

			y++
		}
	}
}

func (cpu *CPU) tick() {
	cpu.fetch()
	cpu.decode()
	cpu.execute()
}
