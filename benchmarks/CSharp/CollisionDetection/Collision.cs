namespace Benchmarks.CollisionDetection;

internal sealed class Collision
{
    public readonly CallSign AircraftA;
    public readonly CallSign AircraftB;
    public readonly Vector3D Position;

    public Collision(CallSign aircraftA, CallSign aircraftB, Vector3D position)
    {
        AircraftA = aircraftA;
        AircraftB = aircraftB;
        Position = position;
    }
}

