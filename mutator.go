package main

import (
	"fmt"
	"math/rand"

	"github.com/thoj/go-galib"
)

// Mutator Basic mutator that will perform one of the following:
// - Change an image to a random image
// - Move an image (x, y)
// - Change the alpha gradients for an image (x, y)
// - Add an image
// - Remove an image
type Mutator struct {
	posStddev float64
}

// Mutate perform mutation
func (m Mutator) Mutate(a ga.GAGenome) ga.GAGenome {
	n := a.Copy().(*Genome)
	idx := (rand.Intn(len(n.Gene)) / 5) * 5
	n.Gene[idx].image = rand.Intn(100)
	n.Gene[idx].x += float32(rand.NormFloat64() * m.posStddev)
	n.Gene[idx].y += float32(rand.NormFloat64() * m.posStddev)
	n.Reset()
	fmt.Println("mutating", n.Score())
	return n
}

func (m Mutator) String() string {
	return "Mutator"
}
