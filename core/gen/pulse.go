package gen

import (
	"github.com/google/gofuzz"
	"github.com/insolar/insolar/core"
)

// PulseNumber generates random pulse number (excluding special cases).
func PulseNumber() (pn core.PulseNumber) {
	fuzz.New().Fuzz(&pn)
	return
}
