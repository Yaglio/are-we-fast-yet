package main

import (
	"fmt"
	"math"

	"awfy/som"
)

// CD is the collision detection benchmark. The Java package cd (Simulator,
// CollisionDetector, RedBlackTree, ...) lives in this file.
type CD struct{}

// benchmark is Java's private overload benchmark(int) (ADR-010).
func (*CD) benchmark(numAircrafts int) int {
	numFrames := 200

	simulator := newSimulator(numAircrafts)
	detector := newCollisionDetector()

	actualCollisions := 0

	for i := 0; i < numFrames; i++ {
		time := float64(i) / 10.0
		collisions := detector.handleNewFrame(simulator.simulate(time))
		actualCollisions += collisions.Size()
	}

	return actualCollisions
}

// InnerBenchmarkLoop overrides the default inner loop (ADR-003).
func (c *CD) InnerBenchmarkLoop(innerIterations int) bool {
	return c.verify(c.benchmark(innerIterations), innerIterations)
}

// verify is Java's overload verifyResult(int, int) (ADR-010).
func (*CD) verify(actualCollisions, numAircrafts int) bool {
	if numAircrafts == 1000 {
		return actualCollisions == 14484
	}
	if numAircrafts == 500 {
		return actualCollisions == 14484
	}
	if numAircrafts == 250 {
		return actualCollisions == 10830
	}
	if numAircrafts == 200 {
		return actualCollisions == 8655
	}
	if numAircrafts == 100 {
		return actualCollisions == 4305
	}
	if numAircrafts == 10 {
		return actualCollisions == 390
	}
	if numAircrafts == 2 {
		return actualCollisions == 42
	}

	fmt.Printf("No verification result for %d found\n", numAircrafts)
	fmt.Printf("Result is: %d\n", actualCollisions)
	return false
}

// Benchmark is never called, InnerBenchmarkLoop is used instead.
func (*CD) Benchmark() any {
	panic("Should never be reached")
}

// VerifyResult is never called, InnerBenchmarkLoop is used instead.
func (*CD) VerifyResult(any) bool {
	panic("Should never be reached")
}

// Java cd.Constants (ADR-012).
const (
	minX            float64 = 0.0
	minY            float64 = 0.0
	maxX            float64 = 1000.0
	maxY            float64 = 1000.0
	minZ            float64 = 0.0
	maxZ            float64 = 10.0
	proximityRadius float64 = 1.0
	goodVoxelSize   float64 = proximityRadius * 2.0
)

type vector2D struct {
	x float64
	y float64
}

func newVector2D(x, y float64) *vector2D {
	return &vector2D{x: x, y: y}
}

func (v *vector2D) plus(other *vector2D) *vector2D {
	return newVector2D(v.x+other.x,
		v.y+other.y)
}

func (v *vector2D) minus(other *vector2D) *vector2D {
	return newVector2D(v.x-other.x,
		v.y-other.y)
}

func (v *vector2D) compareTo(other *vector2D) int {
	result := vector2DCompareNumbers(v.x, other.x)
	if result != 0 {
		return result
	}
	return vector2DCompareNumbers(v.y, other.y)
}

func vector2DCompareNumbers(a, b float64) int {
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}

	// We say that NaN is smaller than non-NaN.
	if a == a {
		return 1
	}
	return -1
}

type vector3D struct {
	x float64
	y float64
	z float64
}

func newVector3D(x, y, z float64) *vector3D {
	return &vector3D{x: x, y: y, z: z}
}

func (v *vector3D) plus(other *vector3D) *vector3D {
	return newVector3D(v.x+other.x,
		v.y+other.y,
		v.z+other.z)
}

func (v *vector3D) minus(other *vector3D) *vector3D {
	return newVector3D(v.x-other.x,
		v.y-other.y,
		v.z-other.z)
}

func (v *vector3D) dot(other *vector3D) float64 {
	return v.x*other.x + v.y*other.y + v.z*other.z
}

