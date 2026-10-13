namespace Benchmarks.CollisionDetection;

internal sealed class Aircraft
{
    public readonly CallSign CallSign;
    public readonly Vector3D Position;

    public Aircraft(CallSign callSign, Vector3D position)
    {
        CallSign = callSign;
        Position = position;
    }
}

