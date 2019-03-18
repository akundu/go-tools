package main

import (
	"fmt"

	"github.com/akundu/utilities"
	"github.com/akundu/utilities/RTJobRunner"
	"github.com/akundu/utilities/logger"
)

type isPrimeJobResult struct {
	input  int
	result bool
}

type isPrimeJobRequest struct {
	name string
	num  int
}

func (this isPrimeJobRequest) GetName() string {
	return this.name
}

type worker struct{}

func CreateWorker() RTJobRunner.Worker {
	return &worker{}
}
func (this *worker) PostRun() {}
func (this *worker) PreRun()  {}
func (this *worker) Run(id int, jh *RTJobRunner.JobHandler) {
	//for jobInfo := range jobs {
	for j := jh.GetJob(); j != nil; j = jh.GetJob() {
		job, ok := j.Req.(*isPrimeJobRequest)
		if ok == false {
			j.Resp = &RTJobRunner.BasicResponseResult{
				Err:    utilities.NewBasicError("object cant cast properly"),
				Result: nil,
			}
			jh.DoneJob(j)
			logger.Error.Printf("got error while processing %v\n", job)
			continue
		}
		result := checkPrime(job.num)
		j.Resp = &RTJobRunner.BasicResponseResult{
			Err:    nil,
			Result: &isPrimeJobResult{input: job.num, result: result},
		}
		jh.DoneJob(j)
	}
}

func isPrimeInIterations(numToGet, numToRunSimultaneously int) []string {
	jh := RTJobRunner.NewJobHandler(numToRunSimultaneously, CreateWorker)
	jh.SetPrintIndividualResults(false)
	jh.SetPrintStatistics(false)
	for index := 0; index < numToGet; index++ {
		jh.AddJob(RTJobRunner.NewRTRequestResultObject(&isPrimeJobRequest{
			name: fmt.Sprintf("%d", index),
			num:  index,
		}))
	}
	//mark that there are no more jobs to add
	jh.DoneAddingJobs()
	//wait for the results
	jh.WaitForJobsToComplete()

	resultString := make([]string, 0, numToGet)
	for _, element := range jh.Jobs {
		if element == nil {
			continue
		}
		result := element.Resp.GetResult().(*isPrimeJobResult)
		if result.result == true {
			resultString = append(resultString, fmt.Sprintf("%d %v", result.input, result.result))
		}
	}
	return resultString
}
