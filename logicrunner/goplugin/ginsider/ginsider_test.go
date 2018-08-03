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

	var dataBuf bytes.Buffer
	cbor := cbor.NewEncoder(&dataBuf)
	cbor.Marshal(HelloWorlder{66})

	gi := NewGoInsider(dir)
	req := goplugin.CallReq{
		Object: logicrunner.Object{
			MachineType: logicrunner.MachineTypeGoPlugin,
			Reference:   "secondary.so",
			Code:        code,
			Data:        dataBuf.Bytes(),
		},
		Method: "Hello",
		Args:   logicrunner.Arguments{},
	}
	resp := goplugin.CallResp{}
	gi.Call(req, &resp)

	var newData HelloWorlder
	_, err = cbor.Unmarshal(resp.Data, &newData)
	if err != nil {
		panic(err)
	}

	if newData.Greeted != 67 {
		t.Fatalf("Got unexpected value: %d, 67 is expected", newData.Greeted)
	}
}
