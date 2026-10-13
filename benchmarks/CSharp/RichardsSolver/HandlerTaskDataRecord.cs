namespace Benchmarks.RichardsSolver;

internal sealed class HandlerTaskDataRecord : RbObject
{
    public HandlerTaskDataRecord()
    {
        WorkIn = DeviceIn = NoWork;
    }

    public Packet? DeviceIn { get; set; }

    public void DeviceInAdd(Packet packet) => DeviceIn = Append(packet, DeviceIn);

    public Packet? WorkIn { get; set; }

    public void WorkInAdd(Packet packet) => WorkIn = Append(packet, WorkIn);
}
