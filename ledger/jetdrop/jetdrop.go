package jetdrop

import (
	"golang.org/x/crypto/sha3"
)

type JetDrop struct {
	PrevHash     []byte
	RecordHashes [][]byte
}

func (jd *JetDrop) Hash() ([]byte, error) {
	// TODO: hash records with merkle tree
	encoded, err := EncodeJetDrop(jd)
	if err != nil {
		return nil, err
	}
	h := sha3.New224()
	h.Write(encoded)
	return h.Sum(nil), nil
}
