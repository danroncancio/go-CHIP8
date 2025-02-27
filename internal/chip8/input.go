package chip8

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// CHIP8 Keyboard
// 1 2 3 C
// 4 5 6 D
// 7 8 9 E
// A 0 B F

type Input struct {
	keys   [16]bool
	keyMap map[ebiten.Key]uint8
}

func (i *Input) initInput() {
	i.keyMap = map[ebiten.Key]uint8{
		ebiten.KeyDigit1: 0x1,
		ebiten.KeyDigit2: 0x2,
		ebiten.KeyDigit3: 0x3,
		ebiten.KeyDigit4: 0xC,
		ebiten.KeyQ:      0x4,
		ebiten.KeyW:      0x5,
		ebiten.KeyE:      0x6,
		ebiten.KeyR:      0xD,
		ebiten.KeyA:      0x7,
		ebiten.KeyS:      0x8,
		ebiten.KeyD:      0x9,
		ebiten.KeyF:      0xE,
		ebiten.KeyZ:      0xA,
		ebiten.KeyX:      0x0,
		ebiten.KeyC:      0xB,
		ebiten.KeyV:      0xF,
	}
}

func (i *Input) processInput() {
	// Reset all keys to false before updating
	for k := range i.keys {
		i.keys[k] = false
	}

	// Update keys based on input state
	for key, chip8Key := range i.keyMap {
		if ebiten.IsKeyPressed(key) {
			i.keys[chip8Key] = true
		}
	}
}
