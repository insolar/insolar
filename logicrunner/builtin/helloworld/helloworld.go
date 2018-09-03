package helloworld

import "github.com/insolar/insolar/core"

type HelloWorld struct {
	Greeted int
}

func CodeRef() core.RecordRef {
	var ref core.RecordRef
	ref[core.RecordRefSize-1] = 1
	return ref
}

func NewHelloWorld() *HelloWorld {
	return &HelloWorld{}
}

func (hw *HelloWorld) Greet(name string) string {
	hw.Greeted++
	return "Hello " + name + "'s world"
}
