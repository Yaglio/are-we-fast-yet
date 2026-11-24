package som

import "testing"

type item struct{ n int }

func TestVectorAppendGrowth(t *testing.T) {
	v := NewVector[int]()
	if v.Capacity() != 0 || !v.IsEmpty() {
		t.Fatalf("new vector: capacity %d, empty %v", v.Capacity(), v.IsEmpty())
	}
	for i := 0; i < 11; i++ {
		v.Append(i)
	}
	if v.Size() != 11 || v.Capacity() != 20 {
		t.Fatalf("size %d capacity %d, want 11 and 20", v.Size(), v.Capacity())
	}
	for i := 0; i < 11; i++ {
		if v.At(i) != i {
			t.Fatalf("At(%d) = %d", i, v.At(i))
		}
	}
	if v.At(20) != 0 {
		t.Fatalf("At beyond storage should be zero value")
	}
}

func TestVectorAtPut(t *testing.T) {
	v := NewVector[int]()
	v.AtPut(3, 7)
	if v.Capacity() != 10 || v.Size() != 4 || v.At(3) != 7 {
		t.Fatalf("capacity %d size %d at3 %d", v.Capacity(), v.Size(), v.At(3))
	}
	v.AtPut(45, 1)
	if v.Capacity() != 80 || v.Size() != 46 || v.At(3) != 7 || v.At(45) != 1 {
		t.Fatalf("after growth: capacity %d size %d", v.Capacity(), v.Size())
	}
	w := NewVector[int]()
	w.AtPut(12, 1)
	if w.Capacity() != 13 {
		t.Fatalf("max(idx+1, initialSize): capacity %d", w.Capacity())
	}
}

func TestVectorFirstRemoveFirst(t *testing.T) {
	a, b := &item{1}, &item{2}
	v := VectorWith(a)
	if v.Capacity() != 1 {
		t.Fatalf("VectorWith capacity %d", v.Capacity())
	}
	v.Append(b)
	if v.First() != a || v.RemoveFirst() != a || v.First() != b || v.Size() != 1 {
		t.Fatalf("first/removeFirst order wrong")
	}
	if v.RemoveFirst() != b || !v.IsEmpty() || v.RemoveFirst() != nil || v.First() != nil {
		t.Fatalf("empty vector should return nil")
	}
}

func TestVectorRemove(t *testing.T) {
	a, b, c := &item{1}, &item{2}, &item{3}
	v := NewVector[*item]()
	if v.Remove(a) {
		t.Fatalf("remove on empty vector")
	}
	v.Append(a)
	v.Append(b)
	v.Append(c)
	v.Append(b)
	v.RemoveFirst()
	if !v.Remove(b) || v.Size() != 1 || v.At(0) != c || v.Capacity() != 10 {
		t.Fatalf("remove: size %d", v.Size())
	}
	if v.Remove(a) || v.Size() != 1 {
		t.Fatalf("remove of absent element")
	}
}

func TestVectorRemoveAll(t *testing.T) {
	v := NewVectorSize[int](4)
	v.Append(1)
	v.Append(2)
	v.RemoveAll()
	if !v.IsEmpty() || v.Capacity() != 4 || v.At(0) != 0 {
		t.Fatalf("removeAll: size %d capacity %d", v.Size(), v.Capacity())
	}
}

func TestVectorHasSomeGetOne(t *testing.T) {
	v := NewVector[int]()
	sum := 0
	for i := 1; i <= 5; i++ {
		v.Append(i)
	}
	v.ForEach(func(e int) { sum += e })
	if sum != 15 {
		t.Fatalf("forEach sum %d", sum)
	}
	if !v.HasSome(func(e int) bool { return e == 4 }) || v.HasSome(func(e int) bool { return e > 5 }) {
		t.Fatalf("hasSome")
	}
	if v.GetOne(func(e int) bool { return e > 2 }) != 3 || v.GetOne(func(e int) bool { return e > 5 }) != 0 {
		t.Fatalf("getOne")
	}
}

func TestVectorSort(t *testing.T) {
	cmp := func(a, b int) int { return a - b }
	v := VectorWith(1)
	v.Sort(cmp) // n <= 1: no-op
	NewVector[int]().Sort(cmp)

	w := NewVector[int]()
	w.Append(1)
	w.Append(2)
	w.Sort(cmp) // already ordered, swap not reached
	if w.At(0) != 1 || w.At(1) != 2 {
		t.Fatalf("sort changed ordered vector")
	}

	defer func() {
		if r := recover(); r != "Not Implemented" {
			t.Fatalf("expected Not Implemented panic, got %v", r)
		}
	}()
	w.Sort(func(a, b int) int { return b - a })
}
