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
