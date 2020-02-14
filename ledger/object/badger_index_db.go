// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package object

import (
	"bytes"
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/pulse"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

// BadgerIndexDB is a db-based storage, that stores a collection of IndexBuckets
type BadgerIndexDB struct {
	lock sync.RWMutex
	db   store.DB

	recordStore RecordStorage
}

type indexKey struct {
	pn    insolar.PulseNumber
	objID insolar.ID
}

func newIndexKey(raw []byte) indexKey {
	ik := indexKey{}
	ik.pn = insolar.NewPulseNumber(raw)
	ik.objID = *insolar.NewIDFromBytes(raw[ik.pn.Size():])

	return ik
}

func (k indexKey) Scope() store.Scope {
	return store.ScopeIndex
}

func (k indexKey) ID() []byte {
	return append(k.pn.Bytes(), k.objID.Bytes()...)
}

type lastKnownIndexPNKey struct {
	objID insolar.ID
}

func (k lastKnownIndexPNKey) Scope() store.Scope {
	return store.ScopeLastKnownIndexPN
}

func (k lastKnownIndexPNKey) ID() []byte {
	id := k.objID
	return bytes.Join([][]byte{id.Pulse().Bytes(), id.Hash()}, nil)
}

// NewBadgerIndexDB creates a new instance of BadgerIndexDB
func NewBadgerIndexDB(db store.DB, recordStore RecordStorage) *BadgerIndexDB {
	return &BadgerIndexDB{db: db, recordStore: recordStore}
}

// SetIndex adds a bucket with provided pulseNumber and ID
func (i *BadgerIndexDB) SetIndex(ctx context.Context, pn insolar.PulseNumber, bucket record.Index) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	err := i.setBucket(pn, bucket.ObjID, &bucket)
	if err != nil {
		return err
	}

	stats.Record(ctx, statIndexesAddedCount.M(1))

	inslogger.FromContext(ctx).Debugf("[SetIndex] bucket for obj - %v was set successfully. Pulse: %d", bucket.ObjID.DebugString(), pn)

	return nil
}

// UpdateLastKnownPulse must be called after updating TopSyncPulse
func (i *BadgerIndexDB) UpdateLastKnownPulse(ctx context.Context, topSyncPulse insolar.PulseNumber) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	inslogger.FromContext(ctx).Info("UpdateLastKnownPulse starts. topSyncPulse: ", topSyncPulse)

	indexes, err := i.ForPulse(ctx, topSyncPulse)
	if err != nil && err != ErrIndexNotFound {
		return errors.Wrapf(err, "failed to get indexes for pulse: %d", topSyncPulse)
	}

	for idx := range indexes {
		inslogger.FromContext(ctx).Debugf("UpdateLastKnownPulse. pulse: %d, object: %s", topSyncPulse, indexes[idx].ObjID.DebugString())
		if err := i.setLastKnownPN(topSyncPulse, indexes[idx].ObjID); err != nil {
			return errors.Wrapf(err, "can't setLastKnownPN. objId: %s. pulse: %d", indexes[idx].ObjID.DebugString(), topSyncPulse)
		}
	}

	return nil
}

// TruncateHead remove all records starting with 'from'
func (i *BadgerIndexDB) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	it := i.db.NewIterator(&indexKey{objID: *insolar.NewID(pulse.MinTimePulse, nil), pn: from}, false)
	defer it.Close()

	var hasKeys bool
	for it.Next() {
		hasKeys = true
		key := newIndexKey(it.Key())
		err := i.db.Delete(&key)
		if err != nil {
			return errors.Wrapf(err, "can't delete key: %+v", key)
		}

		inslogger.FromContext(ctx).Debugf("Erased key. Pulse number: %s. ObjectID: %s", key.pn.String(), key.objID.String())
	}

	if !hasKeys {
		inslogger.FromContext(ctx).Infof("No records. Nothing done. Pulse number: %s", from.String())
	}

	return nil
}

// ForID returns a lifeline from a bucket with provided PN and ObjID
func (i *BadgerIndexDB) ForID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (record.Index, error) {
	var buck *record.Index
	buck, err := i.getBucket(pn, objID)
	if err == ErrIndexNotFound {
		lastPN, err := i.getLastKnownPN(objID)
		if err != nil {
			return record.Index{}, ErrIndexNotFound
		}

		buck, err = i.getBucket(lastPN, objID)
		if err != nil {
			return record.Index{}, err
		}
	} else if err != nil {
		return record.Index{}, err
	}

	return *buck, nil
}

