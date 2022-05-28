package direct

import (
	"image/color"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	timestep = 1e11
)

type Simulation struct {
	p      int
	res    int
	bodies []*Body
	mu     sync.Mutex
}

func NewSimulation(p, n, res int) *Simulation {
	s := &Simulation{p, res, make([]*Body, n), sync.Mutex{}}
	s.startthebodies()
	return s
}

func (s *Simulation) Update(canvasImage *ebiten.Image) error {
	translation := s.res / 2
	for _, body := range s.bodies {
		canvasImage.Set(translation+int(body.rx), translation+int(body.ry), body.color)
		//s.canvasImage.
		//g.fillOval((int) Math.round(bodies[i].rx*250/1e18),(int) Math.round(bodies[i].ry*250/1e18),8,8);
	}
	//go through the Brute Force algorithm (see the function below)
	s.addforces()
	return nil
}

func circlev(rx, ry float64) float64 {
	r2 := math.Sqrt(rx*rx + ry*ry)
	numerator := (6.67e-11) * 1e6 * solarmass
	return math.Sqrt(numerator / r2)
}

//Initialize N bodies with random positions and circular velocities
func (s *Simulation) startthebodies() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	exp := func(lambda float64) float64 {
		return -math.Log(1-r.Float64()) / lambda
	}

	//radius := 1e18 // radius of universe

	for i := range s.bodies {
		px := 1e8 * exp(-1.8) * (.5 - r.Float64())
		py := 1e8 * exp(-1.8) * (.5 - r.Float64())
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

		mass := r.Float64()*solarmass*10 + 1e20
		//Color the masses in green gradients by mass
		red := uint8(mass * 254 / (solarmass*10 + 1e20))
		blue := uint8(mass * 254 / (solarmass*10 + 1e20))
		var green uint8 = 255
		color := color.RGBA{red, green, blue, 0xff}
		s.bodies[i] = NewBody(px, py, vx, vy, mass, color)
	}
	//Put the central mass in
	s.bodies[0] = NewBody(0, 0, 0, 0, 1e6*solarmass, color.RGBA{255, 0, 0, 0xff}) //put a heavy body in the center

}

//Use the method in Body to reset the forces, then add all the new forces
func (s *Simulation) addforces() {
	for i := range s.bodies {
		s.bodies[i].resetForce()
		//Notice-2 loops-->N^2 complexity
		for j := range s.bodies {
			if i != j {
				s.bodies[i].addForce(s.bodies[j])
			}
		}
	}
	//Then, loop again and update the bodies using timestep dt
	for i := range s.bodies {
		s.bodies[i].update(timestep)
	}
}
