package main

import (
	"fmt"
	"math"

	"awfy/som"
)

// Havlak is the loop recognition benchmark. The Java package havlak
// (LoopTesterApp, HavlakLoopFinder, ControlFlowGraph, ...) lives in this file.
type Havlak struct{}

// InnerBenchmarkLoop overrides the default inner loop (ADR-003).
func (h *Havlak) InnerBenchmarkLoop(innerIterations int) bool {
	return h.verify(newLoopTesterApp().main(
		innerIterations, 50, 10 /* was 100 */, 10, 5), innerIterations)
}

// verify is Java's overload verifyResult(int[], int) (ADR-010).
func (*Havlak) verify(r []int, innerIterations int) bool {
	if innerIterations == 15000 {
		return r[0] == 46602 && r[1] == 5213
	}
	if innerIterations == 1500 {
		return r[0] == 6102 && r[1] == 5213
	}
	if innerIterations == 150 {
		return r[0] == 2052 && r[1] == 5213
	}
	if innerIterations == 15 {
		return r[0] == 1647 && r[1] == 5213
	}
	if innerIterations == 1 {
		return r[0] == 1605 && r[1] == 5213
	}

	fmt.Printf("No verification result for %d found\n", innerIterations)
	fmt.Printf("Result is: %d, %d\n", r[0], r[1])
	return false
}

// Benchmark is never called, InnerBenchmarkLoop is used instead.
func (*Havlak) Benchmark() any {
	panic("Should never be reached")
}

// VerifyResult is never called, InnerBenchmarkLoop is used instead.
func (*Havlak) VerifyResult(any) bool {
	panic("Should never be reached")
}

// basicBlock only maintains a vector of in-edges and a vector of out-edges.
type basicBlock struct {
	inEdges  *som.Vector[*basicBlock]
	outEdges *som.Vector[*basicBlock]
	name     int
}

func newBasicBlock(name int) *basicBlock {
	b := &basicBlock{}
	b.name = name
	b.inEdges = som.NewVectorSize[*basicBlock](2)
	b.outEdges = som.NewVectorSize[*basicBlock](2)
	return b
}

func (b *basicBlock) getNumPred() int {
	return b.inEdges.Size()
}

func (b *basicBlock) addOutEdge(to *basicBlock) {
	b.outEdges.Append(to)
}

func (b *basicBlock) addInEdge(from *basicBlock) {
	b.inEdges.Append(from)
}

// CustomHash implements som's custom hash for IdentityDictionary keys.
func (b *basicBlock) CustomHash() int {
	return b.name
}

// basicBlockEdge is an edge between basic blocks. Its fields are written but
// never read, as in Java (ADR-030).
type basicBlockEdge struct {
	from *basicBlock
	to   *basicBlock
}

func newBasicBlockEdge(cfg *controlFlowGraph, fromName, toName int) *basicBlockEdge {
	e := &basicBlockEdge{}
	e.from = cfg.createNode(fromName)
	e.to = cfg.createNode(toName)

	e.from.addOutEdge(e.to)
	e.to.addInEdge(e.from)

	cfg.addEdge(e)
	return e
}

// controlFlowGraph maintains a list of nodes, plus a start node.
type controlFlowGraph struct {
	basicBlockMap *som.Vector[*basicBlock]
	startNode     *basicBlock
	edgeList      *som.Vector[*basicBlockEdge]
}

func newControlFlowGraph() *controlFlowGraph {
	c := &controlFlowGraph{}
	c.startNode = nil
	c.basicBlockMap = som.NewVector[*basicBlock]()
	c.edgeList = som.NewVector[*basicBlockEdge]()
	return c
}

func (c *controlFlowGraph) createNode(name int) *basicBlock {
	var node *basicBlock
	if c.basicBlockMap.At(name) != nil {
		node = c.basicBlockMap.At(name)
	} else {
		node = newBasicBlock(name)
		c.basicBlockMap.AtPut(name, node)
	}

	if c.getNumNodes() == 1 {
		c.startNode = node
	}
	return node
}

func (c *controlFlowGraph) addEdge(edge *basicBlockEdge) {
	c.edgeList.Append(edge)
}

func (c *controlFlowGraph) getNumNodes() int {
	return c.basicBlockMap.Size()
}

