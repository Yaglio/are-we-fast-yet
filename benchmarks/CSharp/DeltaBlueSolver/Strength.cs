using Benchmarks.Som;

namespace Benchmarks.DeltaBlueSolver;

/*
 * Strengths are used to measure the relative importance of constraints. New
 * strengths may be inserted in the strength hierarchy without disrupting
 * current constraints. Strengths cannot be created outside this class, so
 * pointer comparison can be used for value comparison.
 */
internal sealed class Strength
{
    internal sealed class Sym : ICustomHash
    {
        private readonly int _hash;

        public Sym(int hash)
        {
            _hash = hash;
        }

        public int CustomHash()
        {
            return _hash;
        }
    }

    public static readonly Sym AbsoluteStrongestSym = new(0);
    public static readonly Sym RequiredSym = new(1);
    public static readonly Sym StrongPreferredSym = new(2);
    public static readonly Sym PreferredSym = new(3);
    public static readonly Sym StrongDefaultSym = new(4);
    public static readonly Sym DefaultSym = new(5);
    public static readonly Sym WeakDefaultSym = new(6);
    public static readonly Sym AbsoluteWeakestSym = new(7);

    public int ArithmeticValue { get; }

    private readonly Sym _symbolicValue;

    private Strength(Sym symbolicValue)
    {
        _symbolicValue = symbolicValue;
        ArithmeticValue = _strengthTable.At(symbolicValue);
    }

    public bool SameAs(Strength s) => ArithmeticValue == s.ArithmeticValue;

    public bool Stronger(Strength s) => ArithmeticValue < s.ArithmeticValue;

    public bool Weaker(Strength s) => ArithmeticValue > s.ArithmeticValue;

    public Strength Strongest(Strength s) => s.Stronger(this) ? s : this;

    public Strength Weakest(Strength s) => s.Weaker(this) ? s : this;

    public static Strength Of(Sym strength) => _strengthConstant.At(strength)!;

    private static IdentityDictionary<Sym, int> CreateStrengthTable()
    {
        IdentityDictionary<Sym, int> strengthTable = new();
        strengthTable.AtPut(AbsoluteStrongestSym, -10000);
        strengthTable.AtPut(RequiredSym, -800);
        strengthTable.AtPut(StrongPreferredSym, -600);
        strengthTable.AtPut(PreferredSym, -400);
        strengthTable.AtPut(StrongDefaultSym, -200);
        strengthTable.AtPut(DefaultSym, 0);
        strengthTable.AtPut(WeakDefaultSym, 500);
        strengthTable.AtPut(AbsoluteWeakestSym, 10000);
        return strengthTable;
    }

    private static IdentityDictionary<Sym, Strength> CreateStrengthConstants()
    {
        IdentityDictionary<Sym, Strength> strengthConstant = new();
        _strengthTable.GetKeys().ForEach(key =>
            strengthConstant.AtPut(key, new Strength(key))
        );
        return strengthConstant;
    }

    public static Strength AbsoluteWeakest() => _absoluteWeakest;

    public static Strength Required() => _required;

    private static readonly Strength _absoluteWeakest;
    private static readonly Strength _required;

    private static readonly IdentityDictionary<Sym, int> _strengthTable;
    private static readonly IdentityDictionary<Sym, Strength> _strengthConstant;

    static Strength()
    {
        _strengthTable = CreateStrengthTable();
        _strengthConstant = CreateStrengthConstants();

        _absoluteWeakest = Of(AbsoluteWeakestSym);
        _required = Of(RequiredSym);
    }
}
