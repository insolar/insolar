// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
