namespace Benchmarks.RichardsSolver;

internal sealed class Scheduler : RbObject
{
    private TaskControlBlock? _taskList;
    private TaskControlBlock? _currentTask;
    private int _currentTaskIdentity;
    private readonly TaskControlBlock?[] _taskTable;

    private int _queuePacketCount;
    private int _holdCount;

    private int _layout;

    private const bool Tracing = false;

    public Scheduler()
    {
        // init tracing
        _layout = 0;

        // init scheduler
        _queuePacketCount = 0;
        _holdCount = 0;
        _taskTable = new TaskControlBlock[NumTypes];
        Array.Fill(_taskTable, NoTask);
        _taskList = NoTask;
    }


    internal void CreateDevice(int identity, int priority, Packet? workPacket, TaskState state)
    {
        DeviceTaskDataRecord data = new();

        CreateTask(identity, priority, workPacket, state, (workArg, wordArg) =>
            {
                DeviceTaskDataRecord? dataRecord = (DeviceTaskDataRecord?)wordArg;
                Packet? functionWork = workArg;
                if (NoWork == functionWork)
                {
                    if (NoWork == (functionWork = dataRecord!.Pending))
                    {
                        return MarkWaiting();
                    }

                    dataRecord.Pending = NoWork;
                    return QueuePacket(functionWork);
                }

                dataRecord!.Pending = functionWork;

#pragma warning disable CS0162
                if (Tracing)
                {
                    Trace(functionWork.Datum);
                }
#pragma warning restore CS0162

                return HoldSelf();
            },
            data);
    }

    internal void CreateHandler(int identity, int priority, Packet? workPaket, TaskState state)
    {
        HandlerTaskDataRecord data = new();
        CreateTask(identity, priority, workPaket, state, (work, word) =>
            {
                HandlerTaskDataRecord? dataRecord = (HandlerTaskDataRecord?)word;
                if (NoWork != work)
                {
                    if (WorkPacketKind == work.Kind)
                    {
                        dataRecord!.WorkInAdd(work);
                    }
                    else
                    {
                        dataRecord!.DeviceInAdd(work);
                    }
                }

                Packet? workPacket = dataRecord!.WorkIn;
                if (workPacket == NoWork)
                {
                    return MarkWaiting();
                }

                int count = workPacket.Datum;
                if (count >= Packet.DataSize)
                {
                    dataRecord.WorkIn = workPacket.Link;
                    return QueuePacket(workPacket);
                }

                Packet? devicePacket = dataRecord.DeviceIn;
                if (devicePacket == NoWork)
                {
                    return MarkWaiting();
                }

                dataRecord.DeviceIn = devicePacket.Link;
                devicePacket.Datum = workPacket.Data[count];
                workPacket.Datum = count + 1;
                return QueuePacket(devicePacket);
            },
            data);
    }

    internal void CreateIdler(int identity, int priority, Packet? work, TaskState state)
    {
        IdleTaskDataRecord data = new();
        CreateTask(identity, priority, work, state, (_, wordArg) =>
            {
                IdleTaskDataRecord? dataRecord = (IdleTaskDataRecord?)wordArg;
                dataRecord!.Count--;
                if (dataRecord.Count == 0)
                {
                    return HoldSelf();
                }

                if ((dataRecord.Control & 1) == 0)
                {
                    dataRecord.Control /= 2;
                    return Release(DeviceA);
                }

                dataRecord.Control = (dataRecord.Control / 2) ^ 53256;
                return Release(DeviceB);
            },
            data);
    }

    internal Packet CreatePacket(Packet? link, int identity, int kind)
    {
        return new Packet(link, identity, kind);
    }

    internal void CreateTask(int identity, int priority, Packet? work, TaskState state, Func<Packet?, RbObject?, TaskControlBlock?> aBlock,
        RbObject data)
    {
        TaskControlBlock t = new(_taskList!, identity, priority, work, state, aBlock, data);
        _taskList = t;
        _taskTable[identity] = t;
    }

