package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/daroncancio/go-chip8/internal/chip8"
)

func main() {
	app, err := chip8.New()
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(app.Display.WindowWidth, app.Display.WindowHeight)
	ebiten.SetWindowTitle("Go CHIP-8")

	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}
