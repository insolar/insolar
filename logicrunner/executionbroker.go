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
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/executionregistry"
	"github.com/insolar/insolar/logicrunner/metrics"
)

var (
	ErrNotInPending     = errors.New("state is not in pending")
	ErrAlreadyExecuting = errors.New("already executing")
)

const immutableExecutionLimit = 30

//go:generate minimock -i github.com/insolar/insolar/logicrunner.ExecutionBrokerI -o ./ -s _mock.go -g

type ExecutionBrokerI interface {
	HasMoreRequests(ctx context.Context)

	AbandonedRequestsOnLedger(ctx context.Context)

	PendingState() insolar.PendingState
	PrevExecutorStillExecuting(ctx context.Context) error
	PrevExecutorPendingResult(ctx context.Context, prevExecState insolar.PendingState)
	PrevExecutorSentPendingFinished(ctx context.Context) error
	SetNotPending(ctx context.Context)

	OnPulse(ctx context.Context) []payload.Payload
}

type LedgerHasMore int
const (
	LedgerIsEmpty LedgerHasMore = iota + 1
	LedgerHasMoreKnown
	LedgerHasMoreUnknown
)

type ExecutionBroker struct {
	Ref  insolar.Reference
	name string

	stateLock sync.Mutex

	processorActive            bool
	closed                     chan struct{}
	ledgerHasMoreRequests      LedgerHasMore
	probablyMoreSinceLastFetch chan struct{}

	finished []*common.Transcript

	outgoingSender OutgoingRequestSender

	executionRegistry executionregistry.ExecutionRegistry

	pulseAccessor pulse.Accessor

	sender           bus.Sender
	requestsExecutor RequestsExecutor
	artifactsManager artifacts.Client

	pending              insolar.PendingState
	PendingConfirmed     bool
	HasPendingCheckMutex sync.Mutex

	deduplicationTable map[insolar.Reference]bool
}

func NewExecutionBroker(
	ref insolar.Reference,
	_ watermillMsg.Publisher,
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

		ledgerHasMoreRequests:      LedgerIsEmpty,
		closed:                     make(chan struct{}),
		probablyMoreSinceLastFetch: make(chan struct{}, 1),

		outgoingSender:    outgoingSender,
		pulseAccessor:     pulseAccessor,
		requestsExecutor:  requestsExecutor,
		sender:            sender,
		artifactsManager:  artifactsManager,
		executionRegistry: executionRegistry,

		deduplicationTable: make(map[insolar.Reference]bool),
	}
}

func (q *ExecutionBroker) HasMoreRequests(ctx context.Context) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	q.setHasMoreRequests()
	go q.startProcessor(ctx)
}

func (q *ExecutionBroker) AbandonedRequestsOnLedger(ctx context.Context) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if q.pending == insolar.PendingUnknown {
		q.pending = insolar.InPending
		q.PendingConfirmed = false
	}

	q.setHasMoreRequests()
	go q.startProcessor(ctx)
}

func (q *ExecutionBroker) noMoreRequestsOnLedgerInitial(ctx context.Context) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if len(q.probablyMoreSinceLastFetch) == 0 {
		inslogger.FromContext(ctx).Debug("marking that there is no more requests on ledger")
		q.ledgerHasMoreRequests = LedgerIsEmpty
	}
}

func (q *ExecutionBroker) noMoreRequestsOnLedger(ctx context.Context) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if len(q.probablyMoreSinceLastFetch) == 0 {
		inslogger.FromContext(ctx).Debug("marking that there is no more requests on ledger")
		q.ledgerHasMoreRequests = LedgerHasMoreKnown
	}
}


