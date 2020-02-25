// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package handle_test

import (
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/handle"
	"github.com/insolar/insolar/ledger/light/proc"
)

func TestGetPulse_Present(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	var (
		dep  *proc.Dependencies
		meta payload.Meta
	)

	setup := func() {
		dep = proc.NewDependenciesMock()
	}

	t.Run("basic ok", func(t *testing.T) {
		setup()
		defer mc.Finish()

		meta = payload.Meta{
			Polymorph: uint32(payload.TypeMeta),
			Payload: payload.MustMarshal(&payload.GetPulse{
				Polymorph:   uint32(payload.TypeGetPulse),
				PulseNumber: 42,
			}),
			ID: []byte{1, 1, 1},
		}

		handler := handle.NewGetPulse(dep, meta)
		flowMock := flow.NewFlowMock(mc).ProcedureMock.Return(nil)
		err := handler.Present(ctx, flowMock)
		assert.NoError(t, err)
	})
}

func TestGetPulse_NilMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	meta := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   nil,
	}

	handler := handle.NewGetPulse(nil, meta)

	err := handler.Present(ctx, flow.NewFlowMock(mc))
	require.Error(t, err)

	mc.Finish()
}

func TestGetPulse_BadMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	meta := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   []byte{1, 2, 3, 4, 5},
	}

	handler := handle.NewGetPulse(nil, meta)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestGetPulse_IncorrectTypeMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	f := flow.NewFlowMock(mc)

	meta := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		// Incorrect type (SetIncomingRequest instead of GetPulse).
		Payload: payload.MustMarshal(&payload.SetIncomingRequest{
			Polymorph: uint32(payload.TypeSetIncomingRequest),
			Request:   record.Virtual{},
		}),
		ID: []byte{1, 1, 1},
	}

	handler := handle.NewGetPulse(proc.NewDependenciesMock(), meta)

	err := handler.Present(ctx, f)
	require.Error(t, err)

	mc.Finish()
}
