namespace Benchmarks.RichardsSolver;

internal sealed class WorkerTaskDataRecord : RbObject
{
    public int Count { get; set; }

    public int Destination { get; set; } = HandlerA;
}
