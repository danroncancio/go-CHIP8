package chip8

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Chip8 struct {
	cpu           CPU
	memory        Memory
	Display       Display
}

func New() (*Chip8, error) {
	if len(os.Args) != 2 {
		return nil, errors.New("Wrong amount of arguments")
	}

	app := &Chip8{}

	app.Display.WindowWidth = 640
	app.Display.WindowHeight = 320
	app.Display.ResolutionWidth = 64
	app.Display.ResolutionHeight = 32

	app.memory.loadFont()

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

	app.memory.loadRomIntoMemory(rom)

	app.cpu.reset()

	return app, nil
}

func (chip8 *Chip8) Update() error {
	chip8.cpu.tick(&chip8.memory, &chip8.Display)

	return nil
}

func (chip8 *Chip8) Draw(screen *ebiten.Image) {
	for i, bit := range chip8.Display.BinaryBuffer {
		idx := i * 4

		if bit == 1 {
			chip8.Display.RGBABuffer[idx+0] = 255 // R
			chip8.Display.RGBABuffer[idx+1] = 255 // G
			chip8.Display.RGBABuffer[idx+2] = 255 // B
			chip8.Display.RGBABuffer[idx+3] = 255 // A
		} else {
			chip8.Display.RGBABuffer[idx+0] = 0
			chip8.Display.RGBABuffer[idx+1] = 0
			chip8.Display.RGBABuffer[idx+2] = 0
			chip8.Display.RGBABuffer[idx+3] = 0
		}
	}

	screen.WritePixels(chip8.Display.RGBABuffer[:])
}

func (chip8 *Chip8) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return chip8.Display.ResolutionWidth, chip8.Display.ResolutionHeight
}
