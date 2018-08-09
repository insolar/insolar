// Package foundation emulates foundation of types for golang contracts
package foundation

import "github.com/insolar/insolar/logicrunner"

type CBORMarshaler interface {
	Marshal(interface{}) []byte
	Unmarshal(interface{}, []byte)
}

// Call other contract via network dispatcher
func Call(Reference logicrunner.Reference, MethodName string, Arguments []interface{}) ([]interface{}, error) {
	return nil, nil
}

func APICall() { // GetPulsar / GetNodeList / GetValidatorCandidates
	return
}

// ???
type CallContext struct {
	Caller logicrunner.Reference
}
