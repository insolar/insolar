//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package logicrunner

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// TODO: probably it's better to rewrite it using linked list
type TranscriptDequeue struct {
	lock  sync.Mutex
	queue []*Transcript
}

func (d *TranscriptDequeue) Push(els ...*Transcript) {
	d.lock.Lock()
	d.queue = append(d.queue, els...)
	d.lock.Unlock()
}

func (d *TranscriptDequeue) Prepend(els ...*Transcript) {
	d.lock.Lock()
	d.queue = append(els, d.queue...)
	d.lock.Unlock()
}

func (d *TranscriptDequeue) Pop() *Transcript {
	elements := d.Take(1)
	if len(elements) == 0 {
		return nil
	}
	return elements[0]
}

func (d *TranscriptDequeue) Has(ref insolar.Reference) bool {
	d.lock.Lock()
	defer d.lock.Unlock()

	for pos := len(d.queue) - 1; pos >= 0; pos-- {
		if d.queue[pos].RequestRef.Compare(ref) == 0 {
			return true
		}
	}
	return false
}

func (d *TranscriptDequeue) PopByReference(ref *insolar.Reference) *Transcript {
	d.lock.Lock()
	toDelete := -1
	for pos := len(d.queue) - 1; pos >= 0; pos-- {
		if d.queue[pos].RequestRef.Compare(*ref) == 0 {
			toDelete = pos
			break
		}
	}
	var rv *Transcript
	if toDelete != -1 {
		rv = d.queue[toDelete]
		d.queue = append(d.queue[:toDelete], d.queue[toDelete+1:]...)
	}
	d.lock.Unlock()
	return rv
}

func (d *TranscriptDequeue) HasFromLedger() *Transcript {
	d.lock.Lock()
	defer d.lock.Unlock()

	for _, t := range d.queue {
		if t.FromLedger {
			return t
		}
	}
	return nil
}

func (d *TranscriptDequeue) Take(count int) []*Transcript {
	d.lock.Lock()
	defer d.lock.Unlock()

	size := len(d.queue)
	if size < count {
		count = size
	}

	elements := d.queue[:count]
	d.queue = d.queue[count:]

	return elements
}

func (d *TranscriptDequeue) Rotate() []*Transcript {
	d.lock.Lock()
	defer d.lock.Unlock()

	rv := d.queue
	d.queue = make([]*Transcript, 0)

	return rv
}

func (d *TranscriptDequeue) Len() int {
	d.lock.Lock()
	defer d.lock.Unlock()

	return len(d.queue)
}

func NewTranscriptDequeue() *TranscriptDequeue {
	return &TranscriptDequeue{
		lock:  sync.Mutex{},
		queue: make([]*Transcript, 0),
	}
}

type ExecutionBroker struct {
	mutableLock sync.RWMutex
	mutable     *TranscriptDequeue
	immutable   *TranscriptDequeue
	finished    *TranscriptDequeue

	logicRunner    *LogicRunner
	executionState *ExecutionState
	// currently we need to store ES inside, so it looks like circular dependency
	// and it is circular dependency. it will be removed, once Broker will be
	// moved out of ES.

	ledgerChecked sync.Once

	processActive bool
	processLock   sync.Mutex

	deduplicationTable map[insolar.Reference]bool
	deduplicationLock  sync.Mutex
}

type ExecutionBrokerRotationResult struct {
	Requests              []*Transcript
	Finished              []*Transcript
	LedgerHasMoreRequests bool
}

func (q *ExecutionBroker) getImmutableTask(_ context.Context) *Transcript {
	q.mutableLock.RLock()
	defer q.mutableLock.RUnlock()

	transcript := q.immutable.Pop()
	if transcript == nil {
		return nil
	}

	return transcript
}

func (q *ExecutionBroker) getMutableTask(_ context.Context) *Transcript {
	q.mutableLock.RLock()
	defer q.mutableLock.RUnlock()

	transcript := q.mutable.Pop()
	if transcript == nil {
		return nil
	}

	return transcript
}

func (q *ExecutionBroker) releaseTask(_ context.Context, transcript *Transcript) {
	q.mutableLock.RLock()
	defer q.mutableLock.RUnlock()

	if q.finished.Has(*transcript.RequestRef) {
		return
	}

	queue := q.mutable
	if transcript.Request.Immutable {
		queue = q.immutable
	}

	if queue.Has(*transcript.RequestRef) {
		return
	}

	queue.Prepend(transcript)
}

func (q *ExecutionBroker) finishTask(_ context.Context, transcript *Transcript) {
	q.mutableLock.RLock()
	defer q.mutableLock.RUnlock()

	q.finished.Push(transcript)
}

