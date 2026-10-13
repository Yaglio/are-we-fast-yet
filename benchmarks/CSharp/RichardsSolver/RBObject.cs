namespace Benchmarks.RichardsSolver;

internal abstract class RbObject
{
    public Packet Append(Packet packet, Packet? queueHead)
    {
        packet.Link = NoWork;
        if (NoWork == queueHead)
        {
            return packet;
        }

        Packet mouse = queueHead;
        Packet? link;
        while (NoWork != (link = mouse.Link))
        {
            mouse = link;
        }

        mouse.Link = packet;
        return queueHead;
    }

    public const int Idler = 0;
    public const int Worker = 1;
    public const int HandlerA = 2;
    public const int HandlerB = 3;
    public const int DeviceA = 4;
    public const int DeviceB = 5;
    public const int NumTypes = 6;

    public const int DevicePacketKind = 0;
    public const int WorkPacketKind = 1;

    public const Packet? NoWork = null;

    public const TaskControlBlock? NoTask = null;
}
