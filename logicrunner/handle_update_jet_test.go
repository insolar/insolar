package logicrunner

import (
	"context"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/require"
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

	jets.UpdateFunc = func(_ context.Context, pn insolar.PulseNumber, a bool, jets ...insolar.JetID) (r error) {
		require.Equal(t, receivedPayload.Pulse, pn)
		require.Equal(t, true, a)
		require.Equal(t, jets, []insolar.JetID{receivedPayload.JetID})
		return nil
	}
	err = h.Present(ctx, nil)
	require.NoError(t, err)
}
