package integration_test

import (
	"context"
	"crypto/rand"
	"sync"
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
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg)
	require.NoError(t, err)

	t.Run("message before pulse received returns error", func(t *testing.T) {
		p, _ := setCode(ctx, t, s)
		requirePayloadNotError(t, p)
	})

	// First pulse goes in storage then interrupts.
	s.Pulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.Pulse(ctx)

	t.Run("messages after two pulses return result", func(t *testing.T) {
		p, _ := setCode(ctx, t, s)
		requirePayloadNotError(t, p)
	})
}

func Test_BasicOperations(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg)
	require.NoError(t, err)

	// First pulse goes in storage then interrupts.
	s.Pulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.Pulse(ctx)

	runner := func(t *testing.T) {
		// Creating root reason request.
		var reasonID insolar.ID
		{
			p, _ := setIncomingRequest(ctx, t, s, gen.ID(), gen.ID(), record.CTSaveAsChild)
			requirePayloadNotError(t, p)
			reasonID = p.(*payload.RequestInfo).RequestID
		}
		// Save and check code.
		{
			p, sent := setCode(ctx, t, s)
			requirePayloadNotError(t, p)

			p = getCode(ctx, t, s, p.(*payload.ID).ID)
			requirePayloadNotError(t, p)
			material := record.Material{}
			err := material.Unmarshal(p.(*payload.Code).Record)
			require.NoError(t, err)
			require.Equal(t, &sent, material.Virtual)
		}
		var objectID insolar.ID
		// Set, get request.
		{
			p, sent := setIncomingRequest(ctx, t, s, gen.ID(), reasonID, record.CTSaveAsChild)
			requirePayloadNotError(t, p)

			p = getRequest(ctx, t, s, p.(*payload.RequestInfo).RequestID)
			requirePayloadNotError(t, p)
			require.Equal(t, sent, p.(*payload.Request).Request)
			objectID = p.(*payload.Request).RequestID
		}
		// Activate and check object.
		{
			p, state := activateObject(ctx, t, s, objectID)
			requirePayloadNotError(t, p)
			_, material := requireGetObject(ctx, t, s, objectID)
			require.Equal(t, &state, material.Virtual)
		}
		// Amend and check object.
		{
			p, _ := setIncomingRequest(ctx, t, s, objectID, reasonID, record.CTMethod)
			requirePayloadNotError(t, p)
			p, state := amendObject(ctx, t, s, objectID, p.(*payload.RequestInfo).RequestID)
			requirePayloadNotError(t, p)
			_, material := requireGetObject(ctx, t, s, objectID)
			require.Equal(t, &state, material.Virtual)
		}
		// Deactivate and check object.
		{
			p, _ := setIncomingRequest(ctx, t, s, objectID, reasonID, record.CTMethod)
			requirePayloadNotError(t, p)
			deactivateObject(ctx, t, s, objectID, p.(*payload.RequestInfo).RequestID)

			lifeline, _ := getObject(ctx, t, s, objectID)
			_, ok := lifeline.(*payload.Error)
			assert.True(t, ok)
		}
	}

	t.Run("happy basic", runner)

	t.Run("happy concurrent", func(t *testing.T) {
		count := 1000
		var wg sync.WaitGroup
		wg.Add(count)
		for i := 0; i < count; i++ {
			go func() {
				runner(t)
				wg.Done()
			}()
		}

		wg.Wait()
	})
}

func setCode(ctx context.Context, t *testing.T, s *Server) (payload.Payload, record.Virtual) {
	code := make([]byte, 100)
	_, err := rand.Read(code)
	require.NoError(t, err)
	rec := record.Wrap(record.Code{Code: code})
	buf, err := rec.Marshal()
	require.NoError(t, err)
	msg, err := payload.NewMessage(&payload.SetCode{
		Record: buf,
	})
	require.NoError(t, err)
	reps, done := s.Send(ctx, msg)
	defer done()

	rep := <-reps
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	require.NoError(t, err)
	switch pl.(type) {
	case *payload.Error:
		return pl, rec
	case *payload.ID:
		return pl, rec
	default:
		t.Fatalf("received unexpected reply %T", pl)
	}

	return nil, rec
}

func getCode(ctx context.Context, t *testing.T, s *Server, id insolar.ID) payload.Payload {
	msg, err := payload.NewMessage(&payload.GetCode{
		CodeID: id,
	})

	reps, done := s.Send(ctx, msg)
	defer done()

	rep := <-reps
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	require.NoError(t, err)
	switch pl.(type) {
	case *payload.Error:
		return pl
	case *payload.Code:
		return pl
	default:
		t.Fatalf("received unexpected reply %T", pl)
	}
	return nil
}

