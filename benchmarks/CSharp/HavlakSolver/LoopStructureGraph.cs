using Benchmarks.Som;

namespace Benchmarks.HavlakSolver;

internal sealed class LoopStructureGraph
{
    private readonly SimpleLoop _root;
    private readonly Vector<SimpleLoop> _loops;

    private int _loopCounter;

    public LoopStructureGraph()
    {
        _loopCounter = 0;
        _loops = new Vector<SimpleLoop>();
        _root = new SimpleLoop(null, true) { NestingLevel = 0, Counter = _loopCounter };
        _loopCounter += 1;

        _loops.Append(_root);
    }

    public SimpleLoop CreateLoop(BasicBlock bb, bool isReducible)
    {
        SimpleLoop loop = new(bb, isReducible) { Counter = _loopCounter };
        _loopCounter += 1;
        _loops.Append(loop);
        return loop;
    }

    public void CalculateNestingLevel()
    {
        // link up all 1st level loops to artificial root node.
        _loops.ForEach(liter =>
        {
            if (liter is { IsRoot: false, Parent: null })
            {
                liter.SetParent(_root);
            }
        });

        // recursively traverse the tree and assign levels.
        CalculateNestingLevelRec(_root, 0);
    }

    private void CalculateNestingLevelRec(SimpleLoop loop, int depth)
    {
        loop.DepthLevel = depth;
        loop.Children.ForEach(liter =>
        {
            CalculateNestingLevelRec(liter, depth + 1);

            loop.NestingLevel = Math.Max(loop.NestingLevel, 1 + liter.NestingLevel);
        });
    }

    public int GetNumLoops() => _loops.Size();
}
