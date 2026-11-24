package som

import "testing"

func TestRandom(t *testing.T) {
	expected := []int{22896, 34761, 34014, 39231, 52540, 41445, 1546, 5947, 65224}
	rnd := NewRandom()
	for i, want := range expected {
		if got := rnd.Next(); got != want {
			t.Fatalf("value %d: got %d, want %d", i, got, want)
		}
	}
}
