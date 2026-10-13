using Benchmarks.Som;

namespace Benchmarks.DeltaBlueSolver;

// A Plan is an ordered list of constraints to be executed in sequence
// to resatisfy all currently satisfiable constraints in the face of
// one or more changing inputs.
internal sealed class Plan : Vector<AbstractConstraint>
{
    public Plan() : base(15)
    {
    }

    public void Execute()
    {
        ForEach(c => c.Execute());
    }
}
