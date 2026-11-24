package main

import (
	"fmt"
	"time"
)

// Run executes one benchmark and reports its timings.
type Run struct {
	name            string
	benchmarkSuite  func() Benchmark
	numIterations   int
	innerIterations int
	total           int64
}

func newRun(name string) *Run {
	return &Run{
		name:            name,
		benchmarkSuite:  getSuiteFromName(name),
		numIterations:   1,
		innerIterations: 1,
	}
}

func getSuiteFromName(name string) func() Benchmark {
	switch name {
	case "Bounce":
		return func() Benchmark { return &Bounce{} }
	case "List":
		return func() Benchmark { return &List{} }
	case "Mandelbrot":
		return func() Benchmark { return &Mandelbrot{} }
	default:
		panic("No benchmark found with the name: " + name)
	}
}

// RunBenchmark runs all iterations and prints the report.
func (r *Run) RunBenchmark() {
	fmt.Println("Starting " + r.name + " benchmark ...")

	r.doRuns(r.benchmarkSuite())
	r.reportBenchmark()

	fmt.Println()
}

func (r *Run) measure(bench Benchmark) {
	startTime := time.Now()
	if !innerBenchmarkLoop(bench, r.innerIterations) {
		panic("Benchmark failed with incorrect result")
	}
	runTime := time.Since(startTime).Microseconds()

	r.printResult(runTime)

	r.total += runTime
}

func (r *Run) doRuns(bench Benchmark) {
	for i := 0; i < r.numIterations; i++ {
		r.measure(bench)
	}
}

func (r *Run) reportBenchmark() {
	fmt.Printf("%s: iterations=%d average: %dus total: %dus\n\n",
		r.name, r.numIterations, r.total/int64(r.numIterations), r.total)
}

func (r *Run) printResult(runTime int64) {
	fmt.Printf("%s: iterations=1 runtime: %dus\n", r.name, runTime)
}

// PrintTotal prints the accumulated runtime.
func (r *Run) PrintTotal() {
	fmt.Printf("Total Runtime: %dus\n", r.total)
}

// SetNumIterations sets the number of measured iterations.
func (r *Run) SetNumIterations(numIterations int) {
	r.numIterations = numIterations
}

// SetInnerIterations sets the number of benchmark executions per measurement.
func (r *Run) SetInnerIterations(innerIterations int) {
	r.innerIterations = innerIterations
}
