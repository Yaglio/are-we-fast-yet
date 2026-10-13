namespace Benchmarks.CollisionDetection;

internal sealed class RedBlackTree<TK, TV> where TK : IComparable<TK>
{
    private Node? _root;

    private enum Color
    {
        Red,
        Black
    }

    private static Node TreeMinimum(Node x)
    {
        Node current = x;
        while (current.Left != null)
        {
            current = current.Left;
        }

        return current;
    }

    private sealed class Node
    {
        public readonly TK Key;
        public TV Value;
        public Node? Left;
        public Node? Right;
        public Node? Parent;
        public Color Color;

        public Node(TK key, TV value)
        {
            Key = key;
            Value = value;
            Left = null;
            Right = null;
            Parent = null;
            Color = Color.Red;
        }

        public Node? Successor()
        {
            Node x = this;
            if (x.Right != null)
            {
                return TreeMinimum(x.Right);
            }

            Node? y = x.Parent;
            while (y != null && x == y.Right)
            {
                x = y;
                y = y.Parent;
            }

            return y;
        }
    }

    public TV? Put(TK key, TV value)
    {
        InsertResult insertionResult = TreeInsert(key, value);
        if (!insertionResult.IsNewEntry)
        {
            return insertionResult.OldValue;
        }

        Node? x = insertionResult.NewNode;

        while (x != _root && x!.Parent!.Color == Color.Red)
        {
            if (x.Parent == x.Parent.Parent!.Left)
            {
                Node? y = x.Parent.Parent.Right;
                if (y is { Color: Color.Red })
                {
                    // Case 1
                    x.Parent.Color = Color.Black;
                    y.Color = Color.Black;
                    x.Parent.Parent.Color = Color.Red;
                    x = x.Parent.Parent;
                }
                else
                {
                    if (x == x.Parent.Right)
                    {
                        // Case 2
                        x = x.Parent;
                        LeftRotate(x);
                    }

                    // Case 3
                    x.Parent!.Color = Color.Black;
                    x.Parent.Parent!.Color = Color.Red;
                    RightRotate(x.Parent.Parent);
                }
            }
            else
            {
                // Same as "then" clause with "right" and "left" exchanged.
                Node? y = x.Parent.Parent.Left;
                if (y is { Color: Color.Red })
                {
                    // Case 1
                    x.Parent.Color = Color.Black;
                    y.Color = Color.Black;
                    x.Parent.Parent.Color = Color.Red;
                    x = x.Parent.Parent;
                }
                else
                {
                    if (x == x.Parent.Left)
                    {
                        // Case 2
                        x = x.Parent;
                        RightRotate(x);
                    }

                    // Case 3
                    x.Parent!.Color = Color.Black;
                    x.Parent.Parent!.Color = Color.Red;
                    LeftRotate(x.Parent.Parent);
                }
            }
        }

        _root!.Color = Color.Black;
        return default;
    }

    public TV? Remove(TK key)
    {
        Node? z = FindNode(key);
        if (z == null)
        {
            return default;
        }

        // Y is the node to be unlinked from the tree.
        Node? y;
        if (z.Left == null || z.Right == null)
        {
            y = z;
        }
        else
        {
            y = z.Successor();
        }

        // Y is guaranteed to be non-null at this point.
        Node? x = y!.Left != null ? y.Left : y.Right;

        // X is the child of y which might potentially replace y in the tree. X might be null at
        // this point.
        Node? xParent;
        if (x != null)
        {
            x.Parent = y.Parent;
            xParent = x.Parent;
        }
        else
        {
            xParent = y.Parent;
        }

        if (y.Parent == null)
        {
            _root = x;
        }
        else
        {
            if (y == y.Parent.Left)
            {
                y.Parent.Left = x;
            }
            else
            {
                y.Parent.Right = x;
            }
        }

        if (y != z)
        {
            if (y.Color == Color.Black)
            {
                RemoveFixup(x, xParent);
            }

            y.Parent = z.Parent;
            y.Color = z.Color;
            y.Left = z.Left;
            y.Right = z.Right;

            z.Left?.Parent = y;

            z.Right?.Parent = y;

            if (z.Parent != null)
            {
                if (z.Parent.Left == z)
                {
                    z.Parent.Left = y;
                }
                else
                {
                    z.Parent.Right = y;
                }
            }
            else
            {
                _root = y;
            }
        }
        else if (y.Color == Color.Black)
        {
            RemoveFixup(x, xParent);
        }

        return z.Value;
    }

    public TV? Get(TK key)
    {
        Node? node = FindNode(key);

        return node == null ? default : node.Value;
    }

    internal sealed class Entry
    {
        public readonly TK Key;
        public readonly TV Value;

        public Entry(TK key, TV value)
        {
            Key = key;
            Value = value;
        }
    }

    public void ForEach(Action<Entry> fn)
    {
        if (_root == null)
        {
            return;
        }

        Node? current = TreeMinimum(_root);
        while (current != null)
        {
            fn.Invoke(new Entry(current.Key, current.Value));
            current = current.Successor();
        }
    }

