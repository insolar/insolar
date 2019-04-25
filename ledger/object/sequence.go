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

package object

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"sync"
	"unsafe"

	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/internal/ledger/store"
)

var (
	ErrNoSequenceSyncPulses = errors.New("No synchronized pulses")
)

type SequenceRecordPosition struct {
	Index       uint32
	PulseNumber insolar.PulseNumber
}

func (pos SequenceRecordPosition) String() string {
	return fmt.Sprintf("(Index: %d, PulseNumber: %d)", pos.Index, pos.PulseNumber)
}

func (pos SequenceRecordPosition) GreaterThan(other SequenceRecordPosition) bool {
	return pos.PulseNumber > other.PulseNumber ||
		(pos.PulseNumber == other.PulseNumber && pos.Index > other.Index)
}

type SequenceRecord struct {
	Position SequenceRecordPosition
	ID       insolar.ID
	Record   record.MaterialRecord
}

type SequenceRecordAccessor interface {
	// ForPosition returns sequence record for provided record position.
	ForPosition(ctx context.Context, pos SequenceRecordPosition) (SequenceRecord, error)
	// FromPosition returns sequence record stream from provided record position with corresponding limit.
	FromPosition(ctx context.Context, pos SequenceRecordPosition, limit uint32) chan interface{}
}

type SequenceRecordModifier interface {
	// Push saves new record-value in storage and increment topIndex position.
	Push(ctx context.Context, pn insolar.PulseNumber, id insolar.ID) error
}

type SequenceRecordCursor interface {
	// Top returns the greatest position for sequence records.
	Top(ctx context.Context) (SequenceRecordPosition, error)
	// TopSync returns position for sequence records at last synced pulse.
	TopSync(ctx context.Context) (SequenceRecordPosition, error)
}

type SequenceRecordStorage interface {
	SequenceRecordAccessor
	SequenceRecordModifier
	SequenceRecordCursor

	// UpdatePulse makes "upsert" for top pulse number and sets at old pulse descriptor link to next pulse
	UpdatePulse(ctx context.Context, pn insolar.PulseNumber) error
}

type SequenceRecordMemory struct{}

type SequenceRecordDB struct {
	lock    sync.RWMutex
	db      store.DB
	top     SequenceRecordPosition
	records RecordAccessor
}

// Push saves new record-value in storage and increment topIndex position.
func (s *SequenceRecordDB) Push(ctx context.Context, pn insolar.PulseNumber, id insolar.ID) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	pd, err := s.exactPulse(pn)
	if err == nil {
		pd.topIndex++
	} else {
		pd = pulseLabel{curr: pn, next: 0, topIndex: 0}
	}
	err = s.updatePulseDescriptor(pn, pd)
	if err != nil {
		return errors.Wrap(err, "[SequenceRecordIndex] failed to update pulse descriptor")
	}
	pos := SequenceRecordPosition{Index: pd.topIndex, PulseNumber: pn}
	err = s.set(pos, id)
	if err != nil {
		return errors.Wrapf(err, "[SequenceRecordIndex] failed to save record at position: %s", pos)
	}
	if pos.GreaterThan(s.top) {
		s.top = pos
	}
	return nil
}

// ForPosition returns sequence record for provided record position.
func (s *SequenceRecordDB) ForPosition(ctx context.Context, pos SequenceRecordPosition) (SequenceRecord, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.get(ctx, pos)
}

// FromPosition returns sequence record stream from provided record position with corresponding limit.
func (s *SequenceRecordDB) FromPosition(ctx context.Context, pos SequenceRecordPosition, limit uint32) chan interface{} {
	out := make(chan interface{}, limit)
	go func() {
		defer close(out)

		s.lock.RLock()
		defer s.lock.RUnlock()

		capacity := uint32(0)
		for {
			curr, err := s.exactPulse(pos.PulseNumber)
			if err != nil {
				out <- err
				return
			}
			for ; pos.Index <= curr.topIndex && capacity < limit; pos.Index++ {
				rec, err := s.get(ctx, pos)
				if err != nil {
					out <- err
					return
				}
				out <- rec
				capacity++
			}
			if 0 == curr.next || capacity == limit {
				break
			}
			pos = SequenceRecordPosition{Index: 0, PulseNumber: curr.next}
		}
	}()
	return out
}

// UpdatePulse makes "upsert" for top pulse number and sets at old pulse descriptor link to next pulse
func (s *SequenceRecordDB) UpdatePulse(ctx context.Context, pn insolar.PulseNumber) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	prev, err := s.topPulse()
	if err == ErrNotFound {
		err = s.updateFirst(pn)
		if err != nil {
			return errors.Wrap(err, "[SequenceRecordIndex] failed to update first pulse number")
		}
		return s.updateTop(pn)
	} else if err != nil {
		return errors.Wrap(err, "[SequenceRecordIndex] failed to get top pulse number")
	}
	pd, err := s.exactPulse(prev)
	if err != nil {
		return errors.Wrap(err, "[SequenceRecordIndex] failed to get previous pulse descriptor")
	}
	pd.next = pn
	return s.updateTop(pn)
}

func (s *SequenceRecordDB) FirstPulse() (SequenceRecordPosition, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	desc, err := s.firstSyncPulse()
	if err != nil {
		return SequenceRecordPosition{}, err
	}
	return SequenceRecordPosition{Index: desc.topIndex, PulseNumber: desc.curr}, nil
}

// Top returns the greatest position for sequence records.
func (s *SequenceRecordDB) Top(ctx context.Context) (SequenceRecordPosition, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.top, nil
}

// TopSync returns position for sequence records at last synced pulse.
func (s *SequenceRecordDB) TopSync(ctx context.Context) (SequenceRecordPosition, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	desc, err := s.lastSyncPulse()
	if err != nil {
		return SequenceRecordPosition{}, err
	}
	return SequenceRecordPosition{Index: desc.topIndex, PulseNumber: desc.curr}, nil
}

