namespace Benchmarks.RichardsSolver;

internal sealed class DeviceTaskDataRecord : RbObject
{
    internal DeviceTaskDataRecord()
    {
        Pending = NoWork;
    }

    public Packet? Pending { get; set; }
}
