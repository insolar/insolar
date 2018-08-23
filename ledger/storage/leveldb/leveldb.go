/*
 *    Copyright 2018 INS Ecosystem
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

package leveldb

import (
	"os"
	"path/filepath"
	"time"

	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
)

const (
	dbDirPath        = "_db"
	zeroRecordBinary = "" // TODO: Empty ClassActivateRecord serialized
	zeroRecordHash   = "" // TODO: Hash from zeroRecordBinary
)

// The PulseFn type is an adapter to allow set function produces
// pulses for storage.
type PulseFn func() record.PulseNum

// LevelLedger represents ledger's LevelDB storage.
type LevelLedger struct {
	ldb     *leveldb.DB
	pulseFn PulseFn
	zeroRef record.Reference
}

const (
	scopeIDLifeline byte = 1
	scopeIDRecord   byte = 2
	scopeIDJetDrop  byte = 3
)

// InitDB returns LevelDB ledger implementation.
func InitDB(dir string, opts *opt.Options) (*LevelLedger, error) {
	if dir == "" {
		dir = dbDirPath
	}
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	db, err := leveldb.OpenFile(absPath, setOptions(opts))
	if err != nil {
		return nil, err
	}

	var zeroID record.ID
	ledger := LevelLedger{
		ldb: db,
		// FIXME: temporary pulse implementation
		pulseFn: func() record.PulseNum {
			return record.PulseNum(time.Now().Unix() / 10)
		},
		zeroRef: record.Reference{
			Domain: record.ID{}, // TODO: fill domain
			Record: zeroID,
		},
	}
	_, err = db.Get([]byte(zeroRecordHash), nil)
	if err == leveldb.ErrNotFound {
		err = db.Put([]byte(zeroRecordHash), []byte(zeroRecordBinary), nil)
		if err != nil {
			return nil, err
		}
		return &ledger, nil
	}
	if err != nil {
		return nil, err
	}
	return &ledger, nil
}

func prefixkey(prefix byte, key []byte) []byte {
	k := make([]byte, record.RefIDSize+1)
	k[0] = prefix
	_ = copy(k[1:], key)
	return k
}

// GetCurrentPulse return currently stored pulse.
func (ll *LevelLedger) GetCurrentPulse() record.PulseNum {
	return ll.pulseFn()
}

// SetPulseFn allows redefine pulse function.
func (ll *LevelLedger) SetPulseFn(fn PulseFn) {
	ll.pulseFn = fn
}

// GetRecord returns record from leveldb by *record.Reference.
//
// It returns ErrNotFound if the DB does not contains the key.
func (ll *LevelLedger) GetRecord(ref *record.Reference) (record.Record, error) {
	k := prefixkey(scopeIDRecord, ref.Key())
	buf, err := ll.ldb.Get(k, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}
	raw, err := record.DecodeToRaw(buf)
	if err != nil {
		return nil, err
	}
	return raw.ToRecord(), nil
}

// SetRecord stores record in leveldb
func (ll *LevelLedger) SetRecord(rec record.Record) (*record.Reference, error) {
	raw, err := record.EncodeToRaw(rec)
	if err != nil {
		return nil, err
	}
	ref := &record.Reference{
		Domain: rec.Domain().Record,
		Record: record.ID{Pulse: ll.pulseFn(), Hash: raw.Hash()},
	}
	k := prefixkey(scopeIDRecord, ref.Key())
	err = ll.ldb.Put(k, record.MustEncodeRaw(raw), nil)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

// GetClassIndex fetches lifeline index from leveldb
func (ll *LevelLedger) GetClassIndex(ref *record.Reference) (*index.ClassLifeline, error) {
	k := prefixkey(scopeIDLifeline, ref.Key())
	buf, err := ll.ldb.Get(k, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}
	idx, err := index.DecodeClassLifeline(buf)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

// SetClassIndex stores lifeline index into leveldb
func (ll *LevelLedger) SetClassIndex(ref *record.Reference, idx *index.ClassLifeline) error {
	k := prefixkey(scopeIDLifeline, ref.Key())
	encoded, err := index.EncodeClassLifeline(idx)
	if err != nil {
		return err
	}
	return ll.ldb.Put(k, encoded, nil)
}

// GetObjectIndex fetches lifeline index from leveldb
func (ll *LevelLedger) GetObjectIndex(ref *record.Reference) (*index.ObjectLifeline, error) {
	k := prefixkey(scopeIDLifeline, ref.Key())
	buf, err := ll.ldb.Get(k, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}
	idx, err := index.DecodeObjectLifeline(buf)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

// SetObjectIndex stores lifeline index into leveldb
func (ll *LevelLedger) SetObjectIndex(ref *record.Reference, idx *index.ObjectLifeline) error {
	k := prefixkey(scopeIDLifeline, ref.Key())
	encoded, err := index.EncodeObjectLifeline(idx)
	if err != nil {
		return err
	}
	return ll.ldb.Put(k, encoded, nil)
}

// GetPulseKeys returns all record keys from slot after given pulse.
func (ll *LevelLedger) GetPulseKeys(pulse record.PulseNum) ([][]byte, error) {
	// TODO: implement me
	return [][]byte{}, nil
}

// GetDrop returns jet drop for a given pulse number.
func (ll *LevelLedger) GetDrop(pulse record.PulseNum) (*jetdrop.JetDrop, error) {
	k := prefixkey(scopeIDJetDrop, record.EncodePulseNum(pulse))
	buf, err := ll.ldb.Get(k, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}
	drop, err := jetdrop.DecodeJetDrop(buf)
	if err != nil {
		return nil, err
	}
	return drop, nil
}

// SetDrop stores given jet drop for given pulse number.
func (ll *LevelLedger) SetDrop(pulse record.PulseNum, drop *jetdrop.JetDrop) error {
	k := prefixkey(scopeIDJetDrop, record.EncodePulseNum(pulse))
	encoded, err := jetdrop.EncodeJetDrop(drop)
	if err != nil {
		return err
	}
	return ll.ldb.Put(k, encoded, nil)
}

// Close terminates db connection
func (ll *LevelLedger) Close() error {
	return ll.ldb.Close()
}

// DropDB erases all data from storage.
func DropDB() error {
	absPath, err := filepath.Abs(dbDirPath)
	if err != nil {
		return err
	}

	if err = os.RemoveAll(absPath); err != nil {
		return err
	}

	return nil
}
