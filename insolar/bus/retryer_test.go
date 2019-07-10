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
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/payload"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/pulse"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/testutils"
)

func accessorMock(t *testing.T) pulse.Accessor {
	p := pulse.NewAccessorMock(t)
	pulseNumber := 10
	p.LatestFunc = func(p context.Context) (r insolar.Pulse, r1 error) {
		pulseNumber = pulseNumber + 10
		return insolar.Pulse{PulseNumber: insolar.PulseNumber(pulseNumber)}, nil
	}

	return p
}

// Send msg, bus.Sender gets error and closes resp chan
func TestRetryerSend_SendErrored(t *testing.T) {
	sender := &SenderMock{}
	sender.SendRoleFunc = func(p context.Context, p1 *message.Message, p2 insolar.DynamicRole, p3 insolar.Reference) (r <-chan *message.Message, r1 func()) {
		res := make(chan *message.Message)
		close(res)
		return res, func() {}
	}

	msg, err := payload.NewMessage(&payload.State{})
	require.NoError(t, err)
	reps, done := SendRoleWithRetry(context.Background(), sender, accessorMock(t), msg, insolar.DynamicRoleLightExecutor, testutils.RandomRef(), 3)
	defer done()
	for range reps {
		require.Fail(t, "we are not expect any replays")
	}
}

// Send msg, close reply channel by timeout
func TestRetryerSend_Send_Timeout(t *testing.T) {
	once := sync.Once{}
	sender := &SenderMock{}
	innerReps := make(chan *message.Message)
	sender.SendRoleFunc = func(p context.Context, p1 *message.Message, p2 insolar.DynamicRole, p3 insolar.Reference) (r <-chan *message.Message, r1 func()) {
		done := func() {
			once.Do(func() { close(innerReps) })
		}
		go func() {
			time.Sleep(time.Second * 2)
			done()
		}()
		return innerReps, done
	}

	msg, err := payload.NewMessage(&payload.State{})
	require.NoError(t, err)
	reps, _ := SendRoleWithRetry(context.Background(), sender, accessorMock(t), msg, insolar.DynamicRoleLightExecutor, testutils.RandomRef(), 3)
	select {
	case _, ok := <-reps:
		require.False(t, ok, "channel with replies must be closed, without any messages received")
	}
}

// Send msg, client stops waiting for response before request was actually done
func TestRetryerSend_Send_ClientDone(t *testing.T) {
	sender := &SenderMock{}

	msg, err := payload.NewMessage(&payload.State{})
	require.NoError(t, err)
	r := newRetryer(sender, nil)

	r.clientDone()
	r.send(context.Background(), msg, insolar.DynamicRoleLightExecutor, testutils.RandomRef(), 3)

	for range r.replyChan {
		require.Fail(t, "we are not expect any replays")
	}
}

func sendTestReply(pl payload.Payload, ch chan<- *message.Message, isDone chan<- interface{}) {
	msg, _ := payload.NewMessage(pl)
	meta := payload.Meta{
		Payload: msg.Payload,
	}
	buf, _ := meta.Marshal()
	msg.Payload = buf
	ch <- msg
	close(isDone)
}

// Send msg, get one response
func TestRetryerSend(t *testing.T) {
	sender := NewSenderMock(t)
	innerReps := make(chan *message.Message)
	sender.SendRoleFunc = func(p context.Context, p1 *message.Message, p2 insolar.DynamicRole, p3 insolar.Reference) (r <-chan *message.Message, r1 func()) {
		return innerReps, func() { close(innerReps) }
	}

	msg, err := payload.NewMessage(&payload.State{})
	require.NoError(t, err)
	reps, done := SendRoleWithRetry(context.Background(), sender, accessorMock(t), msg, insolar.DynamicRoleLightExecutor, testutils.RandomRef(), 3)

	isDone := make(chan<- interface{})
	go sendTestReply(&payload.Error{Text: "object is deactivated", Code: payload.CodeUnknown}, innerReps, isDone)

	var success bool
	for rep := range reps {
		replyPayload, err := payload.UnmarshalFromMeta(rep.Payload)
		require.Nil(t, err)

		switch p := replyPayload.(type) {
		case *payload.Error:
			switch p.Code {
			case payload.CodeUnknown:
				success = true
			}
		}

		if success {
			break
		}
	}
	done()

	require.True(t, isChannelClosed(innerReps))
}

