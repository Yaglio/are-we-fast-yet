package main

// Sieve counts the primes up to 5000 with the sieve of Eratosthenes.
type Sieve struct{}

// Benchmark runs one iteration of the benchmark.
func (s *Sieve) Benchmark() any {
	flags := make([]bool, 5000)
	for i := range flags {
		flags[i] = true
	}
	return s.sieve(flags, 5000)
}

func (*Sieve) sieve(flags []bool, size int) int {
	primeCount := 0

	for i := 2; i <= size; i++ {
		if flags[i-1] {
			primeCount++
			k := i + i
			for k <= size {
				flags[k-1] = false
				k += i
			}
		}
	}
	return primeCount
}

// VerifyResult checks the benchmark result.
func (*Sieve) VerifyResult(result any) bool {
	return 669 == result.(int)
}
