package power

import (
	"github.com/insolar/insolar/network/consensus/common/capacity"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
)

type Request int16

func NewRequestByLevel(v capacity.Level) Request {
	return -Request(v)
}

func NewRequest(v member.Power) Request {
	return Request(v)
}

func (v Request) AsCapacityLevel() (bool, capacity.Level) {
	return v < 0, capacity.Level(-v)
}

func (v Request) AsMemberPower() (bool, member.Power) {
	return v >= 0, member.Power(v)
}
