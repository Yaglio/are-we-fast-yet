namespace Benchmarks.JsonSolver;

internal sealed class JsonNumber : JsonValue
{
    private readonly string _str;

    public JsonNumber(string str)
    {
        _str = str;
        if (str == null)
        {
            throw new ArgumentNullException(nameof(str), "String must not be null");
        }
    }

    public override string ToString() => _str;

    public override bool IsNumber() => true;
}
