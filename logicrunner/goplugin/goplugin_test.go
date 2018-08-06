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

	err := os.Chdir(d + "/ginsider")
	if err != nil {
		return errors.Wrap(err, "couldn't chdir")
	}

	defer os.Chdir(d) // nolint: errcheck

	err = exec.Command("go", "build", "ginsider.go").Run()
	if err != nil {
		return errors.Wrap(err, "can't build ginsider")
	}

	err = os.Chdir(d + "/testplugins")
	if err != nil {
		return errors.Wrap(err, "couldn't chdir")
	}

	err = exec.Command("make", "secondary.so").Run()
	if err != nil {
		return errors.Wrap(err, "can't build pluigins")
	}
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
	defer os.RemoveAll(dir) // nolint: errcheck

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
	_, err = e.Marshal(HelloWorlder{77})
	if err != nil {
		t.Fatal(err)
	}

	obj := logicrunner.Object{
		MachineType: logicrunner.MachineTypeGoPlugin,
		Reference:   "secondary",
		Data:        buff.Bytes(),
	}

	data, _, err := gp.Exec(obj, "Hello", logicrunner.Arguments{})
	if err != nil {
		t.Fatal(err)
	}

	var newData HelloWorlder
	_, err = e.Unmarshal(data, &newData)
	if err != nil {
		panic(err)
	}
	if newData.Greeted != 78 {
		t.Fatalf("Got unexpected value: %d, 78 is expected", newData.Greeted)
	}

	//TODO: check second returned value
}
