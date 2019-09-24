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

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/executionregistry"
	"github.com/insolar/insolar/logicrunner/metrics"
	"github.com/insolar/insolar/logicrunner/requestsqueue"
)

var (
	ErrNotInPending     = errors.New("state is not in pending")
	ErrAlreadyExecuting = errors.New("already executing")
)

// passToNextLimit - number of requests we pass to next executor on pulse change,
// the rest it should fetch of the ledger
const passToNextLimit = 50

// prefetchLimit - when we reach this number of requests in queue, we start
// pre-fetching requests from ledger
const prefetchLimit = 0

const mutableExecutionLimit = 1
const immutableExecutionLimit = 30

//go:generate minimock -i github.com/insolar/insolar/logicrunner.ExecutionBrokerI -o ./ -s _mock.go -g

type ExecutionBrokerI interface {
	AddFreshRequest(ctx context.Context, transcript *common.Transcript)
	AddRequestsFromPrevExecutor(ctx context.Context, transcripts ...*common.Transcript)
	AddRequestsFromLedger(ctx context.Context, transcripts ...*common.Transcript)
	AddAdditionalRequestFromPrevExecutor(ctx context.Context, transcript *common.Transcript)

	PendingState() insolar.PendingState
	PrevExecutorStillExecuting(ctx context.Context) error
	PrevExecutorPendingResult(ctx context.Context, prevExecState insolar.PendingState)
	PrevExecutorSentPendingFinished(ctx context.Context) error
	SetNotPending(ctx context.Context)

	IsKnownRequest(ctx context.Context, req insolar.Reference) bool

	AbandonedRequestsOnLedger(ctx context.Context)
	MoreRequestsOnLedger(ctx context.Context)
	NoMoreRequestsOnLedger(ctx context.Context)

	OnPulse(ctx context.Context) []payload.Payload
}

type processorState struct {
	name         string
	queue        requestsqueue.RequestsQueue
	workers      int
	workersLimit int
}

// dismissWorker frees worker (decrements currently used worker count)
// one should be certain that he calls that after engageWorker and only once
// that function should be called under ExecutionBroker lock
func (ps *processorState) dismissWorker() {
	if ps.workers <= 0 {
		panic("unreachable")
	}
	ps.workers--
}

// engageWorker occupies worker (increments currently used worker count)
// one should check count of available workers before calling that function
// that function should be called under ExecutionBroker lock
func (ps *processorState) engageWorker() {
	if ps.workers >= ps.workersLimit {
		panic("unreachable")
	}
	ps.workers++
}

// availableWorkers returns number of available workers in the worker queue
// that function should be called under ExecutionBroker lock
func (ps processorState) availableWorkers() int {
	return ps.workersLimit - ps.workers
}

type ExecutionBroker struct {
	Ref  insolar.Reference
	name string

	stateLock sync.Mutex

	mutable   processorState
	immutable processorState

	finished []*common.Transcript

	outgoingSender OutgoingRequestSender

	executionRegistry executionregistry.ExecutionRegistry
	requestsFetcher   RequestFetcher

	pulseAccessor pulse.Accessor

	publisher        watermillMsg.Publisher
	sender           bus.Sender
	requestsExecutor RequestsExecutor
	artifactsManager artifacts.Client

	pending              insolar.PendingState
	PendingConfirmed     bool
	HasPendingCheckMutex sync.Mutex

	ledgerHasMoreRequests bool

	deduplicationTable map[insolar.Reference]bool
}

const (
	ImmutableQueueName = "immutable"
	MutableQueueName   = "mutable"
)

func NewExecutionBroker(
	ref insolar.Reference,
	publisher watermillMsg.Publisher,
	requestsExecutor RequestsExecutor,
	sender bus.Sender,
	artifactsManager artifacts.Client,
	executionRegistry executionregistry.ExecutionRegistry,
	outgoingSender OutgoingRequestSender,
	pulseAccessor pulse.Accessor,
) *ExecutionBroker {
	pulseObject, err := pulseAccessor.Latest(context.Background())
	if err != nil {
		log.Error("failed to create execution broker ", err.Error())
		return nil
	}
	return &ExecutionBroker{
		Ref:  ref,
		name: "executionbroker-" + pulseObject.PulseNumber.String(),

		mutable: processorState{
			name:         MutableQueueName,
			workers:      0,
			workersLimit: mutableExecutionLimit,
			queue:        requestsqueue.New(),
		},
		immutable: processorState{
			name:         ImmutableQueueName,
			workers:      0,
			workersLimit: immutableExecutionLimit,
			queue:        requestsqueue.New(),
		},

		outgoingSender: outgoingSender,
		pulseAccessor:  pulseAccessor,

		publisher:         publisher,
		requestsExecutor:  requestsExecutor,
		sender:            sender,
		artifactsManager:  artifactsManager,
		executionRegistry: executionRegistry,

		deduplicationTable: make(map[insolar.Reference]bool),
	}
}

