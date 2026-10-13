using Benchmarks.Som;

namespace Benchmarks.HavlakSolver;

internal sealed class HavlakLoopFinder
{
    private readonly ControlFlowGraph _cfg; // Control Flow Graph
    private readonly LoopStructureGraph _lsg; // Loop Structure Graph

    // Marker for uninitialized nodes.
    private const int Unvisited = int.MaxValue;

    // Safeguard against pathological algorithm behavior.
    private const int MaxNonBackPreds = 32 * 1024;

    private readonly Vector<Set<int>> _nonBackPreds = new();
    private readonly Vector<Vector<int>> _backPreds = new();
    private readonly IdentityDictionary<BasicBlock, int> _number = new();
    private int _maxSize;
    private int[]? _header;
    private BasicBlockClass[]? _type;
    private int[]? _last;
    private UnionFindNode[]? _nodes;

    public HavlakLoopFinder(ControlFlowGraph cfg, LoopStructureGraph lsg)
    {
        _cfg = cfg;
        _lsg = lsg;
    }

    // Basic Blocks and Loops are being classified as regular, irreducible,
    // and so on. This enum contains a symbolic name for all these
    // classifications.
    private enum BasicBlockClass
    {
        BbTop, // uninitialized
        BbNonHeader, // a regular BB
        BbReducible, // reducible loop
        BbSelf, // single BB loop
        BbIrreducible, // irreducible loop
        BbDead, // a dead BB
        BbLast // Sentinel
    }

    // As described in the paper, determine whether a node 'w' is a
    // "true" ancestor for node 'v'.
    //
    // Dominance can be tested quickly using a pre-order trick
    // for depth-first spanning trees. This is why DFS is the first
    // thing we run below.
    private bool IsAncestor(int w, int v) => w <= v && v <= _last![w];

    // DFS - Depth-First-Search
    //
    // Simple depth first traversal along out edges with node numbering.
    // ReSharper disable once InconsistentNaming
    private int DoDFS(BasicBlock currentNode, int current)
    {
        _nodes![current].InitNode(currentNode, current);
        _number.AtPut(currentNode, current);

        int lastId = current;
        Vector<BasicBlock> outerBlocks = currentNode.OutEdges;

        for (int i = 0; i < outerBlocks.Size(); i++)
        {
            BasicBlock? target = outerBlocks.At(i);
            if (_number.At(target!) == Unvisited)
            {
                lastId = DoDFS(target!, lastId + 1);
            }
        }

        _last![current] = lastId;
        return lastId;
    }

    private void InitAllNodes()
    {
        // Step a:
        //   - initialize all nodes as unvisited.
        //   - depth-first traversal and numbering.
        //   - unreached BB's are marked as dead.
        _cfg.BasicBlocks.ForEach(bb => _number.AtPut(bb, Unvisited));

        DoDFS(_cfg.GetStartBasicBlock()!, 0);
    }

    private void IdentifyEdges(int size)
    {
        // Step b:
        //   - iterate over all nodes.
        //
        //   A backedge comes from a descendant in the DFS tree, and non-backedges
        //   from non-descendants (following Tarjan).
        //
        //   - check incoming edges 'v' and add them to either
        //     - the list of backedges (backPreds) or
        //     - the list of non-backedges (nonBackPreds)
        for (int w = 0; w < size; w++)
        {
            _header![w] = 0;
            _type![w] = BasicBlockClass.BbNonHeader;

            BasicBlock? nodeW = _nodes![w].Bb;
            if (nodeW == null)
            {
                _type[w] = BasicBlockClass.BbDead;
            }
            else
            {
                ProcessEdges(nodeW, w);
            }
        }
    }

    private void ProcessEdges(BasicBlock nodeW, int w)
    {
        if (nodeW.GetNumPred() > 0)
        {
            nodeW.InEdges.ForEach(nodeV =>
            {
                int v = _number.At(nodeV);
                if (v != Unvisited)
                {
                    if (IsAncestor(w, v))
                    {
                        _backPreds.At(w)!.Append(v);
                    }
                    else
                    {
                        _nonBackPreds.At(w)!.Add(v);
                    }
                }
            });
        }
    }

