package primary

import (
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/logicrunner/goplugin/testplugins/foundation"
)

type INFHwRunner interface {
	Run() string
}

// @inscontract
// nolint
type HwRunner struct {
	Runned int
}

func (h *HwRunner) Run() string {
	hw := Hw.GetNewInstance(logicrunner.Reference("#1.#2"))
	return hw.Echo("Ooops")
}

//
type Hw struct {
	Reference logicrunner.Reference
}

func (Hw) GetNewInstance(r logicrunner.Reference) Hw {
	return Hw{Reference: r}
}

func (_self *Hw) Echo(s string) string {
	foundation.APICall(_self.Reference, "Echo", s)
}
