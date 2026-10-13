using Benchmarks.Som;

namespace Benchmarks.JsonSolver;

/// <summary>
/// Represents a JSON object, a set of name/value pairs, where the names are strings and the values
/// are JSON values.
/// </summary>
/// <remarks>
/// Members can be added using the <c>Add(string, ...)</c> methods which accept instances of
/// <see cref="JsonValue"/>, strings, primitive numbers, and bool values. To modify certain values of an
/// object, use the <c>Set(string, ...)</c> methods. Please note that the <c>Add</c>
/// methods are faster than <c>Set</c> as they do not search for existing members. On the other
/// hand, the <c>Add</c> methods do not prevent adding multiple members with the same name.
/// Duplicate names are discouraged but not prohibited by JSON.
///
/// Members can be accessed by their name using <see cref="Get(string)"/>.
///
/// This class is <b>not supposed to be extended</b> by clients.
/// </remarks>
internal sealed class JsonObject : JsonValue
{
    private readonly Vector<string> _names;
    private readonly Vector<JsonValue> _values;
    private readonly HashIndexTable _table;

    /// <summary>
    /// Creates a new empty JsonObject.
    /// </summary>
    public JsonObject()
    {
        _names = new Vector<string>();
        _values = new Vector<JsonValue>();
        _table = new HashIndexTable();
    }

    /// <summary>
    /// Appends a new member to the end of this object, with the specified name and the specified JSON
    /// value.
    /// </summary>
    /// <remarks>
    /// This method <b>does not prevent duplicate names</b>. Calling this method with a name
    /// that already exists in the object will append another member with the same name. In order to
    /// replace existing members, use the method <c>Set(name, value)</c> instead. However,
    /// <b><i>Add</i> is much faster than <i>Set</i></b> (because it does not need to
    /// search for existing members). Therefore <i>Add</i> should be preferred when constructing new
    /// objects.
    /// </remarks>
    /// <param name="name">the name of the member to add</param>
    /// <param name="value">the value of the member to add, must not be <c>null</c></param>
    /// <returns>the object itself, to enable method chaining</returns>
    public JsonObject Add(string name, JsonValue value)
    {
        if (name == null)
        {
            throw new ArgumentNullException(nameof(name), "name must not be null");
        }

        if (value == null)
        {
            throw new ArgumentNullException(nameof(value), "value must not be null");
        }

        _table.Add(name, _names.Size());
        _names.Append(name);
        _values.Append(value);
        return this;
    }

    /// <summary>
    /// Returns the value of the member with the specified name in this object. If this object contains
    /// multiple members with the given name, this method will return the last one.
    /// </summary>
    /// <param name="name">the name of the member whose value is to be returned</param>
    /// <returns>the value of the last member with the specified name, or <c>null</c> if this
    ///         object does not contain a member with that name</returns>
    public JsonValue? Get(string name)
    {
        if (name == null)
        {
            throw new ArgumentNullException(nameof(name), "name must not be null");
        }

        int index = IndexOf(name);
        return index == -1 ? null : _values.At(index);
    }

    /// <summary>
    /// Returns the number of members (name/value pairs) in this object.
    /// </summary>
    /// <returns>the number of members in this object</returns>
    public int Size() => _names.Size();

    /// <summary>
    /// Returns <c>true</c> if this object contains no members.
    /// </summary>
    /// <returns><c>true</c> if this object contains no members</returns>
    public bool IsEmpty() => _names.IsEmpty();

    public override bool IsObject() => true;

    public override JsonObject AsObject() => this;

    private int IndexOf(string name)
    {
        int index = _table.Get(name);
        if (index != -1 && name == _names.At(index))
        {
            return index;
        }

        throw new NotImplementedException(); // Not needed for benchmark
    }

    private sealed class HashIndexTable
    {
        private readonly int[] _hashTable = new int[32]; // must be a power of two

        public void Add(string name, int index)
        {
            int slot = HashSlotFor(name);
            if (index < 0xff)
            {
                // increment by 1, 0 stands for empty
                _hashTable[slot] = (index + 1) & 0xff;
            }
            else
            {
                _hashTable[slot] = 0;
            }
        }

        public int Get(string name)
        {
            int slot = HashSlotFor(name);
            // subtract 1, 0 stands for empty
            return (_hashTable[slot] & 0xff) - 1;
        }

        private int StringHash(string s)
        {
            // this is not a proper hash, but sufficient for the benchmark,
            // and very portable!
            return s.Length * 1402589;
        }

        private int HashSlotFor(string element)
        {
            return StringHash(element) & (_hashTable.Length - 1);
        }
    }
}
