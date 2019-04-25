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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/flow/handler"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type mockReply struct{}

func (r *mockReply) Type() insolar.ReplyType {
	return 88
}

func makeParcelMock(t *testing.T) insolar.Parcel {
	parcelMock := testutils.NewParcelMock(t)
	parcelMock.PulseFunc = func() (r insolar.PulseNumber) {
		return 33
	}

	return parcelMock
}

func TestEmptyHandle(t *testing.T) {
	testReply := &mockReply{}
	hand := handler.NewHandler(
		func(message bus.Message) flow.Handle {
			return func(context context.Context, f flow.Flow) error {
				message.ReplyTo <- bus.Reply{Reply: testReply}
				return nil
			}
		})

	reply, err := hand.WrapBusHandle(context.Background(), makeParcelMock(t))
	require.NoError(t, err)
	require.Equal(t, testReply, reply)
}

type EmptyProcedure struct{}

func (p *EmptyProcedure) Proceed(context.Context) error {
	return nil
}

func TestCallEmptyProcedure(t *testing.T) {
	testReply := &mockReply{}

	hand := handler.NewHandler(
		func(message bus.Message) flow.Handle {
			return func(context context.Context, f flow.Flow) error {
				err := f.Procedure(context, &EmptyProcedure{}, true)
				require.NoError(t, err)
				message.ReplyTo <- bus.Reply{Reply: testReply}
				return nil
			}
		})

	reply, err := hand.WrapBusHandle(context.Background(), makeParcelMock(t))
	require.NoError(t, err)
	require.Equal(t, testReply, reply)

}

type ErrorProcedure struct{}

func (p *ErrorProcedure) Proceed(context.Context) error {
	return errors.New("Errorchik")
}

func TestProcedureReturnError(t *testing.T) {
	testReply := &mockReply{}

	hand := handler.NewHandler(
		func(message bus.Message) flow.Handle {
			return func(context context.Context, f flow.Flow) error {
				err := f.Procedure(context, &ErrorProcedure{}, true)
				require.Error(t, err)
				message.ReplyTo <- bus.Reply{Reply: testReply}
				return nil
			}
		})

	reply, err := hand.WrapBusHandle(context.Background(), makeParcelMock(t))
	require.NoError(t, err)
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

func TestChangePulse(t *testing.T) {
	testReply := &mockReply{}

	procedureStarted := make(chan struct{})

	hand := handler.NewHandler(
		func(message bus.Message) flow.Handle {
			return func(context context.Context, f flow.Flow) error {
				longProcedure := LongProcedure{}
				longProcedure.started = procedureStarted
				err := f.Procedure(context, &longProcedure, true)
				require.Equal(t, flow.ErrCancelled, err)
				message.ReplyTo <- bus.Reply{Reply: testReply}
				return nil
			}
		})

	handleProcessed := make(chan struct{})
	go func() {
		reply, err := hand.WrapBusHandle(context.Background(), makeParcelMock(t))
		require.NoError(t, err)
		require.Equal(t, testReply, reply)
		handleProcessed <- struct{}{}
	}()

	<-procedureStarted
	pulse := pulsar.NewPulse(22, 33, &entropygenerator.StandardEntropyGenerator{})
	hand.ChangePulse(context.Background(), *pulse)
	<-handleProcessed
}

func TestChangePulseAndMigrate(t *testing.T) {
	testReply := &mockReply{}

	firstProcedureStarted := make(chan struct{})
	secondProcedureStarted := make(chan struct{})

	migrateStarted := make(chan struct{})

	hand := handler.NewHandler(
		func(message bus.Message) flow.Handle {
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

					message.ReplyTo <- bus.Reply{Reply: testReply}

					return nil
				})
				return nil
			}
		})

	handleProcessed := make(chan struct{})
	go func() {
		reply, err := hand.WrapBusHandle(context.Background(), makeParcelMock(t))
		require.NoError(t, err)
		require.Equal(t, testReply, reply)
		handleProcessed <- struct{}{}
	}()

	<-firstProcedureStarted
	pulse := pulsar.NewPulse(22, 33, &entropygenerator.StandardEntropyGenerator{})
	hand.ChangePulse(context.Background(), *pulse)
	<-migrateStarted
	<-secondProcedureStarted
	hand.ChangePulse(context.Background(), *pulse)
	<-handleProcessed
}
