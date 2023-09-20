package logicrunner

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
)

func TestHandleUpdateJet_Present(t *testing.T) {
	defer testutils.LeakTester(t)

	tests := []struct {
		name  string
		mocks func(t minimock.Tester) (*HandleUpdateJet, flow.Flow)
		error bool
	}{
		{
			name: "success",
			mocks: func(t minimock.Tester) (*HandleUpdateJet, flow.Flow) {
				receivedPayload := payload.UpdateJet{
					Pulse: gen.PulseNumber(),
					JetID: gen.JetID(),
				}

				buf, err := payload.Marshal(&receivedPayload)
				require.NoError(t, err, "marshal")

				jets := jet.NewStorageMock(t)
				jets.UpdateMock.Inspect(
					func(_ context.Context, pn insolar.PulseNumber, a bool, jets ...insolar.JetID) {
						require.Equal(t, receivedPayload.Pulse, pn)
						require.Equal(t, true, a)
						require.Equal(t, jets, []insolar.JetID{receivedPayload.JetID})
					},
				).Return(nil)

				h := &HandleUpdateJet{
					dep: &Dependencies{JetStorage: jets},
					meta: payload.Meta{
						Payload: buf,
					},
				}

				f := flow.NewFlowMock(t)
				return h, f
			},
		},
		{
			name: "error updating tree",
			mocks: func(t minimock.Tester) (*HandleUpdateJet, flow.Flow) {
				receivedPayload := payload.UpdateJet{
					Pulse: gen.PulseNumber(),
					JetID: gen.JetID(),
				}

				buf, err := payload.Marshal(&receivedPayload)
				require.NoError(t, err, "marshal")

				jets := jet.NewStorageMock(t)
				jets.UpdateMock.Return(errors.New("some"))

				h := &HandleUpdateJet{
					dep: &Dependencies{JetStorage: jets},
					meta: payload.Meta{
						Payload: buf,
					},
				}

				f := flow.NewFlowMock(t)
				return h, f
			},
			error: true,
		},
		{
			name: "error unmarshaling",
			mocks: func(t minimock.Tester) (*HandleUpdateJet, flow.Flow) {
				h := &HandleUpdateJet{
					meta: payload.Meta{
						Payload: []byte{3, 2, 1},
					},
				}

				f := flow.NewFlowMock(t)
				return h, f
			},
			error: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
			mc := minimock.NewController(t)

			h, f := test.mocks(mc)
			err := h.Present(ctx, f)
			if test.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mc.Wait(1 * time.Minute)
			mc.Finish()
		})
	}
}
