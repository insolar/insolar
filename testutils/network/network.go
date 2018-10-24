package network

import (
	"crypto/ecdsa"

	"github.com/insolar/insolar/core"
)

type testNetwork struct {
}

func (n *testNetwork) GetNodeID() core.RecordRef {
	return core.NewRefFromBase58("v1")
}

func (n *testNetwork) SendMessage(nodeID core.RecordRef, method string, msg core.Message) ([]byte, error) {
	return make([]byte, 0), nil
}
func (n *testNetwork) SendCascadeMessage(data core.Cascade, method string, msg core.Message) error {
	return nil
}
func (n *testNetwork) GetAddress() string                                               { return "" }
func (n *testNetwork) RemoteProcedureRegister(name string, method core.RemoteProcedure) {}
func (n *testNetwork) GetPrivateKey() *ecdsa.PrivateKey                                 { return nil }

func GetTestNetwork() core.Network {
	return &testNetwork{}
}
