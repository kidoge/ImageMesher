package main

import (
	"fmt"
	"image"
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
	TargetImage     []image.Image
}

func NewProblem(sourceDir, targetFile string) *Problem {
	p := new(Problem)

	p.imageSourceDir = strings.Replace(sourceDir, "~", homedir(), 1)
	p.targetImageFile = strings.Replace(targetFile, "~", homedir(), 1)
	return p
}

func (prob *Problem) Load() {
	dir, err := os.Stat(prob.imageSourceDir)
	if err != nil || !dir.IsDir() {
		panic("Cannot read source directory: " + prob.imageSourceDir)
	}

	files, _ := ioutil.ReadDir(prob.imageSourceDir)
	for _, f := range files {
		fmt.Println(f.Name())
	}
	prob.SourceImages = make([]image.Image, 16)

	_, err = os.Stat(prob.targetImageFile)
	if err != nil {
		panic("Cannot read target file: " + prob.targetImageFile)
	}
}