func (v *vector3D) squaredMagnitude() float64 {
	return v.dot(v)
}

func (v *vector3D) magnitude() float64 {
	return math.Sqrt(v.squaredMagnitude())
}

func (v *vector3D) times(amount float64) *vector3D {
	return newVector3D(v.x*amount,
		v.y*amount,
		v.z*amount)
}

type callSign struct {
	value int
}

func newCallSign(value int) *callSign {
	return &callSign{value: value}
}

func (c *callSign) compareTo(other *callSign) int {
	if c.value == other.value {
		return 0
	}
	if c.value < other.value {
		return -1
	}
	return 1
}

type aircraft struct {
	callsign *callSign
	position *vector3D
}

func newAircraft(callsign *callSign, position *vector3D) *aircraft {
	return &aircraft{callsign: callsign, position: position}
}

// collision's fields are written but never read (ADR-030).
type collision struct {
	aircraftA *callSign
	aircraftB *callSign
	position  *vector3D
}

func newCollision(aircraftA, aircraftB *callSign, position *vector3D) *collision {
	return &collision{aircraftA: aircraftA, aircraftB: aircraftB, position: position}
}

type motion struct {
	callsign *callSign
	posOne   *vector3D
	posTwo   *vector3D
}

func newMotion(callsign *callSign, posOne, posTwo *vector3D) *motion {
	return &motion{callsign: callsign, posOne: posOne, posTwo: posTwo}
}

func (m *motion) delta() *vector3D {
	return m.posTwo.minus(m.posOne)
}

func (m *motion) findIntersection(other *motion) *vector3D {
	init1 := m.posOne
	init2 := other.posOne
	vec1 := m.delta()
	vec2 := other.delta()
	radius := proximityRadius

	// this test is not geometrical 3-d intersection test, it takes the fact that the aircraft move
	// into account ; so it is more like a 4d test
	// (it assumes that both of the aircraft have a constant speed over the tested interval)

	// we thus have two points, each of them moving on its line segment at constant speed ; we are looking
	// for times when the distance between these two points is smaller than r

	// vec1 is vector of aircraft 1
	// vec2 is vector of aircraft 2

	// a = (V2 - V1)^T * (V2 - V1)
	a := vec2.minus(vec1).squaredMagnitude()

	if a != 0.0 {
		// we are first looking for instances of time when the planes are exactly r from each other
		// at least one plane is moving ; if the planes are moving in parallel, they do not have constant speed

		// if the planes are moving in parallel, then
		//   if the faster starts behind the slower, we can have 2, 1, or 0 solutions
		//   if the faster plane starts in front of the slower, we can have 0 or 1 solutions

		// if the planes are not moving in parallel, then

		// point P1 = I1 + vV1
		// point P2 = I2 + vV2
		//   - looking for v, such that dist(P1,P2) = || P1 - P2 || = r

		// it follows that || P1 - P2 || = sqrt( < P1-P2, P1-P2 > )
		//   0 = -r^2 + < P1 - P2, P1 - P2 >
		//  from properties of dot product
		//   0 = -r^2 + <I1-I2,I1-I2> + v * 2<I1-I2, V1-V2> + v^2 *<V1-V2,V1-V2>
		//   so we calculate a, b, c - and solve the quadratic equation
		//   0 = c + bv + av^2

		// b = 2 * <I1-I2, V1-V2>
		b := 2.0 * init1.minus(init2).dot(vec1.minus(vec2))

		// c = -r^2 + (I2 - I1)^T * (I2 - I1)
		c := -radius*radius + init2.minus(init1).squaredMagnitude()

		discr := b*b - 4.0*a*c
		if discr < 0.0 {
			return nil
		}

		v1 := (-b - math.Sqrt(discr)) / (2.0 * a)
		v2 := (-b + math.Sqrt(discr)) / (2.0 * a)

		if v1 <= v2 && ((v1 <= 1.0 && 1.0 <= v2) ||
			(v1 <= 0.0 && 0.0 <= v2) ||
			(0.0 <= v1 && v2 <= 1.0)) {
			// Pick a good "time" at which to report the collision.
			var v float64
			if v1 <= 0.0 {
				// The collision started before this frame. Report it at the start of the frame.
				v = 0.0
			} else {
				// The collision started during this frame. Report it at that moment.
				v = v1
			}

			result1 := init1.plus(vec1.times(v))
			result2 := init2.plus(vec2.times(v))

			result := result1.plus(result2).times(0.5)
			if result.x >= minX &&
				result.x <= maxX &&
				result.y >= minY &&
				result.y <= maxY &&
				result.z >= minZ &&
				result.z <= maxZ {
				return result
			}
		}

		return nil
	}

	// the planes have the same speeds and are moving in parallel (or they are not moving at all)
	// they  thus have the same distance all the time ; we calculate it from the initial point

	// dist = || i2 - i1 || = sqrt(  ( i2 - i1 )^T * ( i2 - i1 ) )
	dist := init2.minus(init1).magnitude()
	if dist <= radius {
		return init1.plus(init2).times(0.5)
	}

	return nil
}

