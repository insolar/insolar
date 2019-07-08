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
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type TranscriptDequeueElement struct {
	prev  *TranscriptDequeueElement
	next  *TranscriptDequeueElement
	value *Transcript
}

// TODO: probably it's better to rewrite it using linked list
type TranscriptDequeue struct {
	lock   sync.Locker
	first  *TranscriptDequeueElement
	last   *TranscriptDequeueElement
	length int
}

func (d *TranscriptDequeue) pushOne(el *Transcript) {
	newElement := &TranscriptDequeueElement{value: el}
	lastElement := d.last

	if lastElement != nil {
		newElement.prev = lastElement
		lastElement.next = newElement
		d.last = newElement
	} else {
		d.first, d.last = newElement, newElement
	}

	d.length++
}

func (d *TranscriptDequeue) Push(els ...*Transcript) {
	d.lock.Lock()
	defer d.lock.Unlock()

	for _, el := range els {
		d.pushOne(el)
	}
}

func (d *TranscriptDequeue) prependOne(el *Transcript) {
	newElement := &TranscriptDequeueElement{value: el}
	firstElement := d.first

	if firstElement != nil {
		newElement.next = firstElement
		firstElement.prev = newElement
		d.first = newElement
	} else {
		d.first, d.last = newElement, newElement
	}

	d.length++
}

func (d *TranscriptDequeue) Prepend(els ...*Transcript) {
	d.lock.Lock()
	defer d.lock.Unlock()

	for i := len(els) - 1; i >= 0; i-- {
		d.prependOne(els[i])
	}
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

	for elem := d.first; elem != nil; elem = elem.next {
		if elem.value.RequestRef.Compare(ref) == 0 {
			return true
		}
	}
	return false
}

func (d *TranscriptDequeue) PopByReference(ref insolar.Reference) *Transcript {
	d.lock.Lock()
	defer d.lock.Unlock()

	for elem := d.first; elem != nil; elem = elem.next {
		if elem.value.RequestRef.Compare(ref) == 0 {
			if elem.prev != nil {
				elem.prev.next = elem.next
			}
			if elem.next != nil {
				elem.next.prev = elem.prev
			}

			d.length--

			return elem.value
		}
	}

	return nil
}

func (d *TranscriptDequeue) HasFromLedger() *Transcript {
	d.lock.Lock()
	defer d.lock.Unlock()

	for elem := d.first; elem != nil; elem = elem.next {
		if elem.value.FromLedger {
			return elem.value
		}
	}
	return nil
}

func (d *TranscriptDequeue) commonPeek(count int) (*TranscriptDequeueElement, []*Transcript) {
	if d.length < count {
		count = d.length
	}

	rv := make([]*Transcript, count)

	var lastElement *TranscriptDequeueElement
	for i := 0; i < count; i++ {
		if lastElement == nil {
			lastElement = d.first
		} else {
			lastElement = lastElement.next
		}
		rv[i] = lastElement.value
	}

	return lastElement, rv
}

func (d *TranscriptDequeue) take(count int) []*Transcript {
	lastElement, rv := d.commonPeek(count)
	if lastElement != nil {
		if lastElement.next == nil {
			d.first, d.last = nil, nil
		} else {
			lastElement.next.prev, d.first = nil, lastElement.next
			lastElement.next = nil
		}

		d.length -= len(rv)
	}

	return rv
}

func (d *TranscriptDequeue) Peek(count int) []*Transcript {
	d.lock.Lock()
	defer d.lock.Unlock()

	_, rv := d.commonPeek(count)
	return rv
}

func (d *TranscriptDequeue) Take(count int) []*Transcript {
	d.lock.Lock()
	defer d.lock.Unlock()

	return d.take(count)
}

func (d *TranscriptDequeue) Rotate() []*Transcript {
	d.lock.Lock()
	defer d.lock.Unlock()

	return d.take(d.length)
}

func (d *TranscriptDequeue) Length() int {
	d.lock.Lock()
	defer d.lock.Unlock()

	return d.length
}

func NewTranscriptDequeue() *TranscriptDequeue {
	return &TranscriptDequeue{
		lock:   &sync.Mutex{},
		first:  nil,
		last:   nil,
		length: 0,
	}
}

type ExecutionBroker struct {
	stateLock   sync.Locker
	mutable     *TranscriptDequeue
	immutable   *TranscriptDequeue
	finished    *TranscriptDequeue
	currentList *CurrentExecutionList

	logicRunner    *LogicRunner
	executionState *ExecutionState
	// currently we need to store ES inside, so it looks like circular dependency
	// and it is circular dependency. it will be removed, once Broker will be
	// moved out of ES.

	ledgerChecked sync.Once

	processorActive uint32

	deduplicationTable map[insolar.Reference]bool
	deduplicationLock  sync.Mutex
}

func (q *ExecutionBroker) tryTakeProcessor(_ context.Context) bool {
	return atomic.CompareAndSwapUint32(&q.processorActive, 0, 1)
}

func (q *ExecutionBroker) releaseProcessor(_ context.Context) {
	atomic.SwapUint32(&q.processorActive, 0)
}

func (q *ExecutionBroker) isActiveProcessor() bool { //nolint: unused
	return atomic.LoadUint32(&q.processorActive) == 1
}

type ExecutionBrokerRotationResult struct {
	Requests              []*Transcript
	Finished              []*Transcript
	LedgerHasMoreRequests bool
}

func (q *ExecutionBroker) checkCurrent(_ context.Context, transcript *Transcript) {
	if q.currentList.Has(*transcript.RequestRef) {
		panic(fmt.Sprintf("requestRef %s is already in currentList", transcript.RequestRef.String()))
	}
}

