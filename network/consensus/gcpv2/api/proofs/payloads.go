package proofs

import "github.com/insolar/insolar/network/consensus/common/longbits"

type OriginalPulsarPacket interface {
	longbits.FixedReader
	OriginalPulsarPacket()
}