// startProcessors starts one processing goroutine
func (q *ExecutionBroker) startProcessor(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	logger.Debug("considering to start requests processor")

	q.clarifyPendingStateFromLedger(ctx)

	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if q.pending != insolar.NotPending {
		logger.Debug("object is in pending state, not processing")
		return
	}
	if q.processorActive {
		logger.Debug("processor is already active")
		return
	}
	if q.isClosed() {
		logger.Warn("broker is closed, pulse ended, refusing to start processor")
		return
	}

	q.processorActive = true

	// Clear flag if there were more requests before fetching started
	select {
	case <-q.probablyMoreSinceLastFetch:
	default:
	}

	logger.Debug("starting requests processor")

	fetcher := NewRequestsFetcher(q.Ref, q.artifactsManager, q.outgoingSender)
	feedMutable := make(chan *common.Transcript, 10)
	feedImmutable := make(chan *common.Transcript, 10)
	transcriptFeed := fetcher.FetchPendings(ctx)
	fetchedRequests, fetchIteration := 0, 0

	go func() {
		ctx, logger := inslogger.WithFields(context.Background(), map[string]interface{}{
			"broker": q.name,
			"object": q.Ref.String(),
		})

		defer q.stopProcessor(ctx, fetcher, feedMutable, feedImmutable)
		for {
			select {
			case tr, ok := <-transcriptFeed:
				if !ok {
					logger.Debug("fetcher stopped producing")

					select {
					case <-q.probablyMoreSinceLastFetch:
						logger.Debug("had request since last fetch, restarting fetcher")

						transcriptFeed = fetcher.FetchPendings(ctx)
						continue
					case <-q.closed:
						return
					}
				}
				if tr == nil {
					if fetchIteration == 0 && fetchedRequests == 0 {
						q.noMoreRequestsOnLedgerInitial(ctx)
					} else {
						q.noMoreRequestsOnLedger(ctx)
					}
					fetchIteration++

					continue
				}
				if q.upsertToDuplicationTable(ctx, tr) {
					continue
				}
				fetchedRequests++
				if tr.Request.Immutable {
					feedImmutable <- tr
				} else {
					feedMutable <- tr
				}
			case <-q.closed:
				return
			}
		}
	}()

	reader := func(feed chan *common.Transcript) {
		ctx, _ := inslogger.WithFields(context.Background(), map[string]interface{}{"broker": q.name})

		for tr := range feed {
			q.processTranscript(ctx, tr)
		}
	}

	for i := 0; i < immutableExecutionLimit; i++ {
		go reader(feedImmutable)
	}

	go reader(feedMutable)
}

func (q *ExecutionBroker) stopProcessor(ctx context.Context, fetcher RequestFetcher, feeds ...chan *common.Transcript) {
	inslogger.FromContext(ctx).Debug("broker stopped, stopping processor")
	for i := range feeds {
		close(feeds[i])
	}
	fetcher.Abort(ctx)
}

func (q *ExecutionBroker) processTranscript(ctx context.Context, transcript *common.Transcript) {
	stats.Record(ctx, metrics.ExecutionBrokerExecutionStarted.M(1))
	defer stats.Record(ctx, metrics.ExecutionBrokerExecutionFinished.M(1))

	var (
		replyData insolar.Reply
		err       error
	)

	if transcript.Context != nil {
		ctx = transcript.Context
	} else {
		inslogger.FromContext(ctx).Error("context in transcript is nil")
	}

	ctx = instracer.WithParentSpan(ctx, instracer.TraceSpan{
		TraceID: []byte(inslogger.TraceID(ctx)),
		SpanID:  instracer.MakeBinarySpan(transcript.Request.Reason.Bytes()),
	})

	ctx, span := instracer.StartSpan(ctx, "IncomingRequest processing")
	defer span.Finish()

	// ctx, logger := inslogger.WithFields(ctx, map[string]interface{}{})
	logger := inslogger.FromContext(ctx)

	if !q.canProcessTranscript(ctx, transcript) {
		// either closed broker or we're executing this already
		return
	}

	sendReply := true
	defer func() {
		q.finishTranscript(ctx, transcript)

		if sendReply {
			q.requestsExecutor.SendReply(ctx, transcript.RequestRef, *transcript.Request, replyData, err)
		}

		// we're checking here that pulse was changed and we should send
		// a message that we've finished processing tasks
		// note: ideally we should tell here that we've stopped executing
		//       but we only hoped that OnPulse had already told us that
		//       pulse changed and we should stop execution
		logger.Debug("finished request, try to finish pending if needed")
		q.finishPendingIfNeeded(ctx)
	}()

	if transcript.Request.CallType == record.CTMethod {
		logger.Info("processing transcript with method call")

		var objDesc artifacts.ObjectDescriptor
		objDesc, err = q.artifactsManager.GetObject(ctx, *transcript.Request.Object, &transcript.RequestRef)
		if err != nil {
			logger.Error("couldn't get object state: ", err)
			return
		}
		transcript.ObjectDescriptor = objDesc

		if !transcript.Request.Immutable &&
			transcript.ObjectDescriptor.EarliestRequestID() != nil &&
			!transcript.RequestRef.GetLocal().Equal(*transcript.ObjectDescriptor.EarliestRequestID()) {

			logger.Error("got different earliest request from ledger, this shouldn't happen")
			sendReply = false
			return
		}
	}

	var result artifacts.RequestResult
	result, err = q.requestsExecutor.ExecuteAndSave(ctx, transcript)
	if err != nil {
		logger.Warn("contract execution error: ", err)
		return
	}
	logger.Debug("finished executing method")

	objRef := result.ObjectReference()
	replyData = &reply.CallMethod{Result: result.Result(), Object: &objRef}
	// cleanup and reply is in defer
}