type topKey uint8

func (k topKey) Scope() store.Scope {
	return store.ScopeReplicaPosition
}

func (k topKey) ID() []byte {
	return []byte{byte(k)}
}

type pulseLabel struct {
	curr     insolar.PulseNumber
	next     insolar.PulseNumber
	topIndex uint32
}

type pulseKey insolar.PulseNumber

func (k pulseKey) Scope() store.Scope {
	return store.ScopePulseLabel
}

func (k pulseKey) ID() []byte {
	return utils.UInt32ToBytes(uint32(k))
}

type positionKey SequenceRecordPosition

func (k positionKey) Scope() store.Scope {
	return store.ScopeSequenceRecord
}

func (k positionKey) ID() []byte {
	return append(SequenceRecordPosition(k).PulseNumber.Bytes(),
		utils.UInt32ToBytes(SequenceRecordPosition(k).Index)...)
}

func (s *SequenceRecordDB) topPulse() (insolar.PulseNumber, error) {
	key := topKey(42)
	buff, err := s.db.Get(key)
	if err == ErrNotFound {
		return 0, ErrNoSequenceSyncPulses
	} else if err != nil {
		return 0, err
	}
	return insolar.PulseNumber(binary.BigEndian.Uint32(buff)), nil
}

func (s *SequenceRecordDB) firstPulse() (insolar.PulseNumber, error) {
	key := topKey(0)
	buff, err := s.db.Get(key)
	if err == ErrNotFound {
		return 0, ErrNoSequenceSyncPulses
	} else if err != nil {
		return 0, err
	}
	return insolar.PulseNumber(binary.BigEndian.Uint32(buff)), nil
}

func (s *SequenceRecordDB) updateTop(pn insolar.PulseNumber) error {
	key := topKey(42)
	return s.db.Set(key, utils.UInt32ToBytes(uint32(pn)))
}

func (s *SequenceRecordDB) updateFirst(pn insolar.PulseNumber) error {
	key := topKey(0)
	return s.db.Set(key, utils.UInt32ToBytes(uint32(pn)))
}

func encodePulseDescriptor(pd pulseLabel) []byte {
	buff := bytes.NewBuffer(make([]byte, unsafe.Sizeof(pd)))
	enc := codec.NewEncoder(buff, &codec.CborHandle{})
	enc.MustEncode(pd)
	return buff.Bytes()
}

func decodePulseDescriptor(buff []byte) (pulseLabel, error) {
	dec := codec.NewDecoderBytes(buff, &codec.CborHandle{})
	pd := pulseLabel{}
	err := dec.Decode(&pd)
	if err != nil {
		return pulseLabel{}, err
	}
	return pd, nil
}

func (s *SequenceRecordDB) exactPulse(pn insolar.PulseNumber) (pulseLabel, error) {
	key := pulseKey(uint32(pn))
	buff, err := s.db.Get(key)
	if err != nil {
		return pulseLabel{}, err
	}
	return decodePulseDescriptor(buff)
}

func (s *SequenceRecordDB) lastSyncPulse() (pulseLabel, error) {
	top, err := s.topPulse()
	if err != nil {
		return pulseLabel{}, err
	}
	key := pulseKey(uint32(top))
	buff, err := s.db.Get(key)
	if err != nil {
		return pulseLabel{}, err
	}
	return decodePulseDescriptor(buff)
}

func (s *SequenceRecordDB) firstSyncPulse() (pulseLabel, error) {
	firstPulse, err := s.firstPulse()
	if err != nil {
		return pulseLabel{}, err
	}
	key := pulseKey(uint32(firstPulse))
	buff, err := s.db.Get(key)
	if err != nil {
		return pulseLabel{}, err
	}
	return decodePulseDescriptor(buff)
}

func (s *SequenceRecordDB) updatePulseDescriptor(pn insolar.PulseNumber, pd pulseLabel) error {
	key := pulseKey(pn)
	return s.db.Set(key, encodePulseDescriptor(pd))
}

func (s *SequenceRecordDB) set(pos SequenceRecordPosition, ref insolar.ID) error {
	key := positionKey(pos)

	_, err := s.db.Get(key)
	if err == nil {
		return ErrOverride
	}

	return s.db.Set(key, encodeID(ref))
}

func (s *SequenceRecordDB) get(ctx context.Context, pos SequenceRecordPosition) (SequenceRecord, error) {
	buff, err := s.db.Get(positionKey(pos))
	if err == store.ErrNotFound {
		return SequenceRecord{}, ErrNotFound
	}
	if err != nil {
		return SequenceRecord{}, err
	}
	id, err := decodeID(buff)
	if err != nil {
		return SequenceRecord{}, errors.Wrap(err, "[SequenceRecordIndex] failed to decode insolar.ID")
	}

	pure, err := s.records.ForID(ctx, id)
	if err != nil {
		return SequenceRecord{}, errors.Wrap(err, "[SequenceRecordIndex] failed to get record by id")
	}
	return SequenceRecord{Position: pos, ID: id, Record: pure}, nil
}

func encodeID(ref insolar.ID) []byte {
	buff := bytes.NewBuffer(make([]byte, unsafe.Sizeof(ref)))
	enc := codec.NewEncoder(buff, &codec.CborHandle{})
	enc.MustEncode(ref)
	return buff.Bytes()
}

func decodeID(buff []byte) (insolar.ID, error) {
	dec := codec.NewDecoderBytes(buff, &codec.CborHandle{})
	ref := insolar.ID{}
	err := dec.Decode(&ref)
	if err != nil {
		return insolar.ID{}, err
	}
	return ref, nil
}
