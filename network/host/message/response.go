/*
 *    Copyright 2018 INS Ecosystem
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

package message

import (
	"github.com/insolar/insolar/network/host/node"
	"github.com/insolar/insolar/network/host/relay"
)

// ResponseDataFindNode is data for FindNode response.
type ResponseDataFindNode struct {
	Closest []*node.Node
}

// ResponseDataFindValue is data for FindValue response.
type ResponseDataFindValue struct {
	Closest []*node.Node
	Value   []byte
}

// ResponseDataStore is data for Store response.
type ResponseDataStore struct {
	Success bool
}

// ResponseDataRPC is data for RPC response.
type ResponseDataRPC struct {
	Success bool
	Result  []byte
	Error   string
}

// ResponseRelay is data for relay request response
type ResponseRelay struct {
	Success bool
	State   relay.State
}

// ResponseAuth is data for authentication request response
type ResponseAuth struct {
	AuthUniqueKey string
}
