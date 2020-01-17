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
// +build slowtest

package integration_test

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/require"
)

func Test_GetObject_PassingRequestID(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)

	t.Run("Incoming request can't have several different results", func(t *testing.T) {
		var firstReqID insolar.ID
		// Creating root reason request.
		{
			msg, _ := MakeSetIncomingRequest(
				gen.ID(),
				gen.IDWithPulse(s.Pulse()),
				insolar.ID{},
				true,
				true,
				"first",
			)
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			firstReqID = rep.(*payload.RequestInfo).RequestID

			p, _ := CallActivateObject(ctx, s, firstReqID)
			RequireNotError(p)
		}

		s.SetPulse(ctx)

		var secondReqID insolar.ID
		var thirdreqID insolar.ID
		// Register second request
		{
			msg, _ := MakeSetIncomingRequest(
				firstReqID,
				firstReqID,
				insolar.ID{},
				false,
				true,
				"second",
			)
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			secondReqID = rep.(*payload.RequestInfo).RequestID

			msg, _ = MakeSetIncomingRequest(
				firstReqID,
				firstReqID,
				insolar.ID{},
				false,
				true,
				"third",
			)
			rep = SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			thirdreqID = rep.(*payload.RequestInfo).RequestID
		}

		s.SetPulse(ctx)
		// Call get object with second ID (while first isn't closed)
		{
			lifelinePL, statePL := CallGetObject(ctx, s, firstReqID, &thirdreqID)
			RequireNotError(lifelinePL)
			RequireNotError(statePL)

			lifeline := lifelinePL.(*payload.Index)
			require.NotNil(t, lifeline.EarliestRequestID)
			require.Equal(t, secondReqID, *lifeline.EarliestRequestID)
		}
	})
}
