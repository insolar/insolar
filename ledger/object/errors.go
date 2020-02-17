// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package object

import (
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is returned when value not found.
	ErrNotFound = errors.New("object not found")
	// ErrOverride is returned when trying to update existing record with the same id.
	ErrOverride = errors.New("record override is forbidden")

	// ErrIndexNotFound is returned when an index not found.
	ErrIndexNotFound = errors.New("index not found")
)
