package main

import (
	"log"

	"github.com/NikolayDPaev/n-body/direct"
	"github.com/hajimehoshi/ebiten/v2"
)

const RESOLUTION = 600

type Game struct {
	simulation  *direct.Simulation
	canvasImage *ebiten.Image
}

func (g *Game) Update() error {
	return g.simulation.Update(g.canvasImage)
}

func (s *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(s.canvasImage, nil)
}

func (s *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return RESOLUTION, RESOLUTION
}

func main() {
	ebiten.SetWindowSize(RESOLUTION, RESOLUTION)
	ebiten.SetWindowTitle("nBody")
	game := &Game{direct.NewSimulation(1, 100, RESOLUTION), ebiten.NewImage(RESOLUTION, RESOLUTION)}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
