package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/logicrunner/goplugin"
)

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

	gi := NewGoInsider(dir)
	req := goplugin.CallReq{
		logicrunner.Object{
			MachineType: logicrunner.MachineTypeGoPlugin,
			Reference:   "secondary.so",
			Code:        code,
			Data:        []byte{},
		},
		"Method",
		logicrunner.Arguments{},
	}
	resp := goplugin.CallResp{}
	gi.Call(req, &resp)
}
