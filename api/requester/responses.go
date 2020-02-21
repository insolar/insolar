// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package requester

import (
	"time"
)

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

func (e *Error) Error() string {
	return e.Message
}

type Data struct {
	Trace            []string `json:"trace,omitempty"`
	TraceID          string   `json:"traceID,omitempty"`
	RequestReference string   `json:"requestReference,omitempty"`
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

// SeedReply is reply for Seed service requests.
type SeedReply struct {
	Seed    []byte `json:"seed"`
	TraceID string `json:"traceID"`
}

type rpcStatusResponse struct {
	Response
	Result StatusResponse `json:"result"`
}

type Node struct {
	Reference string `json:"reference"`
	Role      string `json:"role"`
	IsWorking bool   `json:"isWorking"`
	ID        uint32 `json:"id"`
}

// StatusResponse represents response from rpc on node.getStatus method
type StatusResponse struct {
	NetworkState       string    `json:"networkState"`
	Origin             Node      `json:"origin"`
	ActiveListSize     int       `json:"activeListSize"`
	WorkingListSize    int       `json:"workingListSize"`
	Nodes              []Node    `json:"nodes"`
	PulseNumber        uint32    `json:"pulseNumber"`
	NetworkPulseNumber uint32    `json:"networkPulseNumber"`
	Entropy            []byte    `json:"entropy"`
	Version            string    `json:"version"`
	Timestamp          time.Time `json:"timestamp"`
	StartTime          time.Time `json:"startTime"`
}
