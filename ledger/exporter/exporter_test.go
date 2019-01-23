/*
 *    Copyright 2019 Insolar
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

package exporter

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/record"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	base58 "github.com/jbenet/go-base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ugorji/go/codec"
)

func TestExporter_Export(t *testing.T) {
	ctx := inslogger.TestContext(t)
	db, clean := storagetest.TmpDB(ctx, t)
	defer clean()
	jetID := core.TODOJetID
	ps := storage.NewPulseStorage(db)
	exporter := NewExporter(db, ps, configuration.Exporter{ExportLag: 0})

	for i := 1; i <= 3; i++ {
		err := db.AddPulse(
			ctx,
			core.Pulse{
				PulseNumber:     core.FirstPulseNumber + 10*core.PulseNumber(i),
				PrevPulseNumber: core.FirstPulseNumber + 10*core.PulseNumber(i-1),
				PulseTimestamp:  10 * int64(i+1),
			},
		)
		require.NoError(t, err)
	}

	type testData struct {
		Field string
		Data  struct {
			Field string
		}
	}
	mem := make([]byte, 0)
	blobData := testData{Field: "objectValue"}
	blobData.Data.Field = "anotherValue"
	codec.NewEncoderBytes(&mem, &codec.CborHandle{}).MustEncode(blobData)
	blobID, err := db.SetBlob(ctx, jetID, core.FirstPulseNumber+10, mem)
	require.NoError(t, err)
	_, err = db.SetRecord(ctx, jetID, core.FirstPulseNumber+10, &record.GenesisRecord{})
	require.NoError(t, err)
	objectID, err := db.SetRecord(ctx, jetID, core.FirstPulseNumber+10, &record.ObjectActivateRecord{
		ObjectStateRecord: record.ObjectStateRecord{
			Memory: blobID,
		},
		IsDelegate: true,
	})
	pl := message.ToBytes(&message.CallConstructor{})
	requestID, err := db.SetRecord(ctx, jetID, core.FirstPulseNumber+10, &record.RequestRecord{
		Payload: pl,
	})
	require.NoError(t, err)

	result, err := exporter.Export(ctx, 0, 15)
	require.NoError(t, err)
	assert.Equal(t, 2, len(result.Data))
	assert.Equal(t, 2, result.Size)
	assert.Nil(t, result.NextFrom)

	result, err = exporter.Export(ctx, 0, 2)
	require.NoError(t, err)
	assert.Equal(t, 2, len(result.Data))
	assert.Equal(t, 2, result.Size)
	assert.Equal(t, core.FirstPulseNumber+20, int(*result.NextFrom))
	_, err = json.Marshal(result)
	assert.NoError(t, err)

	pulse := result.Data[strconv.FormatUint(uint64(core.FirstPulseNumber), 10)].([]*pulseData)[0].Pulse
	assert.Equal(t, core.FirstPulseNumber, int(pulse.PulseNumber))
	assert.Equal(t, int64(0), pulse.PulseTimestamp)
	pulse = result.Data[strconv.FormatUint(uint64(core.FirstPulseNumber+10), 10)].([]*pulseData)[0].Pulse
	assert.Equal(t, core.FirstPulseNumber+10, int(pulse.PulseNumber))
	assert.Equal(t, int64(20), pulse.PulseTimestamp)

	records := result.Data[strconv.FormatUint(uint64(core.FirstPulseNumber+10), 10)].([]*pulseData)[0].Records
	object, ok := records[base58.Encode(objectID[:])]
	if assert.True(t, ok, "object not found by ID") {
		assert.Equal(t, "TypeActivate", object.Type)
		assert.Equal(t, true, object.Data.(*record.ObjectActivateRecord).IsDelegate)
		assert.Equal(t, "objectValue", object.Payload["Memory"].(payload)["Field"])
	}

	request, ok := records[base58.Encode(requestID[:])]
	if assert.True(t, ok, "request not found by ID") {
		assert.Equal(t, "TypeCallRequest", request.Type)
		assert.Equal(t, pl, request.Data.(*record.RequestRecord).Payload)
		assert.Equal(t, core.TypeCallConstructor.String(), request.Payload["Type"])
	}
}
