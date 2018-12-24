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
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func TestGetMessageHash(t *testing.T) {
	pcs := platformpolicy.NewPlatformCryptographyScheme()
	require.Equal(t, 64, len(GetMessageHash(pcs, &message.Parcel{Msg: &message.GenesisRequest{}})))
}

func TestTape_SetReply(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	ctx := inslogger.TestContext(t)
	pn := core.PulseNumber(1)

	rep := reply.Object{Memory: []byte{9, 9, 9}}
	rd, err := reply.Serialize(&rep)
	buff := new(bytes.Buffer)
	_, err = buff.ReadFrom(rd)
	require.NoError(t, err)

	tp := newMemoryTape(pn)
	err = tp.SetReply(ctx, []byte{4, 5, 6}, &rep)
	require.NoError(t, err)
}

func TestTape_GetReply(t *testing.T) {
	// mc := minimock.NewController(t)
	// defer mc.Finish()

	// ctx := inslogger.TestContext(t)
	// pn := core.PulseNumber(1)

	// expectedRep := reply.Object{Memory: []byte{42}}
	// rd, err := reply.Serialize(&expectedRep)
	// buff := new(bytes.Buffer)
	// _, err = buff.ReadFrom(rd)
	// require.NoError(t, err)

	// tp := newMemoryTape(pn)
	// rep, err := tp.GetReply(ctx, []byte{4, 5, 6})
	// require.NoError(t, err)
	// require.Equal(t, expectedRep, *rep.(*reply.Object))
}

// func TestTape_Write(t *testing.T) {
// 	mc := minimock.NewController(t)
// 	defer mc.Finish()

// 	// Prepare test data.
// 	ctx := inslogger.TestContext(t)
// 	pn := core.PulseNumber(core.FirstPulseNumber + 1000)
// 	tp := newMemoryTape(pn)

// 	// Write buffer from storageTape.
// 	buff := bytes.NewBuffer(nil)
// 	expected.GetReply()
// 	err := tp.Write(ctx, buff)
// 	require.NoError(t, err)

// 	r, err := newMemoryTapeFromReader(ctx, bytes.NewReader(buff.Bytes()))
// 	require.NoError(t, err)

// 	// Write expected buffer.
// 	// expectedBuff := bytes.NewBuffer(nil)
// 	// enc := gob.NewEncoder(expectedBuff)
// 	// enc.Encode(tp.pulse)
// 	// enc.Encode(core.KV{K: []byte{1}, V: []byte{2}})
// 	// enc.Encode(core.KV{K: []byte{3}, V: []byte{4}})

// 	require.Equal(t, expectedBuff.Bytes(), buff.Bytes())
// }
