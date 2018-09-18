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

package pulsar

import (
	"crypto/ecdsa"

	"github.com/cenkalti/rpc2"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
)

type Neighbour struct {
	ConnectionType    configuration.ConnectionType
	ConnectionAddress string
	Client            *rpc2.Client
	PublicKey         *ecdsa.PublicKey
}

type RequestType string

const (
	Handshake RequestType = "handshake"
)

func (state RequestType) String() string {
	return string(state)
}

type HandshakePayload struct {
	Entropy core.Entropy
}

type Payload struct {
	PublicKey string
	Signature []byte
	Body      interface{}
}
