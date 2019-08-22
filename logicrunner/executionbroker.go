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
	"sync/atomic"

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/executionregistry"
	"github.com/insolar/insolar/logicrunner/requestsqueue"
)

// passToNextLimit - number of requests we pass to next executor on pulse change,
// the rest it should fetch of the ledger
const passToNextLimit = 10

// prefetchLimit - when we reach this number of requests in queue, we start
// pre-fetching requests from ledger
const prefetchLimit = 3

//go:generate minimock -i github.com/insolar/insolar/logicrunner.ExecutionBrokerI -o ./ -s _mock.go -g

type ExecutionBrokerI interface {
	AddFreshRequest(ctx context.Context, transcript *common.Transcript)
	AddRequestsFromPrevExecutor(ctx context.Context, transcripts ...*common.Transcript)
	AddRequestsFromLedger(ctx context.Context, transcripts ...*common.Transcript)
	AddAdditionalRequestFromPrevExecutor(ctx context.Context, transcript *common.Transcript)

	PendingState() insolar.PendingState
	PrevExecutorStillExecuting(ctx context.Context)
	PrevExecutorPendingResult(ctx context.Context, prevExecState insolar.PendingState)
	PrevExecutorSentPendingFinished(ctx context.Context) error
	SetNotPending(ctx context.Context)

	IsKnownRequest(ctx context.Context, req insolar.Reference) bool

	AbandonedRequestsOnLedger(ctx context.Context)
	MoreRequestsOnLedger(ctx context.Context)
	NoMoreRequestsOnLedger(ctx context.Context)

	OnPulse(ctx context.Context) []payload.Payload
}

type ExecutionBroker struct {
	Ref insolar.Reference

	stateLock sync.Mutex

	mutable   requestsqueue.RequestsQueue
	immutable requestsqueue.RequestsQueue
	finished  []*common.Transcript

	outgoingSender OutgoingRequestSender

	executionRegistry executionregistry.ExecutionRegistry
	requestsFetcher   RequestsFetcher

	pulseAccessor pulse.Accessor

	publisher        watermillMsg.Publisher
	sender           bus.Sender
	requestsExecutor RequestsExecutor
	artifactsManager artifacts.Client

	pending              insolar.PendingState
	PendingConfirmed     bool
	HasPendingCheckMutex sync.Mutex

	ledgerHasMoreRequests bool

	processorActive uint32

	deduplicationTable map[insolar.Reference]bool
}

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
	return &ExecutionBroker{
		Ref: ref,

		mutable:   requestsqueue.New(),
		immutable: requestsqueue.New(),

		outgoingSender: outgoingSender,
		pulseAccessor:  pulseAccessor,

		publisher:         publisher,
		requestsExecutor:  requestsExecutor,
		sender:            sender,
		artifactsManager:  artifactsManager,
		executionRegistry: executionRegistry,

		processorActive: 0,

		deduplicationTable: make(map[insolar.Reference]bool),
	}
}

func (q *ExecutionBroker) tryTakeProcessor(_ context.Context, immutable bool) bool {
	if immutable {
		return true
	}
	return atomic.CompareAndSwapUint32(&q.processorActive, 0, 1)
}

func (q *ExecutionBroker) releaseProcessor(_ context.Context, immutable bool) {
	if immutable {
		return
	}
	atomic.SwapUint32(&q.processorActive, 0)
}

func (q *ExecutionBroker) getTask(ctx context.Context, queue requestsqueue.RequestsQueue) *common.Transcript {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	transcript := queue.TakeFirst(ctx)
	if transcript == nil {
		return nil
	}

	q.executionRegistry.Register(ctx, transcript)

	return transcript
}

func (q *ExecutionBroker) finishTask(ctx context.Context, transcript *common.Transcript) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	logger := inslogger.FromContext(ctx)

	q.finished = append(q.finished, transcript)

	if q.executionRegistry.GetActiveTranscript(transcript.RequestRef) == nil {
		logger.Error("[ ExecutionBroker.FinishTask ] task wasn't executed")
	} else {
		q.executionRegistry.Done(transcript)
	}
}

