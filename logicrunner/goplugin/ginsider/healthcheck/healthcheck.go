package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type One struct {
	foundation.BaseContract
}

var INSATTR_Check_API = true

func (t *One) Check() (bool, error) {
	return true, nil
}
