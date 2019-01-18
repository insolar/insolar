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
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

// GetDrop returns jet drop for a given pulse number and jet id.
func (db *DB) GetDrop(ctx context.Context, jetID core.RecordID, pulse core.PulseNumber) (*jet.JetDrop, error) {
	_, prefix := jet.Jet(jetID)
	k := prefixkey(scopeIDJetDrop, prefix, pulse.Bytes())

	buf, err := db.get(ctx, k)
	if err != nil {
		return nil, err
	}
	drop, err := jet.Decode(buf)
	if err != nil {
		return nil, err
	}
	return drop, nil
}

// CreateDrop creates and stores jet drop for given pulse number.
//
// On success returns saved drop object, slot records, drop size.
func (db *DB) CreateDrop(ctx context.Context, jetID core.RecordID, pulse core.PulseNumber, prevHash []byte) (
	*jet.JetDrop,
	[][]byte,
	uint64,
	error,
) {
	var err error
	db.waitinflight()

	hw := db.PlatformCryptographyScheme.ReferenceHasher()
	_, err = hw.Write(prevHash)
	if err != nil {
		return nil, nil, 0, err
	}

	var messages [][]byte
	_, jetPrefix := jet.Jet(jetID)
	// messagesPrefix := prefixkey(scopeIDMessage, jetPrefix, pulse.Bytes())

	// err = db.db.View(func(txn *badger.Txn) error {
	// 	it := txn.NewIterator(badger.DefaultIteratorOptions)
	// 	defer it.Close()
	//
	// 	for it.Seek(messagesPrefix); it.ValidForPrefix(messagesPrefix); it.Next() {
	// 		val, err := it.Item().ValueCopy(nil)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		messages = append(messages, val)
	// 	}
	// 	return nil
	// })
	// if err != nil {
	// 	return nil, nil, 0, err
	// }

	var dropSize uint64
	recordPrefix := prefixkey(scopeIDRecord, jetPrefix, pulse.Bytes())
	err = db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(recordPrefix); it.ValidForPrefix(recordPrefix); it.Next() {
			val, err := it.Item().ValueCopy(nil)
			if err != nil {
				return err
			}
			_, err = hw.Write(val)
			if err != nil {
				return err
			}
			dropSize += uint64(len(val))
		}
		return nil
	})
	if err != nil {
		return nil, nil, 0, err
	}

	drop := jet.JetDrop{
		Pulse:    pulse,
		PrevHash: prevHash,
		Hash:     hw.Sum(nil),
	}
	return &drop, messages, dropSize, nil
}

// SetDrop saves provided JetDrop in db.
func (db *DB) SetDrop(ctx context.Context, jetID core.RecordID, drop *jet.JetDrop) error {
	fmt.Printf("SetDrop for jet: %v, pulse: %v\n", jetID.JetIDString(), drop.Pulse)

	_, prefix := jet.Jet(jetID)
	k := prefixkey(scopeIDJetDrop, prefix, drop.Pulse.Bytes())
	_, err := db.get(ctx, k)
	if err == nil {
		fmt.Println("override drop for pulse ", drop.Pulse)
		return ErrOverride
	}

	encoded, err := jet.Encode(drop)
	if err != nil {
		return err
	}
	return db.set(ctx, k, encoded)
}

// UpdateJetTree updates jet tree for specified pulse.
func (db *DB) UpdateJetTree(ctx context.Context, pulse core.PulseNumber, setActual bool, ids ...core.RecordID) error {
	db.jetTreeLock.Lock()
	defer db.jetTreeLock.Unlock()

	k := prefixkey(scopeIDSystem, []byte{sysJetTree}, pulse.Bytes())
	tree, err := db.getJetTree(ctx, pulse)
	if err != nil {
		return err
	}
	for _, id := range ids {
		tree.Update(id, setActual)
	}

	return db.set(ctx, k, tree.Bytes())
}

// GetJetTree fetches tree for specified pulse.
func (db *DB) GetJetTree(ctx context.Context, pulse core.PulseNumber) (*jet.Tree, error) {
	db.jetTreeLock.RLock()
	defer db.jetTreeLock.RUnlock()
	return db.getJetTree(ctx, pulse)
}

func (db *DB) getJetTree(ctx context.Context, pulse core.PulseNumber) (*jet.Tree, error) {
	k := prefixkey(scopeIDSystem, []byte{sysJetTree}, pulse.Bytes())
	buff, err := db.get(ctx, k)
	if err == ErrNotFound {
		fmt.Println("NewTree was created with pulse", pulse)
		tree := jet.NewTree(pulse == core.GenesisPulse.PulseNumber)
		err := db.set(ctx, k, tree.Bytes())
		return tree, err
	}
	if err != nil {
		return nil, err
	}

	dec := codec.NewDecoder(bytes.NewReader(buff), &codec.CborHandle{})
	var tree jet.Tree
	err = dec.Decode(&tree)
	if err != nil {
		return nil, err
	}

	return &tree, nil
}

