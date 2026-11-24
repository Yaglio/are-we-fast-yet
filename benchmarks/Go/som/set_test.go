package som

import "testing"

func TestSetAddDeduplicatesValues(t *testing.T) {
	s := NewSet[int]()
	for i := 0; i < 15; i++ {
		s.Add(i)
		s.Add(i)
	}
	if s.Size() != 15 {
		t.Fatalf("size %d, want 15", s.Size())
	}
	next := 0
	s.ForEach(func(e int) {
		if e != next {
			t.Fatalf("ForEach got %d, want %d", e, next)
		}
		next++
	})
	if !s.Contains(14) || s.Contains(15) {
		t.Fatalf("Contains wrong")
	}
	if !s.HasSome(func(e int) bool { return e > 13 }) || s.HasSome(func(e int) bool { return e < 0 }) {
		t.Fatalf("HasSome wrong")
	}
	if s.GetOne(func(e int) bool { return e > 5 }) != 6 {
		t.Fatalf("GetOne wrong")
	}
	s.RemoveAll()
	if s.Size() != 0 || s.Contains(0) {
		t.Fatalf("RemoveAll left elements")
	}
	s.Add(3)
	if s.Size() != 1 {
		t.Fatalf("add after RemoveAll: size %d", s.Size())
	}
}

func TestIdentitySetComparesIdentity(t *testing.T) {
	s := NewIdentitySet[*item]()
	a, b := &item{1}, &item{1}
	s.Add(a)
	s.Add(a)
	s.Add(b)
	if s.Size() != 2 {
		t.Fatalf("size %d, want 2", s.Size())
	}
	if !s.Contains(a) || !s.Contains(b) || s.Contains(&item{1}) {
		t.Fatalf("Contains should use identity")
	}
	if s.GetOne(func(e *item) bool { return e.n == 2 }) != nil {
		t.Fatalf("GetOne without match should be nil")
	}
	for i := 0; i < 20; i++ {
		s.Add(&item{i})
	}
	if s.Size() != 22 {
		t.Fatalf("size after growth %d, want 22", s.Size())
	}
	s.RemoveAll()
	if s.Size() != 0 {
		t.Fatalf("RemoveAll left elements")
	}
}
