package jetdrop

import (
	"golang.org/x/crypto/sha3"
)

type JetDrop struct {
	PrevHash     []byte
	RecordHashes [][]byte // TODO: we should probable store merkle tree here
}

func (jd *JetDrop) Hash() ([]byte, error) {
	// TODO: use merkle tree root instead of records here
	encoded, err := EncodeJetDrop(jd)
	if err != nil {
		return nil, err
	}
	h := sha3.New224()
	h.Write(encoded)
	return h.Sum(nil), nil
}
