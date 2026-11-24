package main

import "awfy/som"

// arrayTree stands for Java's Object[], whose slots hold further Object[] (or null in leaves).
type arrayTree []arrayTree

// Storage builds a tree of arrays to stress allocation.
type Storage struct {
	count int
}

// Benchmark runs one iteration of the benchmark.
func (s *Storage) Benchmark() any {
	random := som.NewRandom()
	s.count = 0
	s.buildTreeDepth(7, random)
	return s.count
}

func (s *Storage) buildTreeDepth(depth int, random *som.Random) arrayTree {
	s.count++
	if depth == 1 {
		return make(arrayTree, random.Next()%10+1)
	}
	arr := make(arrayTree, 4)
	for i := range arr {
		arr[i] = s.buildTreeDepth(depth-1, random)
	}
	return arr
}

// VerifyResult checks the benchmark result.
func (*Storage) VerifyResult(result any) bool {
	return 5461 == result.(int)
}
