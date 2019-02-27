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
	"bytes"
	"context"
	"math"
	"strconv"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/record"
	"github.com/jbenet/go-base58"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

// Exporter provides methods for fetching data view from storage.
type Exporter struct {
	DB            storage.DBContext     `inject:""`
	JetStorage    storage.JetStorage    `inject:""`
	ObjectStorage storage.ObjectStorage `inject:""`
	PulseTracker  storage.PulseTracker  `inject:""`
	PulseStorage  core.PulseStorage     `inject:""`

	cfg configuration.Exporter
}

// NewExporter creates new StorageExporter instance.
func NewExporter(cfg configuration.Exporter) *Exporter {
	return &Exporter{
		cfg: cfg,
	}
}

type payload map[string]interface{}

// MarshalJSON serializes payload into JSON.
func (p payload) MarshalJSON() ([]byte, error) {
	var buf []byte
	err := codec.NewEncoderBytes(&buf, &codec.JsonHandle{}).Encode(&p)
	return buf, err
}

type recordData struct {
	Type    string
	Data    record.Record
	Payload payload
}

type recordsData map[string]recordData

type pulseData struct {
	Records recordsData
	Pulse   core.Pulse
	JetID   core.RecordID
}

// Export returns data view from storage.
func (e *Exporter) Export(ctx context.Context, fromPulse core.PulseNumber, size int) (*core.StorageExportResult, error) {
	result := core.StorageExportResult{Data: map[string]interface{}{}}

	jetIDs, err := e.JetStorage.GetJets(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch jets")
	}

	currentPulse, err := e.PulseStorage.Current(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get current pulse data")
	}

	counter := 0
	fromPulsePN := core.PulseNumber(math.Max(float64(fromPulse), float64(core.GenesisPulse.PulseNumber)))

	if fromPulsePN > currentPulse.PulseNumber {
		return nil, errors.Errorf("failed to fetch data: from-pulse[%v] > current-pulse[%v]",
			fromPulsePN, currentPulse.PulseNumber)
	}

	_, err = e.PulseTracker.GetPulse(ctx, fromPulsePN)
	if err != nil {
		tryPulse, err := e.PulseTracker.GetPulse(ctx, core.GenesisPulse.PulseNumber)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch genesis pulse data")
		}

		for fromPulsePN > *tryPulse.Next {
			tryPulse, err = e.PulseTracker.GetPulse(ctx, *tryPulse.Next)
			if err != nil {
				return nil, errors.Wrap(err, "failed to iterate through first pulses")
			}
		}
		fromPulsePN = *tryPulse.Next
	}

	iterPulse := &fromPulsePN
	for iterPulse != nil && counter < size {
		pulse, err := e.PulseTracker.GetPulse(ctx, *iterPulse)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch pulse data")
		}

		// We don't need data from current pulse, because of
		// not all data for this pulse is persisted at this moment
		// @sergey.morozov 20.01.18 - Blocks are synced to Heavy node with a lag.
		// We can't reliably predict this lag so we add threshold of N seconds.
		if pulse.Pulse.PulseNumber >= (currentPulse.PrevPulseNumber - core.PulseNumber(e.cfg.ExportLag)) {
			iterPulse = nil
			break
		}

		var data []*pulseData
		for jetID := range jetIDs {
			fetchedData, err := e.exportPulse(ctx, jetID, &pulse.Pulse)
			if err != nil {
				return nil, err
			}
			data = append(data, fetchedData)
		}

		result.Data[strconv.FormatUint(uint64(pulse.Pulse.PulseNumber), 10)] = data

		iterPulse = pulse.Next
		counter++
	}

	result.Size = counter
	result.NextFrom = iterPulse

	return &result, nil
}

func (e *Exporter) exportPulse(ctx context.Context, jetID core.RecordID, pulse *core.Pulse) (*pulseData, error) {
	records := recordsData{}
	err := e.DB.IterateRecordsOnPulse(ctx, jetID, pulse.PulseNumber, func(id core.RecordID, rec record.Record) error {
		pl, err := e.getPayload(ctx, jetID, rec)
		if err != nil {
			return errors.Wrap(err, "exportPulse failed to getPayload")
		}
		records[string(base58.Encode(id[:]))] = recordData{
			Type:    recordType(rec),
			Data:    rec,
			Payload: pl,
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "exportPulse failed to IterateRecordsOnPulse")
	}

	data := pulseData{
		Records: records,
		Pulse:   *pulse,
		JetID:   jetID,
	}

	return &data, nil
}

func (e *Exporter) getPayload(ctx context.Context, jetID core.RecordID, rec record.Record) (payload, error) {
	switch r := rec.(type) {
	case record.ObjectState:
		if r.GetMemory() == nil {
			break
		}
		blob, err := e.ObjectStorage.GetBlob(ctx, jetID, r.GetMemory())
		if err != nil {
			return nil, errors.Wrapf(err, "getPayload failed to GetBlob (jet: %s)", jetID.DebugString())
		}
		memory := payload{}
		err = codec.NewDecoderBytes(blob, &codec.CborHandle{}).Decode(&memory)
		if err != nil {
			return payload{"MemoryBinary": blob}, nil
		}
		return payload{"Memory": memory}, nil
	case record.Request:
		if r.GetPayload() == nil {
			break
		}
		parcel, err := message.DeserializeParcel(bytes.NewBuffer(r.GetPayload()))
		if err != nil {
			return payload{"PayloadBinary": r.GetPayload()}, nil
		}

		msg := parcel.Message()
		switch m := parcel.Message().(type) {
		case *message.CallMethod:
			res, err := m.ToMap()
			if err != nil {
				return payload{"Payload": m, "Type": msg.Type().String()}, nil
			}
			return payload{"Payload": res, "Type": msg.Type().String()}, nil
		case *message.CallConstructor:
			res, err := m.ToMap()
			if err != nil {
				return payload{"Payload": m, "Type": msg.Type().String()}, nil
			}
			return payload{"Payload": res, "Type": msg.Type().String()}, nil
		case *message.GenesisRequest:
			return payload{"Payload": m, "Type": msg.Type().String()}, nil
		}

		return payload{"Payload": msg, "Type": msg.Type().String()}, nil
	}

	return nil, nil
}

func recordType(rec record.Record) string {
	switch rec.(type) {
	case *record.GenesisRecord:
		return "TypeGenesis"
	case *record.ChildRecord:
		return "TypeChild"
	case *record.JetRecord:
		return "TypeJet"
	case *record.RequestRecord:
		return "TypeCallRequest"
	case *record.ResultRecord:
		return "TypeResult"
	case *record.TypeRecord:
		return "TypeType"
	case *record.CodeRecord:
		return "TypeCode"
	case *record.ObjectActivateRecord:
		return "TypeActivate"
	case *record.ObjectAmendRecord:
		return "TypeAmend"
	case *record.DeactivationRecord:
		return "TypeDeactivate"
	}

	return record.TypeFromRecord(rec).String()
}
