namespace Benchmarks.DeltaBlueSolver;

// I am an abstract superclass for constraints having two possible
// output variables.
internal abstract class BinaryConstraint : AbstractConstraint
{
    protected Variable? V1;
    protected Variable? V2; // possible output variables
    protected Direction? Direction; // one of the following...

    protected BinaryConstraint(Variable var1, Variable var2, Strength.Sym strength, Planner planner) : base(strength)
    {
        V1 = var1;
        V2 = var2;
        Direction = null;
    }

    // Answer true if this constraint is satisfied in the current solution.
    public override bool IsSatisfied() => Direction != null;

    // Add myself to the constraint graph.
    public override void AddToGraph()
    {
        V1!.AddConstraint(this);
        V2!.AddConstraint(this);
        Direction = null;
    }

    // Remove myself from the constraint graph.
    public override void RemoveFromGraph()
    {
        V1?.RemoveConstraint(this);

        V2?.RemoveConstraint(this);

        Direction = null;
    }

    // Decide if I can be satisfied and which way I should flow based on
    // the relative strength of the variables I relate, and record that
    // decision.
    //
    protected override Direction? ChooseMethod(int mark)
    {
        if (V1!.Mark == mark)
        {
            if (V2!.Mark != mark && Strength.Stronger(V2.WalkStrength))
            {
                Direction = DeltaBlueSolver.Direction.Forward;
                return Direction;
            }

            Direction = null;
            return Direction;
        }

        if (V2!.Mark == mark)
        {
            if (V1.Mark != mark && Strength.Stronger(V1.WalkStrength))
            {
                Direction = DeltaBlueSolver.Direction.Backward;
                return Direction;
            }

            Direction = null;
            return Direction;
        }

        // If we get here, neither variable is marked, so we have a choice.
        if (V1.WalkStrength.Weaker(V2.WalkStrength))
        {
            if (Strength.Stronger(V1.WalkStrength))
            {
                Direction = DeltaBlueSolver.Direction.Backward;
                return Direction;
            }

            Direction = null;
            return Direction;
        }

        if (Strength.Stronger(V2.WalkStrength))
        {
            Direction = DeltaBlueSolver.Direction.Forward;
            return Direction;
        }

        Direction = null;
        return Direction;
    }

    public override void InputsDo(Action<Variable> fn)
    {
        fn.Invoke(Direction == DeltaBlueSolver.Direction.Forward ? V1! : V2!);
    }

    public override bool InputsHasOne(Predicate<Variable> fn)
    {
        return fn.Invoke(Direction == DeltaBlueSolver.Direction.Forward ? V1! : V2!);
    }

    // Record the fact that I am unsatisfied.
    public override void MarkUnsatisfied() => Direction = null;


    // Answer my current output variable.
    public override Variable? GetOutput() => Direction == DeltaBlueSolver.Direction.Forward ? V2 : V1;

    // Calculate the walkabout strength, the stay flag, and, if it is
    // 'stay', the value for the current output of this
    // constraint. Assume this constraint is satisfied.
    //
    public override void Recalculate()
    {
        Variable? input;
        Variable? output;

        if (Direction == DeltaBlueSolver.Direction.Forward)
        {
            input = V1;
            output = V2;
        }
        else
        {
            input = V2;
            output = V1;
        }

        output!.WalkStrength = Strength.Weakest(input!.WalkStrength);
        output.Stay = input.Stay;
        if (output.Stay)
        {
            Execute();
        }
    }
}
