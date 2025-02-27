package chip8

import (
	"fmt"
	"math/rand"
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
	stackIdx  uint8
	opcode    uint16
	nibbles   Nibbles
	wait      bool
	r         *rand.Rand
}

func (cpu *CPU) reset() {
	cpu.pc = 0x200
	cpu.I = 0

	cpu.r = rand.New(rand.NewSource(1337))
}

func (cpu *CPU) fetch(mem *Memory) {
	// Instructions are two bytes (Big endian)
	// First we grab the first byte and shifted 8 bits to the left (hight byte of opcode)
	// Second we get the next byte, no need to shift here

	cpu.opcode = 0
	cpu.opcode |= uint16(mem.memory[cpu.pc]) << 8
	cpu.opcode |= uint16(mem.memory[cpu.pc+1])

	if !cpu.wait {
		cpu.pc += 2
	}
}

func (cpu *CPU) decode() {
	cpu.nibbles.first = (cpu.opcode & 0xF000) >> 12
	cpu.nibbles.second = (cpu.opcode & 0x0F00) >> 8
	cpu.nibbles.third = (cpu.opcode & 0x00F0) >> 4

	cpu.nibbles.n = cpu.opcode & 0x000F
	cpu.nibbles.nn = cpu.opcode & 0x00FF
	cpu.nibbles.nnn = cpu.opcode & 0x0FFF
}

