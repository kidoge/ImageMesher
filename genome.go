package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/thoj/go-galib"
)

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
	cmd.image = rand.Intn(len(problem.SourceImages))
	cmd.x = rand.Intn(problem.TargetImage.Bounds().Dx())
	cmd.y = rand.Intn(problem.TargetImage.Bounds().Dy())
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
	image     *image.RGBA
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

func linearCombine(alpha float64, cImg, cOver color.Color) color.RGBA {
	r1, g1, b1, _ := cImg.RGBA()
	r2, g2, b2, _ := cOver.RGBA()
	return color.RGBA{R: uint8(float64(r2)*alpha + float64(r1)*(1.0-alpha)),
		G: uint8(float64(g2)*alpha + float64(g1)*(1.0-alpha)),
		B: uint8(float64(b2)*alpha + float64(b1)*(1.0-alpha)),
		A: 255}
}

func applyCommand(img *image.RGBA, cmd *Command) {
	imgBounds := img.Bounds()
	subBounds := problem.SourceImages[cmd.image].Bounds()
	subX := 0
	maxX := minInt(int(cmd.x)+subBounds.Dx(), imgBounds.Max.X)
	maxY := minInt(int(cmd.y)+subBounds.Dy(), imgBounds.Max.Y)
	for imgX := cmd.x; imgX < maxX; imgX++ {
		subY := 0
		for imgY := cmd.y; imgY < maxY; imgY++ {
			//noise := problem.Noise.Eval2(float64(subX)/100+cmd.noiseX, float64(subY)/100+cmd.noiseY)
			noise := 1.0
			imgColor := img.RGBAAt(imgX, imgY)
			subColor := problem.SourceImages[cmd.image].At(subX, subY)
			img.SetRGBA(imgX, imgY, linearCombine(noise, imgColor, subColor))
			subY++
		}
		subX++
	}
}
func (g *Genome) calcScore() float64 {
	bounds := problem.TargetImage.Bounds()
	g.image = image.NewRGBA(bounds)
	for _, cmd := range g.Gene {
		applyCommand(g.image, &cmd)
	}

	var score float64
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rs, gs, bs, _ := problem.TargetImage.At(x, y).RGBA()
			ri, gi, bi, _ := g.image.At(x, y).RGBA()
			score += math.Abs(float64(rs-ri)) +
				math.Abs(float64(gs-gi)) +
				math.Abs(float64(bs-bi))
		}
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
