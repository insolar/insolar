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
	"github.com/insolar/insolar/insolar/pulse"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

// Send msg, bus.Sender gets error and closes resp chan
func TestWaitOKSender_SendRole_RetryExceeded(t *testing.T) {
	sender := NewSenderMock(t)

	innerReps := make(chan *message.Message)
	sender.SendRoleMock.Set(func(p context.Context, p1 *message.Message, p2 insolar.DynamicRole, p3 insolar.Reference) (r <-chan *message.Message, r1 func()) {
		innerReps = make(chan *message.Message)
		go sendTestReply(&payload.Error{Text: "test error", Code: payload.CodeFlowCanceled}, innerReps, make(chan<- interface{}))
		return innerReps, func() { close(innerReps) }
	})
	msg, err := payload.NewMessage(&payload.State{})
	require.NoError(t, err)

	tries := 3

	pa := pulse.NewAccessorMock(t)
	pa.LatestMock.Set(accessorMock(t).Latest)
	c := NewWaitOKWithRetrySender(sender, pa, uint(tries))

	c.SendRole(context.Background(), msg, insolar.DynamicRoleLightExecutor, testutils.RandomRef())

	require.EqualValues(t, tries, sender.SendRoleAfterCounter())
}

func sendOK(ch chan<- *message.Message) {
	msg := ReplyAsMessage(context.Background(), &reply.OK{})
	meta := payload.Meta{
		Payload: msg.Payload,
	}
	buf, _ := meta.Marshal()
	msg.Payload = buf
	ch <- msg
}

func TestWaitOKSender_SendRole_RetryOnce(t *testing.T) {
	sender := NewSenderMock(t)

	innerReps := make(chan *message.Message)
	sender.SendRoleMock.Set(func(p context.Context, p1 *message.Message, p2 insolar.DynamicRole, p3 insolar.Reference) (r <-chan *message.Message, r1 func()) {
		innerReps = make(chan *message.Message)
		if sender.SendRoleAfterCounter() == 0 {
			go sendTestReply(&payload.Error{Text: "test error", Code: payload.CodeFlowCanceled}, innerReps, make(chan<- interface{}))
		} else {
			go sendOK(innerReps)
		}
		return innerReps, func() { close(innerReps) }
	})
	msg, err := payload.NewMessage(&payload.State{})
	require.NoError(t, err)

	pa := pulse.NewAccessorMock(t)
	pa.LatestMock.Set(accessorMock(t).Latest)
	c := NewWaitOKWithRetrySender(sender, pa, 3)

	c.SendRole(context.Background(), msg, insolar.DynamicRoleLightExecutor, testutils.RandomRef())

	require.EqualValues(t, 2, sender.SendRoleAfterCounter())
}

func TestWaitOKSender_SendRole_OK(t *testing.T) {
	sender := NewSenderMock(t)

	innerReps := make(chan *message.Message)
	sender.SendRoleMock.Set(func(p context.Context, p1 *message.Message, p2 insolar.DynamicRole, p3 insolar.Reference) (r <-chan *message.Message, r1 func()) {
		return innerReps, func() { close(innerReps) }
	})

	go sendOK(innerReps)

	msg, err := payload.NewMessage(&payload.State{})
	require.NoError(t, err)
	pa := pulse.NewAccessorMock(t)
	pa.LatestMock.Set(accessorMock(t).Latest)
	c := NewWaitOKWithRetrySender(sender, pa, 3)

	c.SendRole(context.Background(), msg, insolar.DynamicRoleLightExecutor, testutils.RandomRef())

	require.EqualValues(t, 1, sender.SendRoleAfterCounter())
}

func TestWaitOKSender_SendRole_NotOK(t *testing.T) {
	sender := NewSenderMock(t)

	innerReps := make(chan *message.Message)
	sender.SendRoleMock.Set(func(p context.Context, p1 *message.Message, p2 insolar.DynamicRole, p3 insolar.Reference) (r <-chan *message.Message, r1 func()) {
		return innerReps, func() { close(innerReps) }
	})

	go sendTestReply(&payload.Error{Text: "test error", Code: payload.CodeUnknown}, innerReps, make(chan<- interface{}))

	msg, err := payload.NewMessage(&payload.State{})
	require.NoError(t, err)
	pa := pulse.NewAccessorMock(t)
	pa.LatestMock.Set(accessorMock(t).Latest)
	c := NewWaitOKWithRetrySender(sender, pa, 3)

	c.SendRole(context.Background(), msg, insolar.DynamicRoleLightExecutor, testutils.RandomRef())

	require.EqualValues(t, 1, sender.SendRoleAfterCounter())
}
