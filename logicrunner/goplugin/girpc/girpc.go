package girpc

import "github.com/insolar/insolar/logicrunner"

// CallReq is a set of arguments for Call RPC in the runner
type CallReq struct {
	Object    logicrunner.Object
	Method    string
	Arguments logicrunner.Arguments
}

// CallResp is response from Call RPC in the runner
type CallResp struct {
	Data []byte
	Ret  logicrunner.Arguments
	Err  error
}
