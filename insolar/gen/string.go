// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package gen

import (
	fuzz "github.com/google/gofuzz"
)

// StringFromBytes generates random id with length from 0 to maxcount randomly filled by provided symbols.
func StringFromBytes(symbols []byte, maxcount int) string {
	if maxcount == 0 {
		return ""
	}
	f := fuzz.New().Funcs(func(b *[]byte, c fuzz.Continue) {
		count := c.Intn(maxcount)
		for i := 0; i < count; i++ {
			*b = append(*b, symbols[c.Intn(len(symbols))])
		}
	})
	var bstr []byte
	f.NilChance(0).Fuzz(&bstr)
	return string(bstr)
}