// simpleLoop is the basic representation of a loop. header, isReducible,
// counter and depthLevel are written but never read, as in Java (ADR-030).
type simpleLoop struct {
	basicBlocks *som.IdentitySet[*basicBlock]
	children    *som.IdentitySet[*simpleLoop]
	parent      *simpleLoop

	header      *basicBlock
	isReducible bool

	isRoot       bool
	nestingLevel int

	counter    int
	depthLevel int
}

func newSimpleLoop(bb *basicBlock, isReducible bool) *simpleLoop {
	l := &simpleLoop{}
	l.isReducible = isReducible
	l.parent = nil
	l.isRoot = false
	l.nestingLevel = 0
	l.depthLevel = 0
	l.basicBlocks = som.NewIdentitySet[*basicBlock]()
	l.children = som.NewIdentitySet[*simpleLoop]()

	if bb != nil {
		l.basicBlocks.Add(bb)
	}
	l.header = bb
	return l
}

func (l *simpleLoop) addNode(bb *basicBlock) {
	l.basicBlocks.Add(bb)
}

func (l *simpleLoop) addChildLoop(loop *simpleLoop) {
	l.children.Add(loop)
}

func (l *simpleLoop) setParent(parent *simpleLoop) {
	l.parent = parent
	l.parent.addChildLoop(l)
}

func (l *simpleLoop) setNestingLevel(level int) {
	l.nestingLevel = level
	if level == 0 {
		l.isRoot = true
	}
}

// loopStructureGraph maintains the loop structure for a given CFG.
type loopStructureGraph struct {
	root        *simpleLoop
	loops       *som.Vector[*simpleLoop]
	loopCounter int
}

func newLoopStructureGraph() *loopStructureGraph {
	g := &loopStructureGraph{}
	g.loopCounter = 0
	g.loops = som.NewVector[*simpleLoop]()
	g.root = newSimpleLoop(nil, true)
	g.root.setNestingLevel(0)
	g.root.counter = g.loopCounter
	g.loopCounter += 1
	g.loops.Append(g.root)
	return g
}

func (g *loopStructureGraph) createNewLoop(bb *basicBlock, isReducible bool) *simpleLoop {
	loop := newSimpleLoop(bb, isReducible)
	loop.counter = g.loopCounter
	g.loopCounter += 1
	g.loops.Append(loop)
	return loop
}

func (g *loopStructureGraph) calculateNestingLevel() {
	// link up all 1st level loops to artificial root node.
	g.loops.ForEach(func(liter *simpleLoop) {
		if !liter.isRoot {
			if liter.parent == nil {
				liter.setParent(g.root)
			}
		}
	})

	// recursively traverse the tree and assign levels.
	g.calculateNestingLevelRec(g.root, 0)
}

func (g *loopStructureGraph) calculateNestingLevelRec(loop *simpleLoop, depth int) {
	loop.depthLevel = depth
	loop.children.ForEach(func(liter *simpleLoop) {
		g.calculateNestingLevelRec(liter, depth+1)

		loop.setNestingLevel(max(loop.nestingLevel, 1+liter.nestingLevel))
	})
}

func (g *loopStructureGraph) getNumLoops() int {
	return g.loops.Size()
}

// unionFindNode is a node of the Union/Find algorithm used to collapse
// complete loops into a single node.
type unionFindNode struct {
	parent    *unionFindNode
	bb        *basicBlock
	loop      *simpleLoop
	dfsNumber int
}

func (u *unionFindNode) initNode(bb *basicBlock, dfsNumber int) {
	u.parent = u
	u.bb = bb
	u.dfsNumber = dfsNumber
	u.loop = nil
}

// findSet is the find routine, with path compression.
func (u *unionFindNode) findSet() *unionFindNode {
	nodeList := som.NewVector[*unionFindNode]()

	node := u
	for node != node.parent {
		if node.parent != node.parent.parent {
			nodeList.Append(node)
		}
		node = node.parent
	}

	// Path Compression, all nodes' parents point to the 1st level parent.
	nodeList.ForEach(func(iter *unionFindNode) { iter.union(u.parent) })
	return node
}

func (u *unionFindNode) union(basicBlock *unionFindNode) {
	u.parent = basicBlock
}

// basicBlockClass is Java's enum BasicBlockClass (ADR-029).
type basicBlockClass int

