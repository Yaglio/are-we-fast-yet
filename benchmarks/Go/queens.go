package main

// Queens solves the eight queens problem by backtracking.
type Queens struct {
	freeMaxs  []bool
	freeRows  []bool
	freeMins  []bool
	queenRows []int
}

// Benchmark runs one iteration of the benchmark.
func (q *Queens) Benchmark() any {
	result := true
	for i := 0; i < 10; i++ {
		result = result && q.queens()
	}
	return result
}

func (q *Queens) queens() bool {
	q.freeRows = make([]bool, 8)
	for i := range q.freeRows {
		q.freeRows[i] = true
	}
	q.freeMaxs = make([]bool, 16)
	for i := range q.freeMaxs {
		q.freeMaxs[i] = true
	}
	q.freeMins = make([]bool, 16)
	for i := range q.freeMins {
		q.freeMins[i] = true
	}
	q.queenRows = make([]int, 8)
	for i := range q.queenRows {
		q.queenRows[i] = -1
	}

	return q.placeQueen(0)
}

func (q *Queens) placeQueen(c int) bool {
	for r := 0; r < 8; r++ {
		if q.getRowColumn(r, c) {
			q.queenRows[r] = c
			q.setRowColumn(r, c, false)

			if c == 7 {
				return true
			}

			if q.placeQueen(c + 1) {
				return true
			}
			q.setRowColumn(r, c, true)
		}
	}
	return false
}

func (q *Queens) getRowColumn(r, c int) bool {
	return q.freeRows[r] && q.freeMaxs[c+r] && q.freeMins[c-r+7]
}

func (q *Queens) setRowColumn(r, c int, v bool) {
	q.freeRows[r] = v
	q.freeMaxs[c+r] = v
	q.freeMins[c-r+7] = v
}

// VerifyResult checks the benchmark result.
func (*Queens) VerifyResult(result any) bool {
	return result.(bool)
}