func (q *ExecutionBroker) getTask(ctx context.Context, state *processorState) *common.Transcript {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	for {
		transcript := state.queue.TakeFirst(ctx)
		if transcript == nil {
			// when we finished processing queue (nothing left in the queue) - while we under lock,
			// release (we've finished processing everything and we should finish current executor)
			state.dismissWorker()
			return nil
		}

		err := q.executionRegistry.Register(ctx, transcript)
		if err != nil {
			if err == executionregistry.ErrAlreadyRegistered {
				stats.Record(ctx, metrics.ExecutionBrokerTranscriptAlreadyRegistered.M(1))
			}
			inslogger.FromContext(ctx).Error("couldn't register transcript, skipping: ", err.Error())
			continue
		}
		return transcript
	}
}

func (q *ExecutionBroker) finishTask(ctx context.Context, transcript *common.Transcript) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	logger := inslogger.FromContext(ctx)

	q.finished = append(q.finished, transcript)

	done := q.executionRegistry.Done(transcript)
	if !done {
		logger.Error("task wasn't in the registry, very bad")
	}
}

func (q *ExecutionBroker) cleanMutableQueue(ctx context.Context) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	q.mutable.queue.Clean(ctx)
	q.deduplicationTable = make(map[insolar.Reference]bool)
}

func (q *ExecutionBroker) processTranscript(ctx context.Context, transcript *common.Transcript) {
	stats.Record(ctx, metrics.ExecutionBrokerExecutionStarted.M(1))
	defer stats.Record(ctx, metrics.ExecutionBrokerExecutionFinished.M(1))

	if transcript.Context != nil {
		ctx = transcript.Context
	} else {
		inslogger.FromContext(ctx).Error("context in transcript is nil")
	}

	ctx, logger := inslogger.WithFields(ctx, map[string]interface{}{
		"request": transcript.RequestRef.String(),
		"broker":  q.name,
	})

	var (
		replyData insolar.Reply
		err       error
	)

	if transcript.Request.CallType == record.CTMethod {
		logger.Info("processing transcript with method")
		var objDesc artifacts.ObjectDescriptor
		objDesc, err = q.artifactsManager.GetObject(ctx, *transcript.Request.Object, &transcript.RequestRef)
		if err == nil {
			transcript.ObjectDescriptor = objDesc

			if !transcript.Request.Immutable &&
				transcript.ObjectDescriptor.EarliestRequestID() != nil &&
				!transcript.RequestRef.GetLocal().Equal(*transcript.ObjectDescriptor.EarliestRequestID()) {
				logger.Info("Got different earliest request from ledger")

				q.cleanMutableQueue(ctx)
				q.ledgerHasMoreRequests = true
			}
		}
	}

	if err == nil {
		var result artifacts.RequestResult
		result, err = q.requestsExecutor.ExecuteAndSave(ctx, transcript)
		if err != nil {
			logger.Warn("contract execution error: ", err)
		} else {
			objRef := result.ObjectReference()
			replyData = &reply.CallMethod{Result: result.Result(), Object: &objRef}
		}
	}

	q.finishTask(ctx, transcript)

	go q.requestsExecutor.SendReply(ctx, transcript.RequestRef, *transcript.Request, replyData, err)

	// we're checking here that pulse was changed and we should send
	// a message that we've finished processing tasks
	// note: ideally we should tell here that we've stopped executing
	//       but we only hoped that OnPulse had already told us that
	//       pulse changed and we should stop execution
	logger.Debug("finished request, try to finish pending if needed")
	q.finishPendingIfNeeded(ctx)
}

