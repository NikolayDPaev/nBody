package body

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const G = 6.673e-11
const SOLARMASS = 1.98892e30

type Body struct {
	rx, ry float64 // position
	vx, vy float64 // velocity
	fx, fy float64 // force
	mass   float64
	color  color.RGBA
}

func NewBody(rx, ry, vx, vy, mass float64, color color.RGBA) *Body {
	return &Body{rx, ry, vx, vy, 0, 0, mass, color}
}

// update the velocity and position using a timestep dt
func (b *Body) Update(dt float64) {
	b.vx += dt * b.fx / b.mass
	b.vy += dt * b.fy / b.mass
	b.rx += dt * b.vx
	b.ry += dt * b.vy
}

// returns the distance between two bodies
func (b *Body) distanceTo(c *Body) float64 {
	dx := b.rx - c.rx
	dy := b.ry - c.ry
	return math.Sqrt(dx*dx + dy*dy)
}

// set the force to 0 for the next iteration
func (b *Body) ResetForce() {
	b.fx = 0.0
	b.fy = 0.0
}

// compute the net force acting between the body a and b, and
// add to the net force acting on a
func (b *Body) AddForce(c *Body) {
	EPS := 3e4 // softening parameter (just to avoid infinities)
	dx := c.rx - b.rx
	dy := c.ry - b.ry
	dist := math.Sqrt(dx*dx + dy*dy)
	F := (G * b.mass * c.mass) / (dist*dist + EPS*EPS)
	b.fx += F * dx / dist
	b.fy += F * dy / dist
}

func (b *Body) ColorPixel(translation int, canvasImage *ebiten.Image) {
	canvasImage.Set(translation+int(b.rx), translation+int(b.ry), b.color)
}
