package main

import (
	"fishing-game/game"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

func main() {
	// Crear el juego
	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	// Configurar ventana
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Fishing Game - Concurrent Programming")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Ejecutar el juego
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}