// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
