using Benchmarks.Som;

namespace Benchmarks.HavlakSolver;

internal sealed class UnionFindNode
{
    private UnionFindNode? _parent;

    public BasicBlock? Bb { get; private set; }
    public SimpleLoop? Loop { get; set; }
    public int DfsNumber { get; private set; }

    // Initialize this node.
    public void InitNode(BasicBlock bb, int dfsNumber)
    {
        _parent = this;
        Bb = bb;
        DfsNumber = dfsNumber;
        Loop = null;
    }

    // Union/Find Algorithm - The find routine.
    //
    // Implemented with Path Compression (inner loops are only
    // visited and collapsed once, however, deep nests would still
    // result in significant traversals).
    public UnionFindNode FindSet()
    {
        Vector<UnionFindNode> nodeList = new();

        UnionFindNode node = this;
        while (node != node._parent)
        {
            if (node._parent != node._parent!._parent)
            {
                nodeList.Append(node);
            }

            node = node._parent;
        }

        // Path Compression, all nodes' parents point to the 1st level parent.
        nodeList.ForEach(iter => iter.Union(_parent!));
        return node;
    }

    // Union/Find Algorithm - The union routine.
    //
    // Trivial. Assigning parent pointer is enough,
    // we rely on path compression.
    public void Union(UnionFindNode basicBlock)
    {
        _parent = basicBlock;
    }
}
