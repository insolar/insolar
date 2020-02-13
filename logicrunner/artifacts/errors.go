// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package artifacts

import (
	"errors"
)

// Custom errors possibly useful to check by artifact manager callers.
var (
	ErrObjectDeactivated = errors.New("object is deactivated")
	ErrNotFound          = errors.New("object not found")
	ErrNoReply           = errors.New("timeout while awaiting reply from watermill")
)
