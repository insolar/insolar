package buitin

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/logicrunner/buitin/helloworld"
)

type Contract interface {
	CodeRef() core.RecordRef
	Call()
}

type BuiltIn struct {
	AM       core.ArtifactManager
	MR       core.MessageRouter
	registry map[string]Contract
}

func NewBuiltIn(am *core.ArtifactManager, mr *core.MessageRouter) *BuiltIn {
	bi := BuiltIn{
		AM:       am,
		MR:       mr,
		registry: make(map[string]Contract),
	}
	hw := helloworld.NewHelloWorld()
	bi.registry[hw.CodeRef().String()] = hw
	return &bi
}

func (b *BuiltIn) Exec(codeRef logicrunner.Reference, data []byte, method string, args logicrunner.Arguments) (newObjectState []byte, methodResults logicrunner.Arguments, err error) {
	panic("implement me")
}
