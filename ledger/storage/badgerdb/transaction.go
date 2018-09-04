package badgerdb

import (
	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/ledger/hash"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
)

type TransactionManager struct {
	store *Store
	txn   *badger.Txn
}

func (m *TransactionManager) Commit() error {
	// TODO: check transaction conflict
	return m.txn.Commit(nil)
}

func (m *TransactionManager) Discard() {
	m.txn.Discard()
}

// GetRecord returns record from BadgerDB by *record.Reference.
//
// It returns storage.ErrNotFound if the DB does not contain the key.
func (m *TransactionManager) GetRecord(ref *record.Reference) (record.Record, error) {
	k := prefixkey(scopeIDRecord, ref.CoreRef()[:])
	item, err := m.txn.Get(k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}
	buf, err := item.Value()
	if err != nil {
		return nil, err
	}
	raw, err := record.DecodeToRaw(buf)

	return raw.ToRecord(), nil
}

// SetRecord stores record in BadgerDB and returns *record.Reference of new record.
func (m *TransactionManager) SetRecord(rec record.Record) (*record.Reference, error) {
	raw, err := record.EncodeToRaw(rec)
	if err != nil {
		return nil, err
	}
	ref := &record.Reference{
		Domain: rec.Domain().Record,
		Record: record.ID{
			Pulse: m.store.GetCurrentPulse(),
			Hash:  raw.Hash(),
		},
	}
	k := prefixkey(scopeIDRecord, ref.CoreRef()[:])
	val := record.MustEncodeRaw(raw)
	err = m.txn.Set(k, val)
	if err != nil {
		return nil, err
	}

	return ref, nil
}

// GetClassIndex fetches class lifeline's index.
func (m *TransactionManager) GetClassIndex(ref *record.Reference) (*index.ClassLifeline, error) {
	k := prefixkey(scopeIDLifeline, ref.CoreRef()[:])
	item, err := m.txn.Get(k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}
	buf, err := item.Value()
	if err != nil {
		return nil, err
	}
	return index.DecodeClassLifeline(buf)
}

// SetClassIndex stores class lifeline index.
func (m *TransactionManager) SetClassIndex(ref *record.Reference, idx *index.ClassLifeline) error {
	k := prefixkey(scopeIDLifeline, ref.CoreRef()[:])
	encoded, err := index.EncodeClassLifeline(idx)
	if err != nil {
		return err
	}
	return m.txn.Set(k, encoded)
}

// GetObjectIndex fetches object lifeline index.
func (m *TransactionManager) GetObjectIndex(ref *record.Reference) (*index.ObjectLifeline, error) {
	k := prefixkey(scopeIDLifeline, ref.CoreRef()[:])
	item, err := m.txn.Get(k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}
	buf, err := item.Value()
	if err != nil {
		return nil, err
	}
	return index.DecodeObjectLifeline(buf)
}

// SetObjectIndex stores object lifeline index.
func (m *TransactionManager) SetObjectIndex(ref *record.Reference, idx *index.ObjectLifeline) error {
	k := prefixkey(scopeIDLifeline, ref.CoreRef()[:])
	encoded, err := index.EncodeObjectLifeline(idx)
	if err != nil {
		return err
	}
	return m.txn.Set(k, encoded)
}

// GetDrop returns jet drop for a given pulse number.
func (m *TransactionManager) GetDrop(pulse record.PulseNum) (*jetdrop.JetDrop, error) {
	k := prefixkey(scopeIDJetDrop, record.EncodePulseNum(pulse))
	item, err := m.txn.Get(k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}
	buf, err := item.Value()
	if err != nil {
		return nil, err
	}
	drop, err := jetdrop.Decode(buf)
	if err != nil {
		return nil, err
	}
	return drop, nil
}

// SetDrop stores jet drop for given pulse number.
// Previous JetDrop should be provided.
// On success returns saved drop hash.
func (m *TransactionManager) SetDrop(pulse record.PulseNum, prevdrop *jetdrop.JetDrop) (*jetdrop.JetDrop, error) {
	k := prefixkey(scopeIDJetDrop, record.EncodePulseNum(pulse))

	hw := hash.NewSHA3()
	err := m.store.ProcessSlotHashes(pulse, func(it HashIterator) error {
		for i := 1; it.Next(); i++ {
			b := it.ShallowHash()
			_, err := hw.Write(b)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	drophash := hw.Sum(nil)

	drop := &jetdrop.JetDrop{
		Pulse:    pulse,
		PrevHash: prevdrop.Hash,
		Hash:     drophash,
	}
	encoded, err := jetdrop.Encode(drop)
	if err != nil {
		return nil, err
	}

	err = m.txn.Set(k, encoded)
	if err != nil {
		drop = nil
	}
	return drop, err
}

// SetEntropy stores given entropy for given pulse in storage.
//
// Entropy is used for calculating node roles.
func (m *TransactionManager) SetEntropy(pulse record.PulseNum, entropy []byte) error {
	k := prefixkey(scopeIDEntropy, record.EncodePulseNum(pulse))
	return m.txn.Set(k, entropy)
}

// GetEntropy returns entropy from storage for given pulse.
//
// Entropy is used for calculating node roles.
func (m *TransactionManager) GetEntropy(pulse record.PulseNum) ([]byte, error) {
	k := prefixkey(scopeIDEntropy, record.EncodePulseNum(pulse))
	item, err := m.txn.Get(k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}
	buf, err := item.Value()
	if err != nil {
		return nil, err
	}
	return buf, nil
}