func (q *ExecutionBroker) storeWithoutDuplication(ctx context.Context, transcript *common.Transcript) bool {
	if _, ok := q.deduplicationTable[transcript.RequestRef]; ok {
		logger := inslogger.FromContext(ctx)
		logger.Infof("Already know about request %s, skipping", transcript.RequestRef.String())

		return true
	}
	q.deduplicationTable[transcript.RequestRef] = true
	return false
}

// One shouldn't mix immutable calls and mutable ones
func (q *ExecutionBroker) add(
	ctx context.Context, source requestsqueue.RequestSource, transcripts ...*common.Transcript,
) {
	for _, transcript := range transcripts {
		queueName := "mutable"
		if transcript.Request.Immutable {
			queueName = "immutable"
		}
		ctx = insmetrics.InsertTag(ctx, metrics.TagExecutionQueueName, queueName)
		if q.storeWithoutDuplication(ctx, transcript) {
			stats.Record(ctx, metrics.ExecutionBrokerTranscriptDuplicate.M(1))
			continue
		}
		if q.executionRegistry.GetActiveTranscript(transcript.RequestRef) != nil {
			stats.Record(ctx, metrics.ExecutionBrokerTranscriptExecuting.M(1))
			inslogger.FromContext(transcript.Context).Warn(
				"this node already executing request, won't add to queue",
			)
			continue
		}

		stats.Record(ctx, metrics.ExecutionBrokerTranscriptRegistered.M(1))

		inslogger.FromContext(transcript.Context).Debug("appending request to queue")

		var list requestsqueue.RequestsQueue
		if transcript.Request.Immutable {
			list = q.immutable.queue
		} else {
			list = q.mutable.queue
		}
		list.Append(ctx, source, transcript)

		inslogger.FromContext(transcript.Context).Debug("request appended to queue")
	}
}

func (q *ExecutionBroker) IsKnownRequest(ctx context.Context, req insolar.Reference) bool {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if _, ok := q.deduplicationTable[req]; ok {
		return true
	}
	return false
}

// startProcessors starts independent processing of mutable and immutable queues
// should be started only under q.stateLock
func (q *ExecutionBroker) startProcessors(ctx context.Context) {
	q.startProcessor(ctx, &q.mutable)

	for i := q.startProcessorCount(ctx, &q.immutable); i > 0; i-- {
		q.startProcessor(ctx, &q.immutable)
	}
}

func (q *ExecutionBroker) startProcessorCount(ctx context.Context, state *processorState) int {
	queueLength := state.queue.Length()
	workerAvailable := state.availableWorkers()

	if queueLength > workerAvailable {
		return workerAvailable
	}
	return queueLength
}

// startProcessor starts processing of queue ensuring that only one processor is active
// at the moment.
func (q *ExecutionBroker) startProcessor(ctx context.Context, state *processorState) bool {
	if state.queue.Length() == 0 {
		return false
	} else if state.availableWorkers() == 0 {
		return false
	}

	state.engageWorker()

	go q.processQueue(ctx, state)
	return true
}

func (q *ExecutionBroker) processQueue(ctx context.Context, state *processorState) {
	q.processQueueInner(ctx, state)
}

func (q *ExecutionBroker) processQueueInner(ctx context.Context, state *processorState) {
	q.clarifyPendingStateFromLedger(ctx)

	ctx, logger := inslogger.WithField(ctx, "processorType", state.name)

	ps := q.PendingState()
	if ps != insolar.NotPending {
		logger.Debug("wont process, pending state is ", ps)

		q.stateLock.Lock()
		state.dismissWorker()
		q.stateLock.Unlock()

		return
	}

	logger.Info("started a new queue processor")

	for elem := q.getTask(ctx, state); elem != nil; elem = q.getTask(ctx, state) {
		q.processTranscript(ctx, elem)
		q.fetchMoreFromLedgerIfNeeded(ctx)
	}

	logger.Debug("queue processor has finished")

	// we shouldn't release lock here, since we're releasing it in getTask function
	// buf it's important to finish ASAP (and do nothing with data) after we exit
	// the for loop.
	// otherwise we should think that multiple queue processor can exist and 'fight'
	// for data
}

func (q *ExecutionBroker) fetchMoreFromLedgerIfNeeded(ctx context.Context) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if !q.ledgerHasMoreRequests {
		return
	}

	if q.mutable.queue.NumberOfOld(ctx)+q.immutable.queue.NumberOfOld(ctx) > prefetchLimit {
		return
	}

	q.startRequestsFetcher(ctx)
}

