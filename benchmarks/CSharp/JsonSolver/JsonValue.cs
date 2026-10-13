namespace Benchmarks.JsonSolver;

/// <summary>
/// Represents a JSON value. This can be a JSON <see cref="JsonObject"/> (object), an <see cref="JsonArray"/> (array),
/// a <b>number</b>, a <b>string</b>, or one of the literals <c>true</c>, <c>false</c>, and <c>null</c>.
/// </summary>
/// <remarks>
/// The literals <c>true</c>, <c>false</c>, and <c>null</c> are represented by the constants <c>TRUE</c>, <c>FALSE</c>, and <c>NULL</c>.
/// JSON objects and arrays are represented by the subtypes <see cref="JsonObject"/> and <see cref="JsonArray"/>. Instances of these types can be created using the
/// public constructors of these classes.
/// 
/// Instances that represent JSON <b>numbers</b>, <b>strings</b> and <b>bool</b> values can be created using the static factory methods
/// <c>valueOf(string)</c>, <c>valueOf(long)</c>, <c>valueOf(double)</c>, etc.
/// 
/// In order to find out whether an instance of this class is of a certain type, the methods
/// <see cref="IsObject()"/>, <see cref="IsArray()"/>, <see cref="IsString()"/>, <see cref="IsNumber()"/> etc. can be
/// used.
/// 
/// If the type of a JSON value is known, the methods <see cref="AsObject()"/>, <see cref="AsArray()"/>,
/// <c>AsString()</c>, <c>AsInt()</c>, etc. can be used to get this value directly in the
/// appropriate target type.
/// 
/// This class is <b>not supposed to be extended</b> by clients.
/// </remarks>
internal abstract class JsonValue
{
    private protected JsonValue()
    {
        // prevent subclasses outside of this package
    }

    /// <summary>
    /// Detects whether this value represents a JSON object. If this is the case, this value is an
    /// instance of <see cref="JsonObject"/>.
    /// </summary>
    /// <returns><c>true</c> if this value is an instance of JsonObject</returns>
    public virtual bool IsObject() => false;

    /// <summary>
    /// Detects whether this value represents a JSON array. If this is the case, this value is an
    /// instance of <see cref="JsonArray"/>.
    /// </summary>
    /// <returns><c>true</c> if this value is an instance of JsonArray</returns>
    public virtual bool IsArray() => false;

    /// <summary>
    /// Detects whether this value represents a JSON number.
    /// </summary>
    /// <returns><c>true</c> if this value represents a JSON number</returns>
    public virtual bool IsNumber() => false;

    /// <summary>
    /// Detects whether this value represents a JSON string.
    /// </summary>
    /// <returns><c>true</c> if this value represents a JSON string</returns>
    public virtual bool IsString() => false;

    /// <summary>
    /// Detects whether this value represents a bool value.
    /// </summary>
    /// <returns><c>true</c> if this value represents either the JSON literal <c>true</c> or <c>false</c></returns>
    public virtual bool IsBool() => false;

    /// <summary>
    /// Detects whether this value represents the JSON literal <c>true</c>.
    /// </summary>
    /// <returns><c>true</c> if this value represents the JSON literal <c>true</c></returns>
    public virtual bool IsTrue() => false;

    /// <summary>
    /// Detects whether this value represents the JSON literal <c>false</c>.
    /// </summary>
    /// <returns><c>true</c> if this value represents the JSON literal <c>false</c></returns>
    public virtual bool IsFalse() => false;

    /// <summary>
    /// Detects whether this value represents the JSON literal <c>null</c>.
    /// </summary>
    /// <returns><c>true</c> if this value represents the JSON literal <c>null</c></returns>
    public virtual bool IsNull() => false;

    /// <summary>
    /// Returns this JSON value as <see cref="JsonObject"/>, assuming that this value represents a JSON
    /// object. If this is not the case, an exception is thrown.
    /// </summary>
    /// <returns>a <see cref="JsonObject"/> for this value</returns>
    /// <exception cref="InvalidOperationException">if this value is not a JSON object</exception>
    public virtual JsonObject AsObject() => throw new InvalidOperationException($"Not an object: {ToString()}");

    /// <summary>
    /// Returns this JSON value as <see cref="JsonArray"/>, assuming that this value represents a JSON array.
    /// If this is not the case, an exception is thrown.
    /// </summary>
    /// <returns>a <see cref="JsonArray"/> for this value</returns>
    /// <exception cref="InvalidOperationException">if this value is not a JSON array</exception>
    public virtual JsonArray AsArray() => throw new InvalidOperationException($"Not an array: {ToString()}");
}
