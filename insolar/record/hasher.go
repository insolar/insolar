package record

import (
	"hash"
)

// HashVirtual returns hash for virtual record.
func HashVirtual(h hash.Hash, rec Virtual) []byte {
	// Signature must not affects material record hash calculating.
	rec.Signature = nil
	buf, err := rec.Marshal()
	if err != nil {
		panic(err)
	}
	_, err = h.Write(buf)
	if err != nil {
		panic(err)
	}
	return h.Sum(nil)
}
