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

package handler

import (
	"context"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestHandleGetJet(t *testing.T) {
	t.Parallel()

	ttable := []struct {
		name   string
		actual bool
	}{
		{name: "Actual", actual: true},
		{name: "NotActual", actual: false},
	}
	for _, tt := range ttable {
		t.Run(tt.name, func(t *testing.T) {
			testCtx := inslogger.TestContext(t)
			testPulseNumber := insolar.PulseNumber(insolar.FirstPulseNumber + 797979)
			testID := gen.ID()
			testJetID := gen.JetID()
			h := New()
			accessorMock := jet.NewAccessorMock(t)
			accessorMock.ForIDFunc = func(ctx context.Context, pn insolar.PulseNumber, id insolar.ID) (r insolar.JetID, r1 bool) {
				require.Equal(t, testPulseNumber, pn)
				require.Equal(t, testID, id)

				return insolar.JetID(testJetID), tt.actual
			}
			h.JetAccessor = accessorMock

			testParcel := testutils.NewParcelMock(t)
			testParcel.MessageFunc = func() (r insolar.Message) {
				return &message.GetJet{Pulse: testPulseNumber, Object: testID}
			}

			rawReply, err := h.handleGetJet(testCtx, testParcel)
			require.NoError(t, err)
			jetReply := rawReply.(*reply.Jet)
			require.Equal(t, insolar.ID(testJetID), jetReply.ID)
			require.Equal(t, tt.actual, jetReply.Actual)
		})
	}
}
