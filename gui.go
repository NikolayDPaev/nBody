package main

import (
	"github.com/NikolayDPaev/n-body/direct"
	"github.com/hajimehoshi/ebiten/v2"
)

const RESOLUTION = 600

type Game struct {
	p, n        int
	simulation  *direct.Simulation
	canvasImage *ebiten.Image
}

func NewGame(p, n int) *Game {
	canvasImage := ebiten.NewImage(RESOLUTION, RESOLUTION)
	return &Game{p, n, direct.NewSimulation(p, 100, RESOLUTION, ebiten.NewImage(RESOLUTION, RESOLUTION)), canvasImage}
}

func (g *Game) Update() error {
	return g.simulation.Update()
}

func (s *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(s.canvasImage, nil)
}

func (s *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return RESOLUTION, RESOLUTION
}
