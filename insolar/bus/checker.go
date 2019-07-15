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
	"bytes"
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

// CheckOKSender allows to send messaged via provided Sender with retries and check that reply.OK was received.
type CheckOKSender struct {
	sender *RetrySender
}

// NewCheckOKSender creates CheckOKSender instance with provided values.
func NewCheckOKSender(sender Sender, pulseAccessor pulse.Accessor, tries uint) *CheckOKSender {
	r := &RetrySender{
		sender:        sender,
		pulseAccessor: pulseAccessor,
		tries:         tries,
	}
	c := &CheckOKSender{
		sender: r,
	}
	return c
}

// SendRole sends message to specified role, using provided Sender.SendRole. It waiting for reply.OK and
// close replies channel after getting it. If error with CodeFlowCanceled was received, it retries request.
func (c *CheckOKSender) SendRole(
	ctx context.Context, msg *message.Message, role insolar.DynamicRole, ref insolar.Reference,
) {
	reps, done := c.sender.SendRole(ctx, msg, role, ref)
	defer done()

	rep, ok := <-reps

	if !ok {
		logger := inslogger.FromContext(ctx)
		logger.Errorf("reply channel was closed before we get any valid replies")
		return
	}

	check(ctx, rep)
}

func check(ctx context.Context, rep *message.Message) {
	logger := inslogger.FromContext(ctx)

	meta := payload.Meta{}
	err := meta.Unmarshal(rep.Payload)
	if err != nil {
		logger.Error(errors.Wrap(err, "can't deserialize message payload"))
		return
	}

	if rep.Metadata.Get(MetaType) == TypeReply {
		r, err := reply.Deserialize(bytes.NewBuffer(meta.Payload))
		if err != nil {
			logger.Error(errors.Wrap(err, "can't deserialize payload to reply"))
			return
		}
		if r.Type() != reply.TypeOK {
			logger.Errorf("expected OK, got %s.", r.Type())
		}
		return
	}

	replyPayload, err := payload.UnmarshalFromMeta(rep.Payload)
	if err != nil {
		logger.Error(errors.Wrap(err, "failed to unmarshal reply"))
		return
	}

	switch p := replyPayload.(type) {
	case *payload.Error:
		logger.Errorf("expected OK, got error: %s, code - %d", p.Text, p.Code)
	default:
		logger.Errorf("got unexpected reply: %#v", p)
	}
	return
}
