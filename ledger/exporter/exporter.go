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

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/jbenet/go-base58"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

// Exporter provides methods for fetching data view from storage.
type Exporter struct {
	db *storage.DB
}

// NewExporter creates new StorageExporter instance.
func NewExporter(db *storage.DB) *Exporter {
	return &Exporter{db: db}
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
}

// Export returns data view from storage.
func (e *Exporter) Export(ctx context.Context, fromPulse core.PulseNumber, size int) (*core.StorageExportResult, error) {
	result := core.StorageExportResult{Data: map[string]interface{}{}}

	counter := 0
	currentPN := core.PulseNumber(math.Max(float64(fromPulse), float64(core.GenesisPulse.PulseNumber)))
	current := &currentPN
	for current != nil && counter < size {
		pulse, err := e.db.GetPulse(ctx, *current)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch pulse data 22222")
		}
		data, err := e.exportPulse(ctx, &pulse.Pulse)
		if err != nil {
			return nil, err
		}
		result.Data[strconv.FormatUint(uint64(pulse.Pulse.PulseNumber), 10)] = *data

		current = pulse.Next
		counter++
	}

	result.Size = counter
	result.NextFrom = current

	return &result, nil
}

func (e *Exporter) exportPulse(ctx context.Context, pulse *core.Pulse) (*pulseData, error) {
	records := recordsData{}
	err := e.db.IterateRecords(ctx, pulse.PulseNumber, func(id core.RecordID, rec record.Record) error {
		pl, err := e.getPayload(ctx, rec)
		if err != nil {
			return err
		}
		records[string(base58.Encode(id[:]))] = recordData{
			Type:    strings.Title(rec.Type().String()),
			Data:    rec,
			Payload: pl,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	data := pulseData{
		Records: records,
		Pulse:   *pulse,
	}

	return &data, nil
}

func (e *Exporter) getPayload(ctx context.Context, rec record.Record) (payload, error) {
	switch r := rec.(type) {
	case record.ObjectState:
		if r.GetMemory() == nil {
			break
		}
		blob, err := e.db.GetBlob(ctx, r.GetMemory())
		if err != nil {
			return nil, err
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
		return payload{"Payload": parcel}, nil
	}

	return nil, nil
}
