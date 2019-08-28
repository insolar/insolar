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

package flow_test

import (
	"context"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus/meta"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
)

type mockReply struct{}

func (r *mockReply) Type() insolar.ReplyType {
	return 88
}

func makeMessage(t *testing.T, ctx context.Context, pn insolar.PulseNumber) *message.Message {
	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)
	msg.Metadata.Set(meta.Pulse, pn.String())
	sp, err := instracer.Serialize(ctx)
	require.NoError(t, err)
	msg.Metadata.Set(meta.SpanData, string(sp))

	return msg
}

func TestEmptyHandle(t *testing.T) {
	testReply := &mockReply{}
	replyChan := make(chan *mockReply, 1)
	currentPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(100)}
	pulseAccessorMock := pulse.NewAccessorMock(t).LatestMock.Return(currentPulse, nil)
	disp := dispatcher.NewDispatcher(pulseAccessorMock, func(message *message.Message) flow.Handle {
		return func(context context.Context, f flow.Flow) error {
			replyChan <- testReply
			return nil
		}
	}, nil, nil)
	ctx := context.Background()

	msg := makeMessage(t, ctx, currentPulse.PulseNumber)

	err := disp.Process(msg)
	require.NoError(t, err)
	reply := <-replyChan
	require.Equal(t, testReply, reply)
}

type EmptyProcedure struct{}

func (p *EmptyProcedure) Proceed(context.Context) error {
	return nil
}

func TestCallEmptyProcedure(t *testing.T) {
	testReply := &mockReply{}

	replyChan := make(chan *mockReply, 1)
	currentPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(100)}
	pulseAccessorMock := pulse.NewAccessorMock(t).LatestMock.Return(currentPulse, nil)
	disp := dispatcher.NewDispatcher(pulseAccessorMock, func(message *message.Message) flow.Handle {
		return func(context context.Context, f flow.Flow) error {
			err := f.Procedure(context, &EmptyProcedure{}, true)
			require.NoError(t, err)
			replyChan <- testReply
			return nil
		}
	}, nil, nil)

	ctx := context.Background()

	msg := makeMessage(t, ctx, currentPulse.PulseNumber)

	err := disp.Process(msg)
	require.NoError(t, err)
	reply := <-replyChan
	require.Equal(t, testReply, reply)

}

type ErrorProcedure struct{}

func (p *ErrorProcedure) Proceed(context.Context) error {
	return errors.New("Errorchik")
}

func TestProcedureReturnError(t *testing.T) {
	testReply := &mockReply{}

	replyChan := make(chan *mockReply, 1)
	currentPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(100)}
	pulseAccessorMock := pulse.NewAccessorMock(t).LatestMock.Return(currentPulse, nil)
	disp := dispatcher.NewDispatcher(pulseAccessorMock, func(message *message.Message) flow.Handle {
		return func(context context.Context, f flow.Flow) error {
			err := f.Procedure(context, &ErrorProcedure{}, true)
			require.Error(t, err)
			replyChan <- testReply
			return nil
		}
	}, nil, nil)

	ctx := context.Background()

	msg := makeMessage(t, ctx, currentPulse.PulseNumber)

	err := disp.Process(msg)
	require.NoError(t, err)
	reply := <-replyChan
	require.Equal(t, testReply, reply)
}

type LongProcedure struct {
	started chan struct{}
}

func (p *LongProcedure) Proceed(context.Context) error {
	p.started <- struct{}{}
	time.Sleep(5 * time.Second)
	return nil
}

func TestClosePulse(t *testing.T) {
	testReply := &mockReply{}

	procedureStarted := make(chan struct{})

	replyChan := make(chan *mockReply, 1)
	currentPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(100)}
	pulseAccessorMock := pulse.NewAccessorMock(t).LatestMock.Return(currentPulse, nil)
	disp := dispatcher.NewDispatcher(pulseAccessorMock, func(message *message.Message) flow.Handle {
		return func(context context.Context, f flow.Flow) error {
			longProcedure := LongProcedure{}
			longProcedure.started = procedureStarted
			err := f.Procedure(context, &longProcedure, true)
			require.Equal(t, flow.ErrCancelled, err)
			replyChan <- testReply
			return nil
		}
	}, nil, nil)

	handleProcessed := make(chan struct{})

	go func() {
		ctx := context.Background()

		msg := makeMessage(t, ctx, currentPulse.PulseNumber)

		err := disp.Process(msg)
		require.NoError(t, err)
		reply := <-replyChan
		require.Equal(t, testReply, reply)
		handleProcessed <- struct{}{}
	}()

	<-procedureStarted
	p := pulsar.NewPulse(22, 33, &entropygenerator.StandardEntropyGenerator{})
	disp.ClosePulse(context.Background(), *p)
	<-handleProcessed
}

func TestChangePulseAndMigrate(t *testing.T) {
	testReply := &mockReply{}

	firstProcedureStarted := make(chan struct{})
	secondProcedureStarted := make(chan struct{})

	migrateStarted := make(chan struct{})

	replyChan := make(chan *mockReply, 1)
	currentPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(100)}
	pulseAccessorMock := pulse.NewAccessorMock(t).LatestMock.Return(currentPulse, nil)
	disp := dispatcher.NewDispatcher(pulseAccessorMock, func(message *message.Message) flow.Handle {
		return func(ctx context.Context, f1 flow.Flow) error {
			longProcedure := LongProcedure{}
			longProcedure.started = firstProcedureStarted
			err := f1.Procedure(ctx, &longProcedure, true)
			require.Equal(t, flow.ErrCancelled, err)

			f1.Migrate(ctx, func(c context.Context, f2 flow.Flow) error {
				migrateStarted <- struct{}{}
				longProcedure := LongProcedure{}
				longProcedure.started = secondProcedureStarted
				err := f2.Procedure(ctx, &longProcedure, true)
				require.Equal(t, flow.ErrCancelled, err)

				replyChan <- testReply

				return nil
			})
			return nil
		}
	}, nil, nil)

	handleProcessed := make(chan struct{})

	go func() {
		ctx := context.Background()

		msg := makeMessage(t, ctx, currentPulse.PulseNumber)

		err := disp.Process(msg)
		require.NoError(t, err)
		reply := <-replyChan
		require.Equal(t, testReply, reply)
		handleProcessed <- struct{}{}
	}()

	<-firstProcedureStarted
	pulse := pulsar.NewPulse(22, 33, &entropygenerator.StandardEntropyGenerator{})
	disp.ClosePulse(context.Background(), *pulse)
	disp.BeginPulse(context.Background(), *pulse)
	<-migrateStarted
	<-secondProcedureStarted
	disp.ClosePulse(context.Background(), *pulse)
	disp.BeginPulse(context.Background(), *pulse)
	<-handleProcessed
}
