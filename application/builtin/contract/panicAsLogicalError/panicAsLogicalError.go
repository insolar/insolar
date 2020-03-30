// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package panicAsLogicalError

import "github.com/insolar/insolar/logicrunner/builtin/foundation"

type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

var INSATTR_Panic_API = true

func (r *One) Panic() error {
	panic("AAAAAAAA!")
	return nil
}
func NewPanic() (*One, error) {
	panic("BBBBBBBB!")
}
