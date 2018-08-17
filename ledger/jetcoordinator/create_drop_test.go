package jetcoordinator

import (
	"testing"

	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage/leveldb"
	"github.com/stretchr/testify/assert"
)

func TestCreateJetDrop_CreatesCorrectDrop(t *testing.T) {
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