// compareKey is Java's bound K extends Comparable<K> (ADR-031).
type compareKey[K any] interface {
	compareTo(other K) int
}

// color is RedBlackTree's private enum Color (ADR-029).
type color int

const (
	colorRed color = iota
	colorBlack
)

type redBlackTreeNode[K, V any] struct {
	key    K
	value  V
	left   *redBlackTreeNode[K, V]
	right  *redBlackTreeNode[K, V]
	parent *redBlackTreeNode[K, V]
	color  color
}

func newRedBlackTreeNode[K, V any](key K, value V) *redBlackTreeNode[K, V] {
	return &redBlackTreeNode[K, V]{
		key:    key,
		value:  value,
		left:   nil,
		right:  nil,
		parent: nil,
		color:  colorRed,
	}
}

func (n *redBlackTreeNode[K, V]) successor() *redBlackTreeNode[K, V] {
	x := n
	if x.right != nil {
		return redBlackTreeTreeMinimum(x.right)
	}
	y := x.parent
	for y != nil && x == y.right {
		x = y
		y = y.parent
	}
	return y
}

func redBlackTreeTreeMinimum[K, V any](x *redBlackTreeNode[K, V]) *redBlackTreeNode[K, V] {
	current := x
	for current.left != nil {
		current = current.left
	}
	return current
}

type redBlackTreeEntry[K, V any] struct {
	key   K
	value V
}

func newRedBlackTreeEntry[K, V any](key K, value V) *redBlackTreeEntry[K, V] {
	return &redBlackTreeEntry[K, V]{key: key, value: value}
}

type redBlackTreeInsertResult[K, V any] struct {
	isNewEntry bool
	newNode    *redBlackTreeNode[K, V]
	oldValue   V
}

func newRedBlackTreeInsertResult[K, V any](isNewEntry bool, newNode *redBlackTreeNode[K, V],
	oldValue V) *redBlackTreeInsertResult[K, V] {
	return &redBlackTreeInsertResult[K, V]{isNewEntry: isNewEntry, newNode: newNode, oldValue: oldValue}
}

type redBlackTree[K compareKey[K], V any] struct {
	root *redBlackTreeNode[K, V]
}

func newRedBlackTree[K compareKey[K], V any]() *redBlackTree[K, V] {
	return &redBlackTree[K, V]{root: nil}
}

