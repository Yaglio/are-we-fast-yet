package main

// Benchmark is implemented by every AWFY benchmark.
type Benchmark interface {
	Benchmark() any
	VerifyResult(result any) bool
}

// innerLooper is implemented by benchmarks that override the default inner
// benchmark loop (Java: overriding Benchmark.innerBenchmarkLoop).
type innerLooper interface {
	InnerBenchmarkLoop(innerIterations int) bool
}

// innerBenchmarkLoop is the default inner loop (Java: Benchmark.innerBenchmarkLoop).
func innerBenchmarkLoop(bench Benchmark, innerIterations int) bool {
	if looper, ok := bench.(innerLooper); ok {
		return looper.InnerBenchmarkLoop(innerIterations)
	}
	for i := 0; i < innerIterations; i++ {
		if !bench.VerifyResult(bench.Benchmark()) {
			return false
		}
	}
	return true
}
