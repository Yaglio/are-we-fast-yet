namespace Benchmarks.DeltaBlueSolver;

// I am an abstract superclass for constraints having a single
// possible output variable.
internal abstract class UnaryConstraint : AbstractConstraint
{
    protected readonly Variable? Output; // possible output variable
    protected bool Satisfied; // true if I am currently satisfied

    protected UnaryConstraint(Variable v, Strength.Sym strength, Planner planner) : base(strength)
    {
        Output = v;
        AddConstraint(planner);
    }

    // Answer true if this constraint is satisfied in the current solution.
    public override bool IsSatisfied() => Satisfied;

    // Add myself to the constraint graph.
    public override void AddToGraph()
    {
        Output!.AddConstraint(this);
        Satisfied = false;
    }

    // Remove myself from the constraint graph.
    public override void RemoveFromGraph()
    {
        Output!.RemoveConstraint(this);
        Satisfied = false;
    }

    // Decide if I can be satisfied and record that decision.
    protected override Direction? ChooseMethod(int mark)
    {
        Satisfied = Output!.Mark != mark && Strength.Stronger(Output.WalkStrength);
        return null;
    }

    public abstract override void Execute();

    public override void InputsDo(Action<Variable> fn)
    {
        // I have no input variables
    }

    public override bool InputsHasOne(Predicate<Variable> fn) => false;

    // Record the fact that I am unsatisfied.
    public override void MarkUnsatisfied() => Satisfied = false;

    // Answer my current output variable.
    public override Variable? GetOutput() => Output;

    // Calculate the walkabout strength, the stay flag, and, if it is
    // 'stay', the value for the current output of this
    // constraint. Assume this constraint is satisfied.
    public override void Recalculate()
    {
        Output!.WalkStrength = Strength;
        Output.Stay = !IsInput();
        if (Output.Stay)
        {
            Execute(); // stay optimization
        }
    }
}