func (t *redBlackTree[K, V]) put(key K, value V) V {
	insertionResult := t.treeInsert(key, value)
	if !insertionResult.isNewEntry {
		return insertionResult.oldValue
	}
	x := insertionResult.newNode

	for x != t.root && x.parent.color == colorRed {
		if x.parent == x.parent.parent.left {
			y := x.parent.parent.right
			if y != nil && y.color == colorRed {
				// Case 1
				x.parent.color = colorBlack
				y.color = colorBlack
				x.parent.parent.color = colorRed
				x = x.parent.parent
			} else {
				if x == x.parent.right {
					// Case 2
					x = x.parent
					t.leftRotate(x)
				}
				// Case 3
				x.parent.color = colorBlack
				x.parent.parent.color = colorRed
				t.rightRotate(x.parent.parent)
			}
		} else {
			// Same as "then" clause with "right" and "left" exchanged.
			y := x.parent.parent.left
			if y != nil && y.color == colorRed {
				// Case 1
				x.parent.color = colorBlack
				y.color = colorBlack
				x.parent.parent.color = colorRed
				x = x.parent.parent
			} else {
				if x == x.parent.left {
					// Case 2
					x = x.parent
					t.rightRotate(x)
				}
				// Case 3
				x.parent.color = colorBlack
				x.parent.parent.color = colorRed
				t.leftRotate(x.parent.parent)
			}
		}
	}

	t.root.color = colorBlack
	var zero V
	return zero
}

func (t *redBlackTree[K, V]) remove(key K) V {
	z := t.findNode(key)
	if z == nil {
		var zero V
		return zero
	}

	// Y is the node to be unlinked from the tree.
	var y *redBlackTreeNode[K, V]
	if z.left == nil || z.right == nil {
		y = z
	} else {
		y = z.successor()
	}

	// Y is guaranteed to be non-null at this point.
	var x *redBlackTreeNode[K, V]
	if y.left != nil {
		x = y.left
	} else {
		x = y.right
	}

	// X is the child of y which might potentially replace y in the tree. X might be null at
	// this point.
	var xParent *redBlackTreeNode[K, V]
	if x != nil {
		x.parent = y.parent
		xParent = x.parent
	} else {
		xParent = y.parent
	}
	if y.parent == nil {
		t.root = x
	} else {
		if y == y.parent.left {
			y.parent.left = x
		} else {
			y.parent.right = x
		}
	}

	if y != z {
		if y.color == colorBlack {
			t.removeFixup(x, xParent)
		}

		y.parent = z.parent
		y.color = z.color
		y.left = z.left
		y.right = z.right

		if z.left != nil {
			z.left.parent = y
		}
		if z.right != nil {
			z.right.parent = y
		}
		if z.parent != nil {
			if z.parent.left == z {
				z.parent.left = y
			} else {
				z.parent.right = y
			}
		} else {
			t.root = y
		}
	} else if y.color == colorBlack {
		t.removeFixup(x, xParent)
	}

	return z.value
}

func (t *redBlackTree[K, V]) get(key K) V {
	node := t.findNode(key)
	if node == nil {
		var zero V
		return zero
	}
	return node.value
}

func (t *redBlackTree[K, V]) forEach(fn func(*redBlackTreeEntry[K, V])) {
	if t.root == nil {
		return
	}
	current := redBlackTreeTreeMinimum(t.root)
	for current != nil {
		fn(newRedBlackTreeEntry(current.key, current.value))
		current = current.successor()
	}
}

func (t *redBlackTree[K, V]) findNode(key K) *redBlackTreeNode[K, V] {
	current := t.root
	for current != nil {
		comparisonResult := key.compareTo(current.key)
		if comparisonResult == 0 {
			return current
		}
		if comparisonResult < 0 {
			current = current.left
		} else {
			current = current.right
		}
	}
	return nil
}

