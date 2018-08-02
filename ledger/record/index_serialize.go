package record

import (
	"bytes"

	"golang.org/x/crypto/sha3"
	"github.com/ugorji/go/codec"
)

// HashLifelineIndex returns 28 bytes of SHA3 hash
func HashLifelineIndex(buf []byte) Hash {
	return sha3.Sum224(buf)
}

// EncodeLifelineIndex converts lifeline index into binary format
func EncodeLifelineIndex(index *LifelineIndex) []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	err := enc.Encode(index)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// DecodeLifelineIndex converts byte array into lifeline index struct
func DecodeLifelineIndex(buf []byte) LifelineIndex {
	dec := codec.NewDecoder(bytes.NewReader(buf), &codec.CborHandle{})
	var index LifelineIndex
	err := dec.Decode(&index)
	if err != nil {
		panic(err)
	}
	return index
}