func setIncomingRequest(
	ctx context.Context, t *testing.T, s *Server, objectID, reasonID insolar.ID, ct record.CallType,
) (payload.Payload, record.Virtual) {
	args := make([]byte, 100)
	_, err := rand.Read(args)
	require.NoError(t, err)
	rec := record.Wrap(record.IncomingRequest{
		Object:    insolar.NewReference(objectID),
		Arguments: args,
		CallType:  ct,
		Reason:    *insolar.NewReference(reasonID),
	})
	msg, err := payload.NewMessage(&payload.SetIncomingRequest{
		Request: rec,
	})
	require.NoError(t, err)
	reps, done := s.Send(ctx, msg)
	defer done()

	rep := <-reps
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	require.NoError(t, err)
	switch pl.(type) {
	case *payload.Error:
		return pl, rec
	case *payload.RequestInfo:
		return pl, rec
	default:
		t.Fatalf("received unexpected reply %T", pl)
	}

	return insolar.ID{}, record.Virtual{}
}

func getRequest(ctx context.Context, t *testing.T, s *Server, requestID insolar.ID) payload.Payload {
	msg, err := payload.NewMessage(&payload.GetRequest{
		RequestID: requestID,
	})
	require.NoError(t, err)
	reps, done := s.Send(ctx, msg)
	defer done()

	rep := <-reps
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	require.NoError(t, err)
	switch pl.(type) {
	case *payload.Error:
		return pl
	case *payload.Request:
		return pl
	default:
		t.Fatalf("received unexpected reply %T", pl)
	}

	return nil
}

func activateObject(ctx context.Context, t *testing.T, s *Server, objectID insolar.ID) (payload.Payload, record.Virtual) {
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

	msg, err := payload.NewMessage(&payload.Activate{
		Record: buf,
		Result: resBuf,
	})
	require.NoError(t, err)
	reps, done := s.Send(ctx, msg)
	defer done()

	rep := <-reps
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	require.NoError(t, err)
	switch pl.(type) {
	case *payload.Error:
		return pl, rec
	case *payload.ResultInfo:
		return pl, rec
	default:
		t.Fatalf("received unexpected reply %T", pl)
	}
	return nil, rec
}

func amendObject(ctx context.Context, t *testing.T, s *Server, objectID, requestID insolar.ID) (payload.Payload, record.Virtual) {
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
		Request: *insolar.NewReference(requestID),
		Object:  objectID,
		Payload: res,
	})
	resBuf, err := resultRecord.Marshal()
	require.NoError(t, err)

	msg, err := payload.NewMessage(&payload.Update{
		Record: buf,
		Result: resBuf,
	})
	require.NoError(t, err)

	reps, done := s.Send(ctx, msg)
	defer done()

	rep := <-reps
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	require.NoError(t, err)
	switch pl.(type) {
	case *payload.ResultInfo:
		return pl, rec
	case *payload.Error:
		return pl, rec
	default:
		t.Fatalf("received unexpected reply %T", pl)
	}
	return nil, rec
}

func deactivateObject(ctx context.Context, t *testing.T, s *Server, objectID, requestID insolar.ID) (payload.Payload, record.Virtual) {
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
		Request: *insolar.NewReference(requestID),
		Object:  objectID,
		Payload: res,
	})
	resBuf, err := resultRecord.Marshal()
	require.NoError(t, err)

	msg, err := payload.NewMessage(&payload.Deactivate{
		Record: buf,
		Result: resBuf,
	})
	require.NoError(t, err)
	reps, done := s.Send(ctx, msg)
	defer done()

	rep := <-reps
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	require.NoError(t, err)
	switch pl.(type) {
	case *payload.ResultInfo:
		return pl, rec
	case *payload.Error:
		return pl, rec
	default:
		t.Fatalf("received unexpected reply %T", pl)
	}
	return pl, rec
}

func getObject(ctx context.Context, t *testing.T, s *Server, objectID insolar.ID) (payload.Payload, payload.Payload) {
	msg, err := payload.NewMessage(&payload.GetObject{
		ObjectID: objectID,
	})
	require.NoError(t, err)
	reps, d := s.Send(ctx, msg)
	defer d()

	var (
		lifeline, state payload.Payload
	)
	done := func() bool {
		return lifeline != nil && state != nil
	}
	for rep := range reps {
		pl, err := payload.UnmarshalFromMeta(rep.Payload)
		require.NoError(t, err)
		switch pl.(type) {
		case *payload.Index:
			lifeline = pl
		case *payload.State:
			state = pl
		case *payload.Error:
			return pl, nil
		default:
			t.Fatalf("received unexpected reply %T", pl)
		}

		if done() {
			break
		}
	}
	require.True(t, done())
	return lifeline, state
}

func requirePayloadNotError(t *testing.T, pl payload.Payload) {
	if err, ok := pl.(*payload.Error); ok {
		t.Fatal(err)
	}
}

func requireGetObject(ctx context.Context, t *testing.T, s *Server, objectID insolar.ID) (record.Lifeline, record.Material) {
	lifelinePL, statePL := getObject(ctx, t, s, objectID)
	requirePayloadNotError(t, lifelinePL)
	requirePayloadNotError(t, statePL)

	lifeline := record.Lifeline{}
	err := lifeline.Unmarshal(lifelinePL.(*payload.Index).Index)
	require.NoError(t, err)

	state := record.Material{}
	err = state.Unmarshal(statePL.(*payload.State).Record)
	require.NoError(t, err)

	return lifeline, state
}