// Send msg, get "flow cancelled" error, than get one response
func TestRetryerSend_FlowCancelled_Once(t *testing.T) {
	sender := NewSenderMock(t)
	innerReps := make(chan *message.Message)
	sender.SendRoleFunc = func(p context.Context, p1 *message.Message, p2 insolar.DynamicRole, p3 insolar.Reference) (r <-chan *message.Message, r1 func()) {
		innerReps = make(chan *message.Message)
		if sender.SendRoleCounter == 0 {
			go sendTestReply(&payload.Error{Text: "test error", Code: payload.CodeFlowCanceled}, innerReps, make(chan<- interface{}))
		} else {
			go sendTestReply(&payload.State{}, innerReps, make(chan<- interface{}))
		}
		return innerReps, func() { close(innerReps) }
	}

	var success bool
	msg, err := payload.NewMessage(&payload.State{})
	require.NoError(t, err)
	reps, done := SendRoleWithRetry(context.Background(), sender, accessorMock(t), msg, insolar.DynamicRoleLightExecutor, testutils.RandomRef(), 3)
	defer done()
	for rep := range reps {
		replyPayload, _ := payload.UnmarshalFromMeta(rep.Payload)

		switch replyPayload.(type) {
		case *payload.State:
			success = true
		}

		if success {
			break
		}
	}
	done()

	require.True(t, isChannelClosed(innerReps))
}

// Send msg, get "flow cancelled" error, than get two responses
func TestRetryerSend_FlowCancelled_Once_SeveralReply(t *testing.T) {
	sender := NewSenderMock(t)
	innerReps := make(chan *message.Message)
	sender.SendRoleFunc = func(p context.Context, p1 *message.Message, p2 insolar.DynamicRole, p3 insolar.Reference) (r <-chan *message.Message, r1 func()) {
		innerReps = make(chan *message.Message)
		if sender.SendRoleCounter == 0 {
			go sendTestReply(&payload.Error{Text: "test error", Code: payload.CodeFlowCanceled}, innerReps, make(chan<- interface{}))
		} else {
			go sendTestReply(&payload.State{}, innerReps, make(chan<- interface{}))
			go sendTestReply(&payload.State{}, innerReps, make(chan<- interface{}))
		}
		return innerReps, func() { close(innerReps) }
	}

	var success int
	msg, err := payload.NewMessage(&payload.State{})
	require.NoError(t, err)
	reps, done := SendRoleWithRetry(context.Background(), sender, accessorMock(t), msg, insolar.DynamicRoleLightExecutor, testutils.RandomRef(), 3)
	for rep := range reps {
		replyPayload, _ := payload.UnmarshalFromMeta(rep.Payload)

		switch replyPayload.(type) {
		case *payload.State:
			success = success + 1
		}

		if success == 2 {
			break
		}
	}
	done()

	require.True(t, isChannelClosed(innerReps))
}

// Send msg, get "flow cancelled" error on every tries
func TestRetryerSend_FlowCancelled_RetryExceeded(t *testing.T) {
	sender := NewSenderMock(t)
	innerReps := make(chan *message.Message)
	sender.SendRoleFunc = func(p context.Context, p1 *message.Message, p2 insolar.DynamicRole, p3 insolar.Reference) (r <-chan *message.Message, r1 func()) {
		innerReps = make(chan *message.Message)
		go sendTestReply(&payload.Error{Text: "test error", Code: payload.CodeFlowCanceled}, innerReps, make(chan<- interface{}))
		return innerReps, func() { close(innerReps) }
	}

	var success bool
	msg, err := payload.NewMessage(&payload.State{})
	require.NoError(t, err)
	reps, done := SendRoleWithRetry(context.Background(), sender, accessorMock(t), msg, insolar.DynamicRoleLightExecutor, testutils.RandomRef(), 3)
	for range reps {
		success = true
		break
	}
	require.False(t, success)

	done()

	require.True(t, isChannelClosed(innerReps))
}

// Send msg, get response, than get "flow cancelled" error, than get two responses
func TestRetryerSend_FlowCancelled_Between(t *testing.T) {
	sender := NewSenderMock(t)
	innerReps := make(chan *message.Message)
	sender.SendRoleFunc = func(p context.Context, p1 *message.Message, p2 insolar.DynamicRole, p3 insolar.Reference) (r <-chan *message.Message, r1 func()) {
		innerReps = make(chan *message.Message)
		if sender.SendRoleCounter == 0 {
			go func() {
				isDone := make(chan interface{})
				go sendTestReply(&payload.State{}, innerReps, isDone)
				<-isDone
				go sendTestReply(&payload.Error{Text: "test error", Code: payload.CodeFlowCanceled}, innerReps, make(chan<- interface{}))
			}()
		} else {
			go func() {
				isDone := make(chan interface{})
				go sendTestReply(&payload.State{}, innerReps, isDone)
				<-isDone
				go sendTestReply(&payload.State{}, innerReps, make(chan<- interface{}))
			}()
		}
		return innerReps, func() { close(innerReps) }
	}

	var success int
	msg, err := payload.NewMessage(&payload.State{})
	require.NoError(t, err)
	reps, done := SendRoleWithRetry(context.Background(), sender, accessorMock(t), msg, insolar.DynamicRoleLightExecutor, testutils.RandomRef(), 3)
	for rep := range reps {
		replyPayload, _ := payload.UnmarshalFromMeta(rep.Payload)

		switch replyPayload.(type) {
		case *payload.State:
			success = success + 1
		default:
		}

		if success == 3 {
			break
		}
	}

	done()

	require.True(t, isChannelClosed(innerReps))
}
