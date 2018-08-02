package goplugin

import (
	"testing"

	"plugin"

	"reflect"

	"github.com/insolar/insolar/logicrunner"

	"bytes"

	"github.com/2tvenom/cbor"
)

type HelloWorlder struct {
	Greeted int
}

func TestHelloWorld(t *testing.T) {
	gp, err := NewGoPlugin("localhost:7777", "localhost:7778")
	defer gp.Stop()
	if err != nil {
		t.Fatal(err)
	}
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

var PATH = "./testplugins"

func TestConfigLoad(t *testing.T) {
	plugin, err := plugin.Open(PATH + "/secondary.so")
	if err != nil {
		t.Fatal(err)
	}

	hw, err := plugin.Lookup("EXP")
	r := reflect.ValueOf(hw)
	m := r.MethodByName("Hello")
	ret := m.Call([]reflect.Value{})
	t.Logf("%+v", ret)
}
