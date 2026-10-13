namespace Benchmarks.JsonSolver;

internal sealed class JsonString : JsonValue
{
    private readonly string _str;

    public JsonString(string str)
    {
        _str = str;
    }

    public override bool IsString() => true;
}
