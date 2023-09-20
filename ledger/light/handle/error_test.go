package handle_test

import (
	"testing"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/handle"
)

func TestError_BadMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	msg := message.NewMessage(watermill.NewUUID(), []byte{1, 2, 3, 4, 5})
	handler := handle.NewError(msg)
	err := handler.Present(ctx, flow.NewFlowMock(t))

	// We get error inside error-handler,
	// but only print log message for this,
	// without error returning.
	require.NoError(t, err)
}

func TestError_IncorrectTypeMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

	meta := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		// Incorrect type (SetIncomingRequest instead of Error).
		Payload: payload.MustMarshal(&payload.SetIncomingRequest{
			Polymorph: uint32(payload.TypeSetIncomingRequest),
		}),
	}

	p, err := meta.Marshal()
	require.NoError(t, err)

	msg := message.NewMessage(watermill.NewUUID(), p)
	handler := handle.NewError(msg)
	err = handler.Present(ctx, f)

	// We get error inside error-handler,
	// but only print log message for this,
	// without error returning.
	require.NoError(t, err)
}

func TestError_HappyPath(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

	meta := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		// Good error payload.
		Payload: payload.MustMarshal(&payload.Error{
			Polymorph: uint32(payload.TypeError),
			Code:      payload.CodeUnknown,
			Text:      "something good",
		}),
	}

	p, err := meta.Marshal()
	require.NoError(t, err)

	msg := message.NewMessage(watermill.NewUUID(), p)
	handler := handle.NewError(msg)
	err = handler.Present(ctx, f)

	// We get error inside error-handler,
	// but only print log message for this,
	// without error returning.
	require.NoError(t, err)
}
