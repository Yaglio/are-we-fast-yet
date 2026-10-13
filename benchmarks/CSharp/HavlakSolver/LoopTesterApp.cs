namespace Benchmarks.HavlakSolver;

internal sealed class LoopTesterApp
{
    private readonly ControlFlowGraph _cfg;
    private readonly LoopStructureGraph _lsg;

    public LoopTesterApp()
    {
        _cfg = new ControlFlowGraph();
        _lsg = new LoopStructureGraph();
        _cfg.CreateNode(0);
    }

    // Create 4 basic blocks, corresponding to and if/then/else clause
    // with a CFG that looks like a diamond
    private int BuildDiamond(int start)
    {
        int bb0 = start;
        new BasicBlockEdge(_cfg, bb0, bb0 + 1);
        new BasicBlockEdge(_cfg, bb0, bb0 + 2);
        new BasicBlockEdge(_cfg, bb0 + 1, bb0 + 3);
        new BasicBlockEdge(_cfg, bb0 + 2, bb0 + 3);

        return bb0 + 3;
    }

    // Connect two existing nodes
    private void BuildConnect(int start, int end) => new BasicBlockEdge(_cfg, start, end);

    // Form a straight connected sequence of n basic blocks
    private int BuildStraight(int start, int n)
    {
        for (int i = 0; i < n; i++)
        {
            BuildConnect(start + i, start + i + 1);
        }

        return start + n;
    }

    // Construct a simple loop with two diamonds in it
    private int BuildBaseLoop(int from)
    {
        int header = BuildStraight(from, 1);
        int diamond1 = BuildDiamond(header);
        int d11 = BuildStraight(diamond1, 1);
        int diamond2 = BuildDiamond(d11);
        int footer = BuildStraight(diamond2, 1);
        BuildConnect(diamond2, d11);
        BuildConnect(diamond1, header);

        BuildConnect(footer, from);
        footer = BuildStraight(footer, 1);
        return footer;
    }

    public int[] Main(int numDummyLoops, int findLoopIterations,
        int parLoops, int pparLoops, int ppparLoops)
    {
        ConstructSimpleCFG();
        AddDummyLoops(numDummyLoops);
        ConstructCFG(parLoops, pparLoops, ppparLoops);

        // Performing Loop Recognition, 1 Iteration, then findLoopIteration
        FindLoops(_lsg);
        for (int i = 0; i < findLoopIterations; i++)
        {
            FindLoops(new LoopStructureGraph());
        }

        _lsg.CalculateNestingLevel();
        return [_lsg.GetNumLoops(), _cfg.GetNumNodes()];
    }

    // ReSharper disable once InconsistentNaming
    private void ConstructCFG(int parLoops, int pparLoops, int ppparLoops)
    {
        int n = 2;

        for (int parlooptrees = 0; parlooptrees < parLoops; parlooptrees++)
        {
            _cfg.CreateNode(n + 1);
            BuildConnect(2, n + 1);
            n += 1;

            for (int i = 0; i < pparLoops; i++)
            {
                int top = n;
                n = BuildStraight(n, 1);
                for (int j = 0; j < ppparLoops; j++)
                {
                    n = BuildBaseLoop(n);
                }

                int bottom = BuildStraight(n, 1);
                BuildConnect(n, top);
                n = bottom;
            }

            BuildConnect(n, 1);
        }
    }

    private void AddDummyLoops(int numDummyLoops)
    {
        for (int dummyloop = 0; dummyloop < numDummyLoops; dummyloop++)
        {
            FindLoops(_lsg);
        }
    }

    private void FindLoops(LoopStructureGraph loopStructure)
    {
        HavlakLoopFinder finder = new(_cfg, loopStructure);
        finder.FindLoops();
    }

    // ReSharper disable once InconsistentNaming
    private void ConstructSimpleCFG()
    {
        _cfg.CreateNode(0);
        BuildBaseLoop(0);
        _cfg.CreateNode(1);
        new BasicBlockEdge(_cfg, 0, 2);
    }
}
