namespace Benchmarks.RichardsSolver;

internal sealed class TaskControlBlock : TaskState
{
    private Packet? _input;
    private readonly Func<Packet?, RbObject?, TaskControlBlock?> _function;
    private readonly RbObject _handle;

    internal TaskControlBlock(TaskControlBlock aLink, int anIdentity, int aPriority, Packet? anInitialWorkQueue,
        TaskState anInitialState, Func<Packet?, RbObject?, TaskControlBlock?> aBlock, RbObject aPrivateData)
    {
        Link = aLink;
        Identity = anIdentity;
        Priority = aPriority;
        _input = anInitialWorkQueue;
        IsPacketPending = anInitialState.IsPacketPending;
        IsTaskWaiting = anInitialState.IsTaskWaiting;
        IsTaskHolding = anInitialState.IsTaskHolding;
        _function = aBlock;
        _handle = aPrivateData;
    }

    public int Identity { get; }
    public TaskControlBlock Link { get; }
    public int Priority { get; }

    public TaskControlBlock? AddInputAndCheckPriority(Packet packet, TaskControlBlock? oldTask)
    {
        if (NoWork == _input)
        {
            _input = packet;
            IsPacketPending = true;
            if (Priority > oldTask!.Priority)
            {
                return this;
            }
        }
        else
        {
            _input = Append(packet, _input);
        }

        return oldTask;
    }

    public TaskControlBlock? RunTask()
    {
        Packet? message;
        if (IsWaitingWithPacket)
        {
            message = _input;
            _input = message!.Link;
            if (NoWork == _input)
            {
                SetRunning();
            }
            else
            {
                SetPacketPending();
            }
        }
        else
        {
            message = NoWork;
        }

        return _function(message, _handle);
    }
}
