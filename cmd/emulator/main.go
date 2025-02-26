package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/daroncancio/go-chip8/internal/chip8"
)

func main() {
	ebiten.SetWindowSize(chip8.WndWidth, chip8.WndHeight)
	ebiten.SetWindowTitle("Go CHIP-8")

	app, err := chip8.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}
