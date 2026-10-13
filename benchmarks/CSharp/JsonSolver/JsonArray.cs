using Benchmarks.Som;

namespace Benchmarks.JsonSolver;

/// <summary>
/// Represents a JSON array, an ordered collection of JSON values.
/// </summary>
/// <remarks>
/// Elements can be added using the <c>Add(...)</c> methods which accept instances of
/// <see cref="JsonValue"/>, strings, primitive numbers, and bool values. To replace an element of an
/// array, use the <c>Set(int, ...)</c> methods.
///
/// Note that this class is <b>not thread-safe</b>. If multiple threads access a
/// <see cref="JsonArray"/> instance concurrently, while at least one of these threads modifies the
/// contents of this array, access to the instance must be synchronized externally. Failure to do so
/// may lead to an inconsistent state.
///
/// This class is <b>not supposed to be extended</b> by clients.
/// </remarks>
internal sealed class JsonArray : JsonValue
{
    private readonly Vector<JsonValue> _values;

    /// <summary>
    /// Creates a new empty JsonArray.
    /// </summary>
    public JsonArray()
    {
        _values = new Vector<JsonValue>();
    }

    /// <summary>
    /// Appends the specified JSON value to the end of this array.
    /// </summary>
    /// <param name="value">the JsonValue to add to the array, must not be <c>null</c></param>
    /// <returns>the array itself, to enable method chaining</returns>
    public JsonArray Add(JsonValue value)
    {
        if (value == null)
        {
            throw new ArgumentNullException(nameof(value), "JsonValue must not be null");
        }

        _values.Append(value);
        return this;
    }

    /// <summary>
    /// Returns the number of elements in this array.
    /// </summary>
    /// <returns>the number of elements in this array</returns>
    public int Size() => _values.Size();

    /// <summary>
    /// Returns the value of the element at the specified position in this array.
    /// </summary>
    /// <param name="index">the index of the array element to return</param>
    /// <returns>the value of the element at the specified position</returns>
    /// <exception cref="IndexOutOfRangeException">if the index is out of range, i.e. <c>index &lt; 0</c> or <c>index &gt;= size</c></exception>
    public JsonValue? Get(int index) => _values.At(index);

    public override bool IsArray() => true;

    public override JsonArray AsArray() => this;
}
