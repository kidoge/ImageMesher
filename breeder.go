package main

import (
	"math/rand"

	"github.com/thoj/go-galib"
)

type Breeder struct{}

func (breeder *Breeder) Breed(a, b ga.GAGenome) (ca, cb ga.GAGenome) {
	minLen := a.Len()
	if b.Len() < minLen {
		minLen = b.Len()
	}

	p1 := rand.Intn(minLen)
	p2 := rand.Intn(minLen)
	if p1 > p2 {
		p1, p2 = p2, p1
	}
	ca, cb = a.Crossover(b, p1, p2)
	return
}

func (b *Breeder) String() string { return "Breeder" }
