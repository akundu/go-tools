package main

import (
	"fmt"

	"github.com/akundu/utilities/jobs"
	"github.com/akundu/utilities/logger"
)

type taskCheckPrime struct {
	num int
}

func newTaskCheckPrime(num int) *taskCheckPrime {
	return &taskCheckPrime{num}
}

func (this *taskCheckPrime) Run() *jobs.GoTaskResult {
	return &jobs.GoTaskResult{
		Value:  checkPrime(this.num),
		Status: nil,
	}
}

func runJobsAllInOneShot(numToRunTill int) []string {
	defer func() {
		if r := recover(); r != nil {
			logger.Error.Println("got error", r)
		}
	}()
	resultStrings := make([]string, 0, numToRunTill)

	gj := jobs.GoJob{}
	for index := 0; index < numToRunTill; index++ {
		gj.AddTask(newTaskCheckPrime(index))
	}
	results, deltaTime := gj.Run()

	failureCount := len(results)
	for _, taskResult := range results {
		//extract the time taken by the task
		if timeTaken, err := taskResult.At(1); err != nil {
			logger.Error.Println("error = ", err)
			continue
		} else {
			logger.Trace.Printf("took %s", timeTaken)
		}

		//extract the task result
		taskR, err := taskResult.At(0)
		if err != nil {
			continue
		}

		//extract the status code
		goTaskResult := taskR.(*jobs.GoTaskResult)
		if goTaskResult.Status != nil {
			logger.Error.Println(" with error = ", goTaskResult.Status.(error))
			continue
		}

		//extract the response
		if goTaskResult.Value == nil {
			logger.Error.Println("body value is nil")
			continue
		}

		taskInput, err := taskResult.At(2)
		if err != nil {
			logger.Info.Println(goTaskResult.Value)
		} else {
			originalTask, ok := taskInput.(*taskCheckPrime)
			if ok == false {
				logger.Error.Println("invalid task type passed in")
			} else {
				if goTaskResult.Value == true {
					resultStrings = append(resultStrings, fmt.Sprintf("%d %v", originalTask.num, goTaskResult.Value))
				}
			}
		}
		failureCount--
	}
	//fmt.Println("total time taken: ", int64(deltaTime/time.Millisecond), "ms")
	logger.Info.Println("total time taken: ", deltaTime)
	logger.Info.Println("succ %: ", int64((len(results)-failureCount)/len(results)*100))

	return resultStrings
}