    // findLoops
    //
    // Find loops and build loop forest using Havlak's algorithm, which
    // is derived from Tarjan. Variable names and step numbering has
    // been chosen to be identical to the nomenclature in Havlak's
    // paper (which, in turn, is similar to the one used by Tarjan).
    public void FindLoops()
    {
        if (_cfg.GetStartBasicBlock() == null)
        {
            return;
        }

        int size = _cfg.GetNumNodes();

        _nonBackPreds.RemoveAll();
        _backPreds.RemoveAll();
        _number.RemoveAll();
        if (size > _maxSize)
        {
            _header = new int[size];
            _type = new BasicBlockClass[size];
            _last = new int[size];
            _nodes = new UnionFindNode[size];
            _maxSize = size;
        }

        for (int i = 0; i < size; ++i)
        {
            _nonBackPreds.Append(new Set<int>());
            _backPreds.Append(new Vector<int>());
            _nodes![i] = new UnionFindNode();
        }

        InitAllNodes();
        IdentifyEdges(size);

        // Start node is root of all other loops.
        _header![0] = 0;

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
        for (int w = size - 1; w >= 0; w--)
        {
            // this is 'P' in Havlak's paper
            Vector<UnionFindNode> nodePool = new();

            BasicBlock? nodeW = _nodes![w].Bb;
            if (nodeW != null)
            {
                StepD(w, nodePool);

                // Copy nodePool to workList.
                Vector<UnionFindNode> workList = new();
                nodePool.ForEach(workList.Append);

                if (nodePool.Size() != 0)
                {
                    _type![w] = BasicBlockClass.BbReducible;
                }

                // work the list...
                while (!workList.IsEmpty())
                {
                    UnionFindNode x = workList.RemoveFirst()!;

                    // Step e:
                    //
                    // Step e represents the main difference from Tarjan's method.
                    // Chasing upwards from the sources of a node w's backedges. If
                    // there is a node y' that is not a descendant of w, w is marked
                    // the header of an irreducible loop, there is another entry
                    // into this loop that avoids w.

                    // The algorithm has degenerated. Break and
                    // return in this case.
                    int nonBackSize = _nonBackPreds.At(x.DfsNumber)!.Size();
                    if (nonBackSize > MaxNonBackPreds)
                    {
                        return;
                    }

                    StepEProcessNonBackPreds(w, nodePool, workList, x);
                }

                // Collapse/Unionize nodes in a SCC to a single node
                // For every SCC found, create a loop descriptor and link it in.
                if (nodePool.Size() > 0 || _type![w] == BasicBlockClass.BbSelf)
                {
                    SimpleLoop loop = _lsg.CreateLoop(nodeW, _type![w] != BasicBlockClass.BbIrreducible);
                    SetLoopAttributes(w, nodePool, loop);
                }
            }
        } // Step c
    } // findLoops

    private void StepEProcessNonBackPreds(int w, Vector<UnionFindNode> nodePool,
        Vector<UnionFindNode> workList, UnionFindNode x)
    {
        _nonBackPreds.At(x.DfsNumber)!.ForEach(iter =>
        {
            UnionFindNode y = _nodes![iter];
            UnionFindNode ydash = y.FindSet();

            if (!IsAncestor(w, ydash.DfsNumber))
            {
                _type![w] = BasicBlockClass.BbIrreducible;
                _nonBackPreds.At(w)!.Add(ydash.DfsNumber);
            }
            else
            {
                if (ydash.DfsNumber != w)
                {
                    if (!nodePool.HasSome(e => e == ydash))
                    {
                        workList.Append(ydash);
                        nodePool.Append(ydash);
                    }
                }
            }
        });
    }

    private void SetLoopAttributes(int w, Vector<UnionFindNode> nodePool, SimpleLoop loop)
    {
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
        _nodes![w].Loop = loop;

        nodePool.ForEach(node =>
        {
            // Add nodes to loop descriptor.
            _header![node.DfsNumber] = w;
            node.Union(_nodes[w]);

            // Nested loops are not added, but linked together.
            if (node.Loop != null)
            {
                node.Loop.SetParent(loop);
            }
            else
            {
                loop.AddNode(node.Bb!);
            }
        });
    }

    private void StepD(int w, Vector<UnionFindNode> nodePool)
    {
        _backPreds.At(w)!.ForEach(v =>
        {
            if (v != w)
            {
                nodePool.Append(_nodes![v].FindSet());
            }
            else
            {
                _type![w] = BasicBlockClass.BbSelf;
            }
        });
    }
}
