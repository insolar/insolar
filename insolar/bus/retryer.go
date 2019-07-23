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
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

// RetrySender allows to send messaged via provided Sender with retries.
type RetrySender struct {
	sender Sender
	tries  uint
}

// NewRetrySender creates RetrySender instance with provided values.
func NewRetrySender(sender Sender, tries uint) *RetrySender {
	r := &RetrySender{
		sender: sender,
		tries:  tries,
	}
	return r
}

func (r *RetrySender) SendTarget(ctx context.Context, msg *message.Message, target insolar.Reference) (<-chan *message.Message, func()) {
	panic("not implemented")
}

func (r *RetrySender) Reply(ctx context.Context, origin payload.Meta, reply *message.Message) {
	panic("not implemented")
}

func (r *RetrySender) LatestPulse(ctx context.Context) (insolar.Pulse, error) {
	return r.sender.LatestPulse(ctx)
}

// SendRole sends message to specified role, using provided Sender.SendRole. If error with CodeFlowCanceled
// was received, it retries request after pulse on current node will be changed.
// Replies will be written to the returned channel. Always read from the channel using multiple assignment
// (rep, ok := <-ch) because the channel will be closed on timeout.
func (r *RetrySender) SendRole(
	ctx context.Context, msg *message.Message, role insolar.DynamicRole, ref insolar.Reference,
) (<-chan *message.Message, func()) {
	tries := r.tries
	once := sync.Once{}
	done := make(chan struct{})
	replyChan := make(chan *message.Message)

	go func() {
		defer close(replyChan)
		logger := inslogger.FromContext(ctx)
		var lastPulse insolar.PulseNumber

		select {
		case <-done:
			return
		default:
		}

		received := false
		for tries > 0 && !received {
			var err error
			lastPulse, err = r.waitForPulseChange(ctx, lastPulse)
			if err != nil {
				logger.Error(errors.Wrap(err, "can't wait for pulse change"))
				break
			}

			reps, d := r.sender.SendRole(ctx, msg, role, ref)
			received = tryReceive(ctx, reps, done, replyChan)
			tries--
			d()
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
		currentPulse, err := r.sender.LatestPulse(ctx)
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
func tryReceive(ctx context.Context, reps <-chan *message.Message, done chan struct{}, receiver chan<- *message.Message) bool {
	for {
		select {
		case <-done:
			return true
		case rep, ok := <-reps:
			if !ok {
				return true
			}
			if isRetryableError(ctx, rep) {
				return false
			}

			select {
			case <-done:
				return true
			case receiver <- rep:
			}
		}
	}
}

func isRetryableError(ctx context.Context, rep *message.Message) bool {
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
