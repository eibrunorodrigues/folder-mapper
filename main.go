package main

import (
	"log"
	"os"
	"runtime"

	"github.com/eibrunorodrigues/folder-mapper/publisher"
)

func initApplication() {
	fatalErrors := make(chan string)

	validateArgs()
	workersOchestrer(&fatalErrors)

	for errors := range fatalErrors {
		log.Fatalln(errors)
	}
}

func workersOchestrer(fatalErrors *chan string) {
	numOfCpus := runtime.NumCPU()

	for cpuNum := 0; cpuNum < numOfCpus; cpuNum++ {
		go publisher.Worker(fatalErrors, cpuNum)
	}
}

func validateArgs() {
	arguments := os.Args[1:]

	if len(arguments) == 0 {
		log.Fatalln("You should pass somethig")
	}

	if _, err := os.Stat(arguments[0]); os.IsNotExist(err) {
		log.Fatalln("The argument " + arguments[0] + " is not a valid folder")
	}
}

func main() {
	initApplication()
}
