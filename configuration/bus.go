// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

import (
	"time"
)

// Bus holds some timeout options
type Bus struct {
	ReplyTimeout time.Duration
}

func NewBus() Bus {
	return Bus{
		ReplyTimeout: 15 * time.Second,
	}
}
