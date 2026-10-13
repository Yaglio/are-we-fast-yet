using Benchmarks.Som;

namespace Benchmarks.CollisionDetection;

internal sealed class Simulator
{
    private readonly Vector<CallSign> _aircraft;

    public Simulator(int numAircraft)
    {
        _aircraft = new Vector<CallSign>();
        for (int i = 0; i < numAircraft; i++)
        {
            _aircraft.Append(new CallSign(i));
        }
    }

    public Vector<Aircraft> Simulate(double time)
    {
        Vector<Aircraft> frame = new();
        for (int i = 0; i < _aircraft.Size(); i += 2)
        {
            frame.Append(new Aircraft(_aircraft.At(i)!, new Vector3D(time, (Math.Cos(time) * 2) + (i * 3), 10)));
            frame.Append(new Aircraft(_aircraft.At(i + 1)!, new Vector3D(time, (Math.Sin(time) * 2) + (i * 3), 10)));
        }

        return frame;
    }
}
