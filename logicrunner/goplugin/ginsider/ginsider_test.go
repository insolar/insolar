package main

import (
	//"github.com/insolar/insolar/logicrunner/goplugin"
	"testing"

	"github.com/insolar/insolar/logicrunner/goplugin"
)

func TestHelloWorld(t *testing.T) {
	gi := GoInsider{"."}
	req := goplugin.CallReq{}
	resp := goplugin.CallResp{}
	gi.Call(req, &resp)

}
