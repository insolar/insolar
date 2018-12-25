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
	"fmt"
	"testing"

	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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

func TestTape_SetGet(t *testing.T) {
	ctx := inslogger.TestContext(t)
	pn := core.PulseNumber(1)

	rep := reply.Object{Memory: []byte{9, 9, 9}}
	rd, err := reply.Serialize(&rep)
	var buf bytes.Buffer
	_, err = buf.ReadFrom(rd)
	require.NoError(t, err)

	msgHash := []byte{4, 5, 6}
	tp := newMemoryTape(pn)
	err = tp.Set(ctx, msgHash, &rep, nil)
	require.NoError(t, err)

	item, err := tp.Get(ctx, msgHash)
	require.NoError(t, err)

	assert.Equal(t, &rep, item.Reply)
	assert.Nil(t, item.Error)
}

func TestTape_SetGet_WithError(t *testing.T) {
	ctx := inslogger.TestContext(t)
	pn := core.PulseNumber(1)

	expectedErr := fmt.Errorf("Error")
	gotErr := fmt.Errorf("Error")

	msgHash := []byte{4, 5, 6}
	tp := newMemoryTape(pn)
	err := tp.Set(ctx, msgHash, nil, gotErr)
	require.NoError(t, err)

	item, err := tp.Get(ctx, msgHash)
	require.NoError(t, err)

	assert.Nil(t, item.Reply)
	assert.Equal(t, expectedErr, item.Error)
}

func TestTape_Write(t *testing.T) {
	// 	mc := minimock.NewController(t)
	// 	defer mc.Finish()

	// 	// Prepare test data.
	ctx := inslogger.TestContext(t)
	pn := core.PulseNumber(core.FirstPulseNumber + 1000)

	tp := newMemoryTape(pn)

	// 	// Write buffer from storageTape.
	// 	expected.GetReply()
	expected := []struct {
		msgHash []byte
		item    TapeItem
	}{
		{
			msgHash: []byte{4, 5, 6},
			item: TapeItem{
				Reply: &reply.Object{Memory: []byte{9, 9, 9}},
			},
		},
		{
			msgHash: []byte{4, 5, 7},
			item: TapeItem{
				Error: errors.New("send failed"),
			},
		},
	}

	for _, tCase := range expected {
		err := tp.Set(ctx, tCase.msgHash, tCase.item.Reply, tCase.item.Error)
		require.NoError(t, err)
	}
	var buf bytes.Buffer
	err := tp.Write(ctx, &buf)
	require.NoError(t, err)

	rTape, err := newMemoryTapeFromReader(ctx, bytes.NewReader(buf.Bytes()))
	require.NoError(t, err)

	for _, tCase := range expected {
		gotItem, err := rTape.Get(ctx, tCase.msgHash)
		// fmt.Println("err =>", err)
		require.NoError(t, err)
		assert.Equal(t, tCase.item.Reply, gotItem.Reply)
		if tCase.item.Error == nil {
			assert.Nil(t, gotItem.Error)
		} else {
			assert.Equal(t, tCase.item.Error.Error(), gotItem.Error.Error())
		}
		// fmt.Printf("gotItem => %+v\n", gotItem)
	}
}
