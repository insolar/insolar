// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor

import (
	"github.com/pkg/errors"
)

var (
	ErrWaiterNotLocked = errors.New("unlocked waiter unlock attempt")
	ErrWriteClosed     = errors.New("requested pulse is closed for writing")
)
