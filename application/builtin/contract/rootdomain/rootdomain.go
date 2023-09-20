package rootdomain

import (
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

// RootDomain is smart contract representing entrance point to system.
type RootDomain struct {
	foundation.BaseContract
}

func (r *RootDomain) Test() error {
	return nil
}
