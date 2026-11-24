package som

// initialCapacity is the default bucket count (Java Dictionary.INITIAL_CAPACITY).
const initialCapacity int = 16

// CustomHash is implemented by dictionary keys (Java Dictionary.CustomHash).
type CustomHash interface {
	CustomHash() int
}

// hashKey is the key constraint: Java's K extends CustomHash, plus comparable
// for the identity check in match.
type hashKey interface {
	comparable
	CustomHash
}

type idEntry[K comparable, V any] struct {
	hash  int
	key   K
	value V
	next  *idEntry[K, V]
}

func (e *idEntry[K, V]) match(hash int, key K) bool {
	return e.hash == hash && e.key == key
}

// IdentityDictionary is the hash table of the SOM class library that compares
// keys by identity.
//
// Porting notes:
//   - Java's base class Dictionary is never instantiated by a benchmark, so its
//     code lives here directly, with IdEntry's match (ADR-021).
//   - getValues is omitted: no benchmark uses it and it would force V comparable.
//   - hash >>> 16 is written as >>: all custom hashes are non-negative.
type IdentityDictionary[K hashKey, V any] struct {
	buckets []*idEntry[K, V]
	size    int
}

// NewIdentityDictionary returns an empty dictionary with the default capacity.
func NewIdentityDictionary[K hashKey, V any]() *IdentityDictionary[K, V] {
	return NewIdentityDictionarySize[K, V](initialCapacity)
}

// NewIdentityDictionarySize returns an empty dictionary with size buckets.
// size must be a power of two.
func NewIdentityDictionarySize[K hashKey, V any](size int) *IdentityDictionary[K, V] {
	return &IdentityDictionary[K, V]{buckets: make([]*idEntry[K, V], size)}
}

func dictionaryHash[K hashKey](key K) int {
	var zero K
	if key == zero {
		return 0
	}
	hash := key.CustomHash()
	return hash ^ hash>>16
}

// Size returns the number of entries.
func (d *IdentityDictionary[K, V]) Size() int {
	return d.size
}

// IsEmpty reports whether the dictionary has no entries.
func (d *IdentityDictionary[K, V]) IsEmpty() bool {
	return d.size == 0
}

func (d *IdentityDictionary[K, V]) getBucketIdx(hash int) int {
	return (len(d.buckets) - 1) & hash
}

func (d *IdentityDictionary[K, V]) getBucket(hash int) *idEntry[K, V] {
	return d.buckets[d.getBucketIdx(hash)]
}

// At returns the value for key, or the zero value if key is absent.
func (d *IdentityDictionary[K, V]) At(key K) V {
	hash := dictionaryHash(key)
	e := d.getBucket(hash)

	for e != nil {
		if e.match(hash, key) {
			return e.value
		}
		e = e.next
	}
	var zero V
	return zero
}

// ContainsKey reports whether key has an entry.
func (d *IdentityDictionary[K, V]) ContainsKey(key K) bool {
	hash := dictionaryHash(key)
	e := d.getBucket(hash)

	for e != nil {
		if e.match(hash, key) {
			return true
		}
		e = e.next
	}
	return false
}

// AtPut stores value under key, replacing an existing value.
func (d *IdentityDictionary[K, V]) AtPut(key K, value V) {
	hash := dictionaryHash(key)
	i := d.getBucketIdx(hash)

	current := d.buckets[i]

	if current == nil {
		d.buckets[i] = d.newEntry(key, value, hash)
		d.size++
	} else {
		d.insertBucketEntry(key, value, hash, current)
	}

	if d.size > len(d.buckets) {
		d.resize()
	}
}

func (d *IdentityDictionary[K, V]) newEntry(key K, value V, hash int) *idEntry[K, V] {
	return &idEntry[K, V]{hash: hash, key: key, value: value, next: nil}
}

func (d *IdentityDictionary[K, V]) insertBucketEntry(key K, value V, hash int, head *idEntry[K, V]) {
	current := head

	for {
		if current.match(hash, key) {
			current.value = value
			return
		}
		if current.next == nil {
			d.size++
			current.next = d.newEntry(key, value, hash)
			return
		}
		current = current.next
	}
}

func (d *IdentityDictionary[K, V]) resize() {
	oldStorage := d.buckets

	newStorage := make([]*idEntry[K, V], len(oldStorage)*2)
	d.buckets = newStorage
	d.transferEntries(oldStorage)
}

func (d *IdentityDictionary[K, V]) transferEntries(oldStorage []*idEntry[K, V]) {
	for i := 0; i < len(oldStorage); i++ {
		current := oldStorage[i]
		if current != nil {
			oldStorage[i] = nil

			if current.next == nil {
				d.buckets[current.hash&(len(d.buckets)-1)] = current
			} else {
				d.splitBucket(oldStorage, i, current)
			}
		}
	}
}

func (d *IdentityDictionary[K, V]) splitBucket(oldStorage []*idEntry[K, V], i int, head *idEntry[K, V]) {
	var loHead, loTail, hiHead, hiTail *idEntry[K, V]
	current := head

	for current != nil {
		if (current.hash & len(oldStorage)) == 0 {
			if loTail == nil {
				loHead = current
			} else {
				loTail.next = current
			}
			loTail = current
		} else {
			if hiTail == nil {
				hiHead = current
			} else {
				hiTail.next = current
			}
			hiTail = current
		}
		current = current.next
	}

	if loTail != nil {
		loTail.next = nil
		d.buckets[i] = loHead
	}
	if hiTail != nil {
		hiTail.next = nil
		d.buckets[i+len(oldStorage)] = hiHead
	}
}

// RemoveAll removes all entries and keeps the bucket count.
func (d *IdentityDictionary[K, V]) RemoveAll() {
	d.buckets = make([]*idEntry[K, V], len(d.buckets))
	d.size = 0
}

// GetKeys returns the keys in bucket order.
func (d *IdentityDictionary[K, V]) GetKeys() *Vector[K] {
	keys := NewVectorSize[K](d.size)
	for i := 0; i < len(d.buckets); i++ {
		current := d.buckets[i]
		for current != nil {
			keys.Append(current.key)
			current = current.next
		}
	}
	return keys
}
