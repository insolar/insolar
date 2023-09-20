package gen

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
)

func TestGen_StringFromBytes(t *testing.T) {
	symbolFuzzer := fuzz.New().NilChance(0).NumElements(1, 100)
	var symbols []byte
	symbolFuzzer.Fuzz(&symbols)
	for i := 0; i < 100; i++ {
		s := StringFromBytes(symbols, i)
		require.GreaterOrEqualf(t, i, len(s), "string length should not be greater than `maxcount`")
		for _, sym := range []byte(s) {
			require.Contains(t, symbols, sym, "byte should be in range")
		}
	}
}