const (
	basicBlockClassBBTop         basicBlockClass = iota // uninitialized
	basicBlockClassBBNonHeader                          // a regular BB
	basicBlockClassBBReducible                          // reducible loop
	basicBlockClassBBSelf                               // single BB loop
	basicBlockClassBBIrreducible                        // irreducible loop
	basicBlockClassBBDead                               // a dead BB
	basicBlockClassBBLast                               // Sentinel
)

const (
	// Marker for uninitialized nodes (Java Integer.MAX_VALUE).
	unvisited int = math.MaxInt32
	// Safeguard against pathological algorithm behavior.
	maxNonBackPreds int = 32 * 1024
)

// havlakLoopFinder encapsulates the complete Havlak loop finding algorithm.
type havlakLoopFinder struct {
	cfg *controlFlowGraph
	lsg *loopStructureGraph

	nonBackPreds *som.Vector[*som.Set[int]]
	backPreds    *som.Vector[*som.Vector[int]]
	number       *som.IdentityDictionary[*basicBlock, int]
	maxSize      int
	header       []int
	typ          []basicBlockClass // Java: type
	last         []int
	nodes        []*unionFindNode
}

func newHavlakLoopFinder(cfg *controlFlowGraph, lsg *loopStructureGraph) *havlakLoopFinder {
	f := &havlakLoopFinder{}
	f.nonBackPreds = som.NewVector[*som.Set[int]]()
	f.backPreds = som.NewVector[*som.Vector[int]]()
	f.number = som.NewIdentityDictionary[*basicBlock, int]()
	f.maxSize = 0
	f.cfg = cfg
	f.lsg = lsg
	return f
}

// isAncestor determines whether node w is a "true" ancestor for node v.
func (f *havlakLoopFinder) isAncestor(w, v int) bool {
	return w <= v && v <= f.last[w]
}

// doDFS is a depth first traversal along out edges with node numbering.
func (f *havlakLoopFinder) doDFS(currentNode *basicBlock, current int) int {
	f.nodes[current].initNode(currentNode, current)
	f.number.AtPut(currentNode, current)

	lastID := current
	outerBlocks := currentNode.outEdges

	for i := 0; i < outerBlocks.Size(); i++ {
		target := outerBlocks.At(i)
		if f.number.At(target) == unvisited {
			lastID = f.doDFS(target, lastID+1)
		}
	}

	f.last[current] = lastID
	return lastID
}

func (f *havlakLoopFinder) initAllNodes() {
	// Step a:
	//   - initialize all nodes as unvisited.
	//   - depth-first traversal and numbering.
	//   - unreached BB's are marked as dead.
	f.cfg.basicBlockMap.ForEach(func(bb *basicBlock) {
		f.number.AtPut(bb, unvisited)
	})

	f.doDFS(f.cfg.startNode, 0)
}

func (f *havlakLoopFinder) identifyEdges(size int) {
	// Step b:
	//   - iterate over all nodes.
	//
	//   A backedge comes from a descendant in the DFS tree, and non-backedges
	//   from non-descendants (following Tarjan).
	//
	//   - check incoming edges 'v' and add them to either
	//     - the list of backedges (backPreds) or
	//     - the list of non-backedges (nonBackPreds)
	for w := 0; w < size; w++ {
		f.header[w] = 0
		f.typ[w] = basicBlockClassBBNonHeader

		nodeW := f.nodes[w].bb
		if nodeW == nil {
			f.typ[w] = basicBlockClassBBDead
		} else {
			f.processEdges(nodeW, w)
		}
	}
}

func (f *havlakLoopFinder) processEdges(nodeW *basicBlock, w int) {
	if nodeW.getNumPred() > 0 {
		nodeW.inEdges.ForEach(func(nodeV *basicBlock) {
			v := f.number.At(nodeV)
			if v != unvisited {
				if f.isAncestor(w, v) {
					f.backPreds.At(w).Append(v)
				} else {
					f.nonBackPreds.At(w).Add(v)
				}
			}
		})
	}
}

