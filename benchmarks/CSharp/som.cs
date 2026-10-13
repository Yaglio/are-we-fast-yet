namespace Benchmarks;

internal static class CollectionConstants
{
  public const int InitialSize = 10;
}

internal class Vector<TE>
{
  private TE[]? _storage;
  private int _firstIdx;
  private int _lastIdx;

  public static Vector<TE> With(TE elem)
  {
    Vector<TE> v = new(1);
    v.Append(elem);
    return v;
  }

  public Vector(int size)
  {
    _storage = size == 0 ? null : new TE[size];
  }

  public Vector()
  {
  }

  public TE? At(int idx)
  {
    if (_storage == null || idx >= _storage.Length)
    {
      return default;
    }

    return _storage[idx];
  }

  public void AtPut(int idx, TE val)
  {
    if (_storage == null)
    {
      _storage = new TE[Math.Max(idx + 1, CollectionConstants.InitialSize)];
    }
    else if (idx >= _storage.Length)
    {
      int newLength = _storage.Length;
      while (newLength <= idx)
      {
        newLength *= 2;
      }

      Array.Resize(ref _storage, newLength);
    }

    _storage[idx] = val;
    if (_lastIdx < idx + 1)
    {
      _lastIdx = idx + 1;
    }
  }

  public void Append(TE elem)
  {
    if (_storage == null)
    {
      _storage = new TE[CollectionConstants.InitialSize];
    }
    else if (_lastIdx >= _storage.Length)
    {
      // Need to expand capacity first
      Array.Resize(ref _storage, 2 * _storage.Length);
    }

    _storage[_lastIdx] = elem;
    _lastIdx++;
  }

  public bool IsEmpty() => _lastIdx == _firstIdx;

  public void ForEach(Action<TE> fn)
  {
    for (int i = _firstIdx; i < _lastIdx; i++)
    {
      fn.Invoke(_storage![i]);
    }
  }

  public bool HasSome(Predicate<TE> fn)
  {
    for (int i = _firstIdx; i < _lastIdx; i++)
    {
      if (fn.Invoke(_storage![i]))
      {
        return true;
      }
    }

    return false;
  }

  public TE? GetOne(Predicate<TE> fn)
  {
    for (int i = _firstIdx; i < _lastIdx; i++)
    {
      TE e = _storage![i];
      if (fn.Invoke(e))
      {
        return e;
      }
    }

    return default;
  }

  public TE? First() => IsEmpty() ? default : _storage![_firstIdx];

  public TE? RemoveFirst()
  {
    if (IsEmpty())
    {
      return default;
    }

    _firstIdx++;
    return _storage![_firstIdx - 1];
  }

  public bool Remove(TE obj)
  {
    if (_storage == null || IsEmpty())
    {
      return false;
    }

    TE[] newArray = new TE[Capacity()];
    int newLast = 0;
    bool found = false;

    ForEach(it =>
    {
      if (EqualityComparer<TE>.Default.Equals(it, obj))
      {
        found = true;
      }
      else
      {
        newArray[newLast] = it;
        newLast++;
      }
    });

    _storage = newArray;
    _lastIdx = newLast;
    _firstIdx = 0;
    return found;
  }

  public void RemoveAll()
  {
    _firstIdx = 0;
    _lastIdx = 0;
    if (_storage != null)
    {
      _storage = new TE[_storage.Length];
    }
  }

  public int Size() => _lastIdx - _firstIdx;

  public int Capacity() => _storage?.Length ?? 0;

  public void Sort(Comparison<TE>? c)
  {
    if (Size() > 0)
    {
      Sort(_firstIdx, _lastIdx - 1, c);
    }
  }

