package main

import (
	"fmt"

	"github.com/thoj/go-galib"
	//"github.com/ojrac/opensimplex-go"
)

var scores int

// Boring fitness/score function.
func score(g *ga.GAOrderedIntGenome) float64 {
	var total int
	for i, c := range g.Gene {
		total += i * c
	}
	scores++
	return float64(total)
}

func main() {
	fmt.Println("Start")

	m := ga.NewMultiMutator()
	msh := new(ga.GAShiftMutator)
	msw := new(ga.GASwitchMutator)
	m.Add(msh)
	m.Add(msw)

	param := ga.GAParameter{
		Initializer: new(ga.GARandomInitializer),
		Selector:    ga.NewGATournamentSelector(0.7, 5),
		Breeder:     new(ga.GA2PointBreeder),
		Mutator:     m,
		PMutate:     0.1,
		PBreed:      0.7,
	}

	gao := ga.NewGA(param)

	genome := ga.NewOrderedIntGenome([]int{10, 11, 12, 13, 14, 15, 16, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}, score)

	gao.Init(100, genome)
	gao.OptimizeUntil(func(best ga.GAGenome) bool {
		return best.Score() <= 680
	})
	gao.PrintTop(10)

	fmt.Printf("Best: %f\n", gao.Best().Score())
}
