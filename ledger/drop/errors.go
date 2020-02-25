// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package drop

import (
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is returned when value was not found.
	ErrNotFound = errors.New("value not found")

	// ErrOverride is returned if something tries to update existing record.
	ErrOverride = errors.New("records override is forbidden")
)
