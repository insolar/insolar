// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package merkle

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
)

type GlobulaLeaf struct {
	// ByteSize = 16

	NodeID insolar.ShortNodeID // ByteSize = 4

	// ByteSize = 4
	NodeRole   member.PrimaryRole // 8
	PowerTotal uint32             // 23

	// ByteSize = 4
	NodePower member.Power // 8
	PowerBase uint32       // 23

	// ByteSize = 4
	RoleIndex uint16        // 10
	RoleTotal uint16        // 10
	NodeTotal uint16        // 10
	OpMode    member.OpMode // 4
}
