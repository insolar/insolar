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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
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

//go:generate minimock -i github.com/insolar/insolar/logicrunner.ExecutionBrokerMethods -o ./ -s _mock.go
type ExecutionBrokerMethods interface {
	Check(context.Context) error
	Execute(context.Context, *Transcript) error
}

type ExecutionBroker struct {
	stateLock   sync.Mutex
	mutable     *TranscriptDequeue
	immutable   *TranscriptDequeue
	finished    *TranscriptDequeue
	errored     *TranscriptDequeue
	currentList *CurrentExecutionList

	methods ExecutionBrokerMethods

	processActive bool
	processLock   sync.Mutex

	deduplicationTable map[insolar.Reference]bool
	deduplicationLock  sync.Mutex
}

var ErrRetryLater = errors.New("Failed to start task, retry next time")

type ExecutionBrokerRotationResult struct {
	Requests              []*Transcript
	Finished              []*Transcript
	LedgerHasMoreRequests bool
}

func (q *ExecutionBroker) processImmutable(ctx context.Context, transcript *Transcript) {
	defer q.releaseTask(ctx, transcript)

	logger := inslogger.FromContext(ctx).WithField("RequestReference", transcript.RequestRef)

	if err := q.methods.Check(ctx); err == ErrRetryLater {
		return
	} else if err != nil {
		logger.Error("[ processImmutable ] check function returned error:", err)
		return
	}

	err := q.methods.Execute(ctx, transcript)
	if err == ErrRetryLater {
		return
	}

	q.finishTask(ctx, transcript, err != nil)
	if err != nil {
		logger.Error("[ processImmutable ] Failed to process immutable Transcript:", err)
	}
	return
}

func (q *ExecutionBroker) processMutable(ctx context.Context, transcript *Transcript) {
	defer q.releaseTask(ctx, transcript)

	logger := inslogger.FromContext(ctx).WithField("RequestReference", transcript.RequestRef)

	err := q.methods.Execute(ctx, transcript)
	if err == ErrRetryLater {
		return
	}

	q.finishTask(ctx, transcript, err != nil)
	if err != nil {
		logger.Error("[ processMutable ] Failed to process immutable Transcript:", err)
	}
	return
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
			q.storeCurrent(ctx, transcript)
			go q.processImmutable(ctx, transcript)
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
		if q.isDuplicate(ctx, transcript) {
			continue
		}

		if transcript.Request.Immutable {
			q.storeCurrent(ctx, transcript)
			go q.processImmutable(ctx, transcript)
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

func (q *ExecutionBroker) getImmutableTask(ctx context.Context) *Transcript {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	transcript := q.mutable.Pop()
	if transcript == nil {
		return nil
	}
	inslogger.FromContext(ctx).Error("put immutable: ", transcript.RequestRef.String())
	q.currentList.Set(*transcript.RequestRef, transcript)

	return transcript
}

func (q *ExecutionBroker) getMutableTask(ctx context.Context) *Transcript {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	transcript := q.mutable.Pop()
	if transcript == nil {
		return nil
	}
	inslogger.FromContext(ctx).Error("put mutable: ", transcript.RequestRef.String())
	q.currentList.Set(*transcript.RequestRef, transcript)

	return transcript
}

func (q *ExecutionBroker) finishTask(ctx context.Context, transcript *Transcript, isErrored bool) {
	logger := inslogger.FromContext(ctx)

	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if !isErrored {
		q.finished.Push(transcript)
	} else {
		q.errored.Push(transcript)
	}

	if !q.currentList.Has(*transcript.RequestRef) {
		logger.Error(transcript.RequestRef.String())
		logger.Error("[ ExecutionBroker.FinishTask ] task is not in current")
	} else {
		q.currentList.Delete(*transcript.RequestRef)
	}
}

func (q *ExecutionBroker) releaseTask(ctx context.Context, transcript *Transcript) {
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

func (q *ExecutionBroker) storeCurrent(ctx context.Context, transcript *Transcript) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	q.currentList.Set(*transcript.RequestRef, transcript)
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

	return q.mutable.PopByReference(r)
}

func (q *ExecutionBroker) startProcessor(ctx context.Context) {
	logger := inslogger.FromContext(ctx)

	defer func() {
		q.processLock.Lock()
		q.processActive = false
		q.processLock.Unlock()
	}()

	// сhecking we're eligible to execute contracts
	if err := q.methods.Check(ctx); err == ErrRetryLater {
		return
	} else if err != nil {
		logger.Error("[ processImmutable ] check function returned error:", err)
		return
	}

	// processing immutable queue (it can appear if we were in pending state)
	// run simultaneously all immutable transcripts and forget about them
	for elem := q.getImmutableTask(ctx); elem != nil; elem = q.getImmutableTask(ctx) {
		go q.processImmutable(ctx, elem)
	}

	// processing mutable queue
	for transcript := q.getMutableTask(ctx); transcript != nil; transcript = q.getMutableTask(ctx) {
		logger := logger.WithField("RequestReference", transcript.RequestRef)

		err := q.methods.Execute(ctx, transcript)
		if err == ErrRetryLater {
			q.releaseTask(ctx, transcript)
			break
		}

		q.finishTask(ctx, transcript, err != nil)
		if err != nil {
			logger.Error("[ StartProcessorIfNeeded ] Failed to process transcript:", err)
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

// TODO: probably rotation should wait till processActive == false (??)
func (q *ExecutionBroker) Rotate(count int) *ExecutionBrokerRotationResult {
	// take mutables, then, if we can, take immutables, if something was left -
	q.stateLock.Lock()
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
	_ = q.errored.Rotate()

	q.deduplicationLock.Lock()
	q.deduplicationTable = make(map[insolar.Reference]bool)
	q.deduplicationLock.Unlock()

	q.stateLock.Unlock()

	return rv
}

func NewExecutionBroker(methods ExecutionBrokerMethods) *ExecutionBroker {
	return &ExecutionBroker{
		stateLock:   sync.Mutex{},
		mutable:     NewTranscriptDequeue(),
		immutable:   NewTranscriptDequeue(),
		finished:    NewTranscriptDequeue(),
		errored:     NewTranscriptDequeue(),
		currentList: NewCurrentExecutionList(),

		methods: methods,

		deduplicationLock:  sync.Mutex{},
		deduplicationTable: make(map[insolar.Reference]bool),
	}
}
