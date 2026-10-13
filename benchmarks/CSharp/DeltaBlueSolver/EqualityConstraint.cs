namespace Benchmarks.DeltaBlueSolver;

// I constrain two variables to have the same value: "v1 = v2".
internal sealed class EqualityConstraint : BinaryConstraint
{
    // Install a constraint with the given strength equating the given
    // variables.
    public EqualityConstraint(Variable var1, Variable var2,
        Strength.Sym strength, Planner planner) : base(var1, var2, strength, planner)
    {
        AddConstraint(planner);
    }

    // Enforce this constraint. Assume that it is satisfied.
    public override void Execute()
    {
        if (Direction == DeltaBlueSolver.Direction.Forward)
        {
            V2!.Value = V1!.Value;
        }
        else
        {
            V1!.Value = V2!.Value;
        }
    }
}
