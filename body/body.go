package body

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
)

var Palette = []color.Color{
	color.RGBA{0, 0, 0, 0xff},
	color.RGBA{255, 255, 153, 0xff},
	color.RGBA{255, 214, 122, 0xff},
	color.RGBA{255, 173, 92, 0xff},
	color.RGBA{255, 133, 61, 0xff},
	color.RGBA{255, 92, 31, 0xff},
	color.RGBA{255, 51, 0, 0xff},
	color.RGBA{255, 0, 0, 0xff},
}

const RADIUS = 1e18 // radius of universe
const RESOLUTION = 250
const G = 6.673e-11
const SOLARMASS = 1.98892e30

type Body struct {
	rx, ry float64 // position
	vx, vy float64 // velocity
	fx, fy float64 // force
	mass   float64
	color  color.Color
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

// compute the net force acting between the body c and b, and
// add to the net force acting on b
func (b *Body) AddForce(c *Body) {
	EPS := 3e4 // softening parameter (just to avoid infinities)
	dx := c.rx - b.rx
	dy := c.ry - b.ry
	dist := math.Sqrt(dx*dx + dy*dy)
	F := (G * b.mass * c.mass) / (dist*dist + EPS*EPS)
	b.fx += F * dx / dist
	b.fy += F * dy / dist
}

func drawCircle(xc, yc, r int, img *image.Paletted, color color.Color) {
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				img.Set(xc+x, yc+y, color)
			}
		}
	}
}

func (b *Body) ColorPixel(img *image.Paletted) {
	if img != nil {
		x := int(math.Round(b.rx*RESOLUTION/RADIUS)) + RESOLUTION/2
		y := int(math.Round(b.ry*RESOLUTION/RADIUS)) + RESOLUTION/2

		radius := int(b.mass*3/(SOLARMASS*10+1e20)) + 1
		if radius > 4 {
			radius = 7
		}
		drawCircle(x, y, radius, img, b.color)
	}
}

func circlev(rx, ry float64) float64 {
	r2 := math.Sqrt(rx*rx + ry*ry)
	numerator := (6.67e-11) * 1e6 * SOLARMASS
	return math.Sqrt(numerator / r2)
}

func exp(r *rand.Rand, lambda float64) float64 {
	return -math.Log(1-r.Float64()) / lambda
}

func NewCentralBody(mass float64, color color.RGBA) *Body {
	return &Body{rx: 0, ry: 0, vx: 0, vy: 0, fx: 0, fy: 0, mass: mass, color: color}
}

func NewRandomBody(r *rand.Rand) *Body {
	px := 1e18 * exp(r, -1.8) * (.5 - r.Float64())
	py := 1e18 * exp(r, -1.8) * (.5 - r.Float64())

	magv := circlev(px, py)
	absangle := math.Atan(math.Abs(py / px))
	thetav := math.Pi/2 - absangle
	//phiv := r.Float64() * math.Pi
	vx := -1 * math.Copysign(1.0, py) * math.Cos(thetav) * magv
	vy := math.Copysign(1.0, px) * math.Sin(thetav) * magv
	//Orient a random 2D circular orbit
	if r.Float64() <= .5 {
		vx = -vx
		vy = -vy
	}

	mass := r.Float64()*SOLARMASS*10 + 1e20
	//Color the masses in green gradients by mass
	color := Palette[int(mass*6/(SOLARMASS*10+1e20))+1]

	return &Body{rx: px, ry: py, vx: vx, vy: vy, fx: 0.0, fy: 0.0, mass: mass, color: color}
}

func (b *Body) String() string {
	return fmt.Sprintf("%f rx, %f ry, %f vx, %f vy, %f mass", b.rx, b.ry, b.vx, b.vy, b.mass)
}
