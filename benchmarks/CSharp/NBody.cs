using Benchmarks.NBodySolver;

namespace Benchmarks;

internal sealed class NBody : Benchmark
{
    public override bool InnerBenchmarkLoop(int innerIterations)
    {
        NBodySystem system = new();
        for (int i = 0; i < innerIterations; i++)
        {
            system.Advance(0.01);
        }

        return VerifyResult(system.Energy(), innerIterations);
    }

    private bool VerifyResult(double result, int innerIterations)
    {
        if (innerIterations == 250000)
        {
            // ReSharper disable once CompareOfFloatsByEqualityOperator
            return result == -0.1690859889909308;
        }

        if (innerIterations == 1)
        {
            // ReSharper disable once CompareOfFloatsByEqualityOperator
            return result == -0.16907495402506745;
        }

        Console.WriteLine($"No verification result for {innerIterations} found");
        Console.WriteLine($"Result is: {result}");
        return false;
    }

    public override object Execute() => throw new InvalidOperationException("Execute should not be called for NBody benchmark.");

    public override bool VerifyResult(object result) => throw new InvalidOperationException("VerifyResult(object) should not be called for NBody benchmark.");
}

