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
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
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

func (d *TranscriptDequeue) PopByReference(ref *insolar.Reference) *Transcript {
	d.lock.Lock()
	toDelete := -1
	for pos := len(d.queue) - 1; pos >= 0; pos-- {
		if d.queue[pos].RequestRef.Compare(*ref) == 0 {
			toDelete = pos
			break
		}
	}
	var rv *Transcript = nil
	if toDelete != -1 {
		rv = d.queue[toDelete]
		d.queue = append(d.queue[:toDelete], d.queue[toDelete+1:]...)
	}
	d.lock.Unlock()
	return rv
}

func (d *TranscriptDequeue) HasFromLedger() *Transcript {
	for _, t := range d.queue {
		if t.FromLedger {
			return t
		}
	}
	return nil
}

func (d *TranscriptDequeue) Take(count int) []*Transcript {
	d.lock.Lock()

	size := len(d.queue)
	if size < count {
		count = size
	}

	elements := d.queue[:count]
	d.queue = d.queue[count:]

	d.lock.Unlock()
	return elements
}

func (d *TranscriptDequeue) Rotate() []*Transcript {
	d.lock.Lock()
	rv := d.queue
	d.queue = make([]*Transcript, 0)
	d.lock.Unlock()
	return rv
}

func (d *TranscriptDequeue) Len() int {
	d.lock.Lock()
	rv := len(d.queue)
	d.lock.Unlock()
	return rv
}

func NewTranscriptDequeue() *TranscriptDequeue {
	return &TranscriptDequeue{
		lock:  sync.Mutex{},
		queue: make([]*Transcript, 0),
	}
}

type Transcript struct {
	State interface{} // Shows current execution status of task

	Context          context.Context
	LogicContext     *insolar.LogicCallContext
	Request          *record.Request
	RequestRef       *insolar.Reference
	RequesterNode    *insolar.Reference
	SentResult       bool
	Nonce            uint64
	Deactivate       bool
	OutgoingRequests []OutgoingRequest

	Parcel     insolar.Parcel
	FromLedger bool
}

func NewTranscript(ctx context.Context, parcel insolar.Parcel, requestRef *insolar.Reference,
	pulse *insolar.Pulse, callee insolar.Reference) *Transcript {

	msg := parcel.Message().(*message.CallMethod)

	logicalContext := &insolar.LogicCallContext{
		Mode:            insolar.ExecuteCallMode,
		Caller:          msg.GetCaller(),
		Callee:          &callee,
		Request:         requestRef,
		Time:            time.Now(), // TODO: probably we should take it earlier
		Pulse:           *pulse,
		TraceID:         inslogger.TraceID(ctx),
		CallerPrototype: &msg.CallerPrototype,
		Immutable:       msg.Immutable,
	}
	sender := parcel.GetSender()

	return &Transcript{
		Context:       ctx,
		LogicContext:  logicalContext,
		Request:       &msg.Request,
		RequestRef:    requestRef,
		RequesterNode: &sender,
		SentResult:    false,
		Nonce:         0,
		Deactivate:    false,

		Parcel:     parcel,
		FromLedger: false,
	}
}

type OutgoingRequest struct {
	Request   record.Request
	NewObject *Ref
	Response  []byte
	Error     error
}

func (t *Transcript) AddOutgoingRequest(
	ctx context.Context, request record.Request, result []byte, newObject *Ref, err error,
) {
	rec := OutgoingRequest{
		Request:   request,
		Response:  result,
		NewObject: newObject,
		Error:     err,
	}
	t.OutgoingRequests = append(t.OutgoingRequests, rec)
}

type CurrentExecutionList struct {
	lock       sync.RWMutex
	executions map[insolar.Reference]*Transcript
}

func (ces *CurrentExecutionList) Get(requestRef insolar.Reference) *Transcript {
	ces.lock.RLock()
	rv := ces.executions[requestRef]
	ces.lock.RUnlock()
	return rv
}

func (ces *CurrentExecutionList) Set(requestRef insolar.Reference, ce *Transcript) {
	ces.lock.Lock()
	ces.executions[requestRef] = ce
	ces.lock.Unlock()
}

func (ces *CurrentExecutionList) Delete(requestRef insolar.Reference) {
	ces.lock.Lock()
	delete(ces.executions, requestRef)
	ces.lock.Unlock()
}

func (ces *CurrentExecutionList) GetByTraceID(traceid string) *Transcript {
	ces.lock.RLock()
	for _, ce := range ces.executions {
		if ce.LogicContext.TraceID == traceid {
			ces.lock.RUnlock()
			return ce
		}
	}
	ces.lock.RUnlock()
	return nil
}

func (ces *CurrentExecutionList) GetMutable() *Transcript {
	ces.lock.RLock()
	for _, ce := range ces.executions {
		if !ce.LogicContext.Immutable {
			ces.lock.RUnlock()
			return ce
		}
	}
	ces.lock.RUnlock()
	return nil
}

func (ces *CurrentExecutionList) Cleanup() {
	ces.lock.Lock()
	ces.executions = make(map[insolar.Reference]*Transcript)
	ces.lock.Unlock()
}

func (ces *CurrentExecutionList) Length() int {
	ces.lock.RLock()
	rv := len(ces.executions)
	ces.lock.RUnlock()
	return rv
}

