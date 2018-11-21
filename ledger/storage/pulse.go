package storage

import (
	"bytes"
	"context"

	"github.com/insolar/insolar/core"
	"github.com/ugorji/go/codec"
)

// Pulse is a record containing pulse info.
type Pulse struct {
	Prev  *core.PulseNumber
	Next  *core.PulseNumber
	Pulse core.Pulse
}

// Bytes serializes pulse.
func (p *Pulse) Bytes() []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	enc.MustEncode(p)
	return buf.Bytes()
}

// GetLatestPulseNumber returns current pulse number.
func (m *TransactionManager) GetLatestPulseNumber(ctx context.Context) (core.PulseNumber, error) {
	buf, err := m.get(ctx, prefixkey(scopeIDSystem, []byte{sysLatestPulse}))
	if err != nil {
		return 0, err
	}
	return core.NewPulseNumber(buf), nil
}

// GetPulse returns pulse for provided pulse number.
func (m *TransactionManager) GetPulse(ctx context.Context, num core.PulseNumber) (*Pulse, error) {
	buf, err := m.get(ctx, prefixkey(scopeIDPulse, num.Bytes()))
	if err != nil {
		return nil, err
	}

	dec := codec.NewDecoder(bytes.NewReader(buf), &codec.CborHandle{})
	var rec Pulse
	err = dec.Decode(&rec)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

// AddPulse saves new pulse data and updates index.
func (db *DB) AddPulse(ctx context.Context, pulse core.Pulse) error {
	return db.Update(ctx, func(tx *TransactionManager) error {
		var previous core.PulseNumber
		previous, err := tx.GetLatestPulseNumber(ctx)
		if err != nil && err != ErrNotFound {
			return err
		}

		// Set next on previous pulse if it exists.
		if err == nil {
			prevPulse, err := tx.GetPulse(ctx, previous)
			if err != nil {
				return err
			}
			prevPulse.Next = &pulse.PulseNumber
			err = tx.set(ctx, prefixkey(scopeIDPulse, previous.Bytes()), prevPulse.Bytes())
			if err != nil {
				return err
			}
		}

		// Save new pulse.
		p := Pulse{
			Prev:  &previous,
			Pulse: pulse,
		}
		err = tx.set(ctx, prefixkey(scopeIDPulse, pulse.PulseNumber.Bytes()), p.Bytes())
		if err != nil {
			return err
		}

		return tx.set(ctx, prefixkey(scopeIDSystem, []byte{sysLatestPulse}), pulse.PulseNumber.Bytes())
	})
}

// GetPulse returns pulse for provided pulse number.
func (db *DB) GetPulse(ctx context.Context, num core.PulseNumber) (*Pulse, error) {
	var (
		pulse *Pulse
		err   error
	)
	err = db.View(ctx, func(tx *TransactionManager) error {
		pulse, err = tx.GetPulse(ctx, num)
		return err
	})
	if err != nil {
		return nil, err
	}
	return pulse, nil
}

// GetLatestPulseNumber returns current pulse number.
func (db *DB) GetLatestPulseNumber(ctx context.Context) (core.PulseNumber, error) {
	tx := db.BeginTransaction(false)
	defer tx.Discard()

	return tx.GetLatestPulseNumber(ctx)
}
