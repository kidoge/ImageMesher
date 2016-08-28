package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ojrac/opensimplex-go"
)

type Problem struct {
	imageSourceDir  string
	targetImageFile string

	SourceBytes   [][]byte
	SourceWidths  []int
	SourceHeights []int

	TargetBytes  []byte
	TargetWidth  int
	TargetHeight int

	Noise *opensimplex.Noise
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
func imageToBitmap(img image.Image) []byte {
	bitmap := make([]byte, img.Bounds().Dx()*img.Bounds().Dy()*4)
	var idx int
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			bitmap[idx] = byte(r)
			bitmap[idx+1] = byte(g)
			bitmap[idx+2] = byte(b)
			bitmap[idx+3] = byte(a)
			idx += 4
		}
	}

	return bitmap
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
	prob.SourceWidths = make([]int, len(files))
	prob.SourceHeights = make([]int, len(files))
	prob.SourceBytes = make([][]byte, len(files))
	for idx, f := range files {
		fullPath := prob.imageSourceDir + "\\" + f.Name()
		img := loadImage(fullPath)

		prob.SourceWidths[idx] = img.Bounds().Dx()
		prob.SourceHeights[idx] = img.Bounds().Dy()
		prob.SourceBytes[idx] = imageToBitmap(img)
	}
	fmt.Printf("%d images loaded.\n", len(prob.SourceBytes))

	_, err = os.Stat(prob.targetImageFile)
	if err != nil {
		panic(err)
	}

	targetImage := loadImage(prob.targetImageFile)
	prob.TargetWidth = targetImage.Bounds().Dx()
	prob.TargetHeight = targetImage.Bounds().Dy()
	prob.TargetBytes = imageToBitmap(targetImage)
}
