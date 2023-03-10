# N-body

Testing of two parallel algorithms for simulating n bodies in gravitational field

## Running
```go run main.go```
  With command line flags:
 - ```anim``` - wether to generate animation gif or not, default = false
 - ```p``` - number of threads to use, default = 1
 - ```n``` - number of bodies to simulate, default = 1000
 - ```steps``` - number of steps for the simulation, default = 1000
 - ```tests``` - number of tests, to take the min time,  default = 10
 - ```arch``` - type of algorithm to use: "direct" or "pipeline", default = "direct"

The program outputs the time in microseconds

## Algorithms

 - "direct" - Each worker is assigned a number of bodies. On each tick every worker thread computes the new position of their own bodies and sends the new coordinates to the other workers. The computation of the new position is done by adding the forces of the other bodies to the current.

 - "pipeline" - The bodies are send through a pipeline consisting of the workers. Each worker computes the forces on their respective bodies. First the worker takes its bodies and calculates the forces on each other as well as the forces of the other bodies traveling through the pipeline. Then it sends its bodies to the next worker.
 The master updates the positions and sends the bodies again for the next tick.

## Animation
![](https://github.com/NikolayDPaev/nBody/blob/master/nBody.gif)

## Results
Speedup on 16 processors and 3000 bodies  
![](https://github.com/NikolayDPaev/nBody/blob/master/Speedup.PNG)
