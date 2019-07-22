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

package replica

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
)

func TestGrpcTransport_Send(t *testing.T) {
	ctx := context.Background()
	trans := NewGRPCTransport(20111)
	trans.(*grpcTransport).Init(ctx)
	trans.(*grpcTransport).Start(ctx)

	t.Logf("trans: %s", trans.Me())

	trans.Register("test.Test", func(data []byte) ([]byte, error) {
		t.Logf("msg: %s", string(data))
		require.Equal(t, "ping", string(data))
		return []byte("pong"), nil
	})
	trans2 := NewGRPCTransport(20112)
	trans2.(component.Initer).Init(ctx)
	trans2.(component.Starter).Start(ctx)

	t.Logf("trans2: %s", trans2.Me())

	reply, err := trans2.Send(ctx, "127.0.0.1:20111", "test.Test", []byte("ping"))
	require.NoError(t, err)
	t.Logf("reply: %s", string(reply))
	require.Equal(t, "pong", string(reply))
}

func TestGrpcSubscribe(t *testing.T) {
	ctx := context.Background()
	trans := NewGRPCTransport(20111)
	trans.(*grpcTransport).Init(ctx)
	trans.(*grpcTransport).Start(ctx)

	parent := NewParentMock(t)
	parent.SubscribeMock.Return(nil)

	trans.Register("replica.Subscribe", func(data []byte) ([]byte, error) {
		ctx := context.Background()
		sub := Subscription{}
		err := insolar.Deserialize(data, &sub)
		if err != nil {
			return []byte{}, errors.Wrapf(err, "failed to deserialize subscription data")
		}
		target := NewRemoteTarget(trans, sub.Target)
		err = parent.Subscribe(ctx, target, sub.At)
		if err != nil {
			return []byte{}, errors.Wrapf(err, "failed to call parent.Subscribe")
		}
		return []byte{}, nil
	})

	trans2 := NewGRPCTransport(20112)
	trans2.(component.Initer).Init(ctx)
	trans2.(component.Starter).Start(ctx)
	remoteParent := NewRemoteParent(trans2, "127.0.0.1:20111")
	err := remoteParent.Subscribe(ctx, nil, Page{Pulse: insolar.GenesisPulse.PulseNumber})
	require.NoError(t, err)
}
