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

func TestRemoteTarget_Notify(t *testing.T) {
	var (
		ctx     = inslogger.TestContext(t)
		address = "127.0.0.1:8080"
		pulse   = insolar.GenesisPulse.PulseNumber
	)
	transport := NewTransportMock(t)
	reply, _ := insolar.Serialize(GenericReply{Data: nil, Error: nil})
	transport.SendMock.Return(reply, nil)
	target := NewRemoteTarget(transport, address)

	err := target.Notify(ctx, pulse)
	require.NoError(t, err)
}
