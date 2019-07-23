//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package integration_test

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/stretchr/testify/require"
)

func callSetCode(ctx context.Context, t *testing.T, s *Server) (payload.Payload, record.Virtual) {
	code := make([]byte, 100)
	_, err := rand.Read(code)
	require.NoError(t, err)
	rec := record.Wrap(record.Code{Code: code})
	buf, err := rec.Marshal()
	require.NoError(t, err)
	reps, done := s.Send(ctx, &payload.SetCode{
		Record: buf,
	})
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

func callGetCode(ctx context.Context, t *testing.T, s *Server, id insolar.ID) payload.Payload {
	reps, done := s.Send(ctx, &payload.GetCode{
		CodeID: id,
	})
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

func callSetIncomingRequest(
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
	reps, done := s.Send(ctx, &payload.SetIncomingRequest{
		Request: rec,
	})
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

func callGetRequest(ctx context.Context, t *testing.T, s *Server, requestID insolar.ID) payload.Payload {
	reps, done := s.Send(ctx, &payload.GetRequest{
		RequestID: requestID,
	})
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

func callActivateObject(ctx context.Context, t *testing.T, s *Server, objectID insolar.ID) (payload.Payload, record.Virtual) {
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

	reps, done := s.Send(ctx, &payload.Activate{
		Record: buf,
		Result: resBuf,
	})
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

func callAmendObject(ctx context.Context, t *testing.T, s *Server, objectID, requestID insolar.ID) (payload.Payload, record.Virtual) {
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

	reps, done := s.Send(ctx, &payload.Update{
		Record: buf,
		Result: resBuf,
	})
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

func callDeactivateObject(ctx context.Context, t *testing.T, s *Server, objectID, requestID insolar.ID) (payload.Payload, record.Virtual) {
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

	reps, done := s.Send(ctx, &payload.Deactivate{
		Record: buf,
		Result: resBuf,
	})
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

func callGetObject(ctx context.Context, t *testing.T, s *Server, objectID insolar.ID) (payload.Payload, payload.Payload) {
	reps, d := s.Send(ctx, &payload.GetObject{
		ObjectID: objectID,
	})
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
