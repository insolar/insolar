package gen

import (
	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/pulse"
)

// PulseNumber generates random pulse number (excluding special cases).
func PulseNumber() insolar.PulseNumber {
	f := fuzz.New().NilChance(0).Funcs(func(pn *insolar.PulseNumber, c fuzz.Continue) {
		*pn = insolar.PulseNumber(c.Int31n(pulse.MaxTimePulse-pulse.MinTimePulse) + pulse.MinTimePulse)
	})

	var pn insolar.PulseNumber
	f.Fuzz(&pn)
	return pn
}
