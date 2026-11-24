package main

// Permute generates all permutations of a six-element array recursively.
type Permute struct {
	count int
	v     []int
}

// Benchmark runs one iteration of the benchmark.
func (p *Permute) Benchmark() any {
	p.count = 0
	p.v = make([]int, 6)
	p.permute(6)
	return p.count
}

func (p *Permute) permute(n int) {
	p.count++
	if n != 0 {
		n1 := n - 1
		p.permute(n1)
		for i := n1; i >= 0; i-- {
			p.swap(n1, i)
			p.permute(n1)
			p.swap(n1, i)
		}
	}
}

func (p *Permute) swap(i, j int) {
	tmp := p.v[i]
	p.v[i] = p.v[j]
	p.v[j] = tmp
}

// VerifyResult checks the number of permute calls.
func (*Permute) VerifyResult(result any) bool {
	return result.(int) == 8660
}
