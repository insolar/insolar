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

	isDone bool

	replyChan chan *wmmsg.Message

	processingStarted sync.Mutex
	once              sync.Once
	channelClosed     chan interface{}
}

func newRetryer(sender bus.Sender, pulseAccessor pulse.Accessor, ppl payload.Payload, role insolar.DynamicRole, ref insolar.Reference, tries uint) *retryer {
	r := &retryer{
		ppl:           ppl,
		role:          role,
		ref:           ref,
		tries:         tries,
		sender:        sender,
		pulseAccessor: pulseAccessor,
		channelClosed: make(chan interface{}),
		replyChan:     make(chan *wmmsg.Message),
		isDone:        false,
	}
	return r
}

func (r *retryer) clientDone() {
	r.once.Do(func() {
		close(r.channelClosed)
		r.processingStarted.Lock()
		close(r.replyChan)
		r.isDone = true
		r.processingStarted.Unlock()
	})
}

func (r *retryer) send(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	var errText string
	retry := true
	var lastPulse insolar.PulseNumber

	r.processingStarted.Lock()
	for r.tries > 0 && retry && !r.isDone {
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

		retry = false

		msg, err := payload.NewMessage(r.ppl)
		if err != nil {
			logger.Error(errors.Wrap(err, "error while create new message from payload"))
			break
		}

		reps, done := r.sender.SendRole(ctx, msg, r.role, r.ref)
		waitingForReply := true
	F:
		for waitingForReply {
			select {
			case rep, ok := <-reps:
				if !ok {
					waitingForReply = false
					break
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
				case <-r.channelClosed:
					waitingForReply = false
				}
			case <-r.channelClosed:
				waitingForReply = false
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