func (ces *CurrentExecutionList) Empty() bool {
	return ces.Length() == 0
}

type CurrentExecutionPredicate func(*Transcript, interface{}) bool

func (ces *CurrentExecutionList) Check(predicate CurrentExecutionPredicate, args interface{}) bool {
	rv := true
	ces.lock.RLock()
	for _, current := range ces.executions {
		if !predicate(current, args) {
			rv = false
			break
		}
	}
	ces.lock.RUnlock()
	return rv
}

func NewCurrentExecutionList() *CurrentExecutionList {
	rv := &CurrentExecutionList{}
	rv.Cleanup()
	return rv
}

type CheckCallback func(context.Context) error
type ExecuteTranscriptCallback func(context.Context, *Transcript, interface{}) error

type ExecutionBroker struct {
	mutableLock sync.RWMutex
	mutable     *TranscriptDequeue
	immutable   *TranscriptDequeue

	finishedLock sync.RWMutex
	finished     *TranscriptDequeue

	checkFunc       CheckCallback
	processFunc     ExecuteTranscriptCallback
	processFuncArgs interface{}

	processActive bool
	processLock   sync.Mutex

	StartProcessorIfNeededCount int32
}

var ErrRetryLater = errors.New("Failed to start task, retry next time")

func (q *ExecutionBroker) processImmutable(ctx context.Context, transcript *Transcript) {
	logger := inslogger.FromContext(ctx)

	if q.processFunc != nil {
		err := q.processFunc(ctx, transcript, q.processFuncArgs)
		if err == nil {
			// In case when we're in pending - we should store it to execute later
			q.finished.Push(transcript)
			return
		} else if err != ErrRetryLater {
			// TODO: we should consider what should be done in that case
			logger.Error("Failed to process immutable transcript:", err)
			return
		}
	}

	// we're retrying this request later, we're in pending now
	q.immutable.Push(transcript)
}

func (q *ExecutionBroker) Prepend(ctx context.Context, start bool, transcripts ...*Transcript) {
	for _, transcript := range transcripts {
		if transcript.LogicContext.Immutable {
			go q.processImmutable(ctx, transcript)
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
		if transcript.LogicContext.Immutable {
			go q.processImmutable(ctx, transcript)
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

func (q *ExecutionBroker) Get(_ context.Context) *Transcript {
	q.mutableLock.RLock()
	rv := q.mutable.Pop()
	q.mutableLock.RUnlock()
	return rv
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
	rv := q.mutable.PopByReference(r)
	q.mutableLock.RUnlock()
	return rv
}

// StartProcessorIfNeeded processes queue messages in strict order
// We need to start manually execution broker only if we were in pending and now we're not
func (q *ExecutionBroker) StartProcessorIfNeeded(ctx context.Context) {
	logger := inslogger.FromContext(ctx)

	if q.checkFunc == nil {
		return
	}

	q.processLock.Lock()
	// i've removed "if we have tasks here"; we can be there in two cases:
	// 1) we've put a task into queue and automatically start processor
	// 2) we've explicitly ask broker to be here
	// both cases means we knew what we are doing and so it's just an
	// unneeded optimisation
	if !q.processActive {
		q.StartProcessorIfNeededCount += 1
		logger.Info("Starting a new queue processor")

		go func() {
			var err error

			// сhecking we're eligable to execute contracts
			err = q.checkFunc(ctx)
			if err != ErrRetryLater {
				// processing immutable queue (it can appear if we were in pending state)
				// run simultaneously all immutable transcripts and forget about them
				for elem := q.immutable.Pop(); elem != nil; elem = q.immutable.Pop() {
					go q.processImmutable(ctx, elem)
				}

				// processing mutable queue
				for elem := q.Get(ctx); elem != nil; elem = q.Get(ctx) {
					err = q.processFunc(ctx, elem, q.processFuncArgs)
					if err == ErrRetryLater {
						q.Prepend(ctx, false, elem)
						break
					} else if err != nil {
						logger.Error("Failed to process transcript:", err)
					}
					q.finished.Push(elem)
				}
			}

			q.processLock.Lock()
			q.processActive = false
			q.processLock.Unlock()
		}()

		q.processActive = true
	}
	q.processLock.Unlock()
}

// TODO: Probably should be rewritten (all rotation in one place)
func (q *ExecutionBroker) TakeAndRotate(count int) ([]*Transcript, []*Transcript) {
	q.mutableLock.Lock()
	head := q.mutable.Take(count)
	tail := q.mutable.Rotate()
	q.mutableLock.Unlock()
	return head, tail
}

func (q *ExecutionBroker) RotateImmutable() []*Transcript {
	return q.immutable.Rotate()
}

func (q *ExecutionBroker) RotateFinished() []*Transcript {
	return q.finished.Rotate()
}

func NewExecutionBroker(checkFunc CheckCallback, execFunc ExecuteTranscriptCallback, args interface{}) *ExecutionBroker {
	return &ExecutionBroker{
		mutableLock: sync.RWMutex{},
		mutable:     NewTranscriptDequeue(),
		immutable:   NewTranscriptDequeue(),
		finished:    NewTranscriptDequeue(),

		checkFunc:       checkFunc,
		processFunc:     execFunc,
		processFuncArgs: args,
	}
}
