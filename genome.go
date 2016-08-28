package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/thoj/go-galib"
)

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Command struct contains instructions for a draw operation
type Command struct {
	image  int
	x      int
	y      int
	noiseX float64
	noiseY float64
}

func (cmd *Command) Randomize() {
	cmd.image = rand.Intn(len(problem.SourceBytes))
	cmd.x = rand.Intn(problem.TargetWidth)
	cmd.y = rand.Intn(problem.TargetHeight)
	cmd.noiseX = rand.Float64()
	cmd.noiseY = rand.Float64()
}

func (cmd *Command) String() string {
	return fmt.Sprintf("{%d, %f, %f}", cmd.image, cmd.x, cmd.y)
}

// Genome struct contains the genetic information for generating a blended image
type Genome struct {
	Gene      []Command
	score     float64
	scoreFunc func(ga *Genome) float64
	hasscore  bool
	img       []byte
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
		g.Gene[idx].Randomize()
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

func linearCombine(alpha float64, bBack, bOver byte) byte {
	return byte(float64(bOver)*alpha + float64(bBack)*(1.0-alpha))
}

func applyCommand(img []byte, cmd *Command) {
	minImgX := maxInt(0, cmd.x)
	minImgY := maxInt(0, cmd.y)
	maxImgX := minInt(cmd.x+problem.SourceWidths[cmd.image], problem.TargetWidth)
	maxImgY := minInt(cmd.y+problem.SourceHeights[cmd.image], problem.TargetHeight)
	minCmdX := maxInt(0, -cmd.x)
	minCmdY := maxInt(0, -cmd.y)
	for cmdY, imgY := minCmdY, minImgY; imgY < maxImgY; cmdY, imgY = cmdY+1, imgY+1 {
		cmdIdx := (cmdY*problem.SourceWidths[cmd.image] + minCmdX) * 4
		imgIdx := (imgY*problem.TargetWidth + minImgX) * 4
		for imgX := minImgX; imgX < maxImgX; imgX++ {
			//noise := problem.Noise.Eval2(float64(subX)/100+cmd.noiseX, float64(subY)/100+cmd.noiseY)
			noise := 1.0

			for c := 0; c < 4; c++ {
				img[imgIdx+c] = linearCombine(noise, img[imgIdx+c], problem.SourceBytes[cmd.image][cmdIdx+c])
			}
			imgIdx += 4
			cmdIdx += 4
		}
	}
}

func (g *Genome) calcScore() float64 {
	g.img = make([]byte, problem.TargetWidth*problem.TargetHeight*4)
	for _, cmd := range g.Gene {
		applyCommand(g.img, &cmd)
	}

	var score float64
	for b := 0; b < problem.TargetWidth*problem.TargetHeight*4; b++ {
		score += math.Abs(float64(problem.TargetBytes[b] - g.img[b]))
	}
	return score
}

func (g *Genome) Score() float64 {
	if !g.hasscore {
		g.score = g.calcScore()
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
