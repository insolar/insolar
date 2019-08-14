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

type Response struct {
	JSONRPC string `json:"jsonrpc"`
	ID      uint64 `json:"id"`
	Error   *Error `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    Data   `json:"data,omitempty"`
}

type Data struct {
	TraceID string `json:"traceID,omitempty"`
}

type ContractResponse struct {
	Response
	Result *ContractResult `json:"result,omitempty"`
}

type ContractResult struct {
	CallResult       interface{} `json:"callResult,omitempty"`
	RequestReference string      `json:"requestReference,omitempty"`
	TraceID          string      `json:"traceID,omitempty"`
}

type seedResponse struct {
	Seed    string `json:"seed"`
	TraceID string `json:"traceID"`
}
type rpcSeedResponse struct {
	Response
	Result seedResponse `json:"result"`
}

// StatusResponse represents response from rpc on node.getStatus method
type StatusResponse struct {
	NetworkState string `json:"networkState"`
}

type rpcStatusResponse struct {
	Response
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
	Response
	Result InfoResponse `json:"result"`
}
