package jetdrop

import (
	"bytes"

	"github.com/ugorji/go/codec"
)

func EncodeJetDrop(drop *JetDrop) ([]byte, error) {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	err := enc.Encode(drop)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeJetDrop(buf []byte) (*JetDrop, error) {
	dec := codec.NewDecoder(bytes.NewReader(buf), &codec.CborHandle{})
	var drop JetDrop
	err := dec.Decode(&drop)
	if err != nil {
		return nil, err
	}
	return &drop, nil
}
