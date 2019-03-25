//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package network

import (
	"github.com/insolar/insolar/insolar"
)

type testNetwork struct {
}

func (n *testNetwork) SendMessage(nodeID insolar.Reference, method string, msg insolar.Parcel) ([]byte, error) {
	return make([]byte, 0), nil
}
func (n *testNetwork) SendCascadeMessage(data insolar.Cascade, method string, msg insolar.Parcel) error {
	return nil
}
func (n *testNetwork) RemoteProcedureRegister(name string, method insolar.RemoteProcedure) {

}

func GetTestNetwork() insolar.Network {
	return &testNetwork{}
}
