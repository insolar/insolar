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

package requester

type rpcResponse struct {
	RPCVersion string                 `json:"jsonrpc"`
	Error      map[string]interface{} `json:"error"`
}

type seedResponse struct {
	Seed    string `json:"Seed"`
	TraceID string `json:"TraceID"`
}
type rpcSeedResponse struct {
	rpcResponse
	Result seedResponse `json:"result"`
}

// StatusResponse represents response from rpc on node.GetStatus method
type StatusResponse struct {
	NetworkState string `json:"NetworkState"`
}

type rpcStatusResponse struct {
	rpcResponse
	Result StatusResponse `json:"result"`
}

// InfoResponse represents response from rpc on network.GetInfo method
type InfoResponse struct {
	RootDomain             string   `json:"RootDomain"`
	RootMember             string   `json:"RootMember"`
	MigrationAdminMember   string   `json:"MigrationAdminMember"`
	MigrationDaemonMembers []string `json:"MigrationDaemonMembers"`
	NodeDomain             string   `json:"NodeDomain"`
	TraceID                string   `json:"TraceID"`
}

type rpcInfoResponse struct {
	rpcResponse
	Result InfoResponse `json:"result"`
}
