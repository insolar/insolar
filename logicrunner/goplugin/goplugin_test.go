package goplugin

import (
	"testing"

	"plugin"

	"reflect"

	"github.com/insolar/insolar/logicrunner"

	"bytes"

	"os"

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

var PATH = "/Users/vany/go/src/github.com/insolar/insolar/logicrunner/goplugin/ginsider/plugins"

func TestConfigLoad(t *testing.T) {
	t.Fatal(os.Getwd())
	pl, err := plugin.Open(PATH + "/main.so")
	if err != nil {
		t.Fatal(err)
	}

	pl2, err := plugin.Open(PATH + "/secondary/main.so")
	if err != nil {
		t.Fatal(err)
	}

	hw, err := pl.Lookup("EXP")
	r := reflect.ValueOf(hw)
	m := r.MethodByName("Hello")
	ret := m.Call([]reflect.Value{})
	t.Logf("%+v", ret)

	hw2, err := pl2.Lookup("EXP")
	r2 := reflect.ValueOf(hw2)
	m2 := r2.MethodByName("Hello")
	ret2 := m2.Call([]reflect.Value{})
	t.Logf("%+v", ret2)

}