func (i *BadgerIndexDB) ForPulse(ctx context.Context, pn insolar.PulseNumber) ([]record.Index, error) {
	indexes := make([]record.Index, 0)

	key := &indexKey{objID: insolar.ID{}, pn: pn}
	it := i.db.NewIterator(key, false)
	defer it.Close()

	for it.Next() {
		rawKey := it.Key()
		currentKey := newIndexKey(rawKey)
		if currentKey.pn != pn {
			break
		}

		index := record.Index{}
		rawIndex, err := it.Value()
		err = index.Unmarshal(rawIndex)
		if err != nil {
			return nil, errors.Wrap(err, "Can't unmarshal index")
		}
		indexes = append(indexes, index)
	}

	if len(indexes) == 0 {
		return nil, ErrIndexNotFound
	}

	return indexes, nil
}

func (i *BadgerIndexDB) LastKnownForID(ctx context.Context, objID insolar.ID) (record.Index, error) {
	lastPN, err := i.getLastKnownPN(objID)
	if err != nil {
		return record.Index{}, ErrIndexNotFound
	}

	idx, err := i.getBucket(lastPN, objID)
	if err != nil {
		return record.Index{}, err
	}

	return *idx, nil
}

func (i *BadgerIndexDB) setBucket(pn insolar.PulseNumber, objID insolar.ID, bucket *record.Index) error {
	key := indexKey{pn: pn, objID: objID}

	buff, err := bucket.Marshal()
	if err != nil {
		return err
	}

	return i.db.Set(key, buff)
}

func (i *BadgerIndexDB) getBucket(pn insolar.PulseNumber, objID insolar.ID) (*record.Index, error) {
	buff, err := i.db.Get(indexKey{pn: pn, objID: objID})
	if err == store.ErrNotFound {
		return nil, ErrIndexNotFound
	}
	if err != nil {
		return nil, err
	}
	bucket := record.Index{}
	err = bucket.Unmarshal(buff)
	return &bucket, err
}

func (i *BadgerIndexDB) setLastKnownPN(pn insolar.PulseNumber, objID insolar.ID) error {
	key := lastKnownIndexPNKey{objID: objID}
	return i.db.Set(key, pn.Bytes())
}

func (i *BadgerIndexDB) getLastKnownPN(objID insolar.ID) (insolar.PulseNumber, error) {
	buff, err := i.db.Get(lastKnownIndexPNKey{objID: objID})
	if err != nil {
		return pulse.MinTimePulse, err
	}
	return insolar.NewPulseNumber(buff), err
}

func (i *BadgerIndexDB) filament(ctx context.Context, b *record.Index) ([]record.CompositeFilamentRecord, error) {
	tempRes := make([]record.CompositeFilamentRecord, len(b.PendingRecords))
	for idx, metaID := range b.PendingRecords {
		metaRec, err := i.recordStore.ForID(ctx, metaID)
		if err != nil {
			return nil, err
		}
		pend := record.Unwrap(&metaRec.Virtual).(*record.PendingFilament)
		rec, err := i.recordStore.ForID(ctx, pend.RecordID)
		if err != nil {
			return nil, err
		}

		tempRes[idx] = record.CompositeFilamentRecord{
			Meta:     metaRec,
			MetaID:   metaID,
			Record:   rec,
			RecordID: pend.RecordID,
		}
	}

	return tempRes, nil
}

func (i *BadgerIndexDB) nextFilament(ctx context.Context, b *record.Index) (canContinue bool, nextPN insolar.PulseNumber, err error) {
	firstRecord := b.PendingRecords[0]
	metaRec, err := i.recordStore.ForID(ctx, firstRecord)
	if err != nil {
		return false, insolar.PulseNumber(0), err
	}
	pf := record.Unwrap(&metaRec.Virtual).(*record.PendingFilament)
	if pf.PreviousRecord != nil {
		return true, pf.PreviousRecord.Pulse(), nil
	}

	return false, insolar.PulseNumber(0), nil
}

func (i *BadgerIndexDB) Records(ctx context.Context, readFrom insolar.PulseNumber, readUntil insolar.PulseNumber, objID insolar.ID) ([]record.CompositeFilamentRecord, error) {
	currentPN := readFrom
	var res []record.CompositeFilamentRecord

	if readUntil > readFrom {
		return nil, errors.New("readUntil can't be more then readFrom")
	}

	hasFilamentBehind := true
	for hasFilamentBehind && currentPN >= readUntil {
		b, err := i.getBucket(currentPN, objID)
		if err != nil {
			return nil, err
		}
		if len(b.PendingRecords) == 0 {
			return nil, errors.New("can't fetch pendings from index")
		}

		tempRes, err := i.filament(ctx, b)
		if err != nil {
			return nil, err
		}
		if len(tempRes) == 0 {
			return nil, errors.New("can't fetch pendings from index")
		}
		res = append(tempRes, res...)

		hasFilamentBehind, currentPN, err = i.nextFilament(ctx, b)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
