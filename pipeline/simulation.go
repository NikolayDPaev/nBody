package pipeline

import (
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/NikolayDPaev/n-body/body"
)

const (
	timestep = 1e11
)

type Simulation struct {
	p   int
	img *image.Paletted

	// bodies
	bodies       []*body.Body
	pipeChannels []chan *body.Body
}

func NewSimulation(p, n int) *Simulation {
	s := &Simulation{p, nil, make([]*body.Body, n), make([]chan *body.Body, p+1)}
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
func (s *Simulation) Update(image *image.Paletted) error {
	s.img = image
	// send to the pipeline
	for i := range s.bodies {
		s.pipeChannels[0] <- s.bodies[i]
	}

	// update positions
	for range s.bodies {
		body := <-s.pipeChannels[s.p]
		body.Update(timestep)
		if s.img != nil {
			body.ColorPixel(s.img)
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

//Initialize N bodies with random positions and circular velocities
func (s *Simulation) startthebodies() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range s.bodies {
		s.bodies[i] = body.NewRandomBody(r)
	}
	//Put the central mass in
	s.bodies[0] = body.NewCentralBody(1e6*body.SOLARMASS, color.RGBA{255, 0, 0, 0xff}) //put a heavy body in the center
}
