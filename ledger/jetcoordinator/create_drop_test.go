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
	"testing"

	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage/leveldb"
	"github.com/stretchr/testify/assert"
)

func TestCreateJetDrop_CreatesCorrectDrop(t *testing.T) {
	// TODO: remove binding to leveldb here
	ledger, _ := leveldb.InitDB()

	prevDrop := jetdrop.JetDrop{PrevHash: []byte{4, 5}}
	prevHash, _ := prevDrop.Hash()
	ledger.SetDrop(1, &prevDrop)
	ledger.SetRecord(&record.CodeRecord{})
	ledger.SetRecord(&record.ClassActivateRecord{})
	ledger.SetRecord(&record.ObjectActivateRecord{})

	drop, err := CreateJetDrop(ledger, 1, 2)
	assert.NoError(t, err)
	assert.Equal(t, jetdrop.JetDrop{
		PrevHash:     prevHash,
		RecordHashes: [][]byte{}, // TODO: after implementing storage.GetPulseKeys should contain created records
	}, *drop)
}
