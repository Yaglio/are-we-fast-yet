namespace Benchmarks.RichardsSolver;

internal class TaskState : RbObject
{
    public bool IsPacketPending { get; set; }
    public bool IsTaskHolding { get; set; }
    public bool IsTaskWaiting { get; set; }

    public void SetPacketPending()
    {
        IsPacketPending = true;
        IsTaskWaiting = false;
        IsTaskHolding = false;
    }

    public void SetRunning()
    {
        IsPacketPending = IsTaskWaiting = IsTaskHolding = false;
    }

    public void SetWaiting()
    {
        IsPacketPending = IsTaskHolding = false;
        IsTaskWaiting = true;
    }

    public void SetWaitingWithPacket()
    {
        IsTaskHolding = false;
        IsTaskWaiting = IsPacketPending = true;
    }

    public bool IsTaskHoldingOrWaiting
        => IsTaskHolding || (!IsPacketPending && IsTaskWaiting);

    public bool IsWaitingWithPacket
        => IsPacketPending && IsTaskWaiting && !IsTaskHolding;

    public static TaskState CreateRunning()
    {
        TaskState t = new();
        t.SetRunning();
        return t;
    }

    public static TaskState CreateWaiting()
    {
        TaskState t = new();
        t.SetWaiting();
        return t;
    }

    public static TaskState CreateWaitingWithPacket()
    {
        TaskState t = new();
        t.SetWaitingWithPacket();
        return t;
    }
}
