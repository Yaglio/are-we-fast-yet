package main

import "awfy/som"

// DeltaBlue is the incremental constraint solver benchmark. The Java package
// deltablue (Planner, constraints, Strength, Variable, Plan) lives in this file.
type DeltaBlue struct{}

// InnerBenchmarkLoop overrides the default inner loop (ADR-003).
func (*DeltaBlue) InnerBenchmarkLoop(innerIterations int) bool {
	plannerChainTest(innerIterations)
	plannerProjectionTest(innerIterations)
	return true
}

// Benchmark is never called, InnerBenchmarkLoop is used instead.
func (*DeltaBlue) Benchmark() any {
	panic("should never be reached")
}

// VerifyResult is never called, InnerBenchmarkLoop is used instead.
func (*DeltaBlue) VerifyResult(any) bool {
	panic("should never be reached")
}

// direction is Java's enum Direction; directionNone is Java's null (ADR-029).
type direction int

const (
	directionNone direction = iota
	directionForward
	directionBackward
)

// sym is Java's Strength.Sym.
type sym struct {
	hash int
}

func (s *sym) CustomHash() int {
	return s.hash
}

// Java Strength.ABSOLUTE_STRONGEST, ... (ADR-022, named after their type).
var (
	symAbsoluteStrongest = &sym{0}
	symRequired          = &sym{1}
	symStrongPreferred   = &sym{2}
	symPreferred         = &sym{3}
	symStrongDefault     = &sym{4}
	symDefault           = &sym{5}
	symWeakDefault       = &sym{6}
	symAbsoluteWeakest   = &sym{7}
)

type strength struct {
	arithmeticValue int
	symbolicValue   *sym
}

func newStrength(symbolicValue *sym) *strength {
	return &strength{
		symbolicValue:   symbolicValue,
		arithmeticValue: strengthTable.At(symbolicValue),
	}
}

func (s *strength) sameAs(o *strength) bool {
	return s.arithmeticValue == o.arithmeticValue
}

func (s *strength) stronger(o *strength) bool {
	return s.arithmeticValue < o.arithmeticValue
}

func (s *strength) weaker(o *strength) bool {
	return s.arithmeticValue > o.arithmeticValue
}

func (s *strength) weakest(o *strength) *strength {
	if o.weaker(s) {
		return o
	}
	return s
}

func strengthOf(strength *sym) *strength {
	return strengthConstant.At(strength)
}

func strengthCreateStrengthTable() *som.IdentityDictionary[*sym, int] {
	strengthTable := som.NewIdentityDictionary[*sym, int]()
	strengthTable.AtPut(symAbsoluteStrongest, -10000)
	strengthTable.AtPut(symRequired, -800)
	strengthTable.AtPut(symStrongPreferred, -600)
	strengthTable.AtPut(symPreferred, -400)
	strengthTable.AtPut(symStrongDefault, -200)
	strengthTable.AtPut(symDefault, 0)
	strengthTable.AtPut(symWeakDefault, 500)
	strengthTable.AtPut(symAbsoluteWeakest, 10000)
	return strengthTable
}

func strengthCreateStrengthConstants() *som.IdentityDictionary[*sym, *strength] {
	strengthConstant := som.NewIdentityDictionary[*sym, *strength]()
	strengthTable.GetKeys().ForEach(func(key *sym) {
		strengthConstant.AtPut(key, newStrength(key))
	})
	return strengthConstant
}

// Java's static initialiser of Strength. Go orders package initialisation by
// dependency, so the tables are built before the cached strengths (ADR-022).
var (
	strengthTable           = strengthCreateStrengthTable()
	strengthConstant        = strengthCreateStrengthConstants()
	strengthAbsoluteWeakest = strengthOf(symAbsoluteWeakest)
	strengthRequired        = strengthOf(symRequired)
)

// constraint is Java's abstract class AbstractConstraint (ADR-024).
type constraint interface {
	getStrength() *strength
	isInput() bool
	isSatisfied() bool
	addToGraph()
	removeFromGraph()
	chooseMethod(mark int) direction
	execute()
	inputsDo(fn func(*variable))
	inputsHasOne(fn func(*variable) bool) bool
	markUnsatisfied()
	getOutput() *variable
	recalculate()
}

// The concrete methods of AbstractConstraint call overridden methods, so they
// are package functions over the interface (ADR-028).

func abstractConstraintAddConstraint(c constraint, planner *planner) {
	c.addToGraph()
	planner.incrementalAdd(c)
}

func abstractConstraintDestroyConstraint(c constraint, planner *planner) {
	if c.isSatisfied() {
		planner.incrementalRemove(c)
	}
	c.removeFromGraph()
}

