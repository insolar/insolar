package integration

import (
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/require"
)

func Test_Light(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewLightServer(ctx, cfg)
	require.NoError(t, err)

	s.Receive(func(msg *message.Message) {
		meta := payload.Meta{}
		err := meta.Unmarshal(msg.Payload)
		require.NoError(t, err)
		pl, err := payload.Unmarshal(meta.Payload)
		require.NoError(t, err)

		switch pl.(type) {
		}
	})

	require.NoError(t, s.Pulse(ctx))
	require.NoError(t, s.Pulse(ctx))
	require.NoError(t, s.Pulse(ctx))
}
