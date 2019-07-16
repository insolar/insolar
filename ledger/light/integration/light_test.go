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

func Test_AllOperations(t *testing.T) {
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

	// Save and check request.
	{
		codeID, codeRecord := setCode(t, s)
		codeRec := getCode(t, s, codeID)
		assert.Equal(t, &codeRecord, codeRec.Virtual)
	}

	// Set request for object.
	objectID, _ := setIncomingRequest(t, s, record.CTSaveAsChild)
	// Activate and check object.
	{
		state := activateObject(t, s, objectID)
		_, material := getObject(t, s, objectID)
		require.Equal(t, &state, material.Virtual)
	}
	// Amend and check object.
	{
		state := amendObject(t, s, objectID)
		_, material := getObject(t, s, objectID)
		require.Equal(t, &state, material.Virtual)
	}
	// Deactivate and check object.
	{
		deactivateObject(t, s, objectID)
		s.Receive(func(meta payload.Meta, pl payload.Payload) {
			received <- pl
		})
		s.Send(&payload.GetObject{
			ObjectID: objectID,
		})
		pl := <-received
		_, ok := pl.(*payload.Error)
		assert.True(t, ok)
	}
}

func setCode(t *testing.T, s *Server) (insolar.ID, record.Virtual) {
	received := make(chan payload.Payload)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		if _, ok := pl.(*payload.ID); !ok {
			return
		}
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
		if _, ok := pl.(*payload.Code); !ok {
			return
		}

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
		if _, ok := pl.(*payload.RequestInfo); !ok {
			return
		}

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
		Reason:    gen.Reference(),
	})
	s.Send(&payload.SetIncomingRequest{
		Request: rec,
	})
	pl := <-received
	id, ok := pl.(*payload.RequestInfo)
	require.True(t, ok)
	return id.RequestID, rec
}

func activateObject(t *testing.T, s *Server, objectID insolar.ID) record.Virtual {
	received := make(chan payload.Payload)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		if _, ok := pl.(*payload.ID); !ok {
			return
		}

		received <- pl
	})

	mem := make([]byte, 100)
	_, err := rand.Read(mem)
	require.NoError(t, err)
	rec := record.Wrap(record.Activate{
		Request: *insolar.NewReference(objectID),
		Memory:  mem,
	})
	buf, err := rec.Marshal()
	require.NoError(t, err)
	res := make([]byte, 100)
	_, err = rand.Read(res)
	require.NoError(t, err)
	resultRecord := record.Wrap(record.Result{
		Request: *insolar.NewReference(objectID),
		Object:  objectID,
		Payload: res,
	})
	resBuf, err := resultRecord.Marshal()
	require.NoError(t, err)
	s.Send(&payload.Activate{
		Record: buf,
		Result: resBuf,
	})
	pl := <-received
	_, ok := pl.(*payload.ID)
	require.True(t, ok)
	return rec
}

func amendObject(t *testing.T, s *Server, objectID insolar.ID) record.Virtual {
	received := make(chan payload.Payload)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		if _, ok := pl.(*payload.ID); !ok {
			return
		}

		received <- pl
	})

	mem := make([]byte, 100)
	_, err := rand.Read(mem)
	require.NoError(t, err)
	rec := record.Wrap(record.Amend{
		Memory: mem,
	})
	buf, err := rec.Marshal()
	require.NoError(t, err)
	res := make([]byte, 100)
	_, err = rand.Read(res)
	require.NoError(t, err)
	resultRecord := record.Wrap(record.Result{
		Request: *insolar.NewReference(objectID),
		Object:  objectID,
		Payload: res,
	})
	resBuf, err := resultRecord.Marshal()
	require.NoError(t, err)
	s.Send(&payload.Update{
		Record: buf,
		Result: resBuf,
	})
	pl := <-received
	_, ok := pl.(*payload.ID)
	require.True(t, ok)
	return rec
}

func deactivateObject(t *testing.T, s *Server, objectID insolar.ID) record.Virtual {
	received := make(chan payload.Payload)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		if _, ok := pl.(*payload.ID); !ok {
			return
		}

		received <- pl
	})

	mem := make([]byte, 100)
	_, err := rand.Read(mem)
	require.NoError(t, err)
	rec := record.Wrap(record.Deactivate{
		Request: *insolar.NewReference(objectID),
	})
	buf, err := rec.Marshal()
	require.NoError(t, err)
	res := make([]byte, 100)
	_, err = rand.Read(res)
	require.NoError(t, err)
	resultRecord := record.Wrap(record.Result{
		Request: *insolar.NewReference(objectID),
		Object:  objectID,
		Payload: res,
	})
	resBuf, err := resultRecord.Marshal()
	require.NoError(t, err)
	s.Send(&payload.Deactivate{
		Record: buf,
		Result: resBuf,
	})
	pl := <-received
	_, ok := pl.(*payload.ID)
	require.True(t, ok)
	return rec
}

func getObject(t *testing.T, s *Server, objectID insolar.ID) (record.Lifeline, record.Material) {
	received := make(chan payload.Payload)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		received <- pl
	})

	s.Send(&payload.GetObject{
		ObjectID: objectID,
	})
	var (
		lifeline *record.Lifeline
		state    *record.Material
	)
	done := func() bool {
		return lifeline != nil && state != nil
	}
	for pl := range received {
		switch p := pl.(type) {
		case *payload.Index:
			lifeline = &record.Lifeline{}
			err := lifeline.Unmarshal(p.Index)
			require.NoError(t, err)
		case *payload.State:
			state = &record.Material{}
			err := state.Unmarshal(p.Record)
			require.NoError(t, err)
		}

		if done() {
			break
		}
	}
	require.True(t, done())
	return *lifeline, *state
}
