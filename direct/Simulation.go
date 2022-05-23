package direct

import (
	"image/color"
	"math/rand"
	"sync"
	"time"

	"github.com/muesli/gamut"
)

type Body struct {
	x    int32
	y    int32
	mass int32
	u    int32
	v    int32
}

type Simulation struct {
	p       int32
	bodies  []Body
	mu      sync.Mutex
	palette []color.Color
}

func NewSimulation(p, n, res int32) *Simulation {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	bodies := make([]Body, n)
	for i := range bodies {
		bodies[i] = Body{r.Int31n(res), r.Int31n(res), r.Int31n(100), r.Int31n(10), r.Int31n(10)}
	}

	palette, _ := gamut.Generate(8, gamut.PastelGenerator{})
	return &Simulation{p, bodies, sync.Mutex{}, palette}
}

func Worker(id, p int, bodies []Body, mu *sync.Mutex) {
	localLen := len(bodies) / p
	begin := id * localLen
	localBodies := bodies[begin : begin+localLen]

}