    internal void CreateWorker(int identity, int priority, Packet? workPaket, TaskState state)
    {
        WorkerTaskDataRecord dataRecord = new();
        CreateTask(identity, priority, workPaket, state, (work, word) =>
            {
                WorkerTaskDataRecord? data = (WorkerTaskDataRecord?)word;
                if (NoWork == work)
                {
                    return MarkWaiting();
                }

                data!.Destination = HandlerA == data.Destination ? HandlerB : HandlerA;
                work.Identity = data.Destination;
                work.Datum = 0;
                for (int i = 0; i < Packet.DataSize; i++)
                {
                    data.Count++;
                    if (data.Count > 26)
                    {
                        data.Count = 1;
                    }

                    work.Data[i] = 65 + data.Count - 1;
                }

                return QueuePacket(work);
            },
            dataRecord);
    }

    public bool Start()
    {
        CreateIdler(Idler, 0, NoWork, TaskState.CreateRunning());
        Packet workQ = CreatePacket(NoWork, Worker, WorkPacketKind);
        workQ = CreatePacket(workQ, Worker, WorkPacketKind);

        CreateWorker(Worker, 1000, workQ, TaskState.CreateWaitingWithPacket());
        workQ = CreatePacket(NoWork, DeviceA, DevicePacketKind);
        workQ = CreatePacket(workQ, DeviceA, DevicePacketKind);
        workQ = CreatePacket(workQ, DeviceA, DevicePacketKind);

        CreateHandler(HandlerA, 2000, workQ, TaskState.CreateWaitingWithPacket());
        workQ = CreatePacket(NoWork, DeviceB, DevicePacketKind);
        workQ = CreatePacket(workQ, DeviceB, DevicePacketKind);
        workQ = CreatePacket(workQ, DeviceB, DevicePacketKind);

        CreateHandler(HandlerB, 3000, workQ, TaskState.CreateWaitingWithPacket());
        CreateDevice(DeviceA, 4000, NoWork, TaskState.CreateWaiting());
        CreateDevice(DeviceB, 5000, NoWork, TaskState.CreateWaiting());

        Schedule();

        return _queuePacketCount == 23246 && _holdCount == 9297;
    }

    internal TaskControlBlock FindTask(int identity)
    {
        TaskControlBlock? t = _taskTable[identity];
        return NoTask == t ? throw new InvalidOperationException("findTask failed") : t;
    }

    internal TaskControlBlock HoldSelf()
    {
        _holdCount++;
        _currentTask!.IsTaskHolding = true;
        return _currentTask.Link;
    }

    internal TaskControlBlock? QueuePacket(Packet packet)
    {
        TaskControlBlock t = FindTask(packet.Identity);
        if (t == NoTask)
        {
            return NoTask;
        }

        _queuePacketCount++;

        packet.Link = NoWork;
        packet.Identity = _currentTaskIdentity;
        return t.AddInputAndCheckPriority(packet, _currentTask);
    }

    internal TaskControlBlock? Release(int identity)
    {
        TaskControlBlock t = FindTask(identity);
        if (NoTask == t)
        {
            return NoTask;
        }

        t.IsTaskHolding = false;
        return t.Priority > _currentTask!.Priority ? t : _currentTask;
    }

    internal void Trace(int id)
    {
        _layout--;
        if (_layout <= 0)
        {
            Console.WriteLine();
            _layout = 50;
        }

        Console.Write(id);
    }

    internal TaskControlBlock MarkWaiting()
    {
        _currentTask!.IsTaskWaiting = true;
        return _currentTask;
    }

    internal void Schedule()
    {
        _currentTask = _taskList;
        while (NoTask != _currentTask)
        {
            if (_currentTask.IsTaskHoldingOrWaiting)
            {
                _currentTask = _currentTask.Link;
            }
            else
            {
                _currentTaskIdentity = _currentTask.Identity;


#pragma warning disable CS0162
                // ReSharper disable HeuristicUnreachableCode
                if (Tracing)  
                {
                    Trace(_currentTaskIdentity);
                }
                // ReSharper restore HeuristicUnreachableCode
#pragma warning restore CS0162

                _currentTask = _currentTask.RunTask();
            }
        }
    }
}