// finishPendingIfNeeded checks whether last execution was a pending one.
// If this is true as a side effect the function sends a PendingFinished
// message to the current executor
func (q *ExecutionBroker) finishPendingIfNeeded(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if q.pending != insolar.InPending {
		logger.Debug("we aren't in pending")
		return
	}

	// we process mutable and immutable calls in parallel
	// and use one pending state for all of them
	// so pending is finished only when all calls are finished
	if !q.executionRegistry.IsEmpty() {
		count := q.executionRegistry.Length()
		inslogger.FromContext(ctx).Debug("we are in pending and still have ", count, " requests to finish")
		return
	}

	logger.Debug("pending finished")

	pendingMsg, err := payload.NewMessage(&payload.PendingFinished{
		ObjectRef: q.Ref,
	})
	if err != nil {
		logger.Error(errors.Wrap(err, "finishPending: Unable to create PendingFinished message"))
		return
	}

	// ensure OK response because we might catch flow cancelled
	waitOKSender := bus.NewWaitOKWithRetrySender(q.sender, q.pulseAccessor, 1)
	waitOKSender.SendRole(ctx, pendingMsg, insolar.DynamicRoleVirtualExecutor, q.Ref)
}

func (q *ExecutionBroker) OnPulse(ctx context.Context) []payload.Payload {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	logger := inslogger.FromContext(ctx)

	onPulseStart := time.Now()
	defer func(ctx context.Context) {
		stats.Record(ctx, metrics.ExecutionBrokerOnPulseTiming.M(float64(time.Since(onPulseStart).Nanoseconds())/1e6))
	}(ctx)

	defer func() {
		// clean everything, just in case
		q.mutable.queue.Clean(ctx)
		q.immutable.queue.Clean(ctx)
		q.finished = nil
		q.deduplicationTable = make(map[insolar.Reference]bool)

		q.pending = insolar.InPending
		q.PendingConfirmed = true
	}()

	q.stopRequestsFetcher(ctx)

	sendExecResults := false

	requests, rest := requestsqueue.SplitNFromMany(ctx, passToNextLimit, q.mutable.queue, q.immutable.queue)
	if len(rest) > 0 {
		stats.Record(ctx, metrics.ExecutionBrokerTruncatedRequests.M(int64(len(rest))))
	}

	switch {
	case q.isActive():
		q.pending = insolar.InPending
		sendExecResults = true
	case q.notConfirmedPending():
		stats.Record(ctx, metrics.ExecutionBrokerOnPulseNotConfirmed.M(1))
		logger.Warn("looks like pending executor died, continuing execution on next executor")
		q.pending = insolar.NotPending
		sendExecResults = true
		q.ledgerHasMoreRequests = true
	case len(q.finished) > 0 || len(requests) > 0:
		sendExecResults = true
	}

	messages := make([]payload.Payload, 0)

	if sendExecResults {
		// TODO: we also should send when executed something for validation
		// TODO: now validation is disabled
		messagesQueue := common.ConvertQueueToMessageQueue(ctx, requests)
		ledgerHasMoreRequests := q.ledgerHasMoreRequests || len(rest) > 0

		messages = append(messages, &payload.ExecutorResults{
			RecordRef:             q.Ref,
			Pending:               q.pending,
			Queue:                 messagesQueue,
			LedgerHasMoreRequests: ledgerHasMoreRequests,
		})
	}

	return messages
}

// notConfirmedPending checks that we were in pending and waiting
// for previous executor to notify us that he still executes it
func (q *ExecutionBroker) notConfirmedPending() bool {
	return q.pending == insolar.InPending && !q.PendingConfirmed
}

func (q *ExecutionBroker) NoMoreRequestsOnLedger(ctx context.Context) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	q.ledgerHasMoreRequests = false
	q.stopRequestsFetcher(ctx)
}

func (q *ExecutionBroker) PendingState() insolar.PendingState {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	return q.pending
}

