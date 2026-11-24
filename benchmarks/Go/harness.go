// Command awfy is the Go harness for the Are We Fast Yet benchmarks.
package main

import (
	"fmt"
	"os"
	"strconv"
)

func processArguments(args []string) *Run {
	run := newRun(args[0])

	if len(args) > 1 {
		run.SetNumIterations(atoi(args[1]))
		if len(args) > 2 {
			run.SetInnerIterations(atoi(args[2]))
		}
	}

	return run
}

func atoi(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return value
}

func printUsage() {
	fmt.Println("Harness benchmark [num-iterations [inner-iter]]")
	fmt.Println()
	fmt.Println("  benchmark      - benchmark class name ")
	fmt.Println("  num-iterations - number of times to execute benchmark, default: 1")
	fmt.Println("  inner-iter     - number of times the benchmark is executed in an inner loop, ")
	fmt.Println("                   which is measured in total, default: 1")
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	run := processArguments(args)
	run.RunBenchmark()
	run.PrintTotal()
}
