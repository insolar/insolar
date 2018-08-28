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

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage/leveldb/leveltestutils"
)

func TestCreateJetDrop_CreatesCorrectDrop(t *testing.T) {
	ledger, cleaner := leveltestutils.TmpDB(t, "")
	defer cleaner()

	jc := &JetCoordinator{
		storage: ledger,
	}
	var (
		zeropulse record.PulseNum
		pulse1    record.PulseNum = 1
		pulse2    record.PulseNum = 2
	)
	// it references on 'fake' zero
	fakeDrop := jetdrop.JetDrop{
		Hash: []byte{0xFF},
	}
	// save zero drop, which references on 'fake' drop with '0xFF' hash.
	dropz, err := ledger.SetDrop(zeropulse, &fakeDrop)
	assert.NoError(t, err)
	assert.NotNil(t, dropz)

	// save pulse1 records
	ledger.SetPulseFn(func() record.PulseNum { return pulse1 })
	ledger.SetRecord(&record.CodeRecord{})
	ledger.SetRecord(&record.ClassActivateRecord{})
	ledger.SetRecord(&record.ObjectActivateRecord{})
	// trigger new pulse on coordinator
	// (should save non zero Pulse)
	drop1, err := jc.Pulse(pulse2)
	assert.NoError(t, err)
	assert.NotNil(t, drop1)
	assert.Equal(t, dropz.Hash, drop1.PrevHash)
}
