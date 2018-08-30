package helloworld

import "github.com/insolar/insolar/core"

type HelloWorld struct {
	Greeted int
}

func (HelloWorld) CodeRef() core.RecordRef {
	ret := make([]byte, 64)
	ret[63] = 1
	return ret
}

func NewHelloWorld() *HelloWorld {
	return &HelloWorld{}
}

func (hw *HelloWorld) Greet(name string) string {
	hw.Greeted++
	return "Hello " + name + "'s world"
}