func (q *ExecutionBroker) processTranscript(ctx context.Context, transcript *Transcript) bool {
	defer q.releaseTask(ctx, transcript)

	if readyToExecute := q.Check(ctx); !readyToExecute {
		return false
	}

	q.Execute(ctx, transcript)
	q.finishTask(ctx, transcript)
	return true
}

func (q *ExecutionBroker) isDuplicate(_ context.Context, transcript *Transcript) bool {
	q.deduplicationLock.Lock()
	defer q.deduplicationLock.Unlock()

	if _, ok := q.deduplicationTable[*transcript.RequestRef]; ok {
		return true
	}
	q.deduplicationTable[*transcript.RequestRef] = true
	return false
}

func (q *ExecutionBroker) Prepend(ctx context.Context, start bool, transcripts ...*Transcript) {
	for _, transcript := range transcripts {
		if q.isDuplicate(ctx, transcript) {
			continue
		}

		if transcript.Request.Immutable {
			go q.processTranscript(ctx, transcript)
		} else {
			q.mutableLock.RLock()
			q.mutable.Prepend(transcript)
			q.mutableLock.RUnlock()
		}
	}
	if start {
		q.StartProcessorIfNeeded(ctx)
	}
}

// One shouldn't mix immutable calls and mutable ones
func (q *ExecutionBroker) Put(ctx context.Context, start bool, transcripts ...*Transcript) {
	for _, transcript := range transcripts {
		if q.isDuplicate(ctx, transcript) {
			continue
		}

		if transcript.Request.Immutable {
			go q.processTranscript(ctx, transcript)
		} else {
			q.mutableLock.RLock()
			q.mutable.Push(transcript)
			q.mutableLock.RUnlock()
		}
	}
	if start {
		q.StartProcessorIfNeeded(ctx)
	}
}

func (q *ExecutionBroker) get(_ context.Context) *Transcript {
	q.mutableLock.RLock()
	defer q.mutableLock.RUnlock()

	return q.mutable.Pop()
}

func (q *ExecutionBroker) HasLedgerRequest(_ context.Context) *Transcript {
	if obj := q.mutable.HasFromLedger(); obj != nil {
		return obj
	}
	if obj := q.immutable.HasFromLedger(); obj != nil {
		return obj
	}
	return nil
}

func (q *ExecutionBroker) GetByReference(_ context.Context, r *insolar.Reference) *Transcript {
	q.mutableLock.RLock()
	defer q.mutableLock.RUnlock()

	return q.mutable.PopByReference(r)
}

func (q *ExecutionBroker) startProcessor(ctx context.Context) {
	defer func() {
		q.processLock.Lock()
		q.processActive = false
		q.processLock.Unlock()
	}()

	// сhecking we're eligible to execute contracts
	if readyToExecute := q.Check(ctx); !readyToExecute {
		return
	}

	// processing immutable queue (it can appear if we were in pending state)
	// run simultaneously all immutable transcripts and forget about them
	for elem := q.getImmutableTask(ctx); elem != nil; elem = q.getImmutableTask(ctx) {
		go q.processTranscript(ctx, elem)
	}

	// processing mutable queue
	for transcript := q.getMutableTask(ctx); transcript != nil; transcript = q.getMutableTask(ctx) {
		if !q.processTranscript(ctx, transcript) {
			break
		}
	}
}

// StartProcessorIfNeeded processes queue messages in strict order
// We need to start manually execution broker only if we were in pending and now we're not
func (q *ExecutionBroker) StartProcessorIfNeeded(ctx context.Context) {
	logger := inslogger.FromContext(ctx)

	q.processLock.Lock()
	// i've removed "if we have tasks here"; we can be there in two cases:
	// 1) we've put a task into queue and automatically start processor
	// 2) we've explicitly ask broker to be here
	// both cases means we knew what we are doing and so it's just an
	// unneeded optimisation
	if !q.processActive {
		logger.Info("[ StartProcessorIfNeeded ] Starting a new queue processor")

		q.processActive = true
		go q.startProcessor(ctx)
	}
	q.processLock.Unlock()

}

// TODO: locking system (mutableLock) should be reconsidered
// TODO: probably rotation should wait till processActive == false (??)
func (q *ExecutionBroker) Rotate(count int) *ExecutionBrokerRotationResult {
	// take mutables, then, if we can, take immutables, if something was left -
	q.mutableLock.Lock()
	rv := &ExecutionBrokerRotationResult{
		Requests:              q.mutable.Take(count),
		Finished:              q.finished.Rotate(),
		LedgerHasMoreRequests: false,
	}

	if leftCount := count - len(rv.Requests); leftCount > 0 {
		rv.Requests = append(rv.Requests, q.immutable.Take(leftCount)...)
	}

	if len(rv.Requests) > 0 && (q.mutable.Len() > 0 || q.immutable.Len() > 0) {
		rv.LedgerHasMoreRequests = true
	}

	_ = q.mutable.Rotate()
	_ = q.immutable.Rotate()

	q.deduplicationLock.Lock()
	q.deduplicationTable = make(map[insolar.Reference]bool)
	q.deduplicationLock.Unlock()

	q.mutableLock.Unlock()

	return rv
}

