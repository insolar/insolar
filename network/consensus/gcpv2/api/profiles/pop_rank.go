package profiles

import "github.com/insolar/insolar/network/consensus/gcpv2/api/member"

type PopulationRank struct {
	Profile ActiveNode
	Power   member.Power
	OpMode  member.OpMode
}
