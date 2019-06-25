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
)

func TestParent_Subscribe(t *testing.T) {
	var (
		pos     = Position{0, insolar.GenesisPulse.PulseNumber}
		address = "127.0.0.1:8080"
	)
	transport := NewTransportMock(t)
	transport.SendMock.Return(nil, nil)
	parent := NewRemoteParent(transport)
	targetTransport := NewTransportMock(t)
	targetTransport.MeMock.Return(address)
	target := NewRemoteTarget(targetTransport)

	err := parent.Subscribe(pos)
	require.NoError(t, err)
}

func TestParent_Pull(t *testing.T) {
	var (
		pos   = Position{10, insolar.GenesisPulse.PulseNumber}
		limit = uint32(10)
		reply = []byte{1, 2, 3}
	)
	rawReply, _ := insolar.Serialize(Reply{reply, nil})
	transport := NewTransportMock(t)
	transport.SendMock.Return(rawReply, nil)
	parent := NewRemoteParent(transport)

	actualReply, err := parent.Pull(0, pos, limit)
	require.NoError(t, err)
	require.Equal(t, reply, actualReply)
}
