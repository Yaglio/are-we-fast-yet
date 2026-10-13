using Benchmarks.HavlakSolver;

namespace Benchmarks;

internal class Havlak : Benchmark
{
    public override bool InnerBenchmarkLoop(int innerIterations) 
        => VerifyResult(new LoopTesterApp().Main(innerIterations, 50, 10, 10, 5), innerIterations);

    public bool VerifyResult(int[] r, int innerIterations)
    {
        if (innerIterations == 15000)
        {
            return r[0] == 46602 && r[1] == 5213;
        }

        if (innerIterations == 1500)
        {
            return r[0] == 6102 && r[1] == 5213;
        }

        if (innerIterations == 150)
        {
            return r[0] == 2052 && r[1] == 5213;
        }

        if (innerIterations == 15)
        {
            return r[0] == 1647 && r[1] == 5213;
        }

        if (innerIterations == 1)
        {
            return r[0] == 1605 && r[1] == 5213;
        }

        Console.WriteLine($"No verification result for {innerIterations} found");
        Console.WriteLine($"Result is: {r[0]}, {r[1]}");

        return false;
    }

    public override object Execute() => throw new InvalidOperationException("Execute should not be called for Havlak benchmark.");

    public override bool VerifyResult(object result) => throw new InvalidOperationException("VerifyResult(object) should not be called for Havlak benchmark.");
}