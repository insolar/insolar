package goplugin

import (
	"testing"

	"github.com/insolar/insolar/logicrunner"

	"bytes"

	"github.com/2tvenom/cbor"
)

type HelloWorlder struct {
	Greeted int
}

func TestHelloWorld(t *testing.T) {
	gp, err := NewGoPlugin(RunnerOptions{
		APIAddr:     "127.0.0.1:7778",
		InsiderAddr: "127.0.0.1:7777",
		StoragePath: ""})
	if err != nil {
		t.Fatal(err)
	}
	defer gp.Stop()

	var buff bytes.Buffer
	e := cbor.NewEncoder(&buff)
	e.Marshal(HelloWorlder{77})

	obj := logicrunner.Object{
		MachineType: logicrunner.MachineTypeGoPlugin,
		Reference:   "reference",
		Data:        buff.Bytes(),
	}

	data, ret, err := gp.Exec(obj, "Hello", logicrunner.Arguments{})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("len of data == 0")
	}
	//	if ret == logicrunner.Arguments{} // IDK, lets decide what must be here
	t.Log(ret)
}
