package main

import "fmt"

// Richards simulates the task dispatcher of an operating system kernel.
type Richards struct{}

// Benchmark runs the scheduler once.
func (*Richards) Benchmark() any {
	return newScheduler().start()
}

// VerifyResult checks the scheduler's own result.
func (*Richards) VerifyResult(r any) bool {
	return r.(bool)
}

// Java RBObject constants (ADR-012). NO_WORK / NO_TASK are written as nil (ADR-026).
const (
	idler    int  = 0
	worker   int  = 1
	handlerA int  = 2
	handlerB int  = 3
	deviceA  int  = 4
	deviceB  int  = 5
	numTypes int  = 6
	tracing  bool = false // Scheduler.TRACING

	devicePacketKind int = 0
	workPacketKind   int = 1

	packetDataSize int = 4 // Packet.DATA_SIZE
)

// rbObjectAppend is Java's inherited RBObject.append (ADR-025).
func rbObjectAppend(packet, queueHead *packet) *packet {
	packet.link = nil
	if queueHead == nil {
		return packet
	}

	mouse := queueHead
	for link := mouse.link; link != nil; link = mouse.link {
		mouse = link
	}
	mouse.link = packet
	return queueHead
}

type packet struct {
	link     *packet
	identity int
	kind     int
	datum    int
	data     []int
}

func newPacket(link *packet, identity, kind int) *packet {
	return &packet{
		link:     link,
		identity: identity,
		kind:     kind,
		datum:    0,
		data:     make([]int, packetDataSize),
	}
}

type taskState struct {
	packetPending bool
	taskWaiting   bool
	taskHolding   bool
}

func (t *taskState) running() {
	t.packetPending = false
	t.taskWaiting = false
	t.taskHolding = false
}

func (t *taskState) waiting() {
	t.packetPending = false
	t.taskHolding = false
	t.taskWaiting = true
}

func (t *taskState) waitingWithPacket() {
	t.taskHolding = false
	t.taskWaiting = true
	t.packetPending = true
}

func taskStateCreateRunning() *taskState {
	t := &taskState{}
	t.running()
	return t
}

func taskStateCreateWaiting() *taskState {
	t := &taskState{}
	t.waiting()
	return t
}

func taskStateCreateWaitingWithPacket() *taskState {
	t := &taskState{}
	t.waitingWithPacket()
	return t
}

// taskControlBlock extends TaskState in Java; it has its own copy of the
// state fields and of the state methods it executes (ADR-025).
type taskControlBlock struct {
	link          *taskControlBlock
	identity      int
	priority      int
	input         *packet
	packetPending bool
	taskWaiting   bool
	taskHolding   bool
	function      func(work *packet, word any) *taskControlBlock
	handle        any
}

func newTaskControlBlock(link *taskControlBlock, identity, priority int,
	initialWorkQueue *packet, initialState *taskState,
	fn func(work *packet, word any) *taskControlBlock, privateData any) *taskControlBlock {
	t := &taskControlBlock{
		link:     link,
		identity: identity,
		priority: priority,
		input:    initialWorkQueue,
	}
	t.packetPending = initialState.packetPending
	t.taskWaiting = initialState.taskWaiting
	t.taskHolding = initialState.taskHolding
	t.function = fn
	t.handle = privateData
	return t
}

// markPacketPending is Java's TaskState.packetPending().
func (t *taskControlBlock) markPacketPending() {
	t.packetPending = true
	t.taskWaiting = false
	t.taskHolding = false
}

func (t *taskControlBlock) running() {
	t.packetPending = false
	t.taskWaiting = false
	t.taskHolding = false
}

func (t *taskControlBlock) isTaskHoldingOrWaiting() bool {
	return t.taskHolding || (!t.packetPending && t.taskWaiting)
}

func (t *taskControlBlock) isWaitingWithPacket() bool {
	return t.packetPending && t.taskWaiting && !t.taskHolding
}

func (t *taskControlBlock) addInputAndCheckPriority(packet *packet, oldTask *taskControlBlock) *taskControlBlock {
	if t.input == nil {
		t.input = packet
		t.packetPending = true
		if t.priority > oldTask.priority {
			return t
		}
	} else {
		t.input = rbObjectAppend(packet, t.input)
	}
	return oldTask
}

