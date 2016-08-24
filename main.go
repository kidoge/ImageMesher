package main

import (
	"fmt"

	"github.com/thoj/go-galib"
	//"github.com/ojrac/opensimplex-go"
)

var scores int

func main() {
	fmt.Println("Start")

	param := ga.GAParameter{
		Initializer: new(ga.GARandomInitializer),
		Selector:    ga.NewGATournamentSelector(0.7, 5),
		Breeder:     new(ga.GA2PointBreeder),
		Mutator:     new(Mutator),
		PMutate:     0.1,
		PBreed:      0.7,
	}

	gao := ga.NewGA(param)

	genome := NewGenome([]Command{{
		image: 0,
		x:     0,
		y:     0,
	}})

	gao.Init(100, genome)
	gao.OptimizeUntil(func(best ga.GAGenome) bool {
		return best.Score() < 0.1
	})
	gao.PrintTop(10)

	fmt.Printf("Best: %f\n", gao.Best().Score())
}
