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

package jetcoordinator

import (
	"bytes"
	"sort"

	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
)

type sortHashes [][]byte

func (s sortHashes) Less(i, j int) bool {
	return bytes.Compare(s[i], s[j]) < 0
}
func (s sortHashes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s sortHashes) Len() int {
	return len(s)
}

// CreateJetDrop creates new jet drop from records in storage. Previous pulse is required to calculate hash.
func CreateJetDrop(storage storage.LedgerStorer, prevPulse, newPulse record.PulseNum) (*jetdrop.JetDrop, error) {
	prevDrop, err := storage.GetDrop(prevPulse)
	if err != nil {
		return nil, err
	}

	// TODO: implement GetPulseKeys
	recordHashes, err := storage.GetPulseKeys(newPulse)
	if err != nil {
		return nil, err
	}
	sort.Sort(sortHashes(recordHashes))

	prevHash, err := prevDrop.Hash()
	if err != nil {
		return nil, err
	}
	drop := jetdrop.JetDrop{
		PrevHash:     prevHash,
		RecordHashes: recordHashes, // TODO: use merkle tree root hash here
	}

	return &drop, nil
}
