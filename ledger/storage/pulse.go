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

func toPulse(raw []byte) (*Pulse, error) {
	dec := codec.NewDecoder(bytes.NewReader(raw), &codec.CborHandle{})
	var rec Pulse
	err := dec.Decode(&rec)
	if err != nil {
		return nil, err
	}
	return &rec, nil
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
		var previousPulseNumber core.PulseNumber
		previousPulse, err := tx.GetLatestPulse(ctx)
		if err != nil && err != ErrNotFound {
			return err
		}

		// Set next on previousPulseNumber pulse if it exists.
		if err == nil {
			if previousPulse != nil {
				previousPulseNumber = previousPulse.Pulse.PulseNumber
			}

			prevPulse, err := tx.GetPulse(ctx, previousPulseNumber)
			if err != nil {
				return err
			}
			prevPulse.Next = &pulse.PulseNumber
			err = tx.set(ctx, prefixkey(scopeIDPulse, previousPulseNumber.Bytes()), prevPulse.Bytes())
			if err != nil {
				return err
			}
		}

		// Save new pulse.
		p := Pulse{
			Prev:  &previousPulseNumber,
			Pulse: pulse,
		}
		err = tx.set(ctx, prefixkey(scopeIDPulse, pulse.PulseNumber.Bytes()), p.Bytes())
		if err != nil {
			return err
		}

		return tx.set(ctx, prefixkey(scopeIDSystem, []byte{sysLatestPulse}), p.Bytes())
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

// GetPreviousPulse returns pulse for provided pulse number.
func (db *DB) GetPreviousPulse(ctx context.Context, num core.PulseNumber) (*Pulse, error) {
	var (
		pulse *Pulse
		err   error
	)
	err = db.View(ctx, func(tx *TransactionManager) error {
		pulse, err = tx.GetPulse(ctx, num)
		if err != nil {
			return err
		}
		if pulse.Prev == nil {
			pulse = nil
			return nil
		}
		pulse, err = tx.GetPulse(ctx, *pulse.Prev)
		return err
	})
	if err != nil {
		return nil, err
	}

	return pulse, nil
}

// GetLatestPulse returns the latest pulse
func (m *TransactionManager) GetLatestPulse(ctx context.Context) (*Pulse, error) {
	buf, err := m.get(ctx, prefixkey(scopeIDSystem, []byte{sysLatestPulse}))
	if err != nil {
		return nil, err
	}
	return toPulse(buf)
}

// GetLatestPulse returns the latest pulse
func (db *DB) GetLatestPulse(ctx context.Context) (*Pulse, error) {
	tx, err := db.BeginTransaction(false)
	if err != nil {
		return nil, err
	}
	defer tx.Discard()

	return tx.GetLatestPulse(ctx)
}

func pulseNumFromKey(from int, key []byte) core.PulseNumber {
	return core.NewPulseNumber(key[from : from+core.PulseNumberSize])
}
