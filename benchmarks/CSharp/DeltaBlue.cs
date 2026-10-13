using Benchmarks.DeltaBlueSolver;

namespace Benchmarks;

internal sealed class DeltaBlue : Benchmark
{
    public override bool InnerBenchmarkLoop(int innerIterations)
    {
        Planner.ChainTest(innerIterations);
        Planner.ProjectionTest(innerIterations);
        return true;
    }

    public override object Execute() => throw new InvalidOperationException("Execute should not be called for DeltaBlue benchmark.");

    public override bool VerifyResult(object result) => throw new InvalidOperationException("VerifyResult(object) should not be called for DeltaBlue benchmark.");
}

