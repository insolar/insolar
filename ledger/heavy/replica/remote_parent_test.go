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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func TestRemoteParent_Subscribe(t *testing.T) {
	var (
		ctx     = inslogger.TestContext(t)
		pos     = Page{Pulse: insolar.GenesisPulse.PulseNumber}
		address = "127.0.0.1:8080"
	)
	transport := NewTransportMock(t)
	transport.MeMock.Return(address)
	reply, _ := insolar.Serialize(GenericReply{Data: nil, Error: nil})
	transport.SendMock.Return(reply, nil)
	parent := NewRemoteParent(transport, address)
	targetTransport := NewTransportMock(t)
	targetTransport.MeMock.Return(address)
	target := NewRemoteTarget(targetTransport, address)

	err := parent.Subscribe(ctx, target, pos)
	require.NoError(t, err)
}

func TestRemoteParent_Pull(t *testing.T) {
	var (
		ctx     = inslogger.TestContext(t)
		pos     = Page{Pulse: insolar.GenesisPulse.PulseNumber}
		total   = uint32(10)
		reply   = []byte{1, 2, 3}
		address = "127.0.0.1:8080"
	)
	extReply, _ := insolar.Serialize(PullReply{Data: reply, Total: total})
	rawReply, _ := insolar.Serialize(GenericReply{Data: extReply, Error: nil})
	transport := NewTransportMock(t)
	transport.SendMock.Return(rawReply, nil)
	parent := NewRemoteParent(transport, address)

	actualReply, actualTotal, err := parent.Pull(ctx, pos)
	require.NoError(t, err)
	require.Equal(t, reply, actualReply)
	require.Equal(t, total, actualTotal)
}
