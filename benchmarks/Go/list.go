package main

type element struct {
	val  any
	next *element
}

func newElement(v any) *element {
	return &element{val: v}
}

func (e *element) length() int {
	if e.next == nil {
		return 1
	} else {
		return 1 + e.next.length()
	}
}

// List exercises recursion over singly linked lists.
type List struct{}

// Benchmark runs one iteration of the benchmark.
func (l *List) Benchmark() any {
	result := l.tail(l.makeList(15), l.makeList(10), l.makeList(6))
	return result.length()
}

func (l *List) makeList(length int) *element {
	if length == 0 {
		return nil
	} else {
		e := newElement(length)
		e.next = l.makeList(length - 1)
		return e
	}
}

func (l *List) isShorterThan(x, y *element) bool {
	xTail := x
	yTail := y

	for yTail != nil {
		if xTail == nil {
			return true
		}
		xTail = xTail.next
		yTail = yTail.next
	}
	return false
}

func (l *List) tail(x, y, z *element) *element {
	if l.isShorterThan(y, x) {
		return l.tail(l.tail(x.next, y, z),
			l.tail(y.next, z, x),
			l.tail(z.next, x, y))
	} else {
		return z
	}
}

// VerifyResult checks the benchmark result.
func (*List) VerifyResult(result any) bool {
	return 10 == result.(int)
}
