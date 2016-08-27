package main

import (
	"fmt"

	"github.com/thoj/go-galib"
	//"github.com/ojrac/opensimplex-go"
)

const imageDir = "~\\testimages"
const targetFile = "~\\target.png"

var scores int

var problem *Problem

func main() {
	problem = NewProblem(imageDir, targetFile)
	problem.Load()

	param := ga.GAParameter{
		Initializer: new(ga.GARandomInitializer),
		Selector:    ga.NewGATournamentSelector(0.7, 5),
		Breeder:     new(Breeder),
		Mutator: &Mutator{
			PLengthChange: 1.0,
			PosStdev:      0.1,
		},
		PMutate: 0.1,
		PBreed:  0.7,
	}

	gao := ga.NewGA(param)

	genome := NewGenome([]Command{{
		image: 0,
		x:     0,
		y:     0,
	}})

	gao.Init(100, genome)
	gao.Optimize(100)
	gao.PrintTop(10)

	fmt.Printf("Best: %f\n", gao.Best().Score())
}
