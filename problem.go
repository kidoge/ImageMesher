package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

func homedir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

type Problem struct {
	imageSourceDir  string
	targetImageFile string
	SourceImages    []image.Image
	TargetImage     image.Image
}

func NewProblem(sourceDir, targetFile string) *Problem {
	p := new(Problem)

	p.imageSourceDir = strings.Replace(sourceDir, "~", homedir(), 1)
	p.targetImageFile = strings.Replace(targetFile, "~", homedir(), 1)
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
		panic("Cannot read source directory: " + prob.imageSourceDir)
	}

	// import source files
	files, _ := ioutil.ReadDir(prob.imageSourceDir)
	prob.SourceImages = make([]image.Image, len(files))
	for _, f := range files {
		prob.SourceImages = append(prob.SourceImages, loadImage(prob.imageSourceDir+"\\"+f.Name()))
	}

	_, err = os.Stat(prob.targetImageFile)
	if err != nil {
		panic("Cannot read target file: " + prob.targetImageFile)
	}

	prob.TargetImage = loadImage(prob.targetImageFile)
}
