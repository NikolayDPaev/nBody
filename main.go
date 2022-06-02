package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"math"
	"os"
	"time"

	"github.com/NikolayDPaev/n-body/body"
	"github.com/NikolayDPaev/n-body/direct"
	"github.com/NikolayDPaev/n-body/pipeline"
)

type Simulation interface {
	Update(image *image.Paletted) error
}

func createGif(steps int, simulation Simulation) {
	var w, h int = body.RESOLUTION, body.RESOLUTION

	var images []*image.Paletted
	var delays []int

	for step := 0; step < steps; step++ {
		img := image.NewPaletted(image.Rect(0, 0, w, h), body.Palette)
		images = append(images, img)
		delays = append(delays, 10)

		simulation.Update(img)
	}

	f, _ := os.OpenFile("nBody.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}

func test(simulation Simulation, steps int) int64 {
	start := time.Now()
	for i := 0; i < steps; i++ {
		simulation.Update(nil)
	}
	return time.Since(start).Microseconds()
}

func main() {
	animPtr := flag.Bool("anim", false, "animation?")
	pPtr := flag.Int("p", 1, "parallelism")
	nPtr := flag.Int("n", 1000, "bodies count")
	stepsPtr := flag.Int("steps", 1000, "steps count")
	testsPtr := flag.Int("tests", 10, "tests count")
	archPtr := flag.String("arch", "direct", "software architecture: direct or pipeline")
	flag.Parse()

	var simulation Simulation
	if *archPtr == "direct" {
		simulation = direct.NewSimulation(*pPtr, *nPtr)
	} else if *archPtr == "pipeline" {
		simulation = pipeline.NewSimulation(*pPtr, *nPtr)
	} else {
		fmt.Println("software architectures are direct and pipeline")
	}

	if *animPtr {
		createGif(*stepsPtr, simulation)
	} else {

		var min int64 = math.MaxInt64
		for i := 0; i < *testsPtr; i++ {
			now := test(simulation, *stepsPtr)
			if now < min {
				min = now
			}
		}
		fmt.Println(min)
	}
}
