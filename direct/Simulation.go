package direct

import (
	"image"
	"image/color"
	"math/rand"
	"sync"
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
	bodies  []*body.Body
	channel []chan []*body.Body

	// synchronizing the ticks
	signal    chan struct{}
	waitGroup *sync.WaitGroup
}

func NewSimulation(p, n int) *Simulation {
	s := &Simulation{p, nil, make([]*body.Body, n), make([]chan []*body.Body, p), make(chan struct{}, p), &sync.WaitGroup{}}
	s.startthebodies()
	for i := range s.channel {
		s.channel[i] = make(chan []*body.Body, p-1)
	}
	// start the workers
	for id := 0; id < s.p; id++ {
		go s.worker(id)
	}
	return s
}

func (s *Simulation) worker(id int) {
	// take own bodies
	//start := int(math.Floor(float64(len(s.bodies))/float64(s.p))) * id
	start := len(s.bodies) / s.p * id
	var local []*body.Body
	if id == s.p-1 {
		local = s.bodies[start:len(s.bodies)]
	} else {
		local = s.bodies[start : start+len(s.bodies)/s.p]
	}

	//send own bodies to the others
	for i := range s.channel {
		if id != i {
			s.channel[i] <- local
		}
	}

	// loop
	for {
		// wait for signal
		<-s.signal
		// add local forces
		for i := range local {
			for j := range local {
				if i != j {
					local[i].AddForce(local[j])
				}
			}
		}

		// add foreign forces
		for i := 0; i < s.p-1; i++ {
			other := <-s.channel[id]
			for i := range local {
				for j := range other {
					local[i].AddForce(other[j])
				}
			}
		}

		// update local bodies positions
		// for i := range local {
		// 	local[i].Update(timestep)
		// 	if s.img != nil {
		// 		local[i].ColorPixel(s.img)
		// 	}
		// 	local[i].ResetForce()
		// }

		// send the new local bodies positions to the others
		for p := range s.channel {
			if p != id {
				s.channel[p] <- local
			}
		}
		// signal that the work is done
		s.waitGroup.Done()
	}
}

// master
func (s *Simulation) Update(image *image.Paletted) error {
	s.img = image
	// start p workers
	s.waitGroup.Add(s.p)
	for i := 0; i < s.p; i++ {
		// give starting signals to the workers
		s.signal <- struct{}{}
	}
	// wait them all
	s.waitGroup.Wait()

	// update positions
	for _, body := range s.bodies {
		body.Update(timestep)
		if s.img != nil {
			body.ColorPixel(s.img)
		}
		body.ResetForce()
	}
	return nil
}

//Initialize N bodies with random positions and circular velocities
func (s *Simulation) startthebodies() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range s.bodies {
		s.bodies[i] = body.NewRandomBody(r)
	}
	//Put the central mass in
	s.bodies[0] = body.NewCentralBody(1e6*body.SOLARMASS, color.RGBA{255, 0, 0, 0xff}) //put a heavy body.body in the center

}
