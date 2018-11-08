/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package network

import (
	"crypto/ecdsa"

	"github.com/insolar/insolar/core"
	ecdsa2 "github.com/insolar/insolar/cryptoproviders/ecdsa"
)

type testNetwork struct {
}

func (n *testNetwork) GetNodeID() core.RecordRef {
	return core.NewRefFromBase58("v1")
}

func (n *testNetwork) SendMessage(nodeID core.RecordRef, method string, msg core.SignedMessage) ([]byte, error) {
	return make([]byte, 0), nil
}
func (n *testNetwork) SendCascadeMessage(data core.Cascade, method string, msg core.SignedMessage) error {
	return nil
}
func (n *testNetwork) GetAddress() string {
	return ""
}
func (n *testNetwork) RemoteProcedureRegister(name string, method core.RemoteProcedure) {

}
func (n *testNetwork) GetPrivateKey() *ecdsa.PrivateKey {
	key, _ := ecdsa2.GeneratePrivateKey()
	return key
}

func GetTestNetwork() core.Network {
	return &testNetwork{}
}