// findLoops finds loops and builds the loop forest using Havlak's algorithm.
func (f *havlakLoopFinder) findLoops() {
	if f.cfg.startNode == nil {
		return
	}

	size := f.cfg.getNumNodes()

	f.nonBackPreds.RemoveAll()
	f.backPreds.RemoveAll()
	f.number.RemoveAll()
	if size > f.maxSize {
		f.header = make([]int, size)
		f.typ = make([]basicBlockClass, size)
		f.last = make([]int, size)
		f.nodes = make([]*unionFindNode, size)
		f.maxSize = size
	}

	for i := 0; i < size; i++ {
		f.nonBackPreds.Append(som.NewSet[int]())
		f.backPreds.Append(som.NewVector[int]())
		f.nodes[i] = &unionFindNode{}
	}

	f.initAllNodes()
	f.identifyEdges(size)

	// Start node is root of all other loops.
	f.header[0] = 0

	// Step c:
	//
	// The outer loop, unchanged from Tarjan. It does nothing except
	// for those nodes which are the destinations of backedges.
	// For a header node w, we chase backward from the sources of the
	// backedges adding nodes to the set P, representing the body of
	// the loop headed by w.
	//
	// By running through the nodes in reverse of the DFST preorder,
	// we ensure that inner loop headers will be processed before the
	// headers for surrounding loops.
	for w := size - 1; w >= 0; w-- {
		// this is 'P' in Havlak's paper
		nodePool := som.NewVector[*unionFindNode]()

		nodeW := f.nodes[w].bb
		if nodeW != nil {
			f.stepD(w, nodePool)

			// Copy nodePool to workList.
			workList := som.NewVector[*unionFindNode]()
			nodePool.ForEach(func(niter *unionFindNode) { workList.Append(niter) })

			if nodePool.Size() != 0 {
				f.typ[w] = basicBlockClassBBReducible
			}

			// work the list...
			for !workList.IsEmpty() {
				x := workList.RemoveFirst()

				// Step e:
				//
				// Step e represents the main difference from Tarjan's method.
				// Chasing upwards from the sources of a node w's backedges. If
				// there is a node y' that is not a descendant of w, w is marked
				// the header of an irreducible loop, there is another entry
				// into this loop that avoids w.

				// The algorithm has degenerated. Break and
				// return in this case.
				nonBackSize := f.nonBackPreds.At(x.dfsNumber).Size()
				if nonBackSize > maxNonBackPreds {
					return
				}
				f.stepEProcessNonBackPreds(w, nodePool, workList, x)
			}

			// Collapse/Unionize nodes in a SCC to a single node
			// For every SCC found, create a loop descriptor and link it in.
			if (nodePool.Size() > 0) || (f.typ[w] == basicBlockClassBBSelf) {
				loop := f.lsg.createNewLoop(nodeW, f.typ[w] != basicBlockClassBBIrreducible)
				f.setLoopAttributes(w, nodePool, loop)
			}
		}
	} // Step c
}

func (f *havlakLoopFinder) stepEProcessNonBackPreds(w int, nodePool, workList *som.Vector[*unionFindNode],
	x *unionFindNode) {
	f.nonBackPreds.At(x.dfsNumber).ForEach(func(iter int) {
		y := f.nodes[iter]
		ydash := y.findSet()

		if !f.isAncestor(w, ydash.dfsNumber) {
			f.typ[w] = basicBlockClassBBIrreducible
			f.nonBackPreds.At(w).Add(ydash.dfsNumber)
		} else {
			if ydash.dfsNumber != w {
				if !nodePool.HasSome(func(e *unionFindNode) bool { return e == ydash }) {
					workList.Append(ydash)
					nodePool.Append(ydash)
				}
			}
		}
	})
}

func (f *havlakLoopFinder) setLoopAttributes(w int, nodePool *som.Vector[*unionFindNode], loop *simpleLoop) {
	// At this point, one can set attributes to the loop, such as:
	//
	// the bottom node:
	//    iter  = backPreds[w].begin();
	//    loop bottom is: nodes[iter].node);
	//
	// the number of backedges:
	//    backPreds[w].size()
	//
	// whether this loop is reducible:
	//    type[w] != BasicBlockClass.BB_IRREDUCIBLE
	f.nodes[w].loop = loop

	nodePool.ForEach(func(node *unionFindNode) {
		// Add nodes to loop descriptor.
		f.header[node.dfsNumber] = w
		node.union(f.nodes[w])

		// Nested loops are not added, but linked together.
		if node.loop != nil {
			node.loop.setParent(loop)
		} else {
			loop.addNode(node.bb)
		}
	})
}

