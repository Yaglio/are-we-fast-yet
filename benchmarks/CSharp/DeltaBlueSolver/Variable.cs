using Benchmarks.Som;

namespace Benchmarks.DeltaBlueSolver;

// ------------------------------ variables ------------------------------

// I represent a constrained variable. In addition to my value, I
// maintain the structure of the constraint graph, the current
// dataflow graph, and various parameters of interest to the DeltaBlue
// incremental constraint solver.
internal sealed class Variable
{
    public int Value { get; set; } // my value; changed by constraints
    public Vector<AbstractConstraint> Constraints { get; } // normal constraints that reference me

    public AbstractConstraint? DeterminedBy { get; set; } // the constraint that currently determines

    // my value (or null if there isn't one)
    public int Mark { get; set; } // used by the planner to mark constraints
    public Strength WalkStrength { get; set; } // my walkabout strength
    public bool Stay { get; set; } // true if I am a planning-time constant

    public static Variable WithValue(int aValue) => new() { Value = aValue };

    public Variable()
    {
        Constraints = new Vector<AbstractConstraint>(2);
        WalkStrength = Strength.AbsoluteWeakest();
        Stay = true;
    }

    // Add the given constraint to the set of all constraints that refer to me.
    public void AddConstraint(AbstractConstraint c)
    {
        Constraints.Append(c);
    }

    // Remove all traces of c from this variable.
    public void RemoveConstraint(AbstractConstraint c)
    {
        Constraints.Remove(c);
        if (DeterminedBy == c)
        {
            DeterminedBy = null;
        }
    }
}
