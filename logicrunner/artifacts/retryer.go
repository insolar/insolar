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

	wmmsg "github.com/ThreeDotsLabs/watermill/message"

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
	replyChan         chan *wmmsg.Message
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
		replyChan:     make(chan *wmmsg.Message),
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

func isChannelClosed(ch chan *wmmsg.Message) bool {
	select {
	case _, ok := <-ch:
		return !ok
	default:
		return false
	}
}

func (r *retryer) send(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	var errText string
	retry := true
	var lastPulse insolar.PulseNumber

	r.processingStarted.Lock()
	if isChannelClosed(r.replyChan) {
		return
	}
	for r.tries > 0 && retry {
		currentPulse, err := r.pulseAccessor.Latest(ctx)
		if err != nil {
			logger.Error(errors.Wrap(err, "can't get latest pulse"))
			break
		}

		if currentPulse.PulseNumber == lastPulse {
			inslogger.FromContext(ctx).Debugf("wait for pulse change in retryer. Current: %d", currentPulse)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		lastPulse = currentPulse.PulseNumber

		msg, err := payload.NewMessage(r.ppl)
		if err != nil {
			logger.Error(errors.Wrap(err, "error while create new message from payload"))
			break
		}

		reps, done := r.sender.SendRole(ctx, msg, r.role, r.ref)

		retry = false
	F:
		for {
			select {
			case rep, ok := <-reps:
				if !ok {
					break F
				}
				replyPayload, err := payload.UnmarshalFromMeta(rep.Payload)
				if err == nil {
					p, ok := replyPayload.(*payload.Error)
					if ok && (p.Code == payload.CodeFlowCanceled) {
						errText = p.Text
						inslogger.FromContext(ctx).Warnf("flow cancelled, retrying (error message - %s)", errText)
						r.tries--
						retry = true
						break F
					}
				}

				select {
				case r.replyChan <- rep:
				case <-r.done:
					break F
				}
			case <-r.done:
				break F
			}
		}
		done()
	}
	r.processingStarted.Unlock()

	if r.tries == 0 {
		logger.Error(errors.Errorf("flow cancelled, retries exceeded (last error - %s)", errText))
	}
	r.clientDone()
}