func (f *havlakLoopFinder) stepD(w int, nodePool *som.Vector[*unionFindNode]) {
	f.backPreds.At(w).ForEach(func(v int) {
		if v != w {
			nodePool.Append(f.nodes[v].findSet())
		} else {
			f.typ[w] = basicBlockClassBBSelf
		}
	})
}

// loopTesterApp builds the control flow graph and runs the loop finder.
type loopTesterApp struct {
	cfg *controlFlowGraph
	lsg *loopStructureGraph
}

func newLoopTesterApp() *loopTesterApp {
	a := &loopTesterApp{}
	a.cfg = newControlFlowGraph()
	a.lsg = newLoopStructureGraph()
	a.cfg.createNode(0)
	return a
}

// buildDiamond creates 4 basic blocks, corresponding to an if/then/else
// clause with a CFG that looks like a diamond.
func (a *loopTesterApp) buildDiamond(start int) int {
	bb0 := start
	newBasicBlockEdge(a.cfg, bb0, bb0+1)
	newBasicBlockEdge(a.cfg, bb0, bb0+2)
	newBasicBlockEdge(a.cfg, bb0+1, bb0+3)
	newBasicBlockEdge(a.cfg, bb0+2, bb0+3)

	return bb0 + 3
}

// buildConnect connects two existing nodes.
func (a *loopTesterApp) buildConnect(start, end int) {
	newBasicBlockEdge(a.cfg, start, end)
}

// buildStraight forms a straight connected sequence of n basic blocks.
func (a *loopTesterApp) buildStraight(start, n int) int {
	for i := 0; i < n; i++ {
		a.buildConnect(start+i, start+i+1)
	}
	return start + n
}

// buildBaseLoop constructs a simple loop with two diamonds in it.
func (a *loopTesterApp) buildBaseLoop(from int) int {
	header := a.buildStraight(from, 1)
	diamond1 := a.buildDiamond(header)
	d11 := a.buildStraight(diamond1, 1)
	diamond2 := a.buildDiamond(d11)
	footer := a.buildStraight(diamond2, 1)
	a.buildConnect(diamond2, d11)
	a.buildConnect(diamond1, header)

	a.buildConnect(footer, from)
	footer = a.buildStraight(footer, 1)
	return footer
}

func (a *loopTesterApp) main(numDummyLoops, findLoopIterations, parLoops, pparLoops, ppparLoops int) []int {
	a.constructSimpleCFG()
	a.addDummyLoops(numDummyLoops)
	a.constructCFG(parLoops, pparLoops, ppparLoops)

	// Performing Loop Recognition, 1 Iteration, then findLoopIteration
	a.findLoops(a.lsg)
	for i := 0; i < findLoopIterations; i++ {
		a.findLoops(newLoopStructureGraph())
	}

	a.lsg.calculateNestingLevel()
	return []int{a.lsg.getNumLoops(), a.cfg.getNumNodes()}
}

func (a *loopTesterApp) constructCFG(parLoops, pparLoops, ppparLoops int) {
	n := 2

	for parlooptrees := 0; parlooptrees < parLoops; parlooptrees++ {
		a.cfg.createNode(n + 1)
		a.buildConnect(2, n+1)
		n += 1

		for i := 0; i < pparLoops; i++ {
			top := n
			n = a.buildStraight(n, 1)
			for j := 0; j < ppparLoops; j++ {
				n = a.buildBaseLoop(n)
			}
			bottom := a.buildStraight(n, 1)
			a.buildConnect(n, top)
			n = bottom
		}
		a.buildConnect(n, 1)
	}
}

func (a *loopTesterApp) addDummyLoops(numDummyLoops int) {
	for dummyloop := 0; dummyloop < numDummyLoops; dummyloop++ {
		a.findLoops(a.lsg)
	}
}

func (a *loopTesterApp) findLoops(loopStructure *loopStructureGraph) {
	finder := newHavlakLoopFinder(a.cfg, loopStructure)
	finder.findLoops()
}

func (a *loopTesterApp) constructSimpleCFG() {
	a.cfg.createNode(0)
	a.buildBaseLoop(0)
	a.cfg.createNode(1)
	newBasicBlockEdge(a.cfg, 0, 2)
}
