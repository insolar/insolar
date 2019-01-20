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

package exporter

import (
	"bytes"
	"context"
	"math"
	"strconv"
	"strings"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/record"
	base58 "github.com/jbenet/go-base58"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

// Exporter provides methods for fetching data view from storage.
type Exporter struct {
	db  *storage.DB
	ps  *storage.PulseStorage
	cfg configuration.Exporter
}

// NewExporter creates new StorageExporter instance.
func NewExporter(db *storage.DB, ps *storage.PulseStorage, cfg configuration.Exporter) *Exporter {
	return &Exporter{db: db, ps: ps, cfg: cfg}
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
	inslog := inslogger.FromContext(ctx)
	inslog.Debugf("[ API Export ] start")

	jetIDs, err := e.db.GetJets(ctx)
	if err != nil {
		inslog.Debugf("[ API Export ] error getting jets: %s", err.Error())
		return nil, err
	}

	currentPulse, err := e.ps.Current(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get current pulse data")
	}

	counter := 0
	fromPulsePN := core.PulseNumber(math.Max(float64(fromPulse), float64(core.GenesisPulse.PulseNumber)))
	iterPulse := &fromPulsePN
	for iterPulse != nil && counter < size {
		pulse, err := e.db.GetPulse(ctx, *iterPulse)
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
	err := e.db.IterateRecordsOnPulse(ctx, jetID, pulse.PulseNumber, func(id core.RecordID, rec record.Record) error {
		pl, err := e.getPayload(ctx, jetID, rec)
		if err != nil {
			return errors.Wrap(err, "exportPulse failed to getPayload")
		}
		records[string(base58.Encode(id[:]))] = recordData{
			Type:    strings.Title(rec.Type().String()),
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
		blob, err := e.db.GetBlob(ctx, jetID, r.GetMemory())
		if err != nil {
			return nil, errors.Wrapf(err, "getPayload failed to GetBlob (jet: %s)", jetID.JetIDString())
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
		msg, err := message.Deserialize(bytes.NewBuffer(r.GetPayload()))
		if err != nil {
			return payload{"PayloadBinary": r.GetPayload()}, nil
		}
		switch m := msg.(type) {
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
