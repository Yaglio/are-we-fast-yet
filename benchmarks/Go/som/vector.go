package som

// Vector is the growable array of the SOM class library.
//
// Porting notes:
//   - does not use an explicit array bounds check, because Go already does
//     that. Don't see a point in doing it twice.
type Vector[E comparable] struct {
	storage  []E
	firstIdx int
	lastIdx  int
}

// VectorWith returns a new vector holding only elem.
func VectorWith[E comparable](elem E) *Vector[E] {
	v := NewVectorSize[E](1)
	v.Append(elem)
	return v
}

// NewVectorSize returns an empty vector with the given initial capacity.
func NewVectorSize[E comparable](size int) *Vector[E] {
	v := &Vector[E]{}
	if size != 0 {
		v.storage = make([]E, size)
	}
	return v
}

// NewVector returns an empty vector without storage.
func NewVector[E comparable]() *Vector[E] {
	return &Vector[E]{}
}

// At returns the element at idx, or the zero value if idx is beyond the storage.
func (v *Vector[E]) At(idx int) E {
	if v.storage == nil || idx >= len(v.storage) {
		var zero E
		return zero
	}
	return v.storage[idx]
}

// AtPut stores val at idx, growing the storage if needed.
func (v *Vector[E]) AtPut(idx int, val E) {
	if v.storage == nil {
		v.storage = make([]E, max(idx+1, initialSize))
	} else if idx >= len(v.storage) {
		newLength := len(v.storage)
		for newLength <= idx {
			newLength *= 2
		}
		newStorage := make([]E, newLength)
		copy(newStorage, v.storage)
		v.storage = newStorage
	}
	v.storage[idx] = val
	if v.lastIdx < idx+1 {
		v.lastIdx = idx + 1
	}
}

// Append adds elem at the end, growing the storage if needed.
func (v *Vector[E]) Append(elem E) {
	if v.storage == nil {
		v.storage = make([]E, initialSize)
	} else if v.lastIdx >= len(v.storage) {
		// Need to expand capacity first
		newStorage := make([]E, 2*len(v.storage))
		copy(newStorage, v.storage)
		v.storage = newStorage
	}

	v.storage[v.lastIdx] = elem
	v.lastIdx++
}

// IsEmpty reports whether the vector has no elements.
func (v *Vector[E]) IsEmpty() bool {
	return v.lastIdx == v.firstIdx
}

// ForEach calls fn for every element in order.
func (v *Vector[E]) ForEach(fn func(E)) {
	for i := v.firstIdx; i < v.lastIdx; i++ {
		fn(v.storage[i])
	}
}

// HasSome reports whether fn returns true for any element.
func (v *Vector[E]) HasSome(fn func(E) bool) bool {
	for i := v.firstIdx; i < v.lastIdx; i++ {
		if fn(v.storage[i]) {
			return true
		}
	}
	return false
}

// GetOne returns the first element for which fn returns true, or the zero value.
func (v *Vector[E]) GetOne(fn func(E) bool) E {
	for i := v.firstIdx; i < v.lastIdx; i++ {
		e := v.storage[i]
		if fn(e) {
			return e
		}
	}
	var zero E
	return zero
}

// First returns the first element, or the zero value if the vector is empty.
func (v *Vector[E]) First() E {
	if v.IsEmpty() {
		var zero E
		return zero
	}
	return v.storage[v.firstIdx]
}

// RemoveFirst removes and returns the first element, or the zero value if empty.
func (v *Vector[E]) RemoveFirst() E {
	if v.IsEmpty() {
		var zero E
		return zero
	}
	v.firstIdx++
	return v.storage[v.firstIdx-1]
}

// Remove deletes every occurrence of obj and reports whether one was found.
func (v *Vector[E]) Remove(obj E) bool {
	if v.storage == nil || v.IsEmpty() {
		return false
	}

	newArray := make([]E, v.Capacity())
	newLast := 0
	found := false

	v.ForEach(func(it E) {
		if it == obj {
			found = true
		} else {
			newArray[newLast] = it
			newLast++
		}
	})

	v.storage = newArray
	v.lastIdx = newLast
	v.firstIdx = 0
	return found
}

// RemoveAll removes all elements and keeps the capacity.
func (v *Vector[E]) RemoveAll() {
	v.firstIdx = 0
	v.lastIdx = 0
	if v.storage != nil {
		v.storage = make([]E, len(v.storage))
	}
}

// Size returns the number of elements.
func (v *Vector[E]) Size() int {
	return v.lastIdx - v.firstIdx
}

// Capacity returns the length of the storage.
func (v *Vector[E]) Capacity() int {
	if v.storage == nil {
		return 0
	}
	return len(v.storage)
}

// Sort sorts the elements with the comparator c.
func (v *Vector[E]) Sort(c func(a, b E) int) {
	if v.Size() > 0 {
		v.sortRange(v.firstIdx, v.lastIdx-1, c)
	}
}

func (v *Vector[E]) sortRange(i, j int, c func(a, b E) int) {
	if c == nil {
		v.defaultSort(i, j)
	}

	n := j + 1 - i
	if n <= 1 {
		return
	}

	di := v.storage[i]
	dj := v.storage[j]

	if c(di, dj) > 0 {
		vectorSwap(v.storage, i, j)
		tt := di
		di = dj
		dj = tt
	}

	if n > 2 {
		ij := (i + j) / 2
		dij := v.storage[ij]

		if c(di, dij) <= 0 {
			if c(dij, dj) > 0 {
				vectorSwap(v.storage, j, ij)
				dij = dj
			}
		} else {
			vectorSwap(v.storage, i, ij)
			dij = di
		}

		if n > 3 {
			k := i
			l := j - 1

			for {
				for k <= l && c(dij, v.storage[l]) <= 0 {
					l--
				}

				k++
				for k <= l && c(v.storage[k], dij) <= 0 {
					k++
				}

				if k > l {
					break
				}
				vectorSwap(v.storage, k, l)
			}

			v.sortRange(i, l, c)
			v.sortRange(k, j, c)
		}
	}
}

func vectorSwap[E comparable](_ []E, _, _ int) {
	panic("Not Implemented")
}

func (v *Vector[E]) defaultSort(_, _ int) {
	panic("Not Implemented")
}
