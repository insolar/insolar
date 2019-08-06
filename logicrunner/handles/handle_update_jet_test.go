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

package handles

import (
	"context"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func Test_HandleUpdateJet_Present(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	jets := jet.NewStorageMock(mc)

	receivedPayload := payload.UpdateJet{
		Pulse: gen.PulseNumber(),
		JetID: gen.JetID(),
	}
	buf, err := payload.Marshal(&receivedPayload)
	h := HandleUpdateJet{
		dep: &Dependencies{JetStorage: jets},
		meta: payload.Meta{
			Payload: buf,
		},
	}

	jets.UpdateMock.Set(func(_ context.Context, pn insolar.PulseNumber, a bool, jets ...insolar.JetID) (r error) {
		require.Equal(t, receivedPayload.Pulse, pn)
		require.Equal(t, true, a)
		require.Equal(t, jets, []insolar.JetID{receivedPayload.JetID})
		return nil
	})
	err = h.Present(ctx, nil)
	require.NoError(t, err)
}