func abstractConstraintInputsKnown(c constraint, mark int) bool {
	return !c.inputsHasOne(func(v *variable) bool {
		return !(v.mark == mark || v.stay || v.determinedBy == nil)
	})
}

func abstractConstraintSatisfy(c constraint, mark int, planner *planner) constraint {
	var overridden constraint

	c.chooseMethod(mark)

	if c.isSatisfied() {
		c.inputsDo(func(in *variable) { in.mark = mark })

		out := c.getOutput()
		overridden = out.determinedBy
		if overridden != nil {
			overridden.markUnsatisfied()
		}
		out.determinedBy = c
		if !planner.addPropagate(c, mark) {
			panic("Cycle encountered")
		}
		out.mark = mark
	} else {
		overridden = nil
		if c.getStrength().sameAs(strengthRequired) {
			panic("Could not satisfy a required constraint")
		}
	}
	return overridden
}

// editConstraint is Java's EditConstraint extends UnaryConstraint (ADR-024/025).
type editConstraint struct {
	strength  *strength
	output    *variable
	satisfied bool
}

func newEditConstraint(v *variable, strength *sym, planner *planner) *editConstraint {
	c := &editConstraint{strength: strengthOf(strength), output: v}
	abstractConstraintAddConstraint(c, planner)
	return c
}

func (c *editConstraint) getStrength() *strength { return c.strength }

func (c *editConstraint) isInput() bool { return true }

func (c *editConstraint) isSatisfied() bool { return c.satisfied }

func (c *editConstraint) addToGraph() {
	c.output.addConstraint(c)
	c.satisfied = false
}

func (c *editConstraint) removeFromGraph() {
	if c.output != nil {
		c.output.removeConstraint(c)
	}
	c.satisfied = false
}

func (c *editConstraint) chooseMethod(mark int) direction {
	c.satisfied = c.output.mark != mark &&
		c.strength.stronger(c.output.walkStrength)
	return directionNone
}

func (c *editConstraint) execute() {} // Edit constraints do nothing.

func (c *editConstraint) inputsDo(func(*variable)) {}

func (c *editConstraint) inputsHasOne(func(*variable) bool) bool { return false }

func (c *editConstraint) markUnsatisfied() { c.satisfied = false }

func (c *editConstraint) getOutput() *variable { return c.output }

func (c *editConstraint) recalculate() {
	c.output.walkStrength = c.strength
	c.output.stay = !c.isInput()
	if c.output.stay {
		c.execute() // stay optimization
	}
}

// stayConstraint is Java's StayConstraint extends UnaryConstraint (ADR-024/025).
type stayConstraint struct {
	strength  *strength
	output    *variable
	satisfied bool
}

func newStayConstraint(v *variable, strength *sym, planner *planner) *stayConstraint {
	c := &stayConstraint{strength: strengthOf(strength), output: v}
	abstractConstraintAddConstraint(c, planner)
	return c
}

func (c *stayConstraint) getStrength() *strength { return c.strength }

func (c *stayConstraint) isInput() bool { return false }

func (c *stayConstraint) isSatisfied() bool { return c.satisfied }

func (c *stayConstraint) addToGraph() {
	c.output.addConstraint(c)
	c.satisfied = false
}

func (c *stayConstraint) removeFromGraph() {
	if c.output != nil {
		c.output.removeConstraint(c)
	}
	c.satisfied = false
}

func (c *stayConstraint) chooseMethod(mark int) direction {
	c.satisfied = c.output.mark != mark &&
		c.strength.stronger(c.output.walkStrength)
	return directionNone
}

func (c *stayConstraint) execute() {} // Stay constraints do nothing.

func (c *stayConstraint) inputsDo(func(*variable)) {}

func (c *stayConstraint) inputsHasOne(func(*variable) bool) bool { return false }

func (c *stayConstraint) markUnsatisfied() { c.satisfied = false }

func (c *stayConstraint) getOutput() *variable { return c.output }

func (c *stayConstraint) recalculate() {
	c.output.walkStrength = c.strength
	c.output.stay = !c.isInput()
	if c.output.stay {
		c.execute() // stay optimization
	}
}

// equalityConstraint is Java's EqualityConstraint extends BinaryConstraint
// (ADR-024/025).
type equalityConstraint struct {
	strength  *strength
	v1        *variable
	v2        *variable
	direction direction
}

