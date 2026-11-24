package som

// Set is the unordered collection without duplicates of the SOM class library,
// backed by a Vector.
//
// Porting notes:
//   - Java's equals becomes == on the comparable element type: value equality
//     for int (Integer.equals), identity for pointers (Object.equals).
//   - collect is omitted: no benchmark uses it and Go has no generic methods.
type Set[E comparable] struct {
	items *Vector[E]
}

// NewSet returns an empty set with the default initial capacity.
func NewSet[E comparable]() *Set[E] {
	return NewSetSize[E](initialSize)
}

// NewSetSize returns an empty set with the given initial capacity.
func NewSetSize[E comparable](size int) *Set[E] {
	return &Set[E]{items: NewVectorSize[E](size)}
}

// Size returns the number of elements.
func (s *Set[E]) Size() int {
	return s.items.Size()
}

// ForEach calls fn for every element in insertion order.
func (s *Set[E]) ForEach(fn func(E)) {
	s.items.ForEach(fn)
}

// HasSome reports whether fn returns true for any element.
func (s *Set[E]) HasSome(fn func(E) bool) bool {
	return s.items.HasSome(fn)
}

// GetOne returns the first element for which fn returns true, or the zero value.
func (s *Set[E]) GetOne(fn func(E) bool) E {
	return s.items.GetOne(fn)
}

// Add inserts obj unless the set already contains it.
func (s *Set[E]) Add(obj E) {
	if !s.Contains(obj) {
		s.items.Append(obj)
	}
}

// Contains reports whether obj is in the set.
func (s *Set[E]) Contains(obj E) bool {
	return s.HasSome(func(e E) bool { return e == obj })
}

// RemoveAll removes all elements and keeps the capacity.
func (s *Set[E]) RemoveAll() {
	s.items.RemoveAll()
}
