package chip8

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WndWidth    = 640
	WndHeight   = 320
	Chip8Width  = 64
	Chip8Height = 32
)

var memory [4 * 1024]uint8
var opcode uint16 = 0

type Chip8 struct {
	cpu           *CPU
	memory        *Memory
	display       [4 * Chip8Width * Chip8Height]byte
	binaryDisplay [Chip8Width * Chip8Height]byte
}

func New() (*Chip8, error) {
	if len(os.Args) != 2 {
		return nil, errors.New("Wrong number of arguments")
	}

	app := &Chip8{}

	// Read rom file

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	rom, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Memory

	app.memory, err = newMemory()
	if err != nil {
		log.Fatal(err)
	}

	app.memory.loadRomIntoMemory(rom)

	// CPU

	app.cpu, err = newCPU(app.memory, app.binaryDisplay[:], Chip8Width, Chip8Height)
	if err != nil {
		log.Fatal(err)
	}

	return app, nil
}

func (chip8 *Chip8) Update() error {
	chip8.cpu.tick()

	return nil
}

func (chip8 *Chip8) Draw(screen *ebiten.Image) {
	for i, bit := range chip8.binaryDisplay {
		idx := i * 4

		if bit == 1 {
			chip8.display[idx+0] = 255 // R
			chip8.display[idx+1] = 255 // G
			chip8.display[idx+2] = 255 // B
			chip8.display[idx+3] = 255 // A
		} else {
			chip8.display[idx+0] = 0   // R
			chip8.display[idx+1] = 0   // G
			chip8.display[idx+2] = 0   // B
			chip8.display[idx+3] = 255 // A
		}
	}

	screen.WritePixels(chip8.display[:])
}

func (chip8 *Chip8) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return Chip8Width, Chip8Height
}
