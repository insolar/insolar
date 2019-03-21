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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
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
	jetStorage    jet.Storage
	pulseStorage  *storage.PulseStorage

	exporter    *Exporter
	jetRecordID insolar.RecordID
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
	s.jetRecordID = insolar.NewID(core.ZeroJetID)

	db, cleaner := storagetest.TmpDB(s.ctx, s.T())
	s.cleaner = cleaner
	s.pulseTracker = storage.NewPulseTracker()
	s.objectStorage = storage.NewObjectStorage()
	s.jetStorage = jet.NewStore()
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
	var (
		Pulse1  insolar.PulseNumber = insolar.FirstPulseNumber
		Pulse10 insolar.PulseNumber = insolar.FirstPulseNumber + 10
		Pulse20 insolar.PulseNumber = insolar.FirstPulseNumber + 20
	)

	for i := 1; i <= 3; i++ {
		err := s.pulseTracker.AddPulse(
			s.ctx,
			insolar.Pulse{
				PulseNumber:     insolar.FirstPulseNumber + 10*insolar.PulseNumber(i),
				PrevPulseNumber: insolar.FirstPulseNumber + 10*insolar.PulseNumber(i-1),
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

	blobID, err := s.objectStorage.SetBlob(s.ctx, s.jetRecordID, Pulse10, mem)
	require.NoError(s.T(), err)

	_, err = s.objectStorage.SetRecord(s.ctx, s.jetRecordID, Pulse10, &object.GenesisRecord{})
	require.NoError(s.T(), err)

	objectID, err := s.objectStorage.SetRecord(s.ctx, s.jetRecordID, Pulse10, &object.ObjectActivateRecord{
		ObjectStateRecord: object.ObjectStateRecord{
			Memory: blobID,
		},
		IsDelegate: true,
	})
	require.NoError(s.T(), err)
	objectID58 := base58.Encode(objectID[:])

	msg := &message.CallConstructor{}
	var parcel insolar.Parcel = &message.Parcel{Msg: msg}

	msgHash := platformpolicy.NewPlatformCryptographyScheme().IntegrityHasher().Hash(message.ToBytes(msg))
	requestID, err := s.objectStorage.SetRecord(
		s.ctx,
		s.jetRecordID,
		Pulse10,
		&object.RequestRecord{
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
	assert.Equal(s.T(), int(Pulse20), int(*result.NextFrom))

	_, err = json.Marshal(result)
	assert.NoError(s.T(), err)

	pn2str := func(pn insolar.PulseNumber) string {
		return strconv.FormatUint(uint64(pn), 10)
	}

	pulse1data := result.Data[pn2str(Pulse1)].([]*pulseData)
	require.NotEmpty(s.T(), pulse1data, "pulse 1 data should not be empty")
	assert.Equal(s.T(), int(Pulse1), int(pulse1data[0].Pulse.PulseNumber))
	assert.Equal(s.T(), int64(0), pulse1data[0].Pulse.PulseTimestamp)

	pulse10data := result.Data[pn2str(Pulse10)].([]*pulseData)
	require.NotEmpty(s.T(), pulse10data, "pulse 10 data should not be empty")
	assert.Equal(s.T(), int(Pulse10), int(pulse10data[0].Pulse.PulseNumber))
	assert.Equal(s.T(), int64(20), pulse10data[0].Pulse.PulseTimestamp)

	records := pulse10data[0].Records
	requestID58 := base58.Encode(requestID[:])

	obj, ok := records[objectID58]
	if assert.True(s.T(), ok, "object not found by ID") {
		assert.Equal(s.T(), "TypeActivate", obj.Type)
		assert.Equal(s.T(), true, obj.Data.(*object.ObjectActivateRecord).IsDelegate)
		assert.Equal(s.T(), "objectValue", obj.Payload["Memory"].(payload)["Field"])
	}

	request, ok := records[requestID58]
	if assert.True(s.T(), ok, "request not found by ID") {
		assert.Equal(s.T(), "TypeCallRequest", request.Type)
		assert.Equal(s.T(), msgHash, request.Data.(*object.RequestRecord).MessageHash)
		assert.Equal(s.T(), insolar.TypeCallConstructor.String(), request.Payload["Type"])
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
			insolar.Pulse{
				PulseNumber:     insolar.FirstPulseNumber + 10*insolar.PulseNumber(i),
				PrevPulseNumber: insolar.FirstPulseNumber + 10*insolar.PulseNumber(i-1),
				PulseTimestamp:  10 * int64(i+1),
			},
		)
		require.NoError(s.T(), err)
	}

	_, err := s.objectStorage.SetRecord(s.ctx, s.jetRecordID, core.FirstPulseNumber+10, &object.ObjectActivateRecord{
		ObjectStateRecord: object.ObjectStateRecord{
			Memory: &insolar.ID{},
		},
		IsDelegate: true,
	})
	require.NoError(s.T(), err)

	result, err := s.exporter.Export(s.ctx, insolar.FirstPulseNumber+10, 10)
	assert.Equal(s.T(), 1, len(result.Data))
	assert.Equal(s.T(), 1, result.Size)
}
