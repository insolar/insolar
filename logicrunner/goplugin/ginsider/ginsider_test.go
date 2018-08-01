package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/2tvenom/cbor"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/logicrunner/goplugin"
)

type HelloWorlder struct {
	Greeted int
}

func TestHelloWorld(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir) // clean up

	code, err := ioutil.ReadFile("../testplugins/secondary.so")
	if err != nil {
		panic(err)
	}

	var data_buf bytes.Buffer
	cbor := cbor.NewEncoder(&data_buf)
	cbor.Marshal(HelloWorlder{66})

	gi := NewGoInsider(dir)
	req := goplugin.CallReq{
		logicrunner.Object{
			MachineType: logicrunner.MachineTypeGoPlugin,
			Reference:   "secondary.so",
			Code:        code,
			Data:        data_buf.Bytes(),
		},
		"Hello",
		logicrunner.Arguments{},
	}
	resp := goplugin.CallResp{}
	gi.Call(req, &resp)
}