func (t *taskControlBlock) runTask() *taskControlBlock {
	var message *packet
	if t.isWaitingWithPacket() {
		message = t.input
		t.input = message.link
		if t.input == nil {
			t.running()
		} else {
			t.markPacketPending()
		}
	} else {
		message = nil
	}
	return t.function(message, t.handle)
}

type deviceTaskDataRecord struct {
	pending *packet
}

func newDeviceTaskDataRecord() *deviceTaskDataRecord {
	return &deviceTaskDataRecord{pending: nil}
}

type handlerTaskDataRecord struct {
	workIn   *packet
	deviceIn *packet
}

func newHandlerTaskDataRecord() *handlerTaskDataRecord {
	return &handlerTaskDataRecord{workIn: nil, deviceIn: nil}
}

func (h *handlerTaskDataRecord) deviceInAdd(packet *packet) {
	h.deviceIn = rbObjectAppend(packet, h.deviceIn)
}

func (h *handlerTaskDataRecord) workInAdd(packet *packet) {
	h.workIn = rbObjectAppend(packet, h.workIn)
}

type idleTaskDataRecord struct {
	control int
	count   int
}

func newIdleTaskDataRecord() *idleTaskDataRecord {
	return &idleTaskDataRecord{control: 1, count: 10000}
}

type workerTaskDataRecord struct {
	destination int
	count       int
}

func newWorkerTaskDataRecord() *workerTaskDataRecord {
	return &workerTaskDataRecord{destination: handlerA, count: 0}
}

type scheduler struct {
	taskList            *taskControlBlock
	currentTask         *taskControlBlock
	currentTaskIdentity int
	taskTable           []*taskControlBlock

	queuePacketCount int
	holdCount        int

	layout int
}

func newScheduler() *scheduler {
	s := &scheduler{}
	// init tracing
	s.layout = 0

	// init scheduler
	s.queuePacketCount = 0
	s.holdCount = 0
	s.taskTable = make([]*taskControlBlock, numTypes)
	for i := range s.taskTable {
		s.taskTable[i] = nil
	}
	s.taskList = nil
	return s
}

func (s *scheduler) createDevice(identity, priority int, workPacket *packet, state *taskState) {
	data := newDeviceTaskDataRecord()

	s.createTask(identity, priority, workPacket, state,
		func(workArg *packet, wordArg any) *taskControlBlock {
			dataRecord := wordArg.(*deviceTaskDataRecord)
			functionWork := workArg
			if functionWork == nil {
				if functionWork = dataRecord.pending; functionWork == nil {
					return s.markWaiting()
				} else {
					dataRecord.pending = nil
					return s.queuePacket(functionWork)
				}
			} else {
				dataRecord.pending = functionWork
				if tracing {
					s.trace(functionWork.datum)
				}
				return s.holdSelf()
			}
		}, data)
}

func (s *scheduler) createHandler(identity, priority int, workPacket *packet, state *taskState) {
	data := newHandlerTaskDataRecord()
	s.createTask(identity, priority, workPacket, state,
		func(work *packet, word any) *taskControlBlock {
			dataRecord := word.(*handlerTaskDataRecord)
			if work != nil {
				if work.kind == workPacketKind {
					dataRecord.workInAdd(work)
				} else {
					dataRecord.deviceInAdd(work)
				}
			}

			if workPacket := dataRecord.workIn; workPacket == nil {
				return s.markWaiting()
			} else {
				count := workPacket.datum
				if count >= packetDataSize {
					dataRecord.workIn = workPacket.link
					return s.queuePacket(workPacket)
				} else if devicePacket := dataRecord.deviceIn; devicePacket == nil {
					return s.markWaiting()
				} else {
					dataRecord.deviceIn = devicePacket.link
					devicePacket.datum = workPacket.data[count]
					workPacket.datum = count + 1
					return s.queuePacket(devicePacket)
				}
			}
		}, data)
}

func (s *scheduler) createIdler(identity, priority int, work *packet, state *taskState) {
	data := newIdleTaskDataRecord()
	s.createTask(identity, priority, work, state,
		func(workArg *packet, wordArg any) *taskControlBlock {
			dataRecord := wordArg.(*idleTaskDataRecord)
			dataRecord.count = dataRecord.count - 1
			if dataRecord.count == 0 {
				return s.holdSelf()
			} else if (dataRecord.control & 1) == 0 {
				dataRecord.control = dataRecord.control / 2
				return s.release(deviceA)
			} else {
				dataRecord.control = (dataRecord.control / 2) ^ 53256
				return s.release(deviceB)
			}
		}, data)
}

