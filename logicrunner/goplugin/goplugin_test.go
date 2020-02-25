// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package goplugin

import (
	"testing"

	"github.com/insolar/insolar/insolar"
)

func TestTypeCompatibility(t *testing.T) {
	var _ insolar.MachineLogicExecutor = (*GoPlugin)(nil)
}