func newEqualityConstraint(var1, var2 *variable, strength *sym, planner *planner) *equalityConstraint {
	c := &equalityConstraint{
		strength:  strengthOf(strength),
		v1:        var1,
		v2:        var2,
		direction: directionNone,
	}
	abstractConstraintAddConstraint(c, planner)
	return c
}

func (c *equalityConstraint) getStrength() *strength { return c.strength }

func (c *equalityConstraint) isInput() bool { return false }

func (c *equalityConstraint) isSatisfied() bool { return c.direction != directionNone }

func (c *equalityConstraint) addToGraph() {
	c.v1.addConstraint(c)
	c.v2.addConstraint(c)
	c.direction = directionNone
}

func (c *equalityConstraint) removeFromGraph() {
	if c.v1 != nil {
		c.v1.removeConstraint(c)
	}
	if c.v2 != nil {
		c.v2.removeConstraint(c)
	}
	c.direction = directionNone
}

func (c *equalityConstraint) chooseMethod(mark int) direction {
	if c.v1.mark == mark {
		if c.v2.mark != mark && c.strength.stronger(c.v2.walkStrength) {
			c.direction = directionForward
			return c.direction
		} else {
			c.direction = directionNone
			return c.direction
		}
	}

	if c.v2.mark == mark {
		if c.v1.mark != mark && c.strength.stronger(c.v1.walkStrength) {
			c.direction = directionBackward
			return c.direction
		} else {
			c.direction = directionNone
			return c.direction
		}
	}

	if c.v1.walkStrength.weaker(c.v2.walkStrength) {
		if c.strength.stronger(c.v1.walkStrength) {
			c.direction = directionBackward
			return c.direction
		} else {
			c.direction = directionNone
			return c.direction
		}
	} else {
		if c.strength.stronger(c.v2.walkStrength) {
			c.direction = directionForward
			return c.direction
		} else {
			c.direction = directionNone
			return c.direction
		}
	}
}

func (c *equalityConstraint) execute() {
	if c.direction == directionForward {
		c.v2.value = c.v1.value
	} else {
		c.v1.value = c.v2.value
	}
}

func (c *equalityConstraint) inputsDo(fn func(*variable)) {
	if c.direction == directionForward {
		fn(c.v1)
	} else {
		fn(c.v2)
	}
}

func (c *equalityConstraint) inputsHasOne(fn func(*variable) bool) bool {
	if c.direction == directionForward {
		return fn(c.v1)
	} else {
		return fn(c.v2)
	}
}

func (c *equalityConstraint) markUnsatisfied() { c.direction = directionNone }

func (c *equalityConstraint) getOutput() *variable {
	if c.direction == directionForward {
		return c.v2
	}
	return c.v1
}

func (c *equalityConstraint) recalculate() {
	var in, out *variable

	if c.direction == directionForward {
		in = c.v1
		out = c.v2
	} else {
		in = c.v2
		out = c.v1
	}

	out.walkStrength = c.strength.weakest(in.walkStrength)
	out.stay = in.stay
	if out.stay {
		c.execute()
	}
}

// scaleConstraint is Java's ScaleConstraint extends BinaryConstraint
// (ADR-024/025).
type scaleConstraint struct {
	strength  *strength
	v1        *variable
	v2        *variable
	direction direction
	scale     *variable // scale factor input variable
	offset    *variable // offset input variable
}

func newScaleConstraint(src, scale, offset, dest *variable, strength *sym, planner *planner) *scaleConstraint {
	c := &scaleConstraint{
		strength:  strengthOf(strength),
		v1:        src,
		v2:        dest,
		direction: directionNone,
	}
	c.scale = scale
	c.offset = offset
	abstractConstraintAddConstraint(c, planner)
	return c
}

func (c *scaleConstraint) getStrength() *strength { return c.strength }

func (c *scaleConstraint) isInput() bool { return false }

func (c *scaleConstraint) isSatisfied() bool { return c.direction != directionNone }

func (c *scaleConstraint) addToGraph() {
	c.v1.addConstraint(c)
	c.v2.addConstraint(c)
	c.scale.addConstraint(c)
	c.offset.addConstraint(c)
	c.direction = directionNone
}

func (c *scaleConstraint) removeFromGraph() {
	if c.v1 != nil {
		c.v1.removeConstraint(c)
	}
	if c.v2 != nil {
		c.v2.removeConstraint(c)
	}
	if c.scale != nil {
		c.scale.removeConstraint(c)
	}
	if c.offset != nil {
		c.offset.removeConstraint(c)
	}
	c.direction = directionNone
}

