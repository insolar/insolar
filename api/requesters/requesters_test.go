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

package requesters

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
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TESTREFERENCE = "222222"
const TESTSEED = "VGVzdA=="

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

	params, err := api.PreprocessRequest(req)
	if err != nil {
		log.Errorf("Can't read request\n")
		return
	}

	qtype := api.QTypeFromString(params.QueryType)
	answer := map[string]interface{}{}
	if qtype == api.GetSeed {
		answer[api.SEED] = TESTSEED
	} else if params.Method == "CreateMember" {
		answer[api.REFERENCE] = TESTREFERENCE
	} else {
		answer["random_data"] = TESTSEED
	}

	writeReponse(response, answer)
}

const LOCATION = "/api/v1"
const PORT = "12221"
const HOST = "127.0.0.1"
const URL = "http://" + HOST + ":" + PORT + LOCATION

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
	http.HandleFunc(LOCATION, fh)
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
	assert.NoError(t, err)
	decodedSeed, err := base64.StdEncoding.DecodeString(TESTSEED)
	assert.NoError(t, err)
	assert.Equal(t, decodedSeed, seed)
}

func TestGetResponseBodyBadRequest(t *testing.T) {
	_, err := GetResponseBody("test", PostParams{})
	assert.EqualError(t, err, "[ getResponseBody ] Problem with sending request: Post test: unsupported protocol scheme \"\"")
}

func TestGetResponseBodyBadHttpStatus(t *testing.T) {
	_, err := GetResponseBody(URL+"TEST", PostParams{})
	assert.EqualError(t, err, "[ getResponseBody ] Bad http response code: 404")
}

func TestGetResponseBody(t *testing.T) {
	data, err := GetResponseBody(URL, PostParams{})
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"random_data": "VGVzdA=="`)
}

func TestSetVerbose(t *testing.T) {
	assert.False(t, verbose)
	SetVerbose(true)
	assert.True(t, verbose)
}

func readConfigs(t *testing.T) (*UserConfigJSON, *RequestConfigJSON) {
	userConf, err := ReadUserConfigFromFile("testdata/userConfig.json")
	require.NoError(t, err)
	reqConf, err := ReadRequestConfigFromFile("testdata/requestConfig.json")
	require.NoError(t, err)

	return userConf, reqConf
}

func TestSend(t *testing.T) {
	userConf, reqConf := readConfigs(t)
	resp, err := Send(URL, userConf, reqConf)
	assert.NoError(t, err)
	assert.Contains(t, string(resp), TESTREFERENCE)
}

func TestSendWithSeed(t *testing.T) {
	userConf, reqConf := readConfigs(t)
	resp, err := SendWithSeed(URL, userConf, reqConf, []byte(TESTSEED))
	assert.NoError(t, err)
	assert.Contains(t, string(resp), TESTREFERENCE)
}

func TestSendWithSeed_WithBadUrl(t *testing.T) {
	userConf, reqConf := readConfigs(t)
	_, err := SendWithSeed(URL+"TTT", userConf, reqConf, []byte(TESTSEED))
	assert.EqualError(t, err, "[ Send ] Problem with sending target request: [ getResponseBody ] Bad http response code: 404")
}

func TestSendWithSeed_NilConfigs(t *testing.T) {
	_, err := SendWithSeed(URL, nil, nil, []byte(TESTSEED))
	assert.EqualError(t, err, "[ Send ] Configs must be initialized")
}

func TestSend_BadSeedUrl(t *testing.T) {
	userConf, reqConf := readConfigs(t)
	_, err := Send(URL+"TTT", userConf, reqConf)
	assert.EqualError(t, err, "[ Send ] Problem with getting seed: [ getSeed ]: [ getResponseBody ] Bad http response code: 404")
}