func (s *scheduler) createPacket(link *packet, identity, kind int) *packet {
	return newPacket(link, identity, kind)
}

func (s *scheduler) createTask(identity, priority int, work *packet, state *taskState,
	aBlock func(work *packet, word any) *taskControlBlock, data any) {
	t := newTaskControlBlock(s.taskList, identity, priority, work, state, aBlock, data)
	s.taskList = t
	s.taskTable[identity] = t
}

func (s *scheduler) createWorker(identity, priority int, workPacket *packet, state *taskState) {
	dataRecord := newWorkerTaskDataRecord()
	s.createTask(identity, priority, workPacket, state,
		func(work *packet, word any) *taskControlBlock {
			data := word.(*workerTaskDataRecord)
			if work == nil {
				return s.markWaiting()
			} else {
				if data.destination == handlerA {
					data.destination = handlerB
				} else {
					data.destination = handlerA
				}
				work.identity = data.destination
				work.datum = 0
				for i := 0; i < packetDataSize; i++ {
					data.count = data.count + 1
					if data.count > 26 {
						data.count = 1
					}
					work.data[i] = 65 + data.count - 1
				}
				return s.queuePacket(work)
			}
		}, dataRecord)
}

func (s *scheduler) start() bool {
	var workQ *packet

	s.createIdler(idler, 0, nil, taskStateCreateRunning())
	workQ = s.createPacket(nil, worker, workPacketKind)
	workQ = s.createPacket(workQ, worker, workPacketKind)

	s.createWorker(worker, 1000, workQ, taskStateCreateWaitingWithPacket())
	workQ = s.createPacket(nil, deviceA, devicePacketKind)
	workQ = s.createPacket(workQ, deviceA, devicePacketKind)
	workQ = s.createPacket(workQ, deviceA, devicePacketKind)

	s.createHandler(handlerA, 2000, workQ, taskStateCreateWaitingWithPacket())
	workQ = s.createPacket(nil, deviceB, devicePacketKind)
	workQ = s.createPacket(workQ, deviceB, devicePacketKind)
	workQ = s.createPacket(workQ, deviceB, devicePacketKind)

	s.createHandler(handlerB, 3000, workQ, taskStateCreateWaitingWithPacket())
	s.createDevice(deviceA, 4000, nil, taskStateCreateWaiting())
	s.createDevice(deviceB, 5000, nil, taskStateCreateWaiting())

	s.schedule()

	return s.queuePacketCount == 23246 && s.holdCount == 9297
}

func (s *scheduler) findTask(identity int) *taskControlBlock {
	t := s.taskTable[identity]
	if t == nil {
		panic("findTask failed")
	}
	return t
}

func (s *scheduler) holdSelf() *taskControlBlock {
	s.holdCount = s.holdCount + 1
	s.currentTask.taskHolding = true
	return s.currentTask.link
}

func (s *scheduler) queuePacket(packet *packet) *taskControlBlock {
	t := s.findTask(packet.identity)
	if t == nil {
		return nil
	}

	s.queuePacketCount = s.queuePacketCount + 1

	packet.link = nil
	packet.identity = s.currentTaskIdentity
	return t.addInputAndCheckPriority(packet, s.currentTask)
}

func (s *scheduler) release(identity int) *taskControlBlock {
	t := s.findTask(identity)
	if t == nil {
		return nil
	}
	t.taskHolding = false
	if t.priority > s.currentTask.priority {
		return t
	} else {
		return s.currentTask
	}
}

func (s *scheduler) trace(id int) {
	s.layout = s.layout - 1
	if s.layout <= 0 {
		fmt.Println()
		s.layout = 50
	}
	fmt.Print(id)
}

func (s *scheduler) markWaiting() *taskControlBlock {
	s.currentTask.taskWaiting = true
	return s.currentTask
}

func (s *scheduler) schedule() {
	s.currentTask = s.taskList
	for s.currentTask != nil {
		if s.currentTask.isTaskHoldingOrWaiting() {
			s.currentTask = s.currentTask.link
		} else {
			s.currentTaskIdentity = s.currentTask.identity
			if tracing {
				s.trace(s.currentTaskIdentity)
			}
			s.currentTask = s.currentTask.runTask()
		}
	}
}