func (q *ExecutionBroker) PrevExecutorPendingResult(ctx context.Context, prevExecState insolar.PendingState) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	logger := inslogger.FromContext(ctx)

	switch q.pending {
	case insolar.InPending:
		if q.isActive() {
			logger.Debug("execution returned to node that is still executing pending")
		} else if prevExecState == insolar.NotPending {
			logger.Debug("executor we came to thinks that execution pending, but previous said to continue")

			q.pending = insolar.NotPending
		}
	case insolar.PendingUnknown:
		q.pending = prevExecState
		logger.Debug("pending state was unknown, setting from previous executor to ", q.pending)
	default:
		logger.Debug("keeping pending state at ", q.pending)
	}
}

func (q *ExecutionBroker) PrevExecutorStillExecuting(ctx context.Context) error {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	logger := inslogger.FromContext(ctx)
	logger.Debugf("got StillExecuting from previous executor")

	switch q.pending {
	case insolar.NotPending:
		// It might be when StillExecuting comes after PendingFinished
		return ErrNotInPending
	case insolar.InPending:
		q.PendingConfirmed = true
	case insolar.PendingUnknown:
		// we are first, strange, soon ExecuteResults message should come
		q.pending = insolar.InPending
		q.PendingConfirmed = true
	}
	return nil
}

func (q *ExecutionBroker) PrevExecutorSentPendingFinished(ctx context.Context) error {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if q.isActive() {
		return ErrAlreadyExecuting
	}

	q.pending = insolar.NotPending
	q.startProcessors(ctx)

	return nil
}

func (q *ExecutionBroker) SetNotPending(ctx context.Context) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	q.pending = insolar.NotPending
}

func (q *ExecutionBroker) AbandonedRequestsOnLedger(ctx context.Context) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if q.pending == insolar.PendingUnknown {
		q.pending = insolar.InPending
		q.PendingConfirmed = false
	}

	q.ledgerHasMoreRequests = true
	q.startRequestsFetcher(ctx)
}

func (q *ExecutionBroker) MoreRequestsOnLedger(ctx context.Context) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	q.ledgerHasMoreRequests = true
}

func (q *ExecutionBroker) startRequestsFetcher(ctx context.Context) {
	if q.requestsFetcher == nil {
		q.requestsFetcher = NewRequestsFetcher(q.Ref, q.artifactsManager, q, q.outgoingSender)
	}
	q.requestsFetcher.FetchPendings(ctx)
}

func (q *ExecutionBroker) stopRequestsFetcher(ctx context.Context) {
	if q.requestsFetcher != nil {
		q.requestsFetcher.Abort(ctx)
		q.requestsFetcher = nil
	}
}

func (q *ExecutionBroker) AddFreshRequest(
	ctx context.Context, tr *common.Transcript,
) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if tr.Request.CallType != record.CTMethod {
		// It's considered that we are not pending except someone calls a method.
		q.pending = insolar.NotPending
	}

	q.add(ctx, requestsqueue.FromThisPulse, tr)
	q.startProcessors(ctx)
}

func (q *ExecutionBroker) AddRequestsFromPrevExecutor(ctx context.Context, transcripts ...*common.Transcript) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	q.add(ctx, requestsqueue.FromPreviousExecutor, transcripts...)
	q.startProcessors(ctx)
}

func (q *ExecutionBroker) AddRequestsFromLedger(ctx context.Context, transcripts ...*common.Transcript) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	q.add(ctx, requestsqueue.FromLedger, transcripts...)
	q.startProcessors(ctx)
}

func (q *ExecutionBroker) AddAdditionalRequestFromPrevExecutor(
	ctx context.Context, tr *common.Transcript,
) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	q.add(ctx, requestsqueue.FromPreviousExecutor, tr)
	q.startProcessors(ctx)
}

func (q *ExecutionBroker) isActive() bool {
	return !q.executionRegistry.IsEmpty()
}

func (q *ExecutionBroker) clarifyPendingStateFromLedger(ctx context.Context) {
	if q.PendingState() != insolar.PendingUnknown {
		return
	}

	q.HasPendingCheckMutex.Lock()
	defer q.HasPendingCheckMutex.Unlock()

	if q.PendingState() != insolar.PendingUnknown {
		return
	}

	has, err := q.artifactsManager.HasPendings(ctx, q.Ref)
	if err != nil {
		inslogger.FromContext(ctx).Error("couldn't check pending state: ", err.Error())
		return
	}

	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if q.pending == insolar.PendingUnknown {
		if has {
			q.pending = insolar.InPending
		} else {
			q.pending = insolar.NotPending
		}
	}
}