// SplitJetTree performs jet split and returns resulting jet ids.
func (db *DB) SplitJetTree(
	ctx context.Context, pulse core.PulseNumber, jetID core.RecordID,
) (*core.RecordID, *core.RecordID, error) {
	db.jetTreeLock.Lock()
	defer db.jetTreeLock.Unlock()

	k := prefixkey(scopeIDSystem, []byte{sysJetTree}, pulse.Bytes())
	tree, err := db.getJetTree(ctx, pulse)
	if err != nil {
		return nil, nil, err
	}

	left, right, err := tree.Split(jetID)
	if err != nil {
		return nil, nil, err
	}
	err = db.set(ctx, k, tree.Bytes())
	if err != nil {
		return nil, nil, err
	}

	return left, right, nil
}

// CloneJetTree copies tree from one pulse to another. Use it to copy past tree into new pulse.
func (db *DB) CloneJetTree(
	ctx context.Context, from, to core.PulseNumber,
) (*jet.Tree, error) {
	db.jetTreeLock.Lock()
	defer db.jetTreeLock.Unlock()

	k := prefixkey(scopeIDSystem, []byte{sysJetTree}, to.Bytes())
	tree, err := db.getJetTree(ctx, from)
	if err != nil {
		return nil, err
	}

	tree.ResetActual()

	err = db.set(ctx, k, tree.Bytes())
	if err != nil {
		return nil, err
	}

	return tree, nil
}

// AddJets stores a list of jets of the current node.
func (db *DB) AddJets(ctx context.Context, jetIDs ...core.RecordID) error {
	db.addJetLock.Lock()
	defer db.addJetLock.Unlock()

	k := prefixkey(scopeIDSystem, []byte{sysJetList})

	var jets jet.IDSet
	buff, err := db.get(ctx, k)
	if err == nil {
		dec := codec.NewDecoder(bytes.NewReader(buff), &codec.CborHandle{})
		err = dec.Decode(&jets)
		if err != nil {
			return err
		}
	} else if err == ErrNotFound {
		jets = jet.IDSet{}
	} else {
		return err
	}

	for _, id := range jetIDs {
		jets[id] = struct{}{}
	}
	return db.set(ctx, k, jets.Bytes())
}

// GetJets returns jets of the current node
func (db *DB) GetJets(ctx context.Context) (jet.IDSet, error) {
	db.addJetLock.RLock()
	defer db.addJetLock.RUnlock()

	k := prefixkey(scopeIDSystem, []byte{sysJetList})
	buff, err := db.get(ctx, k)
	if err != nil {
		return nil, err
	}

	dec := codec.NewDecoder(bytes.NewReader(buff), &codec.CborHandle{})
	var jets jet.IDSet
	err = dec.Decode(&jets)
	if err != nil {
		return nil, err
	}

	return jets, nil
}

func dropSizesPrefixKey(jetID core.RecordID) []byte {
	return prefixkey(scopeIDSystem, []byte{sysDropSizeHistory}, jetID.Bytes())
}

// AddDropSize adds Jet drop size stats (required for split decision).
func (db *DB) AddDropSize(ctx context.Context, dropSize *jet.DropSize) error {
	inslogger.FromContext(ctx).Debug("DB.AddDropSize starts ...")
	db.addBlockSizeLock.Lock()
	defer db.addBlockSizeLock.Unlock()

	k := dropSizesPrefixKey(dropSize.JetID)
	buff, err := db.get(ctx, k)
	if err != nil && err != ErrNotFound {
		return errors.Wrapf(err, "[ AddDropSize ] Can't get object: %s", string(k))
	}

	var dropSizes = jet.DropSizeHistory{}
	if err != ErrNotFound {
		dropSizes, err = jet.DeserializeJetDropSizeHistory(ctx, buff)
		if err != nil {
			return errors.Wrapf(err, "[ AddDropSize ] Can't decode dropSizes")
		}

		if len([]jet.DropSize(dropSizes)) >= db.jetSizesHistoryDepth {
			dropSizes = dropSizes[1:]
		}
	}

	dropSizes = append(dropSizes, *dropSize)

	return db.set(ctx, k, dropSizes.Bytes())
}

// SetDropSizeHistory saves drop sizes history.
func (db *DB) SetDropSizeHistory(ctx context.Context, jetID core.RecordID, dropSizeHistory jet.DropSizeHistory) error {
	inslogger.FromContext(ctx).Debug("DB.ResetDropSizeHistory starts ...")
	db.addBlockSizeLock.Lock()
	defer db.addBlockSizeLock.Unlock()

	k := dropSizesPrefixKey(jetID)
	err := db.set(ctx, k, dropSizeHistory.Bytes())
	return errors.Wrap(err, "[ ResetDropSizeHistory ] Can't db.set")
}

// GetDropSizeHistory returns last drops sizes.
func (db *DB) GetDropSizeHistory(ctx context.Context, jetID core.RecordID) (jet.DropSizeHistory, error) {
	inslogger.FromContext(ctx).Debug("DB.GetDropSizeHistory starts ...")
	db.addBlockSizeLock.RLock()
	defer db.addBlockSizeLock.RUnlock()

	k := dropSizesPrefixKey(jetID)
	buff, err := db.get(ctx, k)
	if err != nil && err != ErrNotFound {
		return nil, errors.Wrap(err, "[ GetDropSizeHistory ] Can't db.set")
	}

	if err == ErrNotFound {
		return jet.DropSizeHistory{}, nil
	}

	dropSizes, err := jet.DeserializeJetDropSizeHistory(ctx, buff)
	if err != nil {
		return nil, errors.Wrapf(err, "[ GetDropSizeHistory ] Can't decode dropSizes")
	}

	return dropSizes, nil
}
