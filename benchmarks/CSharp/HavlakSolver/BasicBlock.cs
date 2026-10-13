using Benchmarks.Som;

namespace Benchmarks.HavlakSolver;

internal sealed class BasicBlock : ICustomHash
{
    public Vector<BasicBlock> InEdges { get; }
    public Vector<BasicBlock> OutEdges { get; }

    private readonly int _name;

    public BasicBlock(int name)
    {
        _name = name;
        InEdges = new Vector<BasicBlock>(2);
        OutEdges = new Vector<BasicBlock>(2);
    }

    public int GetNumPred() => InEdges.Size();

    public void AddOutEdge(BasicBlock to) => OutEdges.Append(to);

    public void AddInEdge(BasicBlock from) => InEdges.Append(from);

    public int CustomHash() => _name;
}
