using Benchmarks.Som;

namespace Benchmarks.DeltaBlueSolver;

internal sealed class Planner
{
    private int _currentMark = 1;

    /**
     * Attempt to satisfy the given constraint and, if successful,
     * incrementally update the dataflow graph. Details: If satisfying
     * the constraint is successful, it may override a weaker constraint
     * on its output. The algorithm attempts to resatisfy that
     * constraint using some other method. This process is repeated
     * until either a) it reaches a variable that was not previously
     * determined by any constraint or b) it reaches a constraint that
     * is too weak to be satisfied using any of its methods. The
     * variables of constraints that have been processed are marked with
     * a unique mark value so that we know where we've been. This allows
     * the algorithm to avoid getting into an infinite loop even if the
     * constraint graph has an inadvertent cycle.
     */
    public void IncrementalAdd(AbstractConstraint c)
    {
        int mark = NewMark();
        AbstractConstraint? overridden = c.Satisfy(mark, this);

        while (overridden != null)
        {
            overridden = overridden.Satisfy(mark, this);
        }
    }

    /**
    * Entry point for retracting a constraint. Remove the given
    * constraint and incrementally update the dataflow graph.
    * Details: Retracting the given constraint may allow some currently
    * unsatisfiable downstream constraint to be satisfied. We therefore collect
    * a list of unsatisfied downstream constraints and attempt to
    * satisfy each one in turn. This list is traversed by constraint
    * strength, strongest first, as a heuristic for avoiding
    * unnecessarily adding and then overriding weak constraints.
    * Assume: c is satisfied.
    */
    public void IncrementalRemove(AbstractConstraint c)
    {
        Variable? output = c.GetOutput();
        c.MarkUnsatisfied();
        c.RemoveFromGraph();

        Vector<AbstractConstraint> unsatisfied = RemovePropagateFrom(output!);
        unsatisfied.ForEach(IncrementalAdd);
    }

    /**
    * Extract a plan for resatisfaction starting from the outputs of
    * the given constraints, usually a set of input constraints.
    */
    private Plan ExtractPlanFromConstraints(Vector<AbstractConstraint> constraints)
    {
        Vector<AbstractConstraint> sources = new();

        constraints.ForEach(c =>
        {
            if (c.IsInput() && c.IsSatisfied())
            {
                sources.Append(c);
            }
        });

        return MakePlan(sources);
    }

    // Extract a plan for resatisfaction starting from the given source
    // constraints, usually a set of input constraints. This method
    // assumes that stay optimization is desired; the plan will contain
    // only constraints whose output variables are not stay. Constraints
    // that do no computation, such as stay and edit constraints, are
    // not included in the plan.
    // Details: The outputs of a constraint are marked when it is added
    // to the plan under construction. A constraint may be appended to
    // the plan when all its input variables are known. A variable is
    // known if either a) the variable is marked (indicating that has
    // been computed by a constraint appearing earlier in the plan), b)
    // the variable is 'stay' (i.e. it is a constant at plan execution
    // time), or c) the variable is not determined by any
    // constraint. The last provision is for past states of history
    // variables, which are not stay but which are also not computed by
    // any constraint.
    // Assume: sources are all satisfied.
    private Plan MakePlan(Vector<AbstractConstraint> sources)
    {
        int mark = NewMark();
        Plan plan = new();
        Vector<AbstractConstraint> todo = sources;

        while (!todo.IsEmpty())
        {
            AbstractConstraint? c = todo.RemoveFirst();

            if (c!.GetOutput()!.Mark != mark && c.InputsKnown(mark))
            {
                // not in plan already and eligible for inclusion
                plan.Append(c);
                c.GetOutput()!.Mark = mark;
                AddConstraintsConsumingTo(c.GetOutput()!, todo);
            }
        }

        return plan;
    }

    public void PropagateFrom(Variable v)
    {
        Vector<AbstractConstraint> todo = new();
        AddConstraintsConsumingTo(v, todo);

        while (!todo.IsEmpty())
        {
            AbstractConstraint? c = todo.RemoveFirst();
            c!.Execute();
            AddConstraintsConsumingTo(c.GetOutput()!, todo);
        }
    }

    private void AddConstraintsConsumingTo(Variable v, Vector<AbstractConstraint> coll)
    {
        AbstractConstraint? determiningC = v.DeterminedBy;

        v.Constraints.ForEach(c =>
        {
            if (c != determiningC && c.IsSatisfied())
            {
                coll.Append(c);
            }
        });
    }

    // Recompute the walkabout strengths and stay flags of all variables
    // downstream of the given constraint and recompute the actual
    // values of all variables whose stay flag is true. If a cycle is
    // detected, remove the given constraint and answer
    // false. Otherwise, answer true.
    // Details: Cycles are detected when a marked variable is
    // encountered downstream of the given constraint. The sender is
    // assumed to have marked the inputs of the given constraint with
    // the given mark. Thus, encountering a marked node downstream of
    // the output constraint means that there is a path from the
    // constraint's output to one of its inputs.
    //
    public bool AddPropagate(AbstractConstraint c, int mark)
    {
        Vector<AbstractConstraint> todo = Vector<AbstractConstraint>.With(c);

        while (!todo.IsEmpty())
        {
            AbstractConstraint? d = todo.RemoveFirst();

            if (d!.GetOutput()!.Mark == mark)
            {
                IncrementalRemove(c);
                return false;
            }

            d.Recalculate();
            AddConstraintsConsumingTo(d.GetOutput()!, todo);
        }

        return true;
    }

