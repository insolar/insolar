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

package bus

import (
	"context"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

type RetrySender struct {
	sender        Sender
	pulseAccessor pulse.Accessor

	tries uint

	once      sync.Once
	done      chan struct{}
	replyChan chan *message.Message
}

func NewRetrySender(sender Sender, pulseAccessor pulse.Accessor, tries uint) *RetrySender {
	r := &RetrySender{
		sender:        sender,
		pulseAccessor: pulseAccessor,
		tries:         tries,
		done:          make(chan struct{}),
		replyChan:     make(chan *message.Message),
	}
	return r
}

// SendWithRetry sends message to specified role, using provided Sender.SendRole. If error with CodeFlowCanceled
// was received, it retries request after pulse on current node will be changed.
// Replies will be written to the returned channel. Always read from the channel using multiple assignment
// (rep, ok := <-ch) because the channel will be closed on timeout.
func (r *RetrySender) SendRole(
	ctx context.Context, msg *message.Message, role insolar.DynamicRole, ref insolar.Reference,
) (<-chan *message.Message, func()) {
	go func() {
		defer close(r.replyChan)
		logger := inslogger.FromContext(ctx)
		var lastPulse insolar.PulseNumber

		select {
		case <-r.done:
			return
		default:
		}

		received := false
		for r.tries > 0 && !received {
			var err error
			lastPulse, err = r.waitForPulseChange(ctx, lastPulse)
			if err != nil {
				logger.Error(errors.Wrap(err, "can't wait for pulse change"))
				break
			}

			reps, done := r.sender.SendRole(ctx, msg, role, ref)
			received = r.tryReceive(ctx, reps)
			r.tries--
			done()
		}

		if r.tries == 0 {
			logger.Error(errors.Errorf("flow cancelled, retries exceeded"))
		}
	}()
	return r.replyChan, r.clientDone
}

func (r *RetrySender) clientDone() {
	r.once.Do(func() {
		close(r.done)
	})
}

func (r *RetrySender) waitForPulseChange(ctx context.Context, lastPulse insolar.PulseNumber) (insolar.PulseNumber, error) {
	logger := inslogger.FromContext(ctx)
	for {
		currentPulse, err := r.pulseAccessor.Latest(ctx)
		if err != nil {
			return lastPulse, errors.Wrap(err, "can't get latest pulse")
		}

		if currentPulse.PulseNumber == lastPulse {
			logger.Debugf("wait for pulse change in RetrySender. Current: %d", currentPulse)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		return currentPulse.PulseNumber, nil
	}
}

// tryReceive returns false if we get retryable error,
// and true if reply was successfully received or client don't want anymore replies
func (r *RetrySender) tryReceive(ctx context.Context, reps <-chan *message.Message) bool {
	for {
		select {
		case <-r.done:
			return true
		case rep, ok := <-reps:
			if !ok {
				return true
			}
			if r.isRetryableError(ctx, rep) {
				return false
			}

			select {
			case <-r.done:
				return true
			case r.replyChan <- rep:
			}
		}
	}
}

func (r *RetrySender) isRetryableError(ctx context.Context, rep *message.Message) bool {
	replyPayload, err := payload.UnmarshalFromMeta(rep.Payload)
	if err != nil {
		return false
	}
	p, ok := replyPayload.(*payload.Error)
	if ok && (p.Code == payload.CodeFlowCanceled) {
		inslogger.FromContext(ctx).Errorf("flow cancelled, retrying (error message - %s)", p.Text)
		return true
	}
	return false
}
