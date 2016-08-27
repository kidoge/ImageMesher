package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ojrac/opensimplex-go"
)

type Problem struct {
	imageSourceDir  string
	targetImageFile string
	SourceImages    []image.Image
	TargetImage     image.Image
	Noise           *opensimplex.Noise
}

func NewProblem(sourceDir, targetFile string) *Problem {
	p := new(Problem)

	p.imageSourceDir = strings.Replace(sourceDir, "~", Homedir(), 1)
	p.targetImageFile = strings.Replace(targetFile, "~", Homedir(), 1)
	p.Noise = opensimplex.New()
	return p
}

func loadImage(path string) image.Image {
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func (prob *Problem) Load() {
	// check source path is a directory
	dir, err := os.Stat(prob.imageSourceDir)
	if err != nil || !dir.IsDir() {
		panic(err)
	}

	// import source files
	files, _ := ioutil.ReadDir(prob.imageSourceDir)
	fmt.Printf("Found %d files in %s\n", len(files), prob.imageSourceDir)
	prob.SourceImages = make([]image.Image, 0)
	for _, f := range files {
		fullPath := prob.imageSourceDir + "\\" + f.Name()
		prob.SourceImages = append(prob.SourceImages, loadImage(fullPath))
	}
	fmt.Printf("%d images loaded.\n", len(prob.SourceImages))

	_, err = os.Stat(prob.targetImageFile)
	if err != nil {
		panic(err)
	}

	prob.TargetImage = loadImage(prob.targetImageFile)
}
