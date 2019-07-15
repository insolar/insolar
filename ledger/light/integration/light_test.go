package integration_test

import (
	"crypto/rand"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BootstrapCalls(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg)
	require.NoError(t, err)

	received := make(chan payload.Payload)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		received <- pl
	})

	t.Run("message before pulse received returns error", func(t *testing.T) {
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
		setCode(t, s)
	})
}

func Test_ReplicationScenario(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg)
	require.NoError(t, err)
	received := make(chan payload.Payload)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		received <- pl
	})

	// First pulse goes in storage then interrupts.
	s.Pulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.Pulse(ctx)

	codeID, codeRecord := setCode(t, s)
	material := getCode(t, s, codeID)
	assert.Equal(t, &codeRecord, material.Virtual)

	// FIXME: doesn't work with old bus. Move Hot data to watermill to fix.
	// _, _ = setIncomingRequest(t, s, record.CTSaveAsChild)
	//
	// // Activate object.
	// var (
	// 	objectID       insolar.ID
	// 	activateRecord record.Virtual
	// )
	// {
	// 	mem := make([]byte, 100)
	// 	_, err := rand.Read(mem)
	// 	require.NoError(t, err)
	// 	activateRecord = record.Wrap(record.Activate{
	// 		Memory: mem,
	// 	})
	// 	buf, err := activateRecord.Marshal()
	// 	require.NoError(t, err)
	// 	res := make([]byte, 100)
	// 	_, err = rand.Read(res)
	// 	require.NoError(t, err)
	// 	resultRecord := record.Wrap(record.Result{Payload: res})
	// 	resBuf, err := resultRecord.Marshal()
	// 	require.NoError(t, err)
	// 	s.Send(&payload.Activate{
	// 		Record: buf,
	// 		Result: resBuf,
	// 	})
	// 	pl := <-received
	// 	id, ok := pl.(*payload.ID)
	// 	require.True(t, ok)
	// 	objectID = id.ID
	// }

	// _ = objectID
}

func setCode(t *testing.T, s *Server) (insolar.ID, record.Virtual) {
	received := make(chan payload.Payload)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		received <- pl
	})

	code := make([]byte, 100)
	_, err := rand.Read(code)
	require.NoError(t, err)
	rec := record.Wrap(record.Code{Code: code})
	buf, err := rec.Marshal()
	require.NoError(t, err)
	s.Send(&payload.SetCode{
		Record: buf,
	})
	pl := <-received
	id, ok := pl.(*payload.ID)
	require.True(t, ok)
	return id.ID, rec
}

func getCode(t *testing.T, s *Server, id insolar.ID) record.Material {
	received := make(chan payload.Payload)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		received <- pl
	})

	s.Send(&payload.GetCode{
		CodeID: id,
	})
	pl := <-received
	code, ok := pl.(*payload.Code)
	require.True(t, ok)
	material := record.Material{}
	err := material.Unmarshal(code.Record)
	require.NoError(t, err)
	return material
}

func setIncomingRequest(t *testing.T, s *Server, ct record.CallType) (insolar.ID, record.Virtual) {
	received := make(chan payload.Payload)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		received <- pl
	})

	args := make([]byte, 100)
	_, err := rand.Read(args)
	require.NoError(t, err)
	objRef := gen.Reference()
	rec := record.Wrap(record.IncomingRequest{
		Object:    &objRef,
		Arguments: args,
		CallType:  ct,
	})
	s.Send(&payload.SetIncomingRequest{
		Request: rec,
	})
	pl := <-received
	id, ok := pl.(*payload.ID)
	require.True(t, ok)
	return id.ID, rec
}
