package goplugin

import (
	"testing"

	"github.com/insolar/insolar/logicrunner"
)

func TestHelloWorld(t *testing.T) {
	gp, err := NewGoPlugin("localhost:7777")
	defer gp.Stop()
	if err != nil {
		t.Fatal(err)
	}
	obj := logicrunner.Object{
		MachineType: logicrunner.MachineTypeGoPlugin,
	}
	ret, err := gp.Exec(obj, "Hello", logicrunner.Arguments{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ret)

}
