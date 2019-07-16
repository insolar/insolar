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
	t.Parallel()

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

	// Creating root reason request.
	reasonID, _ := setIncomingRequest(t, s, gen.ID(), gen.ID(), record.CTSaveAsChild)

	// Save and check code.
	{
		codeID, sent := setCode(t, s)
		received := getCode(t, s, codeID)
		require.Equal(t, &sent, received.Virtual)
	}
	var objectID insolar.ID
	// Set and get request.
	{
		id, sent := setIncomingRequest(t, s, gen.ID(), reasonID, record.CTSaveAsChild)
		received := getRequest(t, s, id)
		require.Equal(t, sent, received)
		objectID = id
	}
	// Activate and check object.
	{
		state := activateObject(t, s, objectID)
		_, material := getObject(t, s, objectID)
		require.Equal(t, &state, material.Virtual)
	}
	// Amend and check object.
	{
		requestID, _ := setIncomingRequest(t, s, objectID, reasonID, record.CTMethod)
		state := amendObject(t, s, objectID, requestID)
		_, material := getObject(t, s, objectID)
		require.Equal(t, &state, material.Virtual)
	}
	// Deactivate and check object.
	{
		deactivateObject(t, s, objectID)
		received := make(chan payload.Payload)
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

// func Test_Concurrency(t *testing.T) {
// 	t.Skip()
// 	t.Parallel()
//
// 	ctx := inslogger.TestContext(t)
// 	cfg := DefaultLightConfig()
// 	s, err := NewServer(ctx, cfg)
// 	require.NoError(t, err)
//
// 	// First pulse goes in storage then interrupts.
// 	s.Pulse(ctx)
// 	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
// 	s.Pulse(ctx)
//
// 	runner := func() {
// 		// Save and check code.
// 		{
// 			codeID, sent := setCode(t, s)
// 			received := getCode(t, s, codeID)
// 			assert.Equal(t, &sent, received.Virtual)
// 		}
// 		var objectID insolar.ID
// 		// Set and get request.
// 		{
// 			id, sent := setIncomingRequest(t, s, record.CTSaveAsChild)
// 			received := getRequest(t, s, id)
// 			assert.Equal(t, sent, received)
// 			objectID = id
// 		}
// 		// Activate and check object.
// 		{
// 			state := activateObject(t, s, objectID)
// 			_, material := getObject(t, s, objectID)
// 			require.Equal(t, &state, material.Virtual)
// 		}
// 		// Amend and check object.
// 		{
// 			state := amendObject(t, s, objectID, gen.ID())
// 			_, material := getObject(t, s, objectID)
// 			require.Equal(t, &state, material.Virtual)
// 		}
// 		// Deactivate and check object.
// 		{
// 			deactivateObject(t, s, objectID)
// 			received := make(chan payload.Payload)
// 			s.Receive(func(meta payload.Meta, pl payload.Payload) {
// 				received <- pl
// 			})
// 			s.Send(&payload.GetObject{
// 				ObjectID: objectID,
// 			})
// 			pl := <-received
// 			_, ok := pl.(*payload.Error)
// 			assert.True(t, ok)
// 		}
// 	}
//
// 	count := 100
// 	var wg sync.WaitGroup
// 	wg.Add(count)
// 	for i := 0; i < count; i++ {
// 		go func() {
// 			runner()
// 			wg.Done()
// 		}()
// 	}
//
// 	wg.Wait()
// }

func setCode(t *testing.T, s *Server) (insolar.ID, record.Virtual) {
	received := make(chan payload.ID)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		switch p := pl.(type) {
		case *payload.Error:
			panic(p.Text)
		case *payload.ID:
			received <- *p
		}
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
	return pl.ID, rec
}

func getCode(t *testing.T, s *Server, id insolar.ID) record.Material {
	received := make(chan payload.Payload)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		switch p := pl.(type) {
		case *payload.Error:
			panic(p.Text)
		case *payload.Code:
			received <- pl
		}
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

func getRequest(t *testing.T, s *Server, requestID insolar.ID) record.Virtual {
	received := make(chan payload.Request)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		switch p := pl.(type) {
		case *payload.Error:
			panic(p.Text)
		case *payload.Request:
			received <- *p
		}
	})

	s.Send(&payload.GetRequest{
		RequestID: requestID,
	})
	pl := <-received
	return pl.Request
}

func setIncomingRequest(t *testing.T, s *Server, objectID, reasonID insolar.ID, ct record.CallType) (insolar.ID, record.Virtual) {
	received := make(chan payload.RequestInfo)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		switch p := pl.(type) {
		case *payload.Error:
			panic(p.Text)
		case *payload.RequestInfo:
			received <- *p
		}
	})

	args := make([]byte, 100)
	_, err := rand.Read(args)
	require.NoError(t, err)
	rec := record.Wrap(record.IncomingRequest{
		Object:    insolar.NewReference(objectID),
		Arguments: args,
		CallType:  ct,
		Reason:    *insolar.NewReference(reasonID),
	})
	s.Send(&payload.SetIncomingRequest{
		Request: rec,
	})
	pl := <-received
	return pl.RequestID, rec
}

func activateObject(t *testing.T, s *Server, objectID insolar.ID) record.Virtual {
	received := make(chan payload.ResultInfo)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		switch p := pl.(type) {
		case *payload.Error:
			panic(p.Text)
		case *payload.ResultInfo:
			received <- *p
		}
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
	<-received
	return rec
}

func amendObject(t *testing.T, s *Server, objectID, requestID insolar.ID) record.Virtual {
	received := make(chan payload.ResultInfo)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		switch p := pl.(type) {
		case *payload.Error:
			panic(p.Text)
		case *payload.ResultInfo:
			received <- *p
		}
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
		Request: *insolar.NewReference(requestID),
		Object:  objectID,
		Payload: res,
	})
	resBuf, err := resultRecord.Marshal()
	require.NoError(t, err)
	s.Send(&payload.Update{
		Record: buf,
		Result: resBuf,
	})
	<-received
	return rec
}

func deactivateObject(t *testing.T, s *Server, objectID insolar.ID) record.Virtual {
	received := make(chan payload.ResultInfo)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		switch p := pl.(type) {
		case *payload.Error:
			panic(p.Text)
		case *payload.ResultInfo:
			received <- *p
		}
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
	<-received
	return rec
}

func getObject(t *testing.T, s *Server, objectID insolar.ID) (record.Lifeline, record.Material) {
	received := make(chan payload.Payload)
	s.Receive(func(meta payload.Meta, pl payload.Payload) {
		switch p := pl.(type) {
		case *payload.Error:
			panic(p.Text)
		case *payload.Index:
			received <- pl
		case *payload.State:
			received <- pl
		}
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
