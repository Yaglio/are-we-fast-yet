package som

import "testing"

type testKey struct{ hash int }

func (k *testKey) CustomHash() int { return k.hash }

func TestIdentityDictionaryAtPutAt(t *testing.T) {
	d := NewIdentityDictionary[*testKey, int]()
	if !d.IsEmpty() {
		t.Fatalf("new dictionary not empty")
	}
	a, b := &testKey{3}, &testKey{3} // same hash, distinct identity
	d.AtPut(a, 1)
	d.AtPut(b, 2)
	if d.Size() != 2 || d.At(a) != 1 || d.At(b) != 2 {
		t.Fatalf("identity keys mixed up: size %d, a %d, b %d", d.Size(), d.At(a), d.At(b))
	}
	d.AtPut(a, 5)
	if d.Size() != 2 || d.At(a) != 5 {
		t.Fatalf("overwrite: size %d, a %d", d.Size(), d.At(a))
	}
	if d.At(&testKey{3}) != 0 || d.ContainsKey(&testKey{3}) || !d.ContainsKey(b) {
		t.Fatalf("missing key handling wrong")
	}
}

func TestIdentityDictionaryNilKey(t *testing.T) {
	d := NewIdentityDictionary[*testKey, int]()
	d.AtPut(nil, 7)
	if d.At(nil) != 7 || d.buckets[0] == nil {
		t.Fatalf("nil key not stored in bucket 0")
	}
}

func TestIdentityDictionaryResize(t *testing.T) {
	d := NewIdentityDictionary[*testKey, int]()
	keys := make([]*testKey, 100)
	for i := range keys {
		// hashes 0, 16, 32, ... collide in a 16-bucket table and split on resize
		keys[i] = &testKey{(i%10)*16 + i/10}
		d.AtPut(keys[i], i)
	}
	if d.Size() != 100 || len(d.buckets) != 128 {
		t.Fatalf("size %d, buckets %d, want 100, 128", d.Size(), len(d.buckets))
	}
	for i, k := range keys {
		if d.At(k) != i {
			t.Fatalf("key %d lost after resize", i)
		}
	}

	d.RemoveAll()
	if d.Size() != 0 || len(d.buckets) != 128 || d.ContainsKey(keys[0]) {
		t.Fatalf("RemoveAll wrong")
	}
}

func TestIdentityDictionaryGetKeys(t *testing.T) {
	d := NewIdentityDictionary[*testKey, int]()
	k2, k0, k0b := &testKey{2}, &testKey{0}, &testKey{16}
	d.AtPut(k2, 0)
	d.AtPut(k0, 0)
	d.AtPut(k0b, 0)
	keys := d.GetKeys()
	if keys.Size() != 3 || keys.At(0) != k0 || keys.At(1) != k0b || keys.At(2) != k2 {
		t.Fatalf("GetKeys not in bucket/chain order")
	}
}
