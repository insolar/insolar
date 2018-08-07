package goplugin

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/logicrunner"
)

type HelloWorlder struct {
	Greeted int
}

func (r *HelloWorlder) ProxyEcho(gp *GoPlugin, s string) {
	ch := new(codec.CborHandle)
	var data []byte
	err := codec.NewEncoderBytes(&data, ch).Encode(*r)
	if err != nil {
		panic(err)
	}

	obj := logicrunner.Object{
		MachineType: logicrunner.MachineTypeGoPlugin,
		Reference:   "secondary",
		Data:        data,
	}

	args := make([]logicrunner.Argument, 1)

	var arg []byte
	err = codec.NewEncoderBytes(&arg, ch).Encode(s)
	if err != nil {
		panic(err)
	}
	args[0] = arg

	data, _, err = gp.Exec(obj, "Echo", args)
	if err != nil {
		panic(err)
	}

	err = codec.NewDecoderBytes(data, ch).Decode(r)
	if err != nil {
		panic(err)
	}
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

	hw := &HelloWorlder{77}
	hw.ProxyEcho(gp, "hi")
	if hw.Greeted != 78 {
		t.Fatalf("Got unexpected value: %d, 78 is expected", hw.Greeted)
	}

	//TODO: check second returned value
}
