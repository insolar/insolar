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

type Exporter struct {
	db *storage.DB
}

func NewExporter(db *storage.DB) *Exporter {
	return &Exporter{db: db}
}

type payload = map[string]interface{}

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

func (e *Exporter) Export(ctx context.Context, fromPulse core.PulseNumber, size int) (*core.ExportResult, error) {
	result := core.ExportResult{Data: map[string]interface{}{}}

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
