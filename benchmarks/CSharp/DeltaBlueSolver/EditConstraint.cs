namespace Benchmarks.DeltaBlueSolver;

// I am a unary input constraint used to mark a variable that the
// client wishes to change.
internal sealed class EditConstraint : UnaryConstraint
{
    public EditConstraint(Variable v, Strength.Sym strength, Planner planner) : base(v, strength, planner)
    {
    }

    // I indicate that a variable is to be changed by imperative code.
    public override bool IsInput() => true;

    public override void Execute()
    {
    } // Edit constraints do nothing.
}