func (q *ExecutionBroker) Check(ctx context.Context) bool {
	logger := inslogger.FromContext(ctx)
	es := q.executionState

	// check pending state of execution (whether we can process task or not)
	es.Lock()
	if es.pending == message.PendingUnknown {
		logger.Debug("One shouldn't call ExecuteTranscript in case when pending state is unknown")
		es.Unlock()
		return false
	} else if es.pending == message.InPending {
		logger.Debug("Object in pending, wont start queue processor")
		es.Unlock()
		return false
	}
	es.Unlock()

	return true
}

func (q *ExecutionBroker) checkLedgerPendingRequestsBase(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	es := q.executionState
	pub := q.logicRunner.publisher

	wmMessage := makeWMMessage(ctx, es.Ref.Bytes(), getLedgerPendingRequestMsg)
	if err := pub.Publish(InnerMsgTopic, wmMessage); err != nil {
		logger.Warnf("can't send getLedgerPendingRequestMsg: ", err)
	}
}

func (q *ExecutionBroker) checkLedgerPendingRequests(ctx context.Context, transcript *Transcript) {
	if transcript == nil {
		// Ask ledger kindly to give us next pending task and continue execution
		// note: should be done only once
		q.ledgerChecked.Do(func() { q.checkLedgerPendingRequestsBase(ctx) })
	} else if transcript.FromLedger {
		// we've already told ledger that we've processed it's task;
		// trying to take another one
		q.checkLedgerPendingRequestsBase(ctx)
	}
}

func (q *ExecutionBroker) Execute(ctx context.Context, transcript *Transcript) {
	q.checkLedgerPendingRequests(ctx, nil)

	q.executionState.Lock()
	q.executionState.CurrentList.Set(*transcript.RequestRef, transcript)
	q.executionState.Unlock()

	reply, err := q.logicRunner.RequestsExecutor.ExecuteAndSave(ctx, transcript)
	if err != nil {
		inslogger.FromContext(ctx).Warn("contract execution error: ", err)
	}

	q.executionState.Lock()
	q.executionState.CurrentList.Delete(*transcript.RequestRef)
	q.executionState.Unlock()

	go q.logicRunner.RequestsExecutor.SendReply(transcript.Context, transcript, reply, err)

	q.checkLedgerPendingRequests(ctx, transcript)

	// we're checking here that pulse was changed and we should send
	// a message that we've finished processing task
	// note: ideally we should tell here that we've stopped executing
	//       but we only hoped that OnPulse had already told us that
	//       pulse changed and we should stop execution
	q.finishPendingIfNeeded(ctx)
}

func (q *ExecutionBroker) finishPending(ctx context.Context) {
	logger := inslogger.FromContext(ctx)

	msg := message.PendingFinished{Reference: q.executionState.Ref}
	_, err := q.logicRunner.MessageBus.Send(ctx, &msg, nil)
	if err != nil {
		logger.Error("Unable to send PendingFinished message:", err)
	}
}

// finishPendingIfNeeded checks whether last execution was a pending one.
// If this is true as a side effect the function sends a PendingFinished
// message to the current executor
func (q *ExecutionBroker) finishPendingIfNeeded(ctx context.Context) {
	es := q.executionState
	lr := q.logicRunner

	es.Lock()
	defer es.Unlock()

	if es.pending != message.InPending {
		return
	}

	es.pending = message.NotPending
	es.PendingConfirmed = false

	pulseObj := lr.pulse(ctx)
	meCurrent, _ := lr.JetCoordinator.IsAuthorized(
		ctx, insolar.DynamicRoleVirtualExecutor, *es.Ref.Record(), pulseObj.PulseNumber, lr.JetCoordinator.Me(),
	)
	if !meCurrent {
		go q.finishPending(ctx)
	}
}

func NewExecutionBroker(es *ExecutionState) *ExecutionBroker {
	return &ExecutionBroker{
		mutableLock: sync.RWMutex{},
		mutable:     NewTranscriptDequeue(),
		immutable:   NewTranscriptDequeue(),
		finished:    NewTranscriptDequeue(),

		executionState: es,

		ledgerChecked: sync.Once{},

		deduplicationLock:  sync.Mutex{},
		deduplicationTable: make(map[insolar.Reference]bool),
	}
}
