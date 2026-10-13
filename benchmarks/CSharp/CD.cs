using Benchmarks.CollisionDetection;
using Benchmarks.Som;

namespace Benchmarks;

// ReSharper disable once InconsistentNaming
internal sealed class CD : Benchmark
{
    private int CalculateCollisions(int numAircraft)
    {
        const int NumFrames = 200;

        Simulator simulator = new(numAircraft);
        CollisionDetector detector = new();

        int actualCollisions = 0;

        for (int i = 0; i < NumFrames; i++)
        {
            double time = i / 10.0;
            Vector<Collision> collisions = detector.HandleNewFrame(simulator.Simulate(time));
            actualCollisions += collisions.Size();
        }

        return actualCollisions;
    }

    public override bool InnerBenchmarkLoop(int innerIterations) => VerifyResult(CalculateCollisions(innerIterations), innerIterations);

    public bool VerifyResult(int actualCollisions, int numAircraft)
    {
        if (numAircraft == 1000)
        {
            return actualCollisions == 14484;
        }

        if (numAircraft == 500)
        {
            return actualCollisions == 14484;
        }

        if (numAircraft == 250)
        {
            return actualCollisions == 10830;
        }

        if (numAircraft == 200)
        {
            return actualCollisions == 8655;
        }

        if (numAircraft == 100)
        {
            return actualCollisions == 4305;
        }

        if (numAircraft == 10)
        {
            return actualCollisions == 390;
        }

        if (numAircraft == 2)
        {
            return actualCollisions == 42;
        }

        Console.WriteLine($"No verification result for {numAircraft} found");
        Console.WriteLine($"Result is: {actualCollisions}");
        return false;
    }

    public override object Execute() => throw new InvalidOperationException("Execute should not be called for CD benchmark.");

    public override bool VerifyResult(object result) => throw new InvalidOperationException("VerifyResult(object) should not be called for CD benchmark.");
}