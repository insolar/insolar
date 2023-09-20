package bus

import (
	"context"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
)

// RetrySender allows to send messaged via provided Sender with retries.
type RetrySender struct {
	sender        Sender
	pulseAccessor pulse.Accessor
	retries       uint
	responseCount uint
}

// NewRetrySender creates RetrySender instance with provided values.
func NewRetrySender(sender Sender, pulseAccessor pulse.Accessor, retries uint, responseCount uint) *RetrySender {
	return &RetrySender{
		sender:        sender,
		pulseAccessor: pulseAccessor,
		retries:       retries,
		responseCount: responseCount,
	}
}

func (r *RetrySender) SendTarget(ctx context.Context, msg *message.Message, target insolar.Reference) (<-chan *message.Message, func()) {
	return r.retryWrapper(ctx, msg, func(ctx context.Context, msg *message.Message) (<-chan *message.Message, func()) {
		return r.sender.SendTarget(ctx, msg, target)
	})
}

func (r *RetrySender) Reply(ctx context.Context, origin payload.Meta, reply *message.Message) {
	inslogger.FromContext(ctx).Panic("not implemented")
}

// SendRole sends message to specified role, using provided Sender.SendRole. If error with CodeFlowCanceled
// was received, it retries request after pulse on current node will be changed.
// Replies will be written to the returned channel. Always read from the channel using multiple assignment
// (rep, ok := <-ch) because the channel will be closed on timeout.
func (r *RetrySender) SendRole(
	ctx context.Context, msg *message.Message, role insolar.DynamicRole, ref insolar.Reference,
) (<-chan *message.Message, func()) {
	return r.retryWrapper(ctx, msg, func(ctx context.Context, msg *message.Message) (<-chan *message.Message, func()) {
		return r.sender.SendRole(ctx, msg, role, ref)
	})
}

func (r *RetrySender) retryWrapper(ctx context.Context, msg *message.Message, caller func(context.Context, *message.Message) (<-chan *message.Message, func())) (<-chan *message.Message, func()) {
	tries := r.retries + 1
	once := sync.Once{}
	done := make(chan struct{})
	replyChan := make(chan *message.Message)

	go func() {
		defer close(replyChan)
		logger := inslogger.FromContext(ctx)
		var lastPulse insolar.PulseNumber

		received := false
		updateUUID := false
		for tries > 0 && !received {
			var err error
			// this doesn't wait on first iteration due to lastPulse being zero
			lastPulse, err = r.waitForPulseChange(ctx, lastPulse)
			if err != nil {
				logger.Error(errors.Wrap(err, "can't wait for pulse change"))
				break
			}

			if updateUUID {
				msg.UUID = watermill.NewUUID()
			}
			reps, d := caller(ctx, msg)
			received = tryReceive(ctx, reps, done, replyChan, r.responseCount)
			tries--
			updateUUID = true
			d()
		}

		if tries < r.retries {
			mctx := insmetrics.InsertTag(ctx, tagMessageType, messagePayloadTypeName(msg))
			stats.Record(mctx, statRetries.M(int64(r.retries-tries)))
		}

		if tries == 0 && !received {
			logger.Error(errors.Errorf("flow cancelled, retries exceeded"))
		}
	}()

	closeDone := func() {
		once.Do(func() {
			close(done)
		})
	}
	return replyChan, closeDone
}

func (r *RetrySender) waitForPulseChange(ctx context.Context, lastPulse insolar.PulseNumber) (insolar.PulseNumber, error) {
	logger := inslogger.FromContext(ctx)
	for {
		currentPulse, err := r.pulseAccessor.Latest(ctx)
		if err != nil {
			return lastPulse, errors.Wrap(err, "can't get latest pulse")
		}

		if currentPulse.PulseNumber == lastPulse {
			logger.Debugf("wait for pulse change in RetrySender. Current: %d", currentPulse.PulseNumber)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		return currentPulse.PulseNumber, nil
	}
}

type messageType int

const (
	messageTypeNotError messageType = iota
	messageTypeErrorRetryable
	messageTypeErrorNonRetryable
)

// tryReceive returns false if we get retryable error,
// and true if reply was successfully received or client don't want anymore replies
func tryReceive(
	ctx context.Context,
	reps <-chan *message.Message,
	done chan struct{},
	receiver chan<- *message.Message,
	responseCount uint,
) bool {
	for i := uint(0); i < responseCount; i++ {
		rep, ok := <-reps
		if !ok {
			return true
		}

		var leave bool
		switch getErrorType(ctx, rep) {
		case messageTypeErrorRetryable:
			return false
		case messageTypeErrorNonRetryable:
			leave = true
		default:
		}

		select {
		case <-done:
		case receiver <- rep:
		}
		if leave {
			break
		}
	}

	return true
}

func getErrorType(ctx context.Context, rep *message.Message) messageType {
	replyPayload, err := payload.UnmarshalFromMeta(rep.Payload)
	if err != nil {
		return messageTypeNotError
	}

	p, ok := replyPayload.(*payload.Error)
	if ok {
		if p.Code == payload.CodeFlowCanceled {
			inslogger.FromContext(ctx).Infof("flow cancelled, retrying (error message - %s)", p.Text)
			return messageTypeErrorRetryable
		}

		return messageTypeErrorNonRetryable
	}
	return messageTypeNotError
}

func ReplyError(ctx context.Context, sender Sender, meta payload.Meta, err error) {
	errCode := payload.CodeUnknown

	// Throwing custom error code
	cause := errors.Cause(err)
	insError, ok := cause.(*payload.CodedError)
	if ok {
		errCode = insError.GetCode()
	}

	// todo refactor this #INS-3191
	if cause == flow.ErrCancelled {
		errCode = payload.CodeFlowCanceled
	}
	errMsg, newErr := payload.NewMessage(&payload.Error{Text: err.Error(), Code: errCode})
	if newErr != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to reply error"))
	}
	sender.Reply(ctx, meta, errMsg)
}
