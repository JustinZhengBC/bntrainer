package bntrainer

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	Title        string = "BN Trainer"
	ScreenWidth  int    = 960
	ScreenHeight int    = 540
	BoardSize    int    = 8
)

type Game struct{}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	UpdateCursorPosition(image.Point{x, y})
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ProcessClick()
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		ProcessRelease()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		ToggleHelp()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		ReturnToMenu()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		HandleStartNewGameShortcut()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	DrawGui(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
