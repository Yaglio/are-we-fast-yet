package main

type towersDisk struct {
	size int
	next *towersDisk
}

func newTowersDisk(size int) *towersDisk {
	return &towersDisk{size: size}
}

// Towers solves the Towers of Hanoi with linked-list piles.
type Towers struct {
	piles     []*towersDisk
	movesDone int
}

func (t *Towers) pushDisk(disk *towersDisk, pile int) {
	top := t.piles[pile]
	if top != nil && disk.size >= top.size {
		panic("Cannot put a big disk on a smaller one")
	}

	disk.next = top
	t.piles[pile] = disk
}

func (t *Towers) popDiskFrom(pile int) *towersDisk {
	top := t.piles[pile]
	if top == nil {
		panic("Attempting to remove a disk from an empty pile")
	}

	t.piles[pile] = top.next
	top.next = nil
	return top
}

func (t *Towers) moveTopDisk(fromPile, toPile int) {
	t.pushDisk(t.popDiskFrom(fromPile), toPile)
	t.movesDone++
}

func (t *Towers) buildTowerAt(pile, disks int) {
	for i := disks; i >= 0; i-- {
		t.pushDisk(newTowersDisk(i), pile)
	}
}

func (t *Towers) moveDisks(disks, fromPile, toPile int) {
	if disks == 1 {
		t.moveTopDisk(fromPile, toPile)
	} else {
		otherPile := (3 - fromPile) - toPile
		t.moveDisks(disks-1, fromPile, otherPile)
		t.moveTopDisk(fromPile, toPile)
		t.moveDisks(disks-1, otherPile, toPile)
	}
}

// Benchmark runs one iteration of the benchmark.
func (t *Towers) Benchmark() any {
	t.piles = make([]*towersDisk, 3)
	t.buildTowerAt(0, 13)
	t.movesDone = 0
	t.moveDisks(13, 0, 1)
	return t.movesDone
}

// VerifyResult checks the benchmark result.
func (*Towers) VerifyResult(result any) bool {
	return 8191 == result.(int)
}
