namespace Benchmarks.CollisionDetection;

internal class CallSign : IComparable<CallSign>
{
    private readonly int _value;

    public CallSign(int value)
    {
        _value = value;
    }

    public int CompareTo(CallSign? other)
    {
        if (other == null)
        {
            return 1;
        }

        return _value == other._value ? 0 : _value < other._value ? -1 : 1;
    }
}
