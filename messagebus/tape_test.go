/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package messagebus

import (
	"bytes"
	"context"
	"encoding/gob"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
)

func TestGetMessageHash(t *testing.T) {
	pcs := platformpolicy.NewPlatformCryptographyScheme()
	require.Equal(t, 64, len(GetMessageHash(pcs, &message.Parcel{Msg: &message.GenesisRequest{}})))
}

func TestTape_SetReply(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	ctx := inslogger.TestContext(t)
	rep := reply.Object{Memory: []byte{9, 9, 9}}
	rd, err := reply.Serialize(&rep)
	buff := new(bytes.Buffer)
	_, err = buff.ReadFrom(rd)
	require.NoError(t, err)

	id := uuid.UUID{1, 2, 3}
	ls := testutils.NewLocalStorageMock(mc)
	ls.SetMock.Expect(ctx, 1, bytes.Join([][]byte{id[:], {4, 5, 6}}, nil), buff.Bytes()).Return(nil)

	tp := storageTape{ls: ls, pulse: 1, id: id}
	err = tp.SetReply(ctx, []byte{4, 5, 6}, &rep)
	require.NoError(t, err)
}

func TestTape_GetReply(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	ctx := inslogger.TestContext(t)
	expectedRep := reply.Object{Memory: []byte{42}}
	rd, err := reply.Serialize(&expectedRep)
	buff := new(bytes.Buffer)
	_, err = buff.ReadFrom(rd)
	require.NoError(t, err)

	id := uuid.UUID{1, 2, 3}
	ls := testutils.NewLocalStorageMock(mc)
	ls.GetMock.Expect(ctx, 1, bytes.Join([][]byte{id[:], {4, 5, 6}}, nil)).Return(buff.Bytes(), nil)

	tp := storageTape{ls: ls, pulse: 1, id: id}
	rep, err := tp.GetReply(ctx, []byte{4, 5, 6})
	require.NoError(t, err)
	require.Equal(t, expectedRep, *rep.(*reply.Object))
}

func TestTape_Write(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	// Prepare test data.
	ctx := inslogger.TestContext(t)
	ls := testutils.NewLocalStorageMock(mc)
	tp, err := newStorageTape(ls, 42)
	require.NoError(t, err)
	ls.IterateFunc = func(ctx context.Context, pulse core.PulseNumber, prefix []byte, handler func(k, v []byte) error) error {
		err := handler(bytes.Join([][]byte{tp.id[:], {1}}, nil), []byte{2})
		require.NoError(t, err)
		err = handler(bytes.Join([][]byte{tp.id[:], {3}}, nil), []byte{4})
		require.NoError(t, err)
		return nil
	}

	// Write buffer from storageTape.
	buff := bytes.NewBuffer(nil)
	err = tp.Write(ctx, buff)
	require.NoError(t, err)

	// Write expected buffer.
	expectedBuff := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(expectedBuff)
	enc.Encode(tp.pulse)
	enc.Encode(tp.id)
	enc.Encode(couple{Key: []byte{1}, Value: []byte{2}})
	enc.Encode(couple{Key: []byte{3}, Value: []byte{4}})

	require.Equal(t, expectedBuff.Bytes(), buff.Bytes())
}

func TestNewTapeFromReader(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	// Prepare test data.
	ctx := inslogger.TestContext(t)
	id := uuid.UUID{1, 2, 3}
	buff := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buff)
	enc.Encode(core.PulseNumber(42))
	enc.Encode(id)
	enc.Encode(couple{Key: []byte{1}, Value: []byte{2}})
	enc.Encode(couple{Key: []byte{3}, Value: []byte{4}})

	var values []couple
	expectedValues := []couple{
		{Key: bytes.Join([][]byte{id[:], {1}}, nil), Value: []byte{2}},
		{Key: bytes.Join([][]byte{id[:], {3}}, nil), Value: []byte{4}},
	}
	ls := testutils.NewLocalStorageMock(mc)
	ls.SetFunc = func(ctx context.Context, pulse core.PulseNumber, k, v []byte) (r error) {
		values = append(values, couple{Key: k, Value: v})
		return nil
	}

	_, err := newStorageTapeFromReader(ctx, ls, bytes.NewBuffer(buff.Bytes()))
	require.NoError(t, err)
	require.Equal(t, expectedValues, values)
}
