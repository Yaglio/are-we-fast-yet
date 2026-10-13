namespace Benchmarks.DeltaBlueSolver;

// I mark variables that should, with some level of preference, stay
// the same. I have one method with zero inputs and one output, which
// does nothing. Planners may exploit the fact that, if I am
// satisfied, my output will not change during plan execution. This is
// called "stay optimization".
//
internal sealed class StayConstraint : UnaryConstraint
{
    // Install a stay constraint with the given strength on the given variable.
    public StayConstraint(Variable v, Strength.Sym strength, Planner planner) : base(v, strength, planner)
    {
    }

    public override void Execute()
    {
    } // Stay constraints do nothing.
}