func (q *ExecutionBroker) processTranscript(ctx context.Context, transcript *common.Transcript) {
	if transcript.Context != nil {
		ctx = transcript.Context
	} else {
		inslogger.FromContext(ctx).Error("context in transcript is nil")
	}

	ctx, logger := inslogger.WithField(ctx, "request", transcript.RequestRef.String())

	reply, err := q.requestsExecutor.ExecuteAndSave(ctx, transcript)
	if err != nil {
		logger.Warn("contract execution error: ", err)
	}

	q.finishTask(ctx, transcript)

	go q.requestsExecutor.SendReply(ctx, transcript, reply, err)

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
		if q.storeWithoutDuplication(ctx, transcript) {
			continue
		}

		inslogger.FromContext(transcript.Context).Debug("appending request to queue")

		var list requestsqueue.RequestsQueue
		if transcript.Request.Immutable {
			list = q.immutable
		} else {
			list = q.mutable
		}
		list.Append(ctx, source, transcript)
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

func (q *ExecutionBroker) commonStartProcessor(ctx context.Context, immutable bool) {
	defer q.releaseProcessor(ctx, immutable)

	q.clarifyPendingStateFromLedger(ctx)

	// сhecking we're eligible to execute contracts
	if readyToExecute := q.Check(ctx); !readyToExecute {
		return
	}

	q.fetchMoreFromLedgerIfNeeded(ctx)

	if immutable {
		for elem := q.getTask(ctx, q.immutable); elem != nil; elem = q.getTask(ctx, q.immutable) {
			go q.processTranscript(ctx, elem)
		}
	} else {
		for elem := q.getTask(ctx, q.mutable); elem != nil; elem = q.getTask(ctx, q.mutable) {
			q.processTranscript(ctx, elem)
		}
	}
}

func getQueueName(immutable bool) string {
	if immutable {
		return "immutable"
	}
	return "mutable"
}

// StartProcessorIfNeeded processes queue messages in strict order (flag determines which
// one, mutable or immutable)
// We need to start manually execution broker only if we were in pending and now we're not.
func (q *ExecutionBroker) StartProcessorIfNeeded(ctx context.Context, immutable bool) {
	logger := inslogger.FromContext(ctx)

	if q.tryTakeProcessor(ctx, immutable) {
		logger.Info("[ StartProcessorIfNeeded ] Starting a new ", getQueueName(immutable), " queue processor")
		go q.commonStartProcessor(ctx, immutable)
	}
}

func (q *ExecutionBroker) StartProcessorsIfNeeded(ctx context.Context) {
	q.StartProcessorIfNeeded(ctx, false)
	q.StartProcessorIfNeeded(ctx, true)
}

func (q *ExecutionBroker) fetchMoreFromLedgerIfNeeded(ctx context.Context) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if !q.ledgerHasMoreRequests {
		return
	}

	if q.mutable.NumberOfOld(ctx)+q.immutable.NumberOfOld(ctx) > prefetchLimit {
		return
	}

	q.startRequestsFetcher(ctx)
}

func (q *ExecutionBroker) Check(ctx context.Context) bool {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	logger := inslogger.FromContext(ctx)

	// check pending state of execution (whether we can process task or not)
	if q.pending == insolar.PendingUnknown {
		logger.Debug("One shouldn't call ExecuteTranscript in case when pending state is unknown")
		return false
	} else if q.pending == insolar.InPending {
		logger.Debug("Object in pending, wont start queue processor")
		return false
	}

	return true
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

	inslogger.FromContext(ctx).Debug("pending finished")
	q.pending = insolar.NotPending
	q.PendingConfirmed = false

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
	logger := inslogger.FromContext(ctx)
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	defer func() {
		// clean everything, just in case
		q.mutable.Clean(ctx)
		q.immutable.Clean(ctx)
		q.finished = nil
		q.deduplicationTable = make(map[insolar.Reference]bool)
	}()

	q.stopRequestsFetcher(ctx)

	sendExecResults := false

	requests, hasMore := requestsqueue.FirstNFromMany(ctx, passToNextLimit, q.mutable, q.immutable)

	switch {
	case q.isActive():
		q.pending = insolar.InPending
		sendExecResults = true
	case q.notConfirmedPending():
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
		messagesQueue := convertQueueToMessageQueue(ctx, requests)
		ledgerHasMoreRequests := q.ledgerHasMoreRequests || hasMore

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

	select {
	case <-ctx.Done():
		inslogger.FromContext(ctx).Debug("pulse changed, skipping")
		return
	default:
	}

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

			q.pending = insolar.NotPending
			q.PendingConfirmed = false
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

func (q *ExecutionBroker) PrevExecutorStillExecuting(ctx context.Context) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	logger := inslogger.FromContext(ctx)
	logger.Debugf("got StillExecuting from previous executor")

	switch q.pending {
	case insolar.NotPending:
		// It might be when StillExecuting comes after PendingFinished
		logger.Warn("got StillExecuting message, but our state says that it's not in pending")
	case insolar.InPending:
		q.PendingConfirmed = true
	case insolar.PendingUnknown:
		// we are first, strange, soon ExecuteResults message should come
		q.pending = insolar.InPending
		q.PendingConfirmed = true
	}
}

func (q *ExecutionBroker) PrevExecutorSentPendingFinished(ctx context.Context) error {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if q.isActive() {
		return errors.New("already executing")
	}

	q.pending = insolar.NotPending
	q.StartProcessorsIfNeeded(ctx)

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
	q.StartProcessorsIfNeeded(ctx)
}

func (q *ExecutionBroker) AddRequestsFromPrevExecutor(ctx context.Context, transcripts ...*common.Transcript) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	q.add(ctx, requestsqueue.FromPreviousExecutor, transcripts...)
	q.StartProcessorsIfNeeded(ctx)
}

func (q *ExecutionBroker) AddRequestsFromLedger(ctx context.Context, transcripts ...*common.Transcript) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	q.add(ctx, requestsqueue.FromLedger, transcripts...)
	q.StartProcessorsIfNeeded(ctx)
}

func (q *ExecutionBroker) AddAdditionalRequestFromPrevExecutor(
	ctx context.Context, tr *common.Transcript,
) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	q.add(ctx, requestsqueue.FromPreviousExecutor, tr)
	q.StartProcessorsIfNeeded(ctx)
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
