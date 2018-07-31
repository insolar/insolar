package hash

import (
	"io"

	"golang.org/x/crypto/sha3"
)

// Writer is the interface that wraps the WriteHash method.
//
// WriteHash should write all required for proper hashing data to io.Writer.
type Writer interface {
	WriteHash(io.Writer)
}

// SHA3hash224 returns SHA3 hash calculated on data recieved from Writer.
func SHA3hash224(hw Writer) []byte {
	h := sha3.New224()
	hw.WriteHash(h)
	return h.Sum(nil)
}
