package pipeline

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/NikolayDPaev/n-body/body"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	timestep = 1e11
)

type Simulation struct {
	p           int
	res         int
	canvasImage *ebiten.Image

	// bodies
	bodies       []*body.Body
	pipeChannels []chan *body.Body
}

func NewSimulation(p, n, res int, canvasImage *ebiten.Image) *Simulation {
	s := &Simulation{p, res, canvasImage, make([]*body.Body, n), make([]chan *body.Body, p+1)}
	s.startthebodies()
	for i := range s.pipeChannels {
		s.pipeChannels[i] = make(chan *body.Body, n)
	}
	// start the workers
	for id := 0; id < s.p; id++ {
		go s.worker(id)
	}
	return s
}

// master
func (s *Simulation) Update() error {
	// send to the pipeline
	for i := range s.bodies {
		s.pipeChannels[0] <- s.bodies[i]
	}

	// update positions
	translation := s.res / 2
	for range s.bodies {
		body := <-s.pipeChannels[s.p]
		body.Update(timestep)
		if s.canvasImage != nil {
			body.ColorPixel(translation, s.canvasImage)
		}
		body.ResetForce()
	}
	return nil
}

func (s *Simulation) worker(id int) {
	// take own bodies
	start := (len(s.bodies) / s.p) * id
	var end int
	if id == s.p-1 {
		end = len(s.bodies)
	} else {
		end = start + len(s.bodies)%s.p
	}

	local := make([]*body.Body, end-start)
	// loop
	for {
		// take local bodies and pair them with one another
		for i := 0; i < end-start; i++ {
			local[i] = <-s.pipeChannels[id]
			for j := 0; j < i-1; j++ {
				local[i].AddForce(local[j])
				local[j].AddForce(local[i])
			}
		}

		// calculate the forces on the local bodies from the others
		for j := len(local); j < len(s.bodies); j++ {
			body := <-s.pipeChannels[id]
			for i := 0; i < end-start; i++ {
				local[i].AddForce(body)
			}
			s.pipeChannels[id+1] <- body
		}

		// send the local bodies to the next node
		for i := 0; i < end-start; i++ {
			s.pipeChannels[id+1] <- local[i]
		}
	}
}

func circlev(rx, ry float64) float64 {
	r2 := math.Sqrt(rx*rx + ry*ry)
	numerator := (6.67e-11) * 1e6 * body.SOLARMASS
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

		mass := r.Float64()*body.SOLARMASS*10 + 1e20
		//Color the masses in green gradients by mass
		red := uint8(mass * 254 / (body.SOLARMASS*10 + 1e20))
		blue := uint8(mass * 254 / (body.SOLARMASS*10 + 1e20))
		var green uint8 = 255
		color := color.RGBA{red, green, blue, 0xff}
		s.bodies[i] = body.NewBody(px, py, vx, vy, mass, color)
	}
	//Put the central mass in
	s.bodies[0] = body.NewBody(0, 0, 0, 0, 1e6*body.SOLARMASS, color.RGBA{255, 0, 0, 0xff}) //put a heavy body.body in the center

}
