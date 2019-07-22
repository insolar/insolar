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

package integration_test

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/require"
)

func Test_Pending_RequestRegistration(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)

	// First pulse goes in storage then interrupts.
	s.Pulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.Pulse(ctx)

	t.Run("pending was added", func(t *testing.T) {
		p, _ := setIncomingRequest(ctx, t, s, gen.ID(), gen.ID(), record.CTSaveAsChild)
		requirePayloadNotError(t, p)
		reqInfo := p.(*payload.RequestInfo)
		p = fetchPendings(ctx, t, s, reqInfo.RequestID)
		requirePayloadNotError(t, p)

		ids := p.(*payload.IDs)
		require.Equal(t, 1, len(ids.IDs))
		require.Equal(t, reqInfo.RequestID, ids.IDs[0])
	})

	t.Run("pending was added and closed", func(t *testing.T) {
		p, _ := setIncomingRequest(ctx, t, s, gen.ID(), gen.ID(), record.CTSaveAsChild)
		requirePayloadNotError(t, p)
		reqInfo := p.(*payload.RequestInfo)

		p, _ = activateObject(ctx, t, s, reqInfo.RequestID)
		requirePayloadNotError(t, p)

		p = fetchPendings(ctx, t, s, reqInfo.RequestID)

		err := p.(*payload.Error)
		require.Equal(t, insolar.ErrNoPendingRequest.Error(), err.Text)
	})
}
