namespace Benchmarks.HavlakSolver;

internal sealed class BasicBlockEdge
{
    private readonly BasicBlock _from;
    private readonly BasicBlock _to;

    public BasicBlockEdge(ControlFlowGraph cfg, int fromName, int toName)
    {
        _from = cfg.CreateNode(fromName);
        _to = cfg.CreateNode(toName);

        _from.AddOutEdge(_to);
        _to.AddInEdge(_from);

        cfg.AddEdge(this);
    }
}
