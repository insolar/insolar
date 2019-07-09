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

package artifacts

import (
	"context"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

type retryer struct {
	ppl           payload.Payload
	role          insolar.DynamicRole
	ref           insolar.Reference
	tries         uint
	sender        bus.Sender
	pulseAccessor pulse.Accessor

	processingStarted sync.Mutex
	once              sync.Once
	done              chan struct{}
	replyChan         chan *message.Message
}

func newRetryer(sender bus.Sender, pulseAccessor pulse.Accessor, ppl payload.Payload, role insolar.DynamicRole, ref insolar.Reference, tries uint) *retryer {
	r := &retryer{
		ppl:           ppl,
		role:          role,
		ref:           ref,
		tries:         tries,
		sender:        sender,
		pulseAccessor: pulseAccessor,
		done:          make(chan struct{}),
		replyChan:     make(chan *message.Message),
	}
	return r
}

func (r *retryer) clientDone() {
	r.once.Do(func() {
		close(r.done)

		r.processingStarted.Lock()
		close(r.replyChan)
		r.processingStarted.Unlock()
	})
}

func isChannelClosed(ch chan *message.Message) bool {
	select {
	case _, ok := <-ch:
		return !ok
	default:
		return false
	}
}

func (r *retryer) waitForPulseChange(ctx context.Context, lastPulse insolar.PulseNumber) (insolar.PulseNumber, error) {
	logger := inslogger.FromContext(ctx)
	for {
		currentPulse, err := r.pulseAccessor.Latest(ctx)
		if err != nil {
			return lastPulse, errors.Wrap(err, "can't get latest pulse")
		}

		if currentPulse.PulseNumber == lastPulse {
			logger.Debugf("wait for pulse change in retryer. Current: %d", currentPulse)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		return currentPulse.PulseNumber, nil
	}
}

// return params tells, if request must be resend
func (r *retryer) proxyReply(ctx context.Context, reps <-chan *message.Message) bool {
	for {
		select {
		case <-r.done:
			return false
		case rep, ok := <-reps:
			if !ok {
				return false
			}
			if r.isRetryableError(ctx, rep) {
				return true
			}

			select {
			case <-r.done:
				return false
			case r.replyChan <- rep:
			}
		}
	}
}

func (r *retryer) isRetryableError(ctx context.Context, rep *message.Message) bool {
	replyPayload, err := payload.UnmarshalFromMeta(rep.Payload)
	if err != nil {
		return false
	}
	p, ok := replyPayload.(*payload.Error)
	if ok && (p.Code == payload.CodeFlowCanceled) {
		inslogger.FromContext(ctx).Errorf("flow cancelled, retrying (error message - %s)", p.Text)
		r.tries--
		return true
	}
	return false
}

func (r *retryer) send(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	retry := true
	var lastPulse insolar.PulseNumber

	func() {
		r.processingStarted.Lock()
		defer r.processingStarted.Unlock()
		if isChannelClosed(r.replyChan) {
			return
		}
		for r.tries > 0 && retry {
			var err error
			lastPulse, err = r.waitForPulseChange(ctx, lastPulse)
			if err != nil {
				logger.Error(errors.Wrap(err, "can't wait for pulse change"))
				break
			}

			msg, err := payload.NewMessage(r.ppl)
			if err != nil {
				logger.Error(errors.Wrap(err, "error while create new message from payload"))
				break
			}

			reps, done := r.sender.SendRole(ctx, msg, r.role, r.ref)
			retry = r.proxyReply(ctx, reps)
			done()
		}
	}()

	if r.tries == 0 {
		logger.Error(errors.Errorf("flow cancelled, retries exceeded"))
	}
	r.clientDone()
}
