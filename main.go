package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"os/user"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/thoj/go-galib"
)

const imageDir = "~\\testimages"
const targetFile = "~\\target.png"

const dumpIncrement = 2
const genGoal = 100

var scores int

var problem *Problem

func Homedir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func bitmapToImage(bitmap []byte, width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	var idx int
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetRGBA(x, y, color.RGBA{
				R: bitmap[idx],
				G: bitmap[idx+1],
				B: bitmap[idx+2],
				A: bitmap[idx+3]})
			idx += 4
		}
	}
	return img
}

func saveImage(outputPath string, img []byte) {
	absoluteOutFile := strings.Replace(outputPath, "~", Homedir(), 1)
	f, err := os.OpenFile(absoluteOutFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, bitmapToImage(img, problem.TargetWidth, problem.TargetHeight))
	if err != nil {
		panic(err)
	}
	f.Close()
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	rand.Seed(time.Now().UnixNano())

	problem = NewProblem(imageDir, targetFile)
	problem.Load()

	param := ga.GAParameter{
		Initializer: new(ga.GARandomInitializer),
		Selector:    ga.NewGATournamentSelector(0.7, 5),
		Breeder:     new(Breeder),
		Mutator: &Mutator{
			PLengthChange: 0.1,
			PosStdev:      100.0,
			NoiseStdev:    20.0,
		},
		PMutate: 0.3,
		PBreed:  0.7,
	}

	gao := ga.NewGA(param)

	var cmds []Command
	for idx := 0; idx < 40; idx++ {
		cmds = append(cmds, Command{
			image: idx,
			x:     0,
			y:     0,
		})
	}
	genome := NewGenome(cmds)

	gao.Init(20, genome)
	startTime := time.Now().UnixNano()
	for gen := 0; gen < genGoal; gen += dumpIncrement {
		gao.Optimize(dumpIncrement)
		gao.PrintTop(10)
		saveImage(fmt.Sprintf("~/output/g%06d.png", gen), gao.Best().(*Genome).OverlayCmds())
	}
	endTime := time.Now().UnixNano()
	elapsed := float64(endTime-startTime) / float64(time.Second)
	fmt.Printf("Best: %f\n", gao.Best().Score())
	fmt.Printf("Elapsed : %fs (avg: %fs/gen)\n", elapsed, elapsed/genGoal)

}
