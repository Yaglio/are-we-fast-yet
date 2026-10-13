using Benchmarks.RichardsSolver;

namespace Benchmarks;

internal sealed class Richards : Benchmark
{
    public override object Execute() => new Scheduler().Start();

    public override bool VerifyResult(object result) => (bool)result;
}