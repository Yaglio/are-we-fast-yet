namespace Benchmarks.DeltaBlueSolver;

// I relate two variables by the linear scaling relationship: "v2 =
// (v1 * scale) + offset". Either v1 or v2 may be changed to maintain
// this relationship but the scale factor and offset are considered
// read-only.
internal sealed class ScaleConstraint : BinaryConstraint
{
    private readonly Variable? _scale; // scale factor input variable
    private readonly Variable? _offset; // offset input variable

    public ScaleConstraint(Variable src, Variable scale,
        Variable offset, Variable dest, Strength.Sym strength,
        Planner planner) : base(src, dest, strength, planner)
    {
        _scale = scale;
        _offset = offset;
        AddConstraint(planner);
    }

    // Add myself to the constraint graph.
    public override void AddToGraph()
    {
        V1!.AddConstraint(this);
        V2!.AddConstraint(this);
        _scale!.AddConstraint(this);
        _offset!.AddConstraint(this);
        Direction = null;
    }

    // Remove myself from the constraint graph.
    public override void RemoveFromGraph()
    {
        V1?.RemoveConstraint(this);

        V2?.RemoveConstraint(this);

        _scale?.RemoveConstraint(this);

        _offset?.RemoveConstraint(this);

        Direction = null;
    }

    // Enforce this constraint. Assume that it is satisfied.
    public override void Execute()
    {
        if (Direction == DeltaBlueSolver.Direction.Forward)
        {
            V2!.Value = (V1!.Value * _scale!.Value) + _offset!.Value;
        }
        else
        {
            V1!.Value = (V2!.Value - _offset!.Value) / _scale!.Value;
        }
    }

    public override void InputsDo(Action<Variable> fn)
    {
        fn.Invoke(Direction == DeltaBlueSolver.Direction.Forward ? V1! : V2!);

        fn.Invoke(_scale!);
        fn.Invoke(_offset!);
    }

    // Calculate the walkabout strength, the stay flag, and, if it is
    // 'stay', the value for the current output of this
    // constraint. Assume this constraint is satisfied.
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
            output = V1;
            input = V2;
        }

        output!.WalkStrength = Strength.Weakest(input!.WalkStrength);
        output.Stay = input.Stay && _scale!.Stay && _offset!.Stay;
        if (output.Stay)
        {
            Execute(); // stay optimization
        }
    }
}