func (t *redBlackTree[K, V]) treeInsert(key K, value V) *redBlackTreeInsertResult[K, V] {
	var y *redBlackTreeNode[K, V]
	x := t.root

	for x != nil {
		y = x
		comparisonResult := key.compareTo(x.key)
		if comparisonResult < 0 {
			x = x.left
		} else if comparisonResult > 0 {
			x = x.right
		} else {
			oldValue := x.value
			x.value = value
			return newRedBlackTreeInsertResult[K, V](false, nil, oldValue)
		}
	}

	z := newRedBlackTreeNode(key, value)
	z.parent = y
	if y == nil {
		t.root = z
	} else {
		if key.compareTo(y.key) < 0 {
			y.left = z
		} else {
			y.right = z
		}
	}
	var zero V
	return newRedBlackTreeInsertResult(true, z, zero)
}

func (t *redBlackTree[K, V]) leftRotate(x *redBlackTreeNode[K, V]) *redBlackTreeNode[K, V] {
	y := x.right

	// Turn y's left subtree into x's right subtree.
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}

	// Link x's parent to y.
	y.parent = x.parent
	if x.parent == nil {
		t.root = y
	} else {
		if x == x.parent.left {
			x.parent.left = y
		} else {
			x.parent.right = y
		}
	}

	// Put x on y's left.
	y.left = x
	x.parent = y

	return y
}

func (t *redBlackTree[K, V]) rightRotate(y *redBlackTreeNode[K, V]) *redBlackTreeNode[K, V] {
	x := y.left

	// Turn x's right subtree into y's left subtree.
	y.left = x.right
	if x.right != nil {
		x.right.parent = y
	}

	// Link y's parent to x;
	x.parent = y.parent
	if y.parent == nil {
		t.root = x
	} else {
		if y == y.parent.left {
			y.parent.left = x
		} else {
			y.parent.right = x
		}
	}

	x.right = y
	y.parent = x

	return x
}

func (t *redBlackTree[K, V]) removeFixup(x, xParent *redBlackTreeNode[K, V]) {
	for x != t.root && (x == nil || x.color == colorBlack) {
		if x == xParent.left {
			// Note: the text points out that w cannot be null. The reason is not obvious from
			// simply looking at the code; it comes about from the properties of the red-black
			// tree.
			w := xParent.right
			if w.color == colorRed {
				// Case 1
				w.color = colorBlack
				xParent.color = colorRed
				t.leftRotate(xParent)
				w = xParent.right
			}
			if (w.left == nil || w.left.color == colorBlack) &&
				(w.right == nil || w.right.color == colorBlack) {
				// Case 2
				w.color = colorRed
				x = xParent
				xParent = x.parent
			} else {
				if w.right == nil || w.right.color == colorBlack {
					// Case 3
					w.left.color = colorBlack
					w.color = colorRed
					t.rightRotate(w)
					w = xParent.right
				}
				// Case 4
				w.color = xParent.color
				xParent.color = colorBlack
				if w.right != nil {
					w.right.color = colorBlack
				}
				t.leftRotate(xParent)
				x = t.root
				xParent = x.parent
			}
		} else {
			// Same as "then" clause with "right" and "left" exchanged.
			w := xParent.left
			if w.color == colorRed {
				// Case 1
				w.color = colorBlack
				xParent.color = colorRed
				t.rightRotate(xParent)
				w = xParent.left
			}
			if (w.right == nil || w.right.color == colorBlack) &&
				(w.left == nil || w.left.color == colorBlack) {
				// Case 2
				w.color = colorRed
				x = xParent
				xParent = x.parent
			} else {
				if w.left == nil || w.left.color == colorBlack {
					// Case 3
					w.right.color = colorBlack
					w.color = colorRed
					t.leftRotate(w)
					w = xParent.left
				}
				// Case 4
				w.color = xParent.color
				xParent.color = colorBlack
				if w.left != nil {
					w.left.color = colorBlack
				}
				t.rightRotate(xParent)
				x = t.root
				xParent = x.parent
			}
		}
	}
	if x != nil {
		x.color = colorBlack
	}
}

type simulator struct {
	aircraft *som.Vector[*callSign]
}

func newSimulator(numAircraft int) *simulator {
	s := &simulator{aircraft: som.NewVector[*callSign]()}
	for i := 0; i < numAircraft; i++ {
		s.aircraft.Append(newCallSign(i))
	}
	return s
}

