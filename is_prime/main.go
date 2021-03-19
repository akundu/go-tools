package main

import (
	"flag"
	"runtime"

	"github.com/akundu/utilities/logger"
)

var (
	numToRun_till          = flag.Int("numToRun_till", 100, "prime number to compute till")
	numToRun               = flag.Int("numToRun", 1, "num cpus to run on")
	numToRunSimultaneously = flag.Int("numToRunSimultaneously", runtime.NumCPU()-1, "num cpus to run on")
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(*numToRunSimultaneously)

	//resultStrings := runJobsAllInOneShot(*numToRun_till)
	resultStrings := isPrimeInIterations(*numToRun_till, *numToRunSimultaneously)
	for _, element := range resultStrings {
		logger.Info.Println(element)
	}
}

func init() {
	logger.DefaultLoggerInit()
}
