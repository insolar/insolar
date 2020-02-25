// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package node

import (
	"github.com/pkg/errors"
)

var (
	// ErrOverride is returned when trying to set nodes for non-empty pulse.
	ErrOverride = errors.New("node override is forbidden")
	// ErrNoNodes is returned when nodes for specified criteria could not be found.
	ErrNoNodes = errors.New("matching nodes not found")
)