func (s *simulator) simulate(time float64) *som.Vector[*aircraft] {
	frame := som.NewVector[*aircraft]()
	for i := 0; i < s.aircraft.Size(); i += 2 {
		frame.Append(newAircraft(s.aircraft.At(i),
			newVector3D(time, math.Cos(time)*2+float64(i*3), 10)))
		frame.Append(newAircraft(s.aircraft.At(i+1),
			newVector3D(time, math.Sin(time)*2+float64(i*3), 10)))
	}
	return frame
}

type collisionDetector struct {
	state *redBlackTree[*callSign, *vector3D]
}

func newCollisionDetector() *collisionDetector {
	return &collisionDetector{state: newRedBlackTree[*callSign, *vector3D]()}
}

func (d *collisionDetector) handleNewFrame(frame *som.Vector[*aircraft]) *som.Vector[*collision] {
	motions := som.NewVector[*motion]()
	seen := newRedBlackTree[*callSign, bool]()

	frame.ForEach(func(aircraft *aircraft) {
		oldPosition := d.state.put(aircraft.callsign, aircraft.position)
		newPosition := aircraft.position
		seen.put(aircraft.callsign, true)

		if oldPosition == nil {
			// Treat newly introduced aircraft as if they were stationary.
			oldPosition = newPosition
		}

		motions.Append(newMotion(aircraft.callsign, oldPosition, newPosition))
	})

	// Remove aircraft that are no longer present.
	toRemove := som.NewVector[*callSign]()
	d.state.forEach(func(e *redBlackTreeEntry[*callSign, *vector3D]) {
		if !seen.get(e.key) {
			toRemove.Append(e.key)
		}
	})

	toRemove.ForEach(func(e *callSign) { d.state.remove(e) })

	allReduced := collisionDetectorReduceCollisionSet(motions)
	collisions := som.NewVector[*collision]()
	allReduced.ForEach(func(reduced *som.Vector[*motion]) {
		for i := 0; i < reduced.Size(); i++ {
			motion1 := reduced.At(i)
			for j := i + 1; j < reduced.Size(); j++ {
				motion2 := reduced.At(j)
				collision := motion1.findIntersection(motion2)
				if collision != nil {
					collisions.Append(newCollision(motion1.callsign, motion2.callsign, collision))
				}
			}
		}
	})

	return collisions
}

func collisionDetectorIsInVoxel(voxel *vector2D, motion *motion) bool {
	if voxel.x > maxX ||
		voxel.x < minX ||
		voxel.y > maxY ||
		voxel.y < minY {
		return false
	}

	init := motion.posOne
	fin := motion.posTwo

	vS := goodVoxelSize
	r := proximityRadius / 2.0

	vX := voxel.x
	x0 := init.x
	xv := fin.x - init.x

	vY := voxel.y
	y0 := init.y
	yv := fin.y - init.y

	lowX := (vX - r - x0) / xv
	highX := (vX + vS + r - x0) / xv

	if xv < 0.0 {
		tmp := lowX
		lowX = highX
		highX = tmp
	}

	lowY := (vY - r - y0) / yv
	highY := (vY + vS + r - y0) / yv

	if yv < 0.0 {
		tmp := lowY
		lowY = highY
		highY = tmp
	}

	return ((xv == 0.0 && vX <= x0+r && x0-r <= vX+vS) /* no motion in x */ ||
		(lowX <= 1.0 && 1.0 <= highX) || (lowX <= 0.0 && 0.0 <= highX) ||
		(0.0 <= lowX && highX <= 1.0)) &&
		((yv == 0.0 && vY <= y0+r && y0-r <= vY+vS) /* no motion in y */ ||
			((lowY <= 1.0 && 1.0 <= highY) || (lowY <= 0.0 && 0.0 <= highY) ||
				(0.0 <= lowY && highY <= 1.0))) &&
		(xv == 0.0 || yv == 0.0 || /* no motion in x or y or both */
			(lowY <= highX && highX <= highY) ||
			(lowY <= lowX && lowX <= highY) ||
			(lowX <= lowY && highY <= highX))
}