  private void Sort(int i, int j, Comparison<TE>? c)
  {
    if (c == null)
    {
      DefaultSort(i, j);
    }

    int n = j + 1 - i;
    if (n <= 1)
    {
      return;
    }

    TE di = _storage![i];
    TE dj = _storage[j];

    if (c!.Invoke(di, dj) > 0)
    {
      Swap(_storage, i, j);
      TE temp = di;
      di = dj;
      dj = temp;
    }

    if (n > 2)
    {
      int ij = (i + j) / 2;
      TE dij = _storage[ij];

      if (c.Invoke(di, dij) <= 0)
      {
        if (c.Invoke(dij, dj) > 0)
        {
          Swap(_storage, j, ij);
          dij = dj;
        }
      }
      else
      {
        Swap(_storage, i, ij);
        dij = di;
      }

      if (n > 3)
      {
        int k = i;
        int l = j - 1;

        while (true)
        {
          while (k <= l && c.Invoke(dij, _storage[l]) <= 0)
          {
            l -= 1;
          }

          k += 1;
          while (k <= l && c.Invoke(_storage[k], dij) <= 0)
          {
            k += 1;
          }

          if (k > l)
          {
            break;
          }

          Swap(_storage, k, l);
        }

        Sort(i, l, c);
        Sort(k, j, c);
      }
    }
  }

  private static void Swap(TE[] storage2, int i, int j) => throw new NotImplementedException();

  private void DefaultSort(int i, int j) => throw new NotImplementedException();
}

internal class Set<TE>
{
  private readonly Vector<TE> _items;

  public Set() : this(CollectionConstants.InitialSize)
  {
  }

  public Set(int size)
  {
    _items = new Vector<TE>(size);
  }

  public int Size() => _items.Size();

  public void ForEach(Action<TE> fn) => _items.ForEach(fn);

  public bool HasSome(Predicate<TE> fn) => _items.HasSome(fn);

  public TE? GetOne(Predicate<TE> fn) => _items.GetOne(fn);

  public void Add(TE obj)
  {
    if (!Contains(obj))
    {
      _items.Append(obj);
    }
  }

  public Vector<T> Collect<T>(Func<TE, T> fn) where T : class
  {
    Vector<T> coll = new();

    ForEach(e => coll.Append(fn.Invoke(e)));
    return coll;
  }

  public virtual bool Contains(TE obj) => HasSome(e => e!.Equals(obj));

  public void RemoveAll() => _items.RemoveAll();
}

internal sealed class IdentitySet<TE> : Set<TE> where TE : class
{
  public IdentitySet()
  {
  }

  public IdentitySet(int size) : base(size)
  {
  }

  public override bool Contains(TE obj) => HasSome(e => e == obj);
}

internal interface ICustomHash
{
  int CustomHash();
}

internal class Dictionary<TK, TV> where TK : class, ICustomHash
{
  protected const int InitialCapacity = 16;

  private Entry?[] _buckets;
  private int _size;

  protected class Entry
  {
    public readonly int Hash;
    public readonly TK Key;
    public TV Value;
    public Entry? Next;

    public Entry(int hash, TK key, TV value, Entry? next)
    {
      Hash = hash;
      Key = key;
      Value = value;
      Next = next;
    }

    public virtual bool Match(int hash, TK key) => Hash == hash && key.Equals(Key);
  }

  public Dictionary(int size)
  {
    _buckets = new Entry[size];
  }

  public Dictionary() : this(InitialCapacity)
  {
  }

  private static int CalculateHash(TK? key)
  {
    if (key == null)
    {
      return 0;
    }

    int hash = key.CustomHash();
    return hash ^ (hash >>> 16);
  }

  public virtual int Size() => _size;

  public bool IsEmpty() => _size == 0;

  private int GetBucketIdx(int hash) => (_buckets.Length - 1) & hash;

  private Entry? GetBucket(int hash) => _buckets[GetBucketIdx(hash)];

  public TV? At(TK key)
  {
    int hash = CalculateHash(key);
    Entry? e = GetBucket(hash);

    while (e != null)
    {
      if (e.Match(hash, key))
      {
        return e.Value;
      }

      e = e.Next;
    }

    return default;
  }

  public bool ContainsKey(TK key)
  {
    int hash = CalculateHash(key);
    Entry? e = GetBucket(hash);

    while (e != null)
    {
      if (e.Match(hash, key))
      {
        return true;
      }

      e = e.Next;
    }

    return false;
  }

