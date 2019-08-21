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
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	busMeta "github.com/insolar/insolar/insolar/bus/meta"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// WaitOKSender allows to send messaged via provided Sender and wait for reply.OK.
type WaitOKSender struct {
	sender Sender
}

// NewWaitOKWithRetrySender creates WaitOKSender instance with RetrySender as Sender.
func NewWaitOKWithRetrySender(sender Sender, pulseAccessor pulse.Accessor, retries uint) *WaitOKSender {
	r := NewRetrySender(sender, pulseAccessor, retries, 1)
	c := NewWaitOKSender(r)
	return c
}

// NewWaitOKSender creates WaitOKSender instance with provided values.
func NewWaitOKSender(sender Sender) *WaitOKSender {
	c := &WaitOKSender{
		sender: sender,
	}
	return c
}

// SendRole sends message to specified role, using provided Sender.SendRole. It waiting for one reply and
// close replies channel after getting it. If reply is not reply.OK, it logs error message.
func (c *WaitOKSender) SendRole(
	ctx context.Context, msg *message.Message, role insolar.DynamicRole, ref insolar.Reference,
) {
	ctx, _ = inslogger.WithField(ctx, "msg_type", messagePayloadTypeName(msg))

	reps, done := c.sender.SendRole(ctx, msg, role, ref)
	defer done()

	processResult(ctx, reps)
}

// SendTarget sends message to specified target, using provided Sender.SendTarget. It waiting for one reply and
// close replies channel after getting it. If reply is not reply.OK, it logs error message.
func (c *WaitOKSender) SendTarget(
	ctx context.Context, msg *message.Message, target insolar.Reference) {

	ctx, _ = inslogger.WithField(ctx, "msg_type", messagePayloadTypeName(msg))

	reps, done := c.sender.SendTarget(ctx, msg, target)
	defer done()

	processResult(ctx, reps)
}

func processResult(ctx context.Context, responses <-chan *message.Message) {
	rep, ok := <-responses

	if !ok {
		logger := inslogger.FromContext(ctx)
		logger.Errorf("reply channel was closed before we get any valid replies")
		return
	}

	checkReply(ctx, rep)
}

func checkReply(ctx context.Context, rep *message.Message) {
	logger := inslogger.FromContext(ctx)

	meta := payload.Meta{}
	err := meta.Unmarshal(rep.Payload)
	if err != nil {
		logger.Error(errors.Wrap(err, "can't deserialize message payload"))
		return
	}

	if rep.Metadata.Get(busMeta.Type) == busMeta.TypeReply {
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
}
