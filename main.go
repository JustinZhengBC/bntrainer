package main

import (
	"log"

	"bntrainer/bntrainer"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := &bntrainer.Game{}
	ebiten.SetWindowSize(bntrainer.ScreenWidth, bntrainer.ScreenHeight)
	ebiten.SetWindowTitle(bntrainer.Title)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
