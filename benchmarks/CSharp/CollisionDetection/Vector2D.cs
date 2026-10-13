namespace Benchmarks.CollisionDetection;

internal sealed class Vector2D : IComparable<Vector2D>
{
    public readonly double X;
    public readonly double Y;

    public Vector2D(double x, double y)
    {
        X = x;
        Y = y;
    }

    public Vector2D Plus(Vector2D other) => new(X + other.X, Y + other.Y);

    public Vector2D Minus(Vector2D other) => new(X - other.X, Y - other.Y);

    public int CompareTo(Vector2D? other)
    {
        int result = CompareNumbers(X, other!.X);
        return result != 0 ? result : CompareNumbers(Y, other.Y);
    }

    private static int CompareNumbers(double a, double b)
    {
        // ReSharper disable once CompareOfFloatsByEqualityOperator
        if (a == b)
        {
            return 0;
        }

        if (a < b)
        {
            return -1;
        }

        if (a > b)
        {
            return 1;
        }

#pragma warning disable CS1718
        // We say that NaN is smaller than non-NaN.
        // ReSharper disable once CompareOfFloatsByEqualityOperator
        // ReSharper disable once EqualExpressionComparison
        if (a == a)
        {
            return 1;
        }
#pragma warning restore CS1718

        return -1;
    }
}