    private void Change(Variable var, int newValue)
    {
        EditConstraint editC = new(var, Strength.PreferredSym, this);

        Vector<AbstractConstraint> editV = Vector<AbstractConstraint>.With(editC);
        Plan plan = ExtractPlanFromConstraints(editV);
        for (int i = 0; i < 10; i++)
        {
            var.Value = newValue;
            plan.Execute();
        }

        editC.DestroyConstraint(this);
    }

    private void ConstraintsConsuming(Variable? v, Action<AbstractConstraint> fn)
    {
        AbstractConstraint? determiningC = v?.DeterminedBy;
        v?.Constraints.ForEach(c =>
        {
            if (c != determiningC && c.IsSatisfied())
            {
                fn.Invoke(c);
            }
        });
    }

    // Select a previously unused mark value.
    private int NewMark()
    {
        _currentMark++;
        return _currentMark;
    }

    // Update the walkabout strengths and stay flags of all variables
    // downstream of the given constraint. Answer a collection of
    // unsatisfied constraints sorted in order of decreasing strength.
    private Vector<AbstractConstraint> RemovePropagateFrom(Variable output)
    {
        Vector<AbstractConstraint> unsatisfied = new();

        output.DeterminedBy = null;
        output.WalkStrength = Strength.AbsoluteWeakest();
        output.Stay = true;

        Vector<Variable> todo = Vector<Variable>.With(output);

        while (!todo.IsEmpty())
        {
            Variable? v = todo.RemoveFirst();

            v?.Constraints.ForEach(c =>
            {
                if (!c.IsSatisfied())
                {
                    unsatisfied.Append(c);
                }
            });

            ConstraintsConsuming(v, c =>
            {
                c.Recalculate();
                todo.Append(c.GetOutput()!);
            });
        }

        unsatisfied.Sort((c1, c2) => c1.Strength.Stronger(c2.Strength) ? -1 : 1);
        return unsatisfied;
    }

    // This is the standard DeltaBlue benchmark. A long chain of
    // equality constraints is constructed with a stay constraint on
    // one end. An edit constraint is then added to the opposite end
    // and the time is measured for adding and removing this
    // constraint, and extracting and executing a constraint
    // satisfaction plan. There are two cases. In case 1, the added
    // constraint is stronger than the stay constraint and values must
    // propagate down the entire length of the chain. In case 2, the
    // added constraint is weaker than the stay constraint so it cannot
    // be accommodated. The cost in this case is, of course, very
    // low. Typical situations lie somewhere between these two
    // extremes.
    //
    public static void ChainTest(int n)
    {
        Planner planner = new();
        Variable[] vars = new Variable[n + 1];

        for (int i = 0; i < n + 1; i++)
        {
            vars[i] = new Variable();
        }

        // Build chain of n equality constraints
        for (int i = 0; i < n; i++)
        {
            Variable v1 = vars[i];
            Variable v2 = vars[i + 1];
            new EqualityConstraint(v1, v2, Strength.RequiredSym, planner);
        }

        new StayConstraint(vars[n], Strength.StrongDefaultSym, planner);
        AbstractConstraint editC = new EditConstraint(vars[0], Strength.PreferredSym, planner);

        Vector<AbstractConstraint> editV = Vector<AbstractConstraint>.With(editC);
        Plan plan = planner.ExtractPlanFromConstraints(editV);
        for (int i = 0; i < 100; i++)
        {
            vars[0].Value = i;
            plan.Execute();
            if (vars[n].Value != i)
            {
                throw new InvalidOperationException("Chain test failed!");
            }
        }

        editC.DestroyConstraint(planner);
    }

    // This test constructs two sets of variables related to each
    // other by a simple linear transformation (scale and offset). The
    // time is measured to change a variable on either side of the
    // mapping and to change the scale and offset factors.
    //
    public static void ProjectionTest(int n)
    {
        Planner planner = new();

        Vector<Variable> dests = new();

        Variable scale = Variable.WithValue(10);
        Variable offset = Variable.WithValue(1000);

        Variable? src = null;
        Variable? dst = null;
        for (int i = 1; i <= n; i++)
        {
            src = Variable.WithValue(i);
            dst = Variable.WithValue(i);
            dests.Append(dst);
            new StayConstraint(src, Strength.DefaultSym, planner);
            new ScaleConstraint(src, scale, offset, dst, Strength.RequiredSym, planner);
        }

        planner.Change(src!, 17);
        if (dst!.Value != 1170)
        {
            throw new InvalidOperationException("Projection test 1 failed!");
        }

        planner.Change(dst, 1050);
        if (src!.Value != 5)
        {
            throw new InvalidOperationException("Projection test 2 failed!");
        }

        planner.Change(scale, 5);
        for (int i = 0; i < n - 1; ++i)
        {
            if (dests.At(i)!.Value != ((i + 1) * 5) + 1000)
            {
                throw new InvalidOperationException("Projection test 3 failed!");
            }
        }

        planner.Change(offset, 2000);
        for (int i = 0; i < n - 1; ++i)
        {
            if (dests.At(i)!.Value != ((i + 1) * 5) + 2000)
            {
                throw new InvalidOperationException("Projection test 4 failed!");
            }
        }
    }
}