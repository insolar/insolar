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

package requester

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

const TESTREFERENCE = "222222"
const TESTSEED = "VGVzdA=="

var testSeedResponse = seedResponse{Seed: []byte("Test"), TraceID: "testTraceID"}
var testInfoResponse = InfoResponse{RootMember: "root_member_ref", RootDomain: "root_domain_ref", NodeDomain: "node_domain_ref"}
var testStatusResponse = StatusResponse{NetworkState: "OK"}

type rpcRequest struct {
	RPCVersion string `json:"jsonrpc"`
	Method     string `json:"method"`
}

func writeReponse(response http.ResponseWriter, answer map[string]interface{}) {
	serJSON, err := json.MarshalIndent(answer, "", "    ")
	if err != nil {
		log.Errorf("Can't serialize response\n")
	}
	var newLine byte = '\n'
	_, err = response.Write(append(serJSON, newLine))
	if err != nil {
		log.Errorf("Can't write response\n")
	}
}

func FakeHandler(response http.ResponseWriter, req *http.Request) {
	response.Header().Add("Content-Type", "application/json")

	params := api.Request{}
	_, err := api.UnmarshalRequest(req, &params)
	if err != nil {
		log.Errorf("Can't read request\n")
		return
	}

	answer := map[string]interface{}{}
	if params.Method == "CreateMember" {
		answer["reference"] = TESTREFERENCE
	} else {
		answer["random_data"] = TESTSEED
	}

	writeReponse(response, answer)
}

func FakeRPCHandler(response http.ResponseWriter, req *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	answer := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "",
	}
	rpcReq := rpcRequest{}
	_, err := api.UnmarshalRequest(req, &rpcReq)
	if err != nil {
		log.Errorf("Can't read request\n")
		return
	}

	switch rpcReq.Method {
	case "status.Get":
		answer["result"] = testStatusResponse
	case "info.Get":
		answer["result"] = testInfoResponse
	case "seed.Get":
		answer["result"] = testSeedResponse
	}
	writeReponse(response, answer)
}

const callLOCATION = "/api/call"
const rpcLOCATION = "/api/rpc"
const PORT = "12221"
const HOST = "127.0.0.1"
const URL = "http://" + HOST + ":" + PORT + "/api"

var server = &http.Server{Addr: ":" + PORT}

func waitForStart() error {
	numAttempts := 5

	for ; numAttempts > 0; numAttempts-- {
		conn, _ := net.DialTimeout("tcp", net.JoinHostPort(HOST, PORT), time.Millisecond*50)
		if conn != nil {
			conn.Close()
			break
		}
	}
	if numAttempts == 0 {
		return errors.New("Problem with launching test api: couldn't wait more")
	}

	return nil
}

func startServer() error {
	server := &http.Server{}
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 12221})
	if err != nil {
		return errors.Wrap(err, "error creating listener")
	}
	go server.Serve(listener)

	return nil
}

func setup() error {
	fh := FakeHandler
	fRPCh := FakeRPCHandler
	http.HandleFunc(callLOCATION, fh)
	http.HandleFunc(rpcLOCATION, fRPCh)
	log.Info("Starting Test api server ...")

	err := startServer()
	if err != nil {
		log.Error("Problem with starting test server: ", err)
		return errors.Wrap(err, "[ setup ]")
	}

	err = waitForStart()
	if err != nil {
		log.Error("Can't start api: ", err)
		return errors.Wrap(err, "[ setup ]")
	}

	return nil
}

func teardown() {
	const timeOut = 2
	log.Infof("Shutting down test server gracefully ...(waiting for %d seconds)", timeOut)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		fmt.Println("STOPPING TEST SERVER:", err)

	}
}

func testMainWrapper(m *testing.M) int {
	err := setup()
	defer teardown()
	if err != nil {
		fmt.Println("error while setup, skip tests: ", err)
		return 1
	}
	code := m.Run()
	return code
}

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func TestGetSeed(t *testing.T) {
	seed, err := GetSeed(URL)
	require.NoError(t, err)
	decodedSeed, err := base64.StdEncoding.DecodeString(TESTSEED)
	require.NoError(t, err)
	require.Equal(t, decodedSeed, seed)
}

func TestGetResponseBodyEmpty(t *testing.T) {
	_, err := GetResponseBody("test", PostParams{})
	require.EqualError(t, err, "[ getResponseBody ] Problem with sending request: Post test: unsupported protocol scheme \"\"")
}

func TestGetResponseBodyBadHttpStatus(t *testing.T) {
	_, err := GetResponseBody(URL+"TEST", PostParams{})
	require.EqualError(t, err, "[ getResponseBody ] Bad http response code: 404")
}

func TestGetResponseBody(t *testing.T) {
	data, err := GetResponseBody(URL+"/call", PostParams{})
	require.NoError(t, err)
	require.Contains(t, string(data), `"random_data": "VGVzdA=="`)
}

func TestSetVerbose(t *testing.T) {
	require.False(t, verbose)
	SetVerbose(true)
	require.True(t, verbose)
}

func readConfigs(t *testing.T) (*UserConfigJSON, *RequestConfigJSON) {
	userConf, err := ReadUserConfigFromFile("testdata/userConfig.json")
	require.NoError(t, err)
	reqConf, err := ReadRequestConfigFromFile("testdata/requestConfig.json")
	require.NoError(t, err)

	return userConf, reqConf
}

func TestSend(t *testing.T) {
	ctx := inslogger.ContextWithTrace(context.Background(), "TestSend")
	userConf, reqConf := readConfigs(t)
	resp, err := Send(ctx, URL, userConf, reqConf)
	require.NoError(t, err)
	require.Contains(t, string(resp), TESTREFERENCE)
}

func TestSendWithSeed(t *testing.T) {
	ctx := inslogger.ContextWithTrace(context.Background(), "TestSendWithSeed")
	userConf, reqConf := readConfigs(t)
	resp, err := SendWithSeed(ctx, URL+"/call", userConf, reqConf, []byte(TESTSEED))
	require.NoError(t, err)
	require.Contains(t, string(resp), TESTREFERENCE)
}

func TestSendWithSeed_WithBadUrl(t *testing.T) {
	ctx := inslogger.ContextWithTrace(context.Background(), "TestSendWithSeed_WithBadUrl")
	userConf, reqConf := readConfigs(t)
	_, err := SendWithSeed(ctx, URL+"TTT", userConf, reqConf, []byte(TESTSEED))
	require.EqualError(t, err, "[ Send ] Problem with sending target request: [ getResponseBody ] Bad http response code: 404")
}

func TestSendWithSeed_NilConfigs(t *testing.T) {
	ctx := inslogger.ContextWithTrace(context.Background(), "TestSendWithSeed_NilConfigs")
	_, err := SendWithSeed(ctx, URL, nil, nil, []byte(TESTSEED))
	require.EqualError(t, err, "[ Send ] Configs must be initialized")
}

func TestInfo(t *testing.T) {
	resp, err := Info(URL)
	require.NoError(t, err)
	require.Equal(t, resp, &testInfoResponse)
}

func TestStatus(t *testing.T) {
	resp, err := Status(URL)
	require.NoError(t, err)
	require.Equal(t, resp, &testStatusResponse)
}
