// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package requester

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/insolar/x-crypto/ecdsa"
	"github.com/insolar/x-crypto/elliptic"
	"github.com/insolar/x-crypto/rand"
	"github.com/insolar/x-crypto/sha256"

	"github.com/insolar/insolar/logicrunner/builtin/foundation"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
)

const TESTREFERENCE = "insolar:1MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"
const TESTSEED = "VGVzdA=="

var testSeedResponse = seedResponse{Seed: "Test", TraceID: "testTraceID"}
var testStatusResponse = StatusResponse{NetworkState: "OK"}

func writeReponse(response http.ResponseWriter, answer interface{}) {
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

type RPCResponse struct {
	Response
	Result interface{} `json:"result,omitempty"`
}

func FakeRPCHandler(response http.ResponseWriter, req *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	rpcResponse := RPCResponse{}
	request := Request{}
	_, err := unmarshalRequest(req, &request)
	if err != nil {
		log.Errorf("Can't read request\n")
		return
	}

	switch request.Method {
	case "node.getStatus":
		rpcResponse.Result = testStatusResponse
	case "node.getSeed":
		rpcResponse.Result = testSeedResponse
	case "contract.call":
		rpcResponse.Result = TESTREFERENCE
	default:
		rpcResponse.Result = TESTSEED

	}
	writeReponse(response, rpcResponse)
}

const rpcLOCATION = "/admin-api/rpc"
const PORT = "12221"
const HOST = "127.0.0.1"
const URL = "http://" + HOST + ":" + PORT + rpcLOCATION

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
	fRPCh := FakeRPCHandler
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
	require.Equal(t, "Test", seed)
}

func TestGetResponseBodyEmpty(t *testing.T) {
	_, err := GetResponseBodyPlatform("test", "", nil)
	require.EqualError(t, err, "problem with sending request: Post test: unsupported protocol scheme \"\"")
}

func TestGetResponseBodyBadHttpStatus(t *testing.T) {
	_, err := GetResponseBodyPlatform(URL+"TEST", "", nil)
	require.EqualError(t, err, "bad http response code: 404")
}

func TestGetResponseBody(t *testing.T) {
	data, err := GetResponseBodyContract(URL, ContractRequest{}, "")
	response := RPCResponse{}
	_ = json.Unmarshal(data, &response)
	require.NoError(t, err)
	require.Contains(t, response.Result, TESTSEED)
}

func TestSetVerbose(t *testing.T) {
	require.False(t, verbose)
	SetVerbose(true)
	require.True(t, verbose)
	// restore original value for future tests, if -count 10 flag is used
	SetVerbose(false)
}

func readConfigs(t *testing.T) (*UserConfigJSON, *Params) {
	userConf, err := ReadUserConfigFromFile("testdata/userConfig.json")
	require.NoError(t, err)
	reqConf, err := ReadRequestParamsFromFile("testdata/requestConfig.json")
	require.NoError(t, err)

	return userConf, reqConf
}

func TestSend(t *testing.T) {
	ctx := inslogger.ContextWithTrace(context.Background(), "TestSend")
	userConf, reqParams := readConfigs(t)
	reqParams.CallSite = "member.create"
	resp, err := Send(ctx, URL, userConf, reqParams)
	require.NoError(t, err)
	require.Contains(t, string(resp), TESTREFERENCE)
}

func TestSendWithSeed(t *testing.T) {
	ctx := inslogger.ContextWithTrace(context.Background(), "TestSendWithSeed")
	userConf, reqParams := readConfigs(t)
	reqParams.CallSite = "member.create"
	resp, err := SendWithSeed(ctx, URL, userConf, reqParams, TESTSEED)
	require.NoError(t, err)
	require.Contains(t, string(resp), TESTREFERENCE)
}

func TestSendWithSeed_WithBadUrl(t *testing.T) {
	ctx := inslogger.ContextWithTrace(context.Background(), "TestSendWithSeed_WithBadUrl")
	userConf, reqConf := readConfigs(t)
	_, err := SendWithSeed(ctx, URL+"TTT", userConf, reqConf, TESTSEED)
	require.EqualError(t, err, "[ SendWithSeed ] Problem with sending target request: bad http response code: 404")
}

func TestSendWithSeed_NilConfigs(t *testing.T) {
	ctx := inslogger.ContextWithTrace(context.Background(), "TestSendWithSeed_NilConfigs")
	_, err := SendWithSeed(ctx, URL, nil, nil, TESTSEED)
	require.EqualError(t, err, "[ SendWithSeed ] Problem with creating target request: configs must be initialized")
}

func TestStatus(t *testing.T) {
	resp, err := Status(URL)
	require.NoError(t, err)
	require.Equal(t, resp, &testStatusResponse)
}

func TestMarshalSig(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	msg := "test"
	hash := sha256.Sum256([]byte(msg))

	r1, s1, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	require.NoError(t, err)
	derString, err := marshalSig(r1, s1)
	require.NoError(t, err)

	sig, err := base64.StdEncoding.DecodeString(derString)

	r2, s2, err := foundation.UnmarshalSig(sig)
	require.NoError(t, err)

	require.Equal(t, r1, r2, errors.Errorf("Invalid S number"))
	require.Equal(t, s1, s2, errors.Errorf("Invalid R number"))
}

// unmarshalRequest unmarshals request to api
func unmarshalRequest(req *http.Request, params interface{}) ([]byte, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, errors.Wrap(err, "[ unmarshalRequest ] Can't read body. So strange")
	}
	if len(body) == 0 {
		return nil, errors.New("[ unmarshalRequest ] Empty body")
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		return body, errors.Wrap(err, "[ unmarshalRequest ] Can't unmarshal input params")
	}
	return body, nil
}
