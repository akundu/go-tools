package main

import (
	"bytes"
	"flag"
	"fmt"
	"math"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/akundu/utilities/jobs"
	"github.com/akundu/utilities/logger"
)

func cleanup() {
	if transport != nil {
		transport.Close()
	}
}

func runJobs() {
	defer func() {
		if r := recover(); r != nil {
			logger.Error.Println("got error", r)
		}
	}()

	gj := jobs.GoJob{}
	gj.AddTask(&GoTaskTest{})
	results, delta_time := gj.Run()

	failure_count := len(results)
	for _, task_result := range results {
		//extract the time taken by the task
		if time_taken, err := task_result.At(1); err != nil {
			logger.Error.Println("error = ", err)
			continue
		} else {
			logger.Info.Printf("took %s", time_taken)
		}

		//extract the task result
		task_r, err := task_result.At(0)
		if err != nil {
			continue
		}

		//extract the status code
		go_task_result := task_r.(*jobs.GoTaskResult)
		if go_task_result.Status != nil {
			logger.Error.Println(" with error = ", go_task_result.Status.(error))
			continue
		}

		//extract the response
		if go_task_result.Value == nil {
			logger.Error.Println("body value is nil")
			continue
		}
		//response_bytes := go_task_result.Value.([]byte)
		//logger.Info.Println(" with response = ", string(response_bytes[:]))
		failure_count--
	}
	//fmt.Println("total time taken: ", int64(delta_time/time.Millisecond), "ms")
	fmt.Println("total time taken: ", delta_time)
	fmt.Println("succ %: ", int64((len(results)-failure_count)/len(results)*100))
}

func main() {
	defer cleanup()
	flag.Parse()

	runJobs()
}

func init() {
	logger.DefaultLoggerInit()
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
}

var (
	percentThreshold          = Percentiles{}
	num_to_run                = flag.Int("num_to_run", 1, "num cpus to run on")
	num_to_run_simultaneously = flag.Int("num_to_run_simultaneously", 1, "num cpus to run on")
	//prefixGauges     = flag.String("prefixGauges", "gauges.", "Prefix for all gauges stats")
)

type Percentiles []*Percentile
type Percentile struct {
	float float64
	str   string
}

func (a *Percentiles) Set(s string) error {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*a = append(*a, &Percentile{f, strings.Replace(s, ".", "_", -1)})
	return nil
}
func (p *Percentile) String() string {
	return p.str
}
func (a *Percentiles) String() string {
	return fmt.Sprintf("%v", *a)
}

type Percentiles []*Percentile

func processTimers(buffer *bytes.Buffer, now int64, pctls Percentiles) int64 {
	var num int64
	for u, t := range timers {
		num++

		sort.Sort(t)
		min := t[0]
		max := t[len(t)-1]
		median := t[uint64(len(t)/2)]
		maxAtThreshold := max
		count := len(t)

		sum := float64(0)
		for _, value := range t {
			sum += value
		}
		mean := float64(sum) / float64(len(t))

		for _, pct := range pctls {
			if len(t) > 1 {
				var abs float64
				if pct.float >= 0 {
					abs = pct.float
				} else {
					abs = 100 + pct.float
				}
				// poor man's math.Round(x):
				// math.Floor(x + 0.5)
				indexOfPerc := int(math.Floor(((abs / 100.0) * float64(count)) + 0.5))
				if pct.float >= 0 {
					indexOfPerc -= 1 // index offset=0
				}
				maxAtThreshold = t[indexOfPerc]
			}

			var tmpl string
			var pctstr string
			if pct.float >= 0 {
				tmpl = "%s.upper_%s %0.2f %d\n"
				pctstr = pct.str
			} else {
				tmpl = "%s.lower_%s %0.2f %d\n"
				pctstr = pct.str[1:]
			}
			fmt.Fprintf(buffer, tmpl, u, pctstr, maxAtThreshold, now)
		}

		fmt.Fprintf(buffer, "%s.mean %0.2f %d\n", u, mean, now)
		fmt.Fprintf(buffer, "%s.median %0.2f %d\n", u, median, now)
		fmt.Fprintf(buffer, "%s.upper %0.2f %d\n", u, max, now)
		fmt.Fprintf(buffer, "%s.lower %0.2f %d\n", u, min, now)
		fmt.Fprintf(buffer, "%s.count %d %d\n", u, count, now)

		delete(timers, u)
	}
	return num
}
