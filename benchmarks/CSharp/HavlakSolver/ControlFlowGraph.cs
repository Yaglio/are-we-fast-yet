using Benchmarks.Som;

namespace Benchmarks.HavlakSolver;

internal sealed class ControlFlowGraph
{
    public Vector<BasicBlock> BasicBlocks { get; } = new();

    private BasicBlock? _startNode;
    private readonly Vector<BasicBlockEdge> _edgeList = new();

    public BasicBlock CreateNode(int name)
    {
        BasicBlock node;

        if (BasicBlocks.At(name) != null)
        {
            node = BasicBlocks.At(name)!;
        }
        else
        {
            node = new BasicBlock(name);
            BasicBlocks.AtPut(name, node);
        }

        if (GetNumNodes() == 1)
        {
            _startNode = node;
        }

        return node;
    }

    public void AddEdge(BasicBlockEdge edge) => _edgeList.Append(edge);

    public int GetNumNodes() => BasicBlocks.Size();

    public BasicBlock? GetStartBasicBlock() => _startNode;
}