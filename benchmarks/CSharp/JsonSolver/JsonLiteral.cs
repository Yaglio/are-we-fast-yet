namespace Benchmarks.JsonSolver;

internal sealed class JsonLiteral : JsonValue
{
    public static readonly JsonValue Null = new JsonLiteral("null");
    public static readonly JsonValue True = new JsonLiteral("true");
    public static readonly JsonValue False = new JsonLiteral("false");

    private readonly string _value;
    private readonly bool _isNull;
    private readonly bool _isTrue;
    private readonly bool _isFalse;

    public JsonLiteral(string value)
    {
        _value = value;
        _isNull = value == "null";
        _isTrue = value == "true";
        _isFalse = value == "false";
    }

    public override string ToString() => _value;

    public override bool IsNull() => _isNull;

    public override bool IsTrue() => _isTrue;

    public override bool IsFalse() => _isFalse;

    public override bool IsBool() => _isTrue || _isFalse;
}