func (cpu *CPU) execute(mem *Memory, inp *Input, dis *Display) {
	// fmt.Printf("Nibbles: First: %X - Second: %X - Third: %X\n", cpu.nibbles.first, cpu.nibbles.second, cpu.nibbles.third)

	switch cpu.nibbles.first {
	case 0x0:
		switch cpu.nibbles.nn {
		case 0xE0: // (CLS) Clear the display
			clear(dis.BinaryBuffer[:])
		case 0xEE: // (RET) Return from a subroutine
			cpu.pc = cpu.stack[cpu.stackIdx-1]
			cpu.stackIdx -= 1
		}
	case 0x1: // (JP) Jump to adress NNN
		cpu.pc = cpu.nibbles.nnn
	case 0x2: // (CALL) Call subroutine at nnn
		cpu.stack[cpu.stackIdx] = cpu.pc
		cpu.stackIdx += 1
		cpu.pc = cpu.nibbles.nnn
	case 0x3: // (SE) Skip next instruction if VX = NN
		if cpu.registers[cpu.nibbles.second] == uint8(cpu.nibbles.nn) {
			cpu.pc += 2
		}
	case 0x4: // (SNE) Skip next instruction if VX != NN
		if cpu.registers[cpu.nibbles.second] != uint8(cpu.nibbles.nn) {
			cpu.pc += 2
		}
	case 0x5: // (SE) Skip next instruction if VX = VY
		if cpu.registers[cpu.nibbles.second] == cpu.registers[cpu.nibbles.third] {
			cpu.pc += 2
		}
	case 0x6: // (LD) Set register VX to NN
		cpu.registers[cpu.nibbles.second] = uint8(cpu.nibbles.nn)
	case 0x7: // (ADD) Add value NN to VX
		cpu.registers[cpu.nibbles.second] += uint8(cpu.nibbles.nn)
	case 0x8:
		switch cpu.nibbles.n {
		case 0x0:
			cpu.registers[cpu.nibbles.second] = cpu.registers[cpu.nibbles.third]
		case 0x1:
			cpu.registers[cpu.nibbles.second] |= cpu.registers[cpu.nibbles.third]
		case 0x2:
			cpu.registers[cpu.nibbles.second] &= cpu.registers[cpu.nibbles.third]
		case 0x3:
			cpu.registers[cpu.nibbles.second] ^= cpu.registers[cpu.nibbles.third]
		case 0x4:
			var sum uint16 = uint16(cpu.registers[cpu.nibbles.second]) + uint16(cpu.registers[cpu.nibbles.third])

			if sum > 255 {
				cpu.registers[0xF] = 1
			} else {
				cpu.registers[0xF] = 0
			}

			cpu.registers[cpu.nibbles.second] = uint8(sum)
		case 0x5:
			if cpu.registers[cpu.nibbles.second] >= cpu.registers[cpu.nibbles.third] {
				cpu.registers[0xF] = 1
			} else {
				cpu.registers[0xF] = 0
			}

			cpu.registers[cpu.nibbles.second] -= cpu.registers[cpu.nibbles.third]
		case 0x6:
			cpu.registers[0xF] = cpu.registers[cpu.nibbles.second] & 0x1
			cpu.registers[cpu.nibbles.second] >>= 1
		case 0x7:
			if cpu.registers[cpu.nibbles.third] >= cpu.registers[cpu.nibbles.second] {
				cpu.registers[0xF] = 1
			} else {
				cpu.registers[0xF] = 0
			}

			cpu.registers[cpu.nibbles.second] = cpu.registers[cpu.nibbles.third] - cpu.registers[cpu.nibbles.second]
		case 0xE:
			cpu.registers[0xF] = (cpu.registers[cpu.nibbles.second] & 0x80) >> 7
			cpu.registers[cpu.nibbles.second] <<= 1
		}
	case 0x9:
		if cpu.registers[cpu.nibbles.second] != cpu.registers[cpu.nibbles.third] {
			cpu.pc += 2
		}
	case 0xA: // (LD I) Set I index to NNN
		cpu.I = cpu.nibbles.nnn
	case 0xC:
		randByte := uint8(cpu.r.Intn(256))
		cpu.registers[cpu.nibbles.second] = randByte & uint8(cpu.nibbles.nn)
	case 0xD: // (DRW) Draw
		x := cpu.registers[cpu.nibbles.second] % uint8(dis.ResolutionWidth)
		y := cpu.registers[cpu.nibbles.third] % uint8(dis.ResolutionHeight)

		cpu.registers[0xF] = 0

		bytesToRead := cpu.nibbles.n

		for row := 0; uint16(row) < bytesToRead; row++ {
			spriteData := mem.memory[cpu.I+uint16(row)]

			for col := 0; col < 8; col++ {
				pixelBit := (spriteData >> (7 - col)) & 1

				idx := (int(y) * dis.ResolutionWidth) + (int(x) + col)

				if pixelBit == 1 {
					if dis.BinaryBuffer[idx] == 1 {
						cpu.registers[0xF] = 1 // Set VF for collision
					}
					dis.BinaryBuffer[idx] ^= 1
				}
			}

			y++
		}
	case 0xE:
		switch cpu.nibbles.nn {
		case 0x9E:
			if inp.keys[cpu.registers[cpu.nibbles.second]] {
				cpu.pc += 2
				fmt.Println("Key pressed")
			}
		case 0xA1:
			if !inp.keys[cpu.registers[cpu.nibbles.second]] {
				cpu.pc += 2
				fmt.Println("Key released")
			}
		}
	case 0xF:
		switch cpu.nibbles.nn {
		case 0x07:
			cpu.registers[cpu.nibbles.second] = cpu.delay
		case 0x0A:
			cpu.wait = true

			for i := 0; i < len(inp.keys); i++ {
				if inp.keys[i] {
					cpu.registers[cpu.nibbles.second] = uint8(i)
					cpu.wait = false
				}
			}
		case 0x15:
			cpu.delay = cpu.registers[cpu.nibbles.second]
		case 0x18:
			cpu.sound = cpu.registers[cpu.nibbles.second]
		case 0x1E:
			cpu.I += uint16(cpu.registers[cpu.nibbles.second])
		case 0x29:
			cpu.I = 0x50 + (uint16(cpu.registers[cpu.nibbles.second]) * 5)
		case 0x33:
			value := cpu.registers[cpu.nibbles.second]

			mem.memory[cpu.I] = value / 100
			mem.memory[cpu.I+1] = (value / 10) % 10
			mem.memory[cpu.I+2] = value % 10
		case 0x55:
			for i := 0; i <= int(cpu.nibbles.second); i++ {
				mem.memory[int(cpu.I)+i] = cpu.registers[i]
			}
		case 0x65:
			for i := 0; i <= int(cpu.nibbles.second); i++ {
				cpu.registers[i] = mem.memory[int(cpu.I)+i]
			}
		}
	}
}

func (cpu *CPU) updateTimers() {
	if cpu.delay > 0 {
		cpu.delay -= 1
	}

	if cpu.sound > 0 {
		cpu.sound -= 1
	}
}

func (cpu *CPU) tick(mem *Memory, inp *Input, dis *Display) {
	cpu.fetch(mem)
	cpu.decode()
	cpu.execute(mem, inp, dis)
}
