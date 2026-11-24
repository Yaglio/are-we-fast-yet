package som

// IdentitySet is a Set that compares elements by identity.
//
// Porting notes:
//   - Java's IdentitySet extends Set and overrides only contains. Go has no
//     inheritance, so this is a separate type with its own delegating methods.
//   - Add calls IdentitySet.Contains, as Java's virtual dispatch does.
type IdentitySet[E comparable] struct {
	items *Vector[E]
}

// NewIdentitySet returns an empty identity set with the default initial capacity.
func NewIdentitySet[E comparable]() *IdentitySet[E] {
	return NewIdentitySetSize[E](initialSize)
}

// NewIdentitySetSize returns an empty identity set with the given initial capacity.
func NewIdentitySetSize[E comparable](size int) *IdentitySet[E] {
	return &IdentitySet[E]{items: NewVectorSize[E](size)}
}

// Size returns the number of elements.
func (s *IdentitySet[E]) Size() int {
	return s.items.Size()
}

// ForEach calls fn for every element in insertion order.
func (s *IdentitySet[E]) ForEach(fn func(E)) {
	s.items.ForEach(fn)
}

// HasSome reports whether fn returns true for any element.
func (s *IdentitySet[E]) HasSome(fn func(E) bool) bool {
	return s.items.HasSome(fn)
}

// GetOne returns the first element for which fn returns true, or the zero value.
func (s *IdentitySet[E]) GetOne(fn func(E) bool) E {
	return s.items.GetOne(fn)
}

// Add inserts obj unless the set already contains it.
func (s *IdentitySet[E]) Add(obj E) {
	if !s.Contains(obj) {
		s.items.Append(obj)
	}
}

// Contains reports whether obj itself is in the set.
func (s *IdentitySet[E]) Contains(obj E) bool {
	return s.HasSome(func(e E) bool { return e == obj })
}

// RemoveAll removes all elements and keeps the capacity.
func (s *IdentitySet[E]) RemoveAll() {
	s.items.RemoveAll()
}
