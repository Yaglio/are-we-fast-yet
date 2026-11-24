// Package som contains the shared SOM class library used by the benchmarks.
package som

// Random is the deterministic pseudo-random number generator of the SOM library.
type Random struct {
	seed int
}

// NewRandom returns a generator with the standard AWFY seed.
func NewRandom() *Random {
	return &Random{seed: 74755}
}

// Next returns the next pseudo-random number in [0, 65535].
func (r *Random) Next() int {
	r.seed = ((r.seed * 1309) + 13849) & 65535
	return r.seed
}
