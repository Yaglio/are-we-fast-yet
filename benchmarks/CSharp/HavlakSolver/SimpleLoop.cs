using Benchmarks.Som;

namespace Benchmarks.HavlakSolver;

internal sealed class SimpleLoop
{
    private readonly IdentitySet<BasicBlock> _basicBlocks;
    private int _nestingLevel;

    private readonly BasicBlock? _header;

    private readonly bool _isReducible;

    public SimpleLoop? Parent { get; private set; }
    public IdentitySet<SimpleLoop> Children { get; }
    public int Counter { get; set; }
    public int DepthLevel { get; set; }
    public bool IsRoot { get; set; }

    public SimpleLoop(BasicBlock? bb, bool isReducible)
    {
        _isReducible = isReducible;
        Parent = null;
        IsRoot = false;
        _nestingLevel = 0;
        DepthLevel = 0;
        _basicBlocks = new IdentitySet<BasicBlock>();
        Children = new IdentitySet<SimpleLoop>();

        if (bb != null)
        {
            _basicBlocks.Add(bb);
        }

        _header = bb;
    }

    public void AddNode(BasicBlock bb) => _basicBlocks.Add(bb);

    private void AddChildLoop(SimpleLoop loop) => Children.Add(loop);

    public void SetParent(SimpleLoop? p)
    {
        Parent = p;
        Parent!.AddChildLoop(this);
    }

    public int NestingLevel
    {
        get => _nestingLevel;
        set
        {
            _nestingLevel = value;
            if (value == 0)
            {
                IsRoot = true;
            }
        }
    }
}