func (q *ExecutionBroker) canProcessTranscript(ctx context.Context, transcript *common.Transcript) bool {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if q.isClosed() {
		return false
	}

	err := q.executionRegistry.Register(ctx, transcript)
	if err != nil {
		if err == executionregistry.ErrAlreadyRegistered {
			stats.Record(ctx, metrics.ExecutionBrokerTranscriptAlreadyRegistered.M(1))
		}
		inslogger.FromContext(ctx).Error("couldn't register transcript, skipping: ", err.Error())
		return false
	}

	return true
}

func (q *ExecutionBroker) finishTranscript(ctx context.Context, transcript *common.Transcript) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	logger := inslogger.FromContext(ctx)
	logger.Debug("finishing transcript, moving from execution registry to finished")

	done := q.executionRegistry.Done(transcript)
	if !done {
		logger.Error("task wasn't in the registry, very bad")
	}

	q.finished = append(q.finished, transcript)
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
		q.finished = nil
		q.deduplicationTable = make(map[insolar.Reference]bool)
		q.pending = insolar.InPending
		q.PendingConfirmed = true
	}()

	q.close()

	sendExecResults := false

	switch {
	case q.isActive():
		logger.Debug("object is executing, sending result to next executor")
		q.pending = insolar.InPending
		sendExecResults = true
	case q.pending == insolar.InPending && !q.PendingConfirmed:
		stats.Record(ctx, metrics.ExecutionBrokerOnPulseNotConfirmed.M(1))
		logger.Warn("looks like pending executor died, continuing execution on next executor")
		q.pending = insolar.NotPending
		sendExecResults = true
		q.ledgerHasMoreRequests = LedgerHasMoreKnown
	case len(q.finished) > 0:
		logger.Debug("had activity on object, sending result to next executor")
		sendExecResults = true
	case q.ledgerHasMoreRequests > LedgerIsEmpty:
		logger.Debug("object is marked as having requests on ledger, sending result to next executor")
		sendExecResults = true
	default:
		logger.Debug("not sending result to next executor")
	}

	messages := make([]payload.Payload, 0)

	if sendExecResults {
		ledgerHasMoreRequests := false
		if q.ledgerHasMoreRequests > LedgerIsEmpty {
			ledgerHasMoreRequests = true
		}

		messages = append(messages, &payload.ExecutorResults{
			RecordRef:             q.Ref,
			Pending:               q.pending,
			LedgerHasMoreRequests: ledgerHasMoreRequests,
		})
	}

	return messages
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
	go q.startProcessor(ctx)

	return nil
}

func (q *ExecutionBroker) SetNotPending(ctx context.Context) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	q.pending = insolar.NotPending
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

	logger := inslogger.FromContext(ctx)
	logger.Debug("clarifying pending state from ledger")

	has, err := q.artifactsManager.HasPendings(ctx, q.Ref)
	if err != nil {
		logger.Error("couldn't check pending state: ", err.Error())
		return
	}

	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if q.pending == insolar.PendingUnknown {
		if has {
			logger.Debug("ledger has requests older than one pulse")
			q.pending = insolar.InPending
			q.ledgerHasMoreRequests = LedgerHasMoreUnknown
		} else {
			logger.Debug("no requests on ledger older than one pulse")
			q.pending = insolar.NotPending
		}
	}
}

func (q *ExecutionBroker) isClosed() bool {
	select {
	case <-q.closed:
		return true
	default:
		return false
	}
}

func (q *ExecutionBroker) close() {
	if !q.isClosed() {
		close(q.closed)
	}
}

// must be called under lock
func (q *ExecutionBroker) setHasMoreRequests() {
	q.ledgerHasMoreRequests = LedgerHasMoreKnown
	select {
	case q.probablyMoreSinceLastFetch <- struct{}{}:
	default:
	}
}

func (q *ExecutionBroker) upsertToDuplicationTable(ctx context.Context, transcript *common.Transcript) (alreadyInTable bool) {
	q.stateLock.Lock()
	defer q.stateLock.Unlock()

	if _, ok := q.deduplicationTable[transcript.RequestRef]; ok {
		logger := inslogger.FromContext(ctx)
		logger.Infof("Already know about request %s, skipping", transcript.RequestRef.String())

		return true
	}
	q.deduplicationTable[transcript.RequestRef] = true
	return false
}