// Java's private static final Vector2D horizontal/vertical (ADR-022).
var (
	collisionDetectorHorizontal = newVector2D(goodVoxelSize, 0.0)
	collisionDetectorVertical   = newVector2D(0.0, goodVoxelSize)
)

func collisionDetectorPutIntoMap(voxelMap *redBlackTree[*vector2D, *som.Vector[*motion]],
	voxel *vector2D, m *motion) {
	array := voxelMap.get(voxel)
	if array == nil {
		array = som.NewVector[*motion]()
		voxelMap.put(voxel, array)
	}
	array.Append(m)
}

func collisionDetectorRecurse(voxelMap *redBlackTree[*vector2D, *som.Vector[*motion]],
	seen *redBlackTree[*vector2D, bool], nextVoxel *vector2D, motion *motion) {
	if !collisionDetectorIsInVoxel(nextVoxel, motion) {
		return
	}

	if seen.put(nextVoxel, true) {
		return
	}

	collisionDetectorPutIntoMap(voxelMap, nextVoxel, motion)

	collisionDetectorRecurse(voxelMap, seen, nextVoxel.minus(collisionDetectorHorizontal), motion)
	collisionDetectorRecurse(voxelMap, seen, nextVoxel.plus(collisionDetectorHorizontal), motion)
	collisionDetectorRecurse(voxelMap, seen, nextVoxel.minus(collisionDetectorVertical), motion)
	collisionDetectorRecurse(voxelMap, seen, nextVoxel.plus(collisionDetectorVertical), motion)
	collisionDetectorRecurse(voxelMap, seen,
		nextVoxel.minus(collisionDetectorHorizontal).minus(collisionDetectorVertical), motion)
	collisionDetectorRecurse(voxelMap, seen,
		nextVoxel.minus(collisionDetectorHorizontal).plus(collisionDetectorVertical), motion)
	collisionDetectorRecurse(voxelMap, seen,
		nextVoxel.plus(collisionDetectorHorizontal).minus(collisionDetectorVertical), motion)
	collisionDetectorRecurse(voxelMap, seen,
		nextVoxel.plus(collisionDetectorHorizontal).plus(collisionDetectorVertical), motion)
}

func collisionDetectorReduceCollisionSet(motions *som.Vector[*motion]) *som.Vector[*som.Vector[*motion]] {
	voxelMap := newRedBlackTree[*vector2D, *som.Vector[*motion]]()
	motions.ForEach(func(motion *motion) { collisionDetectorDrawMotionOnVoxelMap(voxelMap, motion) })

	result := som.NewVector[*som.Vector[*motion]]()
	voxelMap.forEach(func(e *redBlackTreeEntry[*vector2D, *som.Vector[*motion]]) {
		if e.value.Size() > 1 {
			result.Append(e.value)
		}
	})
	return result
}

func collisionDetectorVoxelHash(position *vector3D) *vector2D {
	xDiv := int(position.x / goodVoxelSize)
	yDiv := int(position.y / goodVoxelSize)

	x := goodVoxelSize * float64(xDiv)
	y := goodVoxelSize * float64(yDiv)

	if position.x < 0 {
		x -= goodVoxelSize
	}
	if position.y < 0 {
		y -= goodVoxelSize
	}

	return newVector2D(x, y)
}

func collisionDetectorDrawMotionOnVoxelMap(voxelMap *redBlackTree[*vector2D, *som.Vector[*motion]],
	motion *motion) {
	seen := newRedBlackTree[*vector2D, bool]()
	collisionDetectorRecurse(voxelMap, seen, collisionDetectorVoxelHash(motion.posOne), motion)
}
