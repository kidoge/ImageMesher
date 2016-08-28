package main

import (
	"math/rand"

	"github.com/thoj/go-galib"
)

// Mutator is a struct that will change the Genome in one of the following ways:
// - Change an image to a random image
// - Move an image (x, y)
// - Change the alpha gradients for an image (x, y)
// - Add an image
// - Remove an image
type Mutator struct {
	PImageChange  float64
	PLengthChange float64
	PosStdev      float64
	NoiseStdev    float64
}

// alter modifies an existing
func (m Mutator) alter(genome *Genome) {
	idx := rand.Intn(len(genome.Gene))
	if rand.Float64() < m.PImageChange {
		genome.Gene[idx].image = rand.Intn(len(problem.SourceBytes))
	}
	genome.Gene[idx].x += int(rand.NormFloat64() * m.PosStdev)
	genome.Gene[idx].y += int(rand.NormFloat64() * m.PosStdev)
	genome.Gene[idx].noiseX += rand.NormFloat64() * m.NoiseStdev
	genome.Gene[idx].noiseY += rand.NormFloat64() * m.NoiseStdev
}

func (m Mutator) addRandom(genome *Genome) {
	idx := rand.Intn(len(genome.Gene) + 1)
	cmd := new(Command)
	cmd.Randomize()

	// insert to slice
	genome.Gene = append(genome.Gene, Command{})
	copy(genome.Gene[idx+1:], genome.Gene[idx:])
	genome.Gene[idx] = *cmd
}

func (m Mutator) removeRandom(genome *Genome) {
	idx := rand.Intn(len(genome.Gene))
	genome.Gene = append(genome.Gene[:idx], genome.Gene[idx+1:]...)
}

// Mutate perform mutation
func (m Mutator) Mutate(a ga.GAGenome) ga.GAGenome {
	newGen := a.Copy().(*Genome)

	if rand.Float64() < m.PLengthChange {
		if newGen.Len() == 1 || rand.Float32() < 0.5 {
			m.addRandom(newGen)
		} else {
			m.removeRandom(newGen)
		}
	} else {
		m.alter(newGen)
	}

	newGen.Reset()
	return newGen
}

func (m Mutator) String() string {
	return "Mutator"
}
