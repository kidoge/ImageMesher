package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
	"os/user"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/thoj/go-galib"
	//"github.com/ojrac/opensimplex-go"
)

const imageDir = "~\\testimages"
const targetFile = "~\\target.png"
const dumpIncrement = 5

var scores int

var problem *Problem

func Homedir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func saveImage(outputPath string, img image.Image) {
	absoluteOutFile := strings.Replace(outputPath, "~", Homedir(), 1)
	f, err := os.OpenFile(absoluteOutFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
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
			PImageChange:  0.01,
			PLengthChange: 0.5,
			PosStdev:      20.0,
			NoiseStdev:    20.0,
		},
		PMutate: 0.2,
		PBreed:  0.8,
	}

	gao := ga.NewGA(param)

	genome := NewGenome([]Command{{
		image: 0,
		x:     0,
		y:     0,
	}})

	gao.Init(20, genome)
	for gen := 0; gen < 100; gen += dumpIncrement {
		gao.Optimize(dumpIncrement)
		gao.PrintTop(10)
		saveImage(fmt.Sprintf("~/output/g%06d.png", gen), gao.Best().(*Genome).image)
	}

	fmt.Printf("Best: %f\n", gao.Best().Score())

}
