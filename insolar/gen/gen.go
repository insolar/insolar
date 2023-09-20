package gen

import (
	fuzz "github.com/google/gofuzz"
)

// Signature generates random non nil bytes sequence of provided size.
func Signature(size int) []byte {
	if size < 0 {
		return nil
	}
	b := make([]byte, size)
	fuzz.New().NilChance(0).NumElements(size, size).Fuzz(&b)
	return b
}
