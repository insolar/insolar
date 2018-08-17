package jetdrop

import (
	"golang.org/x/crypto/sha3"
)

// JetDrop is a blockchain block. It contains hashes from all records from slot.
type JetDrop struct {
	PrevHash     []byte
	RecordHashes [][]byte // TODO: this should be a byte slice that represents the merkle tree root of records
}

// Hash calculates jet drop hash. Raw data for hash should contain previous hash and merkle tree hash from records.
func (jd *JetDrop) Hash() ([]byte, error) {
	encoded, err := EncodeJetDrop(jd)
	if err != nil {
		return nil, err
	}
	h := sha3.New224()
	_, err = h.Write(encoded)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