func (c *scaleConstraint) chooseMethod(mark int) direction {
	if c.v1.mark == mark {
		if c.v2.mark != mark && c.strength.stronger(c.v2.walkStrength) {
			c.direction = directionForward
			return c.direction
		} else {
			c.direction = directionNone
			return c.direction
		}
	}

	if c.v2.mark == mark {
		if c.v1.mark != mark && c.strength.stronger(c.v1.walkStrength) {
			c.direction = directionBackward
			return c.direction
		} else {
			c.direction = directionNone
			return c.direction
		}
	}

	if c.v1.walkStrength.weaker(c.v2.walkStrength) {
		if c.strength.stronger(c.v1.walkStrength) {
			c.direction = directionBackward
			return c.direction
		} else {
			c.direction = directionNone
			return c.direction
		}
	} else {
		if c.strength.stronger(c.v2.walkStrength) {
			c.direction = directionForward
			return c.direction
		} else {
			c.direction = directionNone
			return c.direction
		}
	}
}

func (c *scaleConstraint) execute() {
	if c.direction == directionForward {
		c.v2.value = c.v1.value*c.scale.value + c.offset.value
	} else {
		c.v1.value = (c.v2.value - c.offset.value) / c.scale.value
	}
}

func (c *scaleConstraint) inputsDo(fn func(*variable)) {
	if c.direction == directionForward {
		fn(c.v1)
		fn(c.scale)
		fn(c.offset)
	} else {
		fn(c.v2)
		fn(c.scale)
		fn(c.offset)
	}
}

func (c *scaleConstraint) inputsHasOne(fn func(*variable) bool) bool {
	if c.direction == directionForward {
		return fn(c.v1)
	} else {
		return fn(c.v2)
	}
}

func (c *scaleConstraint) markUnsatisfied() { c.direction = directionNone }

func (c *scaleConstraint) getOutput() *variable {
	if c.direction == directionForward {
		return c.v2
	}
	return c.v1
}

func (c *scaleConstraint) recalculate() {
	var in, out *variable

	if c.direction == directionForward {
		in = c.v1
		out = c.v2
	} else {
		out = c.v1
		in = c.v2
	}

	out.walkStrength = c.strength.weakest(in.walkStrength)
	out.stay = in.stay && c.scale.stay && c.offset.stay
	if out.stay {
		c.execute() // stay optimization
	}
}

// plan is Java's Plan extends Vector<AbstractConstraint>: a defined type over
// som.Vector, so a plan is one object plus its storage, as in Java.
type plan som.Vector[constraint]

func newPlan() *plan {
	return (*plan)(som.NewVectorSize[constraint](15))
}

func (p *plan) execute() {
	(*som.Vector[constraint])(p).ForEach(func(c constraint) { c.execute() })
}

type variable struct {
	value        int                     // my value; changed by constraints
	constraints  *som.Vector[constraint] // normal constraints that reference me
	determinedBy constraint              // the constraint that currently determines my value
	mark         int                     // used by the planner to mark constraints
	walkStrength *strength               // my walkabout strength
	stay         bool                    // true if I am a planning-time constant
}

func newVariable() *variable {
	return &variable{
		value:        0,
		constraints:  som.NewVectorSize[constraint](2),
		determinedBy: nil,
		walkStrength: strengthAbsoluteWeakest,
		stay:         true,
		mark:         0,
	}
}

func variableValue(aValue int) *variable {
	v := newVariable()
	v.value = aValue
	return v
}

func (v *variable) addConstraint(c constraint) {
	v.constraints.Append(c)
}

func (v *variable) removeConstraint(c constraint) {
	v.constraints.Remove(c)
	if v.determinedBy == c {
		v.determinedBy = nil
	}
}

type planner struct {
	currentMark int
}

func newPlanner() *planner {
	return &planner{currentMark: 1}
}

func (p *planner) incrementalAdd(c constraint) {
	mark := p.newMark()
	overridden := abstractConstraintSatisfy(c, mark, p)

	for overridden != nil {
		overridden = abstractConstraintSatisfy(overridden, mark, p)
	}
}

func (p *planner) incrementalRemove(c constraint) {
	out := c.getOutput()
	c.markUnsatisfied()
	c.removeFromGraph()

	unsatisfied := p.removePropagateFrom(out)
	unsatisfied.ForEach(func(u constraint) { p.incrementalAdd(u) })
}

func (p *planner) extractPlanFromConstraints(constraints *som.Vector[constraint]) *plan {
	sources := som.NewVector[constraint]()

	constraints.ForEach(func(c constraint) {
		if c.isInput() && c.isSatisfied() {
			sources.Append(c)
		}
	})

	return p.makePlan(sources)
}

