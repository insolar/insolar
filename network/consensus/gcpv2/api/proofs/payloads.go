package proofs

import (
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/pulse"
)

type OriginalPulsarPacket interface {
	longbits.FixedReader
	pulse.DataHolder
	OriginalPulsarPacket()
}
