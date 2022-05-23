package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	RESOLUTION = 800
)

type Game struct {
	canvasImage *ebiten.Image
}

func NewGame() *Game {
	g := &Game{
		canvasImage: ebiten.NewImage(RESOLUTION, RESOLUTION),
	}
	return g
}

func (g *Game) Update() error {
	//g.canvasImage
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.canvasImage, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return RESOLUTION, RESOLUTION
}
