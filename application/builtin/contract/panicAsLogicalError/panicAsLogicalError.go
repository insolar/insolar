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
