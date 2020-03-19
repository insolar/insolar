// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package testresponse

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
)

type PostParams map[string]interface{}

type RPCResponseInterface interface {
	GetRPCVersion() string
	GetError() map[string]interface{}
}

type RPCResponse struct {
	RPCVersion string                 `json:"jsonrpc"`
	Error      map[string]interface{} `json:"error"`
}

func (r *RPCResponse) GetRPCVersion() string {
	return r.RPCVersion
}

func (r *RPCResponse) GetError() map[string]interface{} {
	return r.Error
}

type getSeedResponse struct {
	RPCResponse
	Result struct {
		Seed    string `json:"seed"`
		TraceID string `json:"traceID"`
	} `json:"result"`
}

type InfoResponse struct {
	NodeDomain string `json:"NodeDomain"`
	TraceID    string `json:"TraceID"`
}

type rpcInfoResponse struct {
	RPCResponse
	Result InfoResponse `json:"result"`
}

type StatusResponse struct {
	NetworkState    string `json:"networkState"`
	WorkingListSize int    `json:"workingListSize"`
}

func CheckConvertRequesterError(t *testing.T, err error) *requester.Error {
	rv, ok := err.(*requester.Error)
	require.Truef(t, ok, "got wrong error %T (expected *requester.Error) with text '%s'", err, err.Error())
	return rv
}

func GetRPSResponseBody(t testing.TB, URL string, postParams map[string]interface{}) []byte {
	jsonValue, _ := json.Marshal(postParams)

	postResp, err := http.Post(URL, "application/json", bytes.NewBuffer(jsonValue))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	require.NoError(t, err)
	return body
}

func GetSeed(t testing.TB) string {
	body := GetRPSResponseBody(t, launchnet.TestRPCUrl, PostParams{
		"jsonrpc": "2.0",
		"method":  "node.getSeed",
		"id":      "",
	})
	getSeedResponse := &getSeedResponse{}
	UnmarshalRPCResponse(t, body, getSeedResponse)
	require.NotNil(t, getSeedResponse.Result)
	return getSeedResponse.Result.Seed
}

func GetInfo(t testing.TB) InfoResponse {
	pp := PostParams{
		"jsonrpc": "2.0",
		"method":  "network.getInfo",
		"id":      1,
		"params":  map[string]string{},
	}
	body := GetRPSResponseBody(t, launchnet.TestRPCUrl, pp)
	rpcInfoResponse := &rpcInfoResponse{}
	UnmarshalRPCResponse(t, body, rpcInfoResponse)
	require.NotNil(t, rpcInfoResponse.Result)
	return rpcInfoResponse.Result
}

func UnmarshalRPCResponse(t testing.TB, body []byte, response RPCResponseInterface) {
	err := json.Unmarshal(body, &response)
	require.NoError(t, err)
	require.Equal(t, "2.0", response.GetRPCVersion())
	require.Nil(t, response.GetError())
}

func UnmarshalCallResponse(t testing.TB, body []byte, response *requester.ContractResponse) {
	err := json.Unmarshal(body, &response)
	require.NoError(t, err)
}

