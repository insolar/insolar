package integration_test

import (
	"crypto/rand"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CallsOnBootstrap(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg)
	require.NoError(t, err)

	t.Run("message before pulse received returns error", func(t *testing.T) {
		received := make(chan payload.Payload)
		s.Receive(func(meta payload.Meta, pl payload.Payload) {
			received <- pl
		})

		s.Send(&payload.SetCode{})
		pl := <-received
		_, ok := pl.(*payload.Error)
		require.True(t, ok)
	})

	// First pulse goes in storage then interrupts.
	s.Pulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.Pulse(ctx)

	t.Run("messages after two pulses return result", func(t *testing.T) {
		assertAllCalls(t, s)
	})
}

func assertAllCalls(t *testing.T, s *Server) {
	received := make(chan payload.Payload)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		received <- pl
	})

	// Save code.
	var (
		codeID     insolar.ID
		codeRecord record.Virtual
	)
	{
		code := make([]byte, 100)
		_, err := rand.Read(code)
		require.NoError(t, err)
		codeRecord = record.Wrap(record.Code{Code: code})
		buf, err := codeRecord.Marshal()
		require.NoError(t, err)
		s.Send(&payload.SetCode{
			Record: buf,
		})
		pl := <-received
		id, ok := pl.(*payload.ID)
		require.True(t, ok)
		codeID = id.ID
	}

	// Get code.
	{
		s.Send(&payload.GetCode{
			CodeID: codeID,
		})
		pl := <-received
		code, ok := pl.(*payload.Code)
		require.True(t, ok)
		material := record.Material{}
		err := material.Unmarshal(code.Record)
		require.NoError(t, err)
		assert.Equal(t, &codeRecord, material.Virtual)
	}
}