  public void AtPut(TK key, TV value)
  {
    int hash = CalculateHash(key);
    int i = GetBucketIdx(hash);

    Entry? current = _buckets[i];

    if (current == null)
    {
      _buckets[i] = NewEntry(key, value, hash);
      _size += 1;
    }
    else
    {
      InsertBucketEntry(key, value, hash, current);
    }

    if (_size > _buckets.Length)
    {
      Resize();
    }
  }

  protected virtual Entry NewEntry(TK key, TV value, int hash) => new(hash, key, value, null);

  private void InsertBucketEntry(TK key, TV value, int hash, Entry head)
  {
    Entry current = head;

    while (true)
    {
      if (current.Match(hash, key))
      {
        current.Value = value;
        return;
      }

      if (current.Next == null)
      {
        _size++;
        current.Next = NewEntry(key, value, hash);
        return;
      }

      current = current.Next;
    }
  }

  private void Resize()
  {
    Entry?[] oldStorage = _buckets;

    Entry?[] newStorage = new Entry?[oldStorage.Length * 2];
    _buckets = newStorage;
    TransferEntries(oldStorage);
  }

  private void TransferEntries(Entry?[] oldStorage)
  {
    for (int i = 0; i < oldStorage.Length; ++i)
    {
      Entry? current = oldStorage[i];
      if (current != null)
      {
        oldStorage[i] = null;

        if (current.Next == null)
        {
          _buckets[current.Hash & (_buckets.Length - 1)] = current;
        }
        else
        {
          SplitBucket(oldStorage, i, current);
        }
      }
    }
  }

  private void SplitBucket(Entry?[] oldStorage, int i, Entry? head)
  {
    Entry? loHead = null;
    Entry? loTail = null;
    Entry? hiHead = null;
    Entry? hiTail = null;
    Entry? current = head;

    while (current != null)
    {
      if ((current.Hash & oldStorage.Length) == 0)
      {
        if (loTail == null)
        {
          loHead = current;
        }
        else
        {
          loTail.Next = current;
        }

        loTail = current;
      }
      else
      {
        if (hiTail == null)
        {
          hiHead = current;
        }
        else
        {
          hiTail.Next = current;
        }

        hiTail = current;
      }

      current = current.Next;
    }

    if (loTail != null)
    {
      loTail.Next = null;
      _buckets[i] = loHead;
    }

    if (hiTail != null)
    {
      hiTail.Next = null;
      _buckets[i + oldStorage.Length] = hiHead;
    }
  }

  public void RemoveAll()
  {
    _buckets = new Entry[_buckets.Length];
    _size = 0;
  }

  public Vector<TK> GetKeys()
  {
    Vector<TK> keys = new(_size);
    foreach (Entry? c in _buckets)
    {
      Entry? current = c;
      while (current != null)
      {
        keys.Append(current.Key);
        current = current.Next;
      }
    }

    return keys;
  }

  public Vector<TV> GetValues()
  {
    Vector<TV> values = new(_size);
    foreach (Entry? c in _buckets)
    {
      Entry? current = c;
      while (current != null)
      {
        values.Append(current.Value);
        current = current.Next;
      }
    }

    return values;
  }
}

internal class IdentityDictionary<TK, TV> : Dictionary<TK, TV> where TK : class, ICustomHash
{
  private sealed class IdEntry : Entry
  {
    public IdEntry(int hash, TK key, TV value, Entry? next) : base(hash, key, value, next)
    {
    }

    public override bool Match(int hash, TK key) => Hash == hash && Key == key;
  }

  public IdentityDictionary(int size) : base(size)
  {
  }

  public IdentityDictionary() : base(InitialCapacity)
  {
  }

  protected override Entry NewEntry(TK key, TV value, int hash) => new IdEntry(hash, key, value, null);
}

internal class Random
{
  private int _seed = 74755;

  public int Next()
  {
    _seed = (_seed * 1309 + 13849) & 65535;
    return _seed;
  }
}

internal sealed class Pair<TK, TV>
{
  public TK Key { get; set; }
  private TV Value { get; set; }

  public Pair(TK key, TV value)
  {
    Key = key;
    Value = value;
  }
}
