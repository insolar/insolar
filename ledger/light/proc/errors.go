// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc

import (
	"errors"
)

var (
	ErrNotExecutor      = errors.New("not executor for jet")
	ErrExecutorMismatch = errors.New("sender isn't executor for object")
	ErrNotActivated     = errors.New("object should be activated")
)