func (q *ExecutionBroker) getImmutableTask(ctx context.Context) *Transcript {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	transcript := q.immutable.Pop()
	if transcript == nil {
		return nil
	}

	q.checkCurrent(ctx, transcript)
	q.currentList.SetTranscript(transcript)
	return transcript
}

func (q *ExecutionBroker) getMutableTask(ctx context.Context) *Transcript {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	transcript := q.mutable.Pop()
	if transcript == nil {
		return nil
	}

	q.checkCurrent(ctx, transcript)
	q.currentList.SetTranscript(transcript)
	return transcript
}

func (q *ExecutionBroker) storeCurrent(ctx context.Context, transcript *Transcript) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	q.checkCurrent(ctx, transcript)
	q.currentList.SetTranscript(transcript)
}

func (q *ExecutionBroker) releaseTask(_ context.Context, transcript *Transcript) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if !q.currentList.Has(*transcript.RequestRef) {
		return
	}
	q.currentList.Delete(*transcript.RequestRef)

	queue := q.mutable
	if transcript.Request.Immutable {
		queue = q.immutable
	}

	queue.Prepend(transcript)
}

func (q *ExecutionBroker) finishTask(ctx context.Context, transcript *Transcript) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	logger := inslogger.FromContext(ctx)

	q.finished.Push(transcript)

	if !q.currentList.Has(*transcript.RequestRef) {
		logger.Error("[ ExecutionBroker.FinishTask ] task '%s' is not in current", transcript.RequestRef.String())
	} else {
		q.currentList.Delete(*transcript.RequestRef)
	}
}

func (q *ExecutionBroker) processTranscript(ctx context.Context, transcript *Transcript) bool {
	defer q.releaseTask(ctx, transcript)

	if readyToExecute := q.Check(ctx); !readyToExecute {
		return false
	}

	q.Execute(ctx, transcript)
	// q.finishTask(ctx, transcript)
	return true
}

func (q *ExecutionBroker) storeWithoutDuplication(_ context.Context, transcript *Transcript) bool {
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
		if q.storeWithoutDuplication(ctx, transcript) {
			inslogger.FromContext(ctx).Info(
				"Already know about request ",
				transcript.RequestRef.String(), ", skipping",
			)
			continue
		}

		if transcript.Request.Immutable {
			q.storeCurrent(ctx, transcript)
			go q.processTranscript(ctx, transcript)
		} else {
			q.stateLock.Lock()
			q.mutable.Prepend(transcript)
			q.stateLock.Unlock()
		}
	}
	if start {
		q.StartProcessorIfNeeded(ctx)
	}
}

// One shouldn't mix immutable calls and mutable ones
func (q *ExecutionBroker) Put(ctx context.Context, start bool, transcripts ...*Transcript) {
	for _, transcript := range transcripts {
		if q.storeWithoutDuplication(ctx, transcript) {
			inslogger.FromContext(ctx).Info(
				"Already know about request ",
				transcript.RequestRef.String(), ", skipping",
			)
			continue
		}

		if transcript.Request.Immutable {
			q.storeCurrent(ctx, transcript)
			go q.processTranscript(ctx, transcript)
		} else {
			q.stateLock.Lock()
			q.mutable.Push(transcript)
			q.stateLock.Unlock()
		}
	}
	if start {
		q.StartProcessorIfNeeded(ctx)
	}
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
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	q.deduplicationLock.Lock()
	defer q.deduplicationLock.Unlock()

	delete(q.deduplicationTable, *r)
	return q.mutable.PopByReference(*r)
}

func (q *ExecutionBroker) startProcessor(ctx context.Context) {
	defer q.releaseProcessor(ctx)

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

	// i've removed "if we have tasks here"; we can be there in two cases:
	// 1) we've put a task into queue and automatically start processor
	// 2) we've explicitly ask broker to be here
	// both cases means we knew what we are doing and so it's just an
	// unneeded optimisation
	if q.tryTakeProcessor(ctx) {
		logger.Info("[ StartProcessorIfNeeded ] Starting a new queue processor")
		go q.startProcessor(ctx)
	}

}

// TODO: probably rotation should wait till processActive == false (??)
func (q *ExecutionBroker) Rotate(count int) *ExecutionBrokerRotationResult {
	// take mutables, then, if we can, take immutables, if something was left -
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	rv := &ExecutionBrokerRotationResult{
		Requests:              q.mutable.Take(count),
		Finished:              q.finished.Rotate(),
		LedgerHasMoreRequests: false,
	}

	if leftCount := count - len(rv.Requests); leftCount > 0 {
		rv.Requests = append(rv.Requests, q.immutable.Take(leftCount)...)
	}

	if len(rv.Requests) > 0 && (q.mutable.Length() > 0 || q.immutable.Length() > 0) {
		rv.LedgerHasMoreRequests = true
	}

	_ = q.mutable.Rotate()
	_ = q.immutable.Rotate()

	q.deduplicationLock.Lock()
	q.deduplicationTable = make(map[insolar.Reference]bool)
	q.deduplicationLock.Unlock()

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

	reply, err := q.logicRunner.RequestsExecutor.ExecuteAndSave(ctx, transcript)
	if err != nil {
		inslogger.FromContext(ctx).Warn("contract execution error: ", err)
	}

	q.finishTask(ctx, transcript) // TODO: hack for now, later that function need to be splitted

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
		stateLock:   &sync.Mutex{},
		mutable:     NewTranscriptDequeue(),
		immutable:   NewTranscriptDequeue(),
		finished:    NewTranscriptDequeue(),
		currentList: NewCurrentExecutionList(),

		executionState: es,

		ledgerChecked: sync.Once{},

		deduplicationLock:  sync.Mutex{},
		deduplicationTable: make(map[insolar.Reference]bool),
	}
}
