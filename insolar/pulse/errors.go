// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulse

import (
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is returned when value was not found.
	ErrNotFound = errors.New("pulse not found")
	// ErrBadPulse is returned when appended Pulse is less than the latest.
	ErrBadPulse = errors.New("pulse should be greater than the latest")
)
