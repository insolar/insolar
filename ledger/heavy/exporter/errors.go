// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package exporter

import (
	"errors"
)

var (
	ErrNilCount          = errors.New("count can't be 0")
	ErrNotFinalPulseData = errors.New("trying to get a non-finalized pulse data")
)
