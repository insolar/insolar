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

package storage

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"

	"github.com/insolar/insolar/core"
	"github.com/ugorji/go/codec"
)

type pulseTracker struct {
	DB DBContext `inject:""`
}

// NewPulseTracker returns new instance PulseTracker with DB-storage realization
func NewPulseTracker() PulseTracker {
	return new(pulseTracker)
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
func (pt *pulseTracker) AddPulse(ctx context.Context, pulse core.Pulse) error {
	return pt.DB.Update(ctx, func(tx *TransactionManager) error {
		var (
			previousPulseNumber  core.PulseNumber
			previousSerialNumber int
		)

		_, err := tx.get(ctx, prefixkey(scopeIDPulse, pulse.PulseNumber.Bytes()))
		if err == nil {
			return ErrOverride
		} else if err != core.ErrNotFound {
			return err
		}

		previousPulse, err := tx.GetLatestPulse(ctx)
		if err != nil && err != core.ErrNotFound {
			return err
		}

		// Set next on previousPulseNumber pulse if it exists.
		if err == nil {
			if previousPulse != nil {
				previousPulseNumber = previousPulse.Pulse.PulseNumber
				previousSerialNumber = previousPulse.SerialNumber
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
			Prev:         &previousPulseNumber,
			SerialNumber: previousSerialNumber + 1,
			Pulse:        pulse,
		}
		err = tx.set(ctx, prefixkey(scopeIDPulse, pulse.PulseNumber.Bytes()), p.Bytes())
		if err != nil {
			return err
		}

		return tx.set(ctx, prefixkey(scopeIDSystem, []byte{sysLatestPulse}), p.Bytes())
	})
}

// GetPulse returns pulse for provided pulse number.
func (pt *pulseTracker) GetPulse(ctx context.Context, num core.PulseNumber) (*Pulse, error) {
	var (
		pulse *Pulse
		err   error
	)
	err = pt.DB.View(ctx, func(tx *TransactionManager) error {
		pulse, err = tx.GetPulse(ctx, num)
		return err
	})
	if err != nil {
		return nil, err
	}
	return pulse, nil
}

// GetPreviousPulse returns pulse for provided pulse number.
func (pt *pulseTracker) GetPreviousPulse(ctx context.Context, num core.PulseNumber) (*Pulse, error) {
	var (
		pulse *Pulse
		err   error
	)
	err = pt.DB.View(ctx, func(tx *TransactionManager) error {
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

// GetNthPrevPulse returns Nth previous pulse from some pulse number
func (pt *pulseTracker) GetNthPrevPulse(ctx context.Context, n uint, num core.PulseNumber) (*Pulse, error) {
	pulse, err := pt.GetPulse(ctx, num)
	if err != nil {
		return nil, err
	}

	err = pt.DB.View(ctx, func(tx *TransactionManager) error {
		for n > 0 {
			if pulse.Prev == nil {
				pulse = nil
				return core.ErrNotFound
			}
			pulse, err = tx.GetPulse(ctx, *pulse.Prev)
			if err != nil {
				return err
			}
			n--
		}
		return nil
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

// Deprecated: use core.PulseStorage.Current() instead (or private getLatestPulse if applicable).
func (pt *pulseTracker) GetLatestPulse(ctx context.Context) (*Pulse, error) {
	return pt.getLatestPulse(ctx)
}

// DeletePulse delete pulse data.
func (pt *pulseTracker) DeletePulse(ctx context.Context, num core.PulseNumber) error {
	return errors.New("DB pulse removal is forbidden")
}

func (pt *pulseTracker) getLatestPulse(ctx context.Context) (*Pulse, error) {
	tx, err := pt.DB.BeginTransaction(false)
	if err != nil {
		return nil, err
	}
	defer tx.Discard()

	return tx.GetLatestPulse(ctx)
}

func pulseNumFromKey(from int, key []byte) core.PulseNumber {
	return core.NewPulseNumber(key[from : from+core.PulseNumberSize])
}

// Key type for wrapping storage binary key.
type Key []byte

// PulseNumber returns pulse number for provided storage binary key.
func (b Key) PulseNumber() core.PulseNumber {
	// by default expect jetID after:
	// offset in this case: is 1 + RecordHashSize (jet length) - 1 minus jet prefix
	from := core.RecordHashSize
	switch b[0] {
	case scopeIDPulse:
		from = 1
	case scopeIDSystem:
		// for specific system records is different rules
		// pulse number could exist or not
		return 0
	}
	return pulseNumFromKey(from, b)
}

// String string hex representation
func (b Key) String() string {
	return hex.EncodeToString(b)
}
