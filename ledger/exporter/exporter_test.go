/*
 *    Copyright 2019 Insolar Technologies
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
	"context"
	"encoding/json"
	"strconv"
	"testing"

	base58 "github.com/jbenet/go-base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/internal/jet"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/platformpolicy"
)

type exporterSuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()

	pulseTracker  storage.PulseTracker
	objectStorage storage.ObjectStorage
	jetStorage    jet.JetStorage
	pulseStorage  *storage.PulseStorage

	exporter *Exporter
	jetID    core.RecordID
}

func NewExporterSuite() *exporterSuite {
	return &exporterSuite{
		Suite: suite.Suite{},
	}
}

// Init and run suite
func TestExporter(t *testing.T) {
	suite.Run(t, NewExporterSuite())
}

func (s *exporterSuite) BeforeTest(suiteName, testName string) {
	s.cm = &component.Manager{}
	s.ctx = inslogger.TestContext(s.T())
	s.jetID = core.TODOJetID

	db, cleaner := storagetest.TmpDB(s.ctx, s.T())
	s.cleaner = cleaner
	s.pulseTracker = storage.NewPulseTracker()
	s.objectStorage = storage.NewObjectStorage()
	s.jetStorage = jet.NewJetStorage()
	s.pulseStorage = storage.NewPulseStorage()
	s.exporter = NewExporter(configuration.Exporter{ExportLag: 0})

	s.cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		db,
		s.pulseTracker,
		s.objectStorage,
		s.jetStorage,
		s.pulseStorage,
		s.exporter,
	)

	err := s.cm.Init(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager init failed", err)
	}
	err = s.cm.Start(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager start failed", err)
	}
}

func (s *exporterSuite) AfterTest(suiteName, testName string) {
	err := s.cm.Stop(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager stop failed", err)
	}
	s.cleaner()
}

func (s *exporterSuite) TestExporter_Export() {
	for i := 1; i <= 3; i++ {
		err := s.pulseTracker.AddPulse(
			s.ctx,
			core.Pulse{
				PulseNumber:     core.FirstPulseNumber + 10*core.PulseNumber(i),
				PrevPulseNumber: core.FirstPulseNumber + 10*core.PulseNumber(i-1),
				PulseTimestamp:  10 * int64(i+1),
			},
		)
		require.NoError(s.T(), err)
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
	blobID, err := s.objectStorage.SetBlob(s.ctx, s.jetID, core.FirstPulseNumber+10, mem)
	require.NoError(s.T(), err)
	_, err = s.objectStorage.SetRecord(s.ctx, s.jetID, core.FirstPulseNumber+10, &object.GenesisRecord{})
	require.NoError(s.T(), err)
	objectID, err := s.objectStorage.SetRecord(s.ctx, s.jetID, core.FirstPulseNumber+10, &object.ObjectActivateRecord{
		ObjectStateRecord: object.ObjectStateRecord{
			Memory: blobID,
		},
		IsDelegate: true,
	})
	msg := &message.CallConstructor{}
	var parcel core.Parcel = &message.Parcel{Msg: msg}

	msgHash := platformpolicy.NewPlatformCryptographyScheme().IntegrityHasher().Hash(message.ToBytes(msg))
	requestID, err := s.objectStorage.SetRecord(s.ctx, s.jetID, core.FirstPulseNumber+10, &object.RequestRecord{
		MessageHash: msgHash,
		Parcel:      message.ParcelToBytes(parcel),
	})
	require.NoError(s.T(), err)

	result, err := s.exporter.Export(s.ctx, 0, 15)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, len(result.Data))
	assert.Equal(s.T(), 2, result.Size)
	assert.Nil(s.T(), result.NextFrom)

	result, err = s.exporter.Export(s.ctx, 0, 2)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, len(result.Data))
	assert.Equal(s.T(), 2, result.Size)
	assert.Equal(s.T(), core.FirstPulseNumber+20, int(*result.NextFrom))
	_, err = json.Marshal(result)
	assert.NoError(s.T(), err)

	pulse := result.Data[strconv.FormatUint(uint64(core.FirstPulseNumber), 10)].([]*pulseData)[0].Pulse
	assert.Equal(s.T(), core.FirstPulseNumber, int(pulse.PulseNumber))
	assert.Equal(s.T(), int64(0), pulse.PulseTimestamp)
	pulse = result.Data[strconv.FormatUint(uint64(core.FirstPulseNumber+10), 10)].([]*pulseData)[0].Pulse
	assert.Equal(s.T(), core.FirstPulseNumber+10, int(pulse.PulseNumber))
	assert.Equal(s.T(), int64(20), pulse.PulseTimestamp)

	records := result.Data[strconv.FormatUint(uint64(core.FirstPulseNumber+10), 10)].([]*pulseData)[0].Records
	obj, ok := records[base58.Encode(objectID[:])]
	if assert.True(s.T(), ok, "object not found by ID") {
		assert.Equal(s.T(), "TypeActivate", obj.Type)
		assert.Equal(s.T(), true, obj.Data.(*object.ObjectActivateRecord).IsDelegate)
		assert.Equal(s.T(), "objectValue", obj.Payload["Memory"].(payload)["Field"])
	}

	request, ok := records[base58.Encode(requestID[:])]
	if assert.True(s.T(), ok, "request not found by ID") {
		assert.Equal(s.T(), "TypeCallRequest", request.Type)
		assert.Equal(s.T(), msgHash, request.Data.(*object.RequestRecord).MessageHash)
		assert.Equal(s.T(), core.TypeCallConstructor.String(), request.Payload["Type"])
	}

	_, err = s.exporter.Export(s.ctx, 100000, 2)
	require.Error(s.T(), err, "From-pulse should be smaller (or equal) current-pulse")

	_, err = s.exporter.Export(s.ctx, 60000, 2)
	require.NoError(s.T(), err, "From-pulse should be smaller (or equal) current-pulse")
}

func (s *exporterSuite) TestExporter_ExportGetBlobFailed() {
	for i := 1; i <= 3; i++ {
		err := s.pulseTracker.AddPulse(
			s.ctx,
			core.Pulse{
				PulseNumber:     core.FirstPulseNumber + 10*core.PulseNumber(i),
				PrevPulseNumber: core.FirstPulseNumber + 10*core.PulseNumber(i-1),
				PulseTimestamp:  10 * int64(i+1),
			},
		)
		require.NoError(s.T(), err)
	}

	_, err := s.objectStorage.SetRecord(s.ctx, s.jetID, core.FirstPulseNumber+10, &object.ObjectActivateRecord{
		ObjectStateRecord: object.ObjectStateRecord{
			Memory: &core.RecordID{},
		},
		IsDelegate: true,
	})
	require.NoError(s.T(), err)

	result, err := s.exporter.Export(s.ctx, core.FirstPulseNumber+10, 10)
	assert.Equal(s.T(), 1, len(result.Data))
	assert.Equal(s.T(), 1, result.Size)
}