func (p *planner) makePlan(sources *som.Vector[constraint]) *plan {
	mark := p.newMark()
	plan := newPlan()
	todo := sources

	for !todo.IsEmpty() {
		c := todo.RemoveFirst()

		if c.getOutput().mark != mark && abstractConstraintInputsKnown(c, mark) {
			(*som.Vector[constraint])(plan).Append(c)
			c.getOutput().mark = mark
			p.addConstraintsConsumingTo(c.getOutput(), todo)
		}
	}
	return plan
}

func (p *planner) addConstraintsConsumingTo(v *variable, coll *som.Vector[constraint]) {
	determiningC := v.determinedBy

	v.constraints.ForEach(func(c constraint) {
		if c != determiningC && c.isSatisfied() {
			coll.Append(c)
		}
	})
}

func (p *planner) addPropagate(c constraint, mark int) bool {
	todo := som.VectorWith(c)

	for !todo.IsEmpty() {
		d := todo.RemoveFirst()

		if d.getOutput().mark == mark {
			p.incrementalRemove(c)
			return false
		}

		d.recalculate()
		p.addConstraintsConsumingTo(d.getOutput(), todo)
	}
	return true
}

func (p *planner) change(v *variable, newValue int) {
	editC := newEditConstraint(v, symPreferred, p)

	editV := som.VectorWith[constraint](editC)
	plan := p.extractPlanFromConstraints(editV)
	for i := 0; i < 10; i++ {
		v.value = newValue
		plan.execute()
	}
	abstractConstraintDestroyConstraint(editC, p)
}

func (p *planner) constraintsConsuming(v *variable, fn func(constraint)) {
	determiningC := v.determinedBy
	v.constraints.ForEach(func(c constraint) {
		if c != determiningC && c.isSatisfied() {
			fn(c)
		}
	})
}

func (p *planner) newMark() int {
	p.currentMark++
	return p.currentMark
}

func (p *planner) removePropagateFrom(out *variable) *som.Vector[constraint] {
	unsatisfied := som.NewVector[constraint]()

	out.determinedBy = nil
	out.walkStrength = strengthAbsoluteWeakest
	out.stay = true

	todo := som.VectorWith(out)

	for !todo.IsEmpty() {
		v := todo.RemoveFirst()

		v.constraints.ForEach(func(c constraint) {
			if !c.isSatisfied() {
				unsatisfied.Append(c)
			}
		})

		p.constraintsConsuming(v, func(c constraint) {
			c.recalculate()
			todo.Append(c.getOutput())
		})
	}

	unsatisfied.Sort(func(c1, c2 constraint) int {
		if c1.getStrength().stronger(c2.getStrength()) {
			return -1
		}
		return 1
	})
	return unsatisfied
}

func plannerChainTest(n int) {
	planner := newPlanner()
	vars := make([]*variable, n+1)
	for i := range vars {
		vars[i] = newVariable()
	}

	// Build chain of n equality constraints
	for i := 0; i < n; i++ {
		v1 := vars[i]
		v2 := vars[i+1]
		newEqualityConstraint(v1, v2, symRequired, planner)
	}

	newStayConstraint(vars[n], symStrongDefault, planner)
	var editC constraint = newEditConstraint(vars[0], symPreferred, planner)

	editV := som.VectorWith(editC)
	plan := planner.extractPlanFromConstraints(editV)
	for i := 0; i < 100; i++ {
		vars[0].value = i
		plan.execute()
		if vars[n].value != i {
			panic("Chain test failed!")
		}
	}
	abstractConstraintDestroyConstraint(editC, planner)
}

func plannerProjectionTest(n int) {
	planner := newPlanner()
	dests := som.NewVector[*variable]()

	scale := variableValue(10)
	offset := variableValue(1000)

	var src, dst *variable
	for i := 1; i <= n; i++ {
		src = variableValue(i)
		dst = variableValue(i)
		dests.Append(dst)
		newStayConstraint(src, symDefault, planner)
		newScaleConstraint(src, scale, offset, dst, symRequired, planner)
	}

	planner.change(src, 17)
	if dst.value != 1170 {
		panic("Projection test 1 failed!")
	}

	planner.change(dst, 1050)
	if src.value != 5 {
		panic("Projection test 2 failed!")
	}

	planner.change(scale, 5)
	for i := 0; i < n-1; i++ {
		if dests.At(i).value != (i+1)*5+1000 {
			panic("Projection test 3 failed!")
		}
	}

	planner.change(offset, 2000)
	for i := 0; i < n-1; i++ {
		if dests.At(i).value != (i+1)*5+2000 {
			panic("Projection test 4 failed!")
		}
	}
}
