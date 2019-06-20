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
	Seed    string `json:"seed"`
	TraceID string `json:"traceID"`
}
type rpcSeedResponse struct {
	rpcResponse
	Result seedResponse `json:"result"`
}

// StatusResponse represents response from rpc on node.getStatus method
type StatusResponse struct {
	NetworkState string `json:"networkState"`
}

type rpcStatusResponse struct {
	rpcResponse
	Result StatusResponse `json:"result"`
}

// InfoResponse represents response from rpc on network.getInfo method
type InfoResponse struct {
	RootDomain             string   `json:"rootDomain"`
	RootMember             string   `json:"rootMember"`
	MigrationAdminMember   string   `json:"migrationAdminMember"`
	MigrationDaemonMembers []string `json:"migrationDaemonMembers"`
	NodeDomain             string   `json:"nodeDomain"`
	TraceID                string   `json:"traceID"`
}

type rpcInfoResponse struct {
	rpcResponse
	Result InfoResponse `json:"result"`
}
