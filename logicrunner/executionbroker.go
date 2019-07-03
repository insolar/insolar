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

type ExecutionBrokerMethods interface {
	Check(context.Context) error
	ExecuteTranscript(context.Context, *Transcript) error
}

type ExecutionBroker struct {
	mutableLock sync.RWMutex
	mutable     *TranscriptDequeue
	immutable   *TranscriptDequeue
	finished    *TranscriptDequeue

	methods ExecutionBrokerMethods

	processActive bool
	processLock   sync.Mutex

	StartProcessorIfNeededCount int32

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
	logger := inslogger.FromContext(ctx).WithField("RequestReference", transcript.RequestRef)

	if err := q.methods.Check(ctx); err == nil {
		if err := q.methods.ExecuteTranscript(ctx, transcript); err == nil {
			// In case when we're in pending - we should store it to execute later
			q.finished.Push(transcript)
			return
		} else if err != ErrRetryLater {
			logger.Error("[ processImmutable ] Failed to process immutable Transcript:", err)
			return
		}
	} else if err != ErrRetryLater {
		logger.Error("[ processImmutable ] check function returned error:", err)
		return
	}

	// we're retrying this request later, we're in pending now
	q.immutable.Push(transcript)
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
		if q.isDuplicate(ctx, transcript) {
			continue
		}

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

	q.processLock.Lock()
	// i've removed "if we have tasks here"; we can be there in two cases:
	// 1) we've put a task into queue and automatically start processor
	// 2) we've explicitly ask broker to be here
	// both cases means we knew what we are doing and so it's just an
	// unneeded optimisation
	if !q.processActive {
		q.StartProcessorIfNeededCount++
		logger.Info("[ StartProcessorIfNeeded ] Starting a new queue processor")

		go func() {
			var err error

			// сhecking we're eligible to execute contracts
			if err = q.methods.Check(ctx); err == nil {
				// processing immutable queue (it can appear if we were in pending state)
				// run simultaneously all immutable transcripts and forget about them
				for elem := q.immutable.Pop(); elem != nil; elem = q.immutable.Pop() {
					go q.processImmutable(ctx, elem)
				}

				// processing mutable queue
				for transcript := q.Get(ctx); transcript != nil; transcript = q.Get(ctx) {
					logger := logger.WithField("RequestReference", transcript.RequestRef)

					err = q.methods.ExecuteTranscript(ctx, transcript)
					if err == ErrRetryLater {
						q.Prepend(ctx, false, transcript)
						break
					} else if err != nil {
						logger.Error("[ StartProcessorIfNeeded ] Failed to process transcript:", err)
					}

					q.finished.Push(transcript)
				}
			} else if err != ErrRetryLater {
				logger.Error("[ StartProcessorIfNeeded ] check function returned error:", err)
			}

			q.processLock.Lock()
			q.processActive = false
			q.processLock.Unlock()
		}()

		q.processActive = true
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

func NewExecutionBroker(methods ExecutionBrokerMethods) *ExecutionBroker {
	return &ExecutionBroker{
		mutableLock: sync.RWMutex{},
		mutable:     NewTranscriptDequeue(),
		immutable:   NewTranscriptDequeue(),
		finished:    NewTranscriptDequeue(),

		methods: methods,

		deduplicationLock:  sync.Mutex{},
		deduplicationTable: make(map[insolar.Reference]bool),
	}
}
