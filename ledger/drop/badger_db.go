// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package drop

import (
	"bytes"
	"context"
	"encoding/base64"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

type BadgerDB struct {
	db store.DB
}

// NewBadgerDB creates a new storage, that holds data in a db.
func NewBadgerDB(db store.DB) *BadgerDB {
	return &BadgerDB{db: db}
}

type dropDbKey struct {
	jetPrefix []byte
	pn        insolar.PulseNumber
}

func (dk *dropDbKey) Scope() store.Scope {
	return store.ScopeJetDrop
}

func (dk *dropDbKey) ID() []byte {
	// order ( pn + jetPrefix ) is important: we use this logic for removing not finalized drops
	return bytes.Join([][]byte{dk.pn.Bytes(), dk.jetPrefix}, nil)
}

func newDropDbKey(raw []byte) dropDbKey {
	dk := dropDbKey{}
	dk.pn = insolar.NewPulseNumber(raw)
	dk.jetPrefix = raw[dk.pn.Size():]

	return dk
}

// ForPulse returns a Drop for a provided pulse, that is stored in a db.
func (ds *BadgerDB) ForPulse(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) (Drop, error) {
	k := dropDbKey{jetID.Prefix(), pulse}

	buf, err := ds.db.Get(&k)
	if err != nil {
		return Drop{}, err
	}

	drop := Drop{}
	err = drop.Unmarshal(buf)
	if err != nil {
		return Drop{}, err
	}
	return drop, nil
}

// Set saves a provided Drop to a db.
func (ds *BadgerDB) Set(ctx context.Context, drop Drop) error {
	k := dropDbKey{drop.JetID.Prefix(), drop.Pulse}

	_, err := ds.db.Get(&k)
	if err == nil {
		return ErrOverride
	}

	encoded, err := drop.Marshal()
	if err != nil {
		return err
	}

	return ds.db.Set(&k, encoded)
}

// TruncateHead remove all records starting with 'from'
func (ds *BadgerDB) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	it := ds.db.NewIterator(&dropDbKey{jetPrefix: []byte{}, pn: from}, false)
	defer it.Close()

	var hasKeys bool
	for it.Next() {
		hasKeys = true
		key := newDropDbKey(it.Key())
		err := ds.db.Delete(&key)
		if err != nil {
			return errors.Wrapf(err, "can't delete key: %+v", key)
		}

		inslogger.FromContext(ctx).Debugf("Erased key. Pulse number: %s. Jet prefix: %s", key.pn.String(), base64.RawURLEncoding.EncodeToString(key.jetPrefix))
	}
	if !hasKeys {
		inslogger.FromContext(ctx).Debug("No records. Nothing done. Pulse number: " + from.String())
	}

	return nil
}
