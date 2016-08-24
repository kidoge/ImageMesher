package main

import (
	"fmt"
	"math/rand"

	"github.com/thoj/go-galib"
)

// Command struct contains instructions for a draw operation
type Command struct {
	image int
	x     float32
	y     float32
}

// Genome struct contains the genetic information for generating a blended image
type Genome struct {
	Gene      []Command
	score     float64
	scoreFunc func(ga *Genome) float64
	hasscore  bool
}

// NewGenome creates a new genome.
func NewGenome(cmds []Command) *Genome {
	g := new(Genome)
	g.Gene = cmds
	return g
}

// Crossover mixes genes from two genomes.
func (a *Genome) Crossover(bi ga.GAGenome, p1, p2 int) (ga.GAGenome, ga.GAGenome) {
	ca := a.Copy().(*Genome)
	b := bi.(*Genome)
	cb := b.Copy().(*Genome)
	copy(ca.Gene[p1:p2+1], b.Gene[p1:p2+1])
	copy(cb.Gene[p1:p2+1], a.Gene[p1:p2+1])
	ca.Reset()
	cb.Reset()
	return ca, cb
}

func (a *Genome) Splice(bi ga.GAGenome, from, to, length int) {
	b := bi.(*Genome)
	copy(a.Gene[to:length+to], b.Gene[from:length+from])
	a.Reset()
}

func (g *Genome) Valid() bool {
	//TODO: Make this
	return true
}

func (g *Genome) Switch(x, y int) {
	g.Gene[x], g.Gene[y] = g.Gene[y], g.Gene[x]
	g.Reset()
}

func (g *Genome) Randomize() {
	l := len(g.Gene)
	for idx := 0; idx < l; idx++ {
		g.Gene[idx].image = rand.Intn(len(problem.SourceImages))
		g.Gene[idx].x = rand.Float32()*200.0 - 100.0
		g.Gene[idx].y = rand.Float32()*200.0 - 100.0
	}
	g.Reset()
}

func (g *Genome) Copy() ga.GAGenome {
	n := new(Genome)
	n.Gene = make([]Command, len(g.Gene))
	copy(n.Gene, g.Gene)
	n.score = g.score
	n.hasscore = g.hasscore
	return n
}

func (g *Genome) Len() int {
	return len(g.Gene)
}

func (g *Genome) Score() float64 {
	if !g.hasscore {
		g.score = rand.Float64()
		g.hasscore = true
	}
	return g.score
}

func (g *Genome) Reset() {
	g.hasscore = false
}

func (g *Genome) String() string {
	return fmt.Sprintf("%v", g.Gene)
}

func (cmd *Command) String() string {
	return fmt.Sprintf("{%d, %f, %f}", cmd.image, cmd.x, cmd.y)
}
