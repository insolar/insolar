package goplugin

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/2tvenom/cbor"

	"os/exec"

	"github.com/insolar/insolar/logicrunner"
	"github.com/pkg/errors"
)

type HelloWorlder struct {
	Greeted int
}

func compileBinaries() error {
	d, _ := os.Getwd()

	os.Chdir(d + "/ginsider")
	err := exec.Command("go", "build", "ginsider.go").Run()
	if err != nil {
		return errors.Wrap(err, "can't build ginsider")
	}

	os.Chdir(d + "/testplugins")
	err = exec.Command("make", "secondary.so").Run()
	if err != nil {
		return errors.Wrap(err, "can't build pluigins")
	}

	os.Chdir(d)
	return nil
}

func TestHelloWorld(t *testing.T) {
	if err := compileBinaries(); err != nil {
		t.Fatal("Can't compile binaries", err)
	}
	dir, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up

	gp, err := NewGoPlugin(
		Options{
			Listen:   "127.0.0.1:7778",
			CodePath: "./testplugins/",
		},
		RunnerOptions{
			Listen:          "127.0.0.1:7777",
			CodeStoragePath: dir,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	defer gp.Stop()

	var buff bytes.Buffer
	e := cbor.NewEncoder(&buff)
	e.Marshal(HelloWorlder{77})

	obj := logicrunner.Object{
		MachineType: logicrunner.MachineTypeGoPlugin,
		Reference:   "secondary",
		Data:        buff.Bytes(),
	}

	data, ret, err := gp.Exec(obj, "Hello", logicrunner.Arguments{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", data)
	if len(data) == 0 {
		t.Fatal("len of data == 0")
	}
	//	if ret == logicrunner.Arguments{} // IDK, lets decide what must be here
	t.Log(ret)
}