    private Node? FindNode(TK key)
    {
        Node? current = _root;
        while (current != null)
        {
            int comparisonResult = key.CompareTo(current.Key);
            if (comparisonResult == 0)
            {
                return current;
            }

            current = comparisonResult < 0 ? current.Left : current.Right;
        }

        return null;
    }

    private sealed class InsertResult
    {
        public readonly bool IsNewEntry;
        public readonly Node? NewNode;
        public readonly TV? OldValue;

        public InsertResult(bool isNewEntry, Node? newNode, TV? oldValue)
        {
            IsNewEntry = isNewEntry;
            NewNode = newNode;
            OldValue = oldValue;
        }
    }

    private InsertResult TreeInsert(TK key, TV value)
    {
        Node? y = null;
        Node? x = _root;

        while (x != null)
        {
            y = x;
            int comparisonResult = key.CompareTo(x.Key);
            if (comparisonResult < 0)
            {
                x = x.Left;
            }
            else if (comparisonResult > 0)
            {
                x = x.Right;
            }
            else
            {
                TV oldValue = x.Value;
                x.Value = value;
                return new InsertResult(false, null, oldValue);
            }
        }

        Node z = new(key, value) { Parent = y };


        if (y == null)
        {
            _root = z;
        }
        else
        {
            if (key.CompareTo(y.Key) < 0)
            {
                y.Left = z;
            }
            else
            {
                y.Right = z;
            }
        }

        return new InsertResult(true, z, default);
    }

    private Node LeftRotate(Node x)
    {
        Node? y = x.Right;

        // Turn y's left subtree into x's right subtree.
        x.Right = y!.Left;
        y.Left?.Parent = x;

        // Link x's parent to y.
        y.Parent = x.Parent;
        if (x.Parent == null)
        {
            _root = y;
        }
        else
        {
            if (x == x.Parent.Left)
            {
                x.Parent.Left = y;
            }
            else
            {
                x.Parent.Right = y;
            }
        }

        // Put x on y's left.
        y.Left = x;
        x.Parent = y;

        return y;
    }

    private Node RightRotate(Node y)
    {
        Node? x = y.Left;

        // Turn x's right subtree into y's left subtree.
        y.Left = x!.Right;
        x.Right?.Parent = y;

        // Link y's parent to x;
        x.Parent = y.Parent;
        if (y.Parent == null)
        {
            _root = x;
        }
        else
        {
            if (y == y.Parent.Left)
            {
                y.Parent.Left = x;
            }
            else
            {
                y.Parent.Right = x;
            }
        }

        x.Right = y;
        y.Parent = x;

        return x;
    }

    private void RemoveFixup(Node? x, Node? xParent)
    {
        while (x != _root && (x == null || x.Color == Color.Black))
        {
            if (x == xParent!.Left)
            {
                // Note: the text points out that w cannot be null. The reason is not obvious from
                // simply looking at the code; it comes about from the properties of the red-black
                // tree.
                Node? w = xParent.Right;
                if (w!.Color == Color.Red)
                {
                    // Case 1
                    w.Color = Color.Black;
                    xParent.Color = Color.Red;
                    LeftRotate(xParent);
                    w = xParent.Right;
                }

                if ((w!.Left == null || w.Left.Color == Color.Black)
                    && (w.Right == null || w.Right.Color == Color.Black))
                {
                    // Case 2
                    w.Color = Color.Red;
                    x = xParent;
                    xParent = x.Parent!;
                }
                else
                {
                    if (w.Right == null || w.Right.Color == Color.Black)
                    {
                        // Case 3
                        w.Left!.Color = Color.Black;
                        w.Color = Color.Red;
                        RightRotate(w);
                        w = xParent.Right;
                    }

                    // Case 4
                    w!.Color = xParent.Color;
                    xParent.Color = Color.Black;
                    w.Right?.Color = Color.Black;

                    LeftRotate(xParent);
                    x = _root;
                    xParent = x!.Parent;
                }
            }
            else
            {
                // Same as "then" clause with "right" and "left" exchanged.
                Node? w = xParent.Left;
                if (w!.Color == Color.Red)
                {
                    // Case 1
                    w.Color = Color.Black;
                    xParent.Color = Color.Red;
                    RightRotate(xParent);
                    w = xParent.Left;
                }

                if ((w!.Right == null || w.Right.Color == Color.Black)
                    && (w.Left == null || w.Left.Color == Color.Black))
                {
                    // Case 2
                    w.Color = Color.Red;
                    x = xParent;
                    xParent = x.Parent!;
                }
                else
                {
                    if (w.Left == null || w.Left.Color == Color.Black)
                    {
                        // Case 3
                        w.Right!.Color = Color.Black;
                        w.Color = Color.Red;
                        LeftRotate(w);
                        w = xParent.Left;
                    }

                    // Case 4
                    w!.Color = xParent.Color;
                    xParent.Color = Color.Black;
                    w.Left?.Color = Color.Black;

                    RightRotate(xParent);
                    x = _root;
                    xParent = x!.Parent;
                }
            }
        }

        x?.Color = Color.Black;
    }
}