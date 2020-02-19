// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
