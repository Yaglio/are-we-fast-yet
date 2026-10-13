namespace Benchmarks.CollisionDetection;

internal sealed class Vector3D
{
    public readonly double X;
    public readonly double Y;
    public readonly double Z;

    public Vector3D(double x, double y, double z)
    {
        X = x;
        Y = y;
        Z = z;
    }

    public Vector3D Plus(Vector3D other) => new(X + other.X, Y + other.Y, Z + other.Z);

    public Vector3D Minus(Vector3D other) => new(X - other.X, Y - other.Y, Z - other.Z);

    public double Dot(Vector3D other) => (X * other.X) + (Y * other.Y) + (Z * other.Z);

    public double SquaredMagnitude() => Dot(this);

    public double Magnitude() => Math.Sqrt(SquaredMagnitude());

    public Vector3D Times(double amount) => new(X * amount, Y * amount, Z * amount);
}