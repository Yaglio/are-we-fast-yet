namespace Benchmarks.DeltaBlueSolver;

// ------------------------ constraints ------------------------------------

// I am an abstract class representing a system-maintainable
// relationship (or "constraint") between a set of variables. I supply
// a strength instance variable; concrete subclasses provide a means
// of storing the constrained variables and other information required
// to represent a constraint.
internal abstract class AbstractConstraint
{
    public Strength Strength { get; } // the strength of this constraint

    protected AbstractConstraint(Strength.Sym strength)
    {
        Strength = Strength.Of(strength);
    }

    // Normal constraints are not input constraints. An input constraint
    // is one that depends on external state, such as the mouse, the
    // keyboard, a clock, or some arbitrary piece of imperative code.
    public virtual bool IsInput() => false;

    // Answer true if this constraint is satisfied in the current solution.
    public abstract bool IsSatisfied();

    // Activate this constraint and attempt to satisfy it.
    protected void AddConstraint(Planner planner)
    {
        AddToGraph();
        planner.IncrementalAdd(this);
    }

    // Add myself to the constraint graph.
    public abstract void AddToGraph();

    // Deactivate this constraint, remove it from the constraint graph,
    // possibly causing other constraints to be satisfied, and destroy
    // it.
    public void DestroyConstraint(Planner planner)
    {
        if (IsSatisfied())
        {
            planner.IncrementalRemove(this);
        }

        RemoveFromGraph();
    }

    // Remove myself from the constraint graph.
    public abstract void RemoveFromGraph();

    // Decide if I can be satisfied and record that decision. The output
    // of the chosen method must not have the given mark and must have
    // a walkabout strength less than that of this constraint.
    protected abstract Direction? ChooseMethod(int mark);

    // Enforce this constraint. Assume that it is satisfied.
    public abstract void Execute();

    public abstract void InputsDo(Action<Variable> fn);
    public abstract bool InputsHasOne(Predicate<Variable> fn);

    // Assume that I am satisfied. Answer true if all my current inputs
    // are known. A variable is known if either a) it is 'stay' (i.e. it
    // is a constant at plan execution time), b) it has the given mark
    // (indicating that it has been computed by a constraint appearing
    // earlier in the plan), or c) it is not determined by any
    // constraint.
    public bool InputsKnown(int mark)
    {
        return !InputsHasOne(v => !(v.Mark == mark || v.Stay || v.DeterminedBy == null));
    }

    // Record the fact that I am unsatisfied.
    public abstract void MarkUnsatisfied();

    // Answer my current output variable. Raise an error if I am not
    // currently satisfied.
    public abstract Variable? GetOutput();

    // Calculate the walkabout strength, the stay flag, and, if it is
    // 'stay', the value for the current output of this
    // constraint. Assume this constraint is satisfied.
    public abstract void Recalculate();

    // Attempt to find a way to enforce this constraint. If successful,
    // record the solution, perhaps modifying the current dataflow
    // graph. Answer the constraint that this constraint overrides, if
    // there is one, or nil, if there isn't.
    // Assume: I am not already satisfied.
    //
    public AbstractConstraint? Satisfy(int mark, Planner planner)
    {
        AbstractConstraint? overridden;

        ChooseMethod(mark);

        if (IsSatisfied())
        {
            // constraint can be satisfied
            // mark inputs to allow cycle detection in addPropagate
            InputsDo(input =>
            {
                input.Mark = mark;
            });

            Variable? output = GetOutput();
            overridden = output!.DeterminedBy;
            overridden?.MarkUnsatisfied();

            output.DeterminedBy = this;
            if (!planner.AddPropagate(this, mark))
            {
                throw new InvalidOperationException("Cycle encountered");
            }

            output.Mark = mark;
        }
        else
        {
            overridden = null;
            if (Strength.SameAs(Strength.Required()))
            {
                throw new InvalidOperationException("Could not satisfy a required constraint");
            }
        }

        return overridden;
    }
}
