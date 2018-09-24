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

package functest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const HOST = "http://localhost:19191"
const TestUrl = HOST + "/api/v1"
const insolarImportPath = "github.com/insolar/insolar"

var cmd *exec.Cmd
var stdin io.WriteCloser
var stdout io.ReadCloser
var insolardPath = filepath.Join(testdataPath(), "insolard")

func testdataPath() string {
	p, err := build.Default.Import("github.com/insolar/insolar", "", build.FindOnly)
	if err != nil {
		panic(err)
	}
	return filepath.Join(p.Dir, "testdata", "functional")
}

func functestPath() string {
	p, err := build.Default.Import("github.com/insolar/insolar", "", build.FindOnly)
	if err != nil {
		panic(err)
	}
	return filepath.Join(p.Dir, "functest")
}

func buildInsolard() error {
	_, err := exec.Command(
		"go", "build",
		"-o", insolardPath,
		insolarImportPath+"/cmd/insolard/",
	).CombinedOutput()
	return err
}

func createDirForContracts() error {
	return os.MkdirAll(filepath.Join(functestPath(), "contractstorage"), 0777)
}

func deleteDirForContracts() error {
	return os.RemoveAll(filepath.Join(functestPath(), "contractstorage"))
}

func deleteDirForData() error {
	return os.RemoveAll(filepath.Join(functestPath(), "data"))
}

func buildGinsiderCLI() error {
	_, _, err := testutil.Build()
	return err
}

func waitForLaunch(stdout io.ReadCloser) error {
	done := make(chan bool, 1)
	timeout := 30 * time.Second

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			if strings.Contains(line, "======= Host info ======") {
				done <- true
			}
		}
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return errors.Errorf("could't wait for launch: timeout of %s was exceeded", timeout)
	}

}

func startInsolard() error {
	cmd = exec.Command(
		insolardPath,
	)
	var err error

	stdin, err = cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}
	err = waitForLaunch(stdout)
	if err != nil {
		return err
	}
	return nil
}

func stopInsolard() error {
	if stdin != nil {
		defer stdin.Close()
	}
	if stdout != nil {
		defer stdout.Close()
	}
	if cmd.Process == nil {
		return nil
	}
	io.WriteString(stdin, "exit\n")
	err := cmd.Wait()
	if err != nil {
		fmt.Println("try to kill, wait done with error: ", err)
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Println("failed to kill process: ", err)
		}
	}
	return nil
}

func setup() error {
	err := createDirForContracts()
	if err != nil {
		return err
	}

	err = buildGinsiderCLI()
	if err != nil {
		return err
	}

	err = buildInsolard()
	if err != nil {
		return err
	}

	err = startInsolard()
	if err != nil {
		return err
	}

	return nil
}

func teardown() {
	err := stopInsolard()
	if err != nil {
		fmt.Println("failed to stop insolard: ", err)
	}

	err = deleteDirForData()
	if err != nil {
		fmt.Println("failed to remove data directory for func tests: ", err)
	}
	err = deleteDirForContracts()
	if err != nil {
		fmt.Println("failed to remove contractstorage directory for func tests: ", err)
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

type postParams map[string]interface{}

type errorResponse struct {
	Code  int    `json:"code"`
	Event string `json:"event"`
}

type responseInterface interface {
	getError() *errorResponse
}

type baseResponse struct {
	Qid string         `json:"qid"`
	Err *errorResponse `json:"error"`
}

func (r *baseResponse) getError() *errorResponse {
	return r.Err
}

type createMemberResponse struct {
	baseResponse
	Reference string `json:"reference"`
}

type sendMoneyResponse struct {
	baseResponse
	Success bool `json:"success"`
}

type getBalanceResponse struct {
	baseResponse
	Amount   uint   `json:"amount"`
	Currency string `json:"currency"`
}

type userInfo struct {
	baseResponse
	Member string `json:"member"`
	Wallet uint   `json:"wallet"`
}

type dumpAllUsersResponse struct {
	baseResponse
	DumpInfo []userInfo `json:"dump_info"`
}

func getResponseBody(t *testing.T, postParams map[string]interface{}) []byte {
	jsonValue, _ := json.Marshal(postParams)
	postResp, err := http.Post(TestUrl, "application/json", bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	assert.NoError(t, err)
	return body
}

func unmarshalResponse(t *testing.T, body []byte, response responseInterface) {
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Nil(t, response.getError())
}

func unmarshalResponseWithError(t *testing.T, body []byte, response responseInterface) {
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.getError())
}

func TestInsolardResponseNotErr(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "dump_all_users",
	})

	response := &dumpAllUsersResponse{}
	unmarshalResponse(t, body, response)

	assert.Nil(t, response.Err)
}

func TestTransferMoney(t *testing.T) {
	// Create member which balance will increase
	body := getResponseBody(t, postParams{
		"query_type": "create_member",
		"name":       "First",
	})

	firstMemberResponse := &createMemberResponse{}
	unmarshalResponse(t, body, firstMemberResponse)

	firstMemberRef := firstMemberResponse.Reference
	assert.NotEqual(t, "", firstMemberRef)

	// Create member which balance will decrease
	body = getResponseBody(t, postParams{
		"query_type": "create_member",
		"name":       "Second",
	})

	secondMemberResponse := &createMemberResponse{}
	unmarshalResponse(t, body, secondMemberResponse)

	secondMemberRef := secondMemberResponse.Reference
	assert.NotEqual(t, "", secondMemberRef)

	// Transfer money from one member to another
	body = getResponseBody(t, postParams{
		"query_type": "send_money",
		"from":       secondMemberRef,
		"to":         firstMemberRef,
		"amount":     111,
	})

	transferResponse := &sendMoneyResponse{}
	unmarshalResponse(t, body, transferResponse)

	assert.Equal(t, true, transferResponse.Success)

	// Check balance of first member
	body = getResponseBody(t, postParams{
		"query_type": "get_balance",
		"reference":  firstMemberRef,
	})

	firstBalanceResponse := &getBalanceResponse{}
	unmarshalResponse(t, body, firstBalanceResponse)

	assert.Equal(t, uint(1111), firstBalanceResponse.Amount)
	assert.Equal(t, "RUB", firstBalanceResponse.Currency)

	// Check balance of second member
	body = getResponseBody(t, postParams{
		"query_type": "get_balance",
		"reference":  secondMemberRef,
	})

	secondBalanceResponse := &getBalanceResponse{}
	unmarshalResponse(t, body, secondBalanceResponse)

	assert.Equal(t, uint(889), secondBalanceResponse.Amount)
	assert.Equal(t, "RUB", secondBalanceResponse.Currency)
}

func TestWrongUrl(t *testing.T) {
	jsonValue, _ := json.Marshal(postParams{
		"query_type": "dump_all_users",
	})
	testUrl := HOST + "/not_api/v1"
	postResp, err := http.Post(testUrl, "application/json", bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, postResp.StatusCode)
}

func TestGetRequest(t *testing.T) {
	postResp, err := http.Get(TestUrl)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	assert.NoError(t, err)

	getResponse := &baseResponse{}
	unmarshalResponseWithError(t, body, getResponse)

	assert.Equal(t, api.BadRequest, getResponse.Err.Code)
	assert.Equal(t, "Bad request", getResponse.Err.Event)
}

func TestWrongJson(t *testing.T) {
	postResp, err := http.Post(TestUrl, "application/json", bytes.NewBuffer([]byte("some not json value")))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	assert.NoError(t, err)

	response := &baseResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.BadRequest, response.Err.Code)
	assert.Equal(t, "Bad request", response.Err.Event)
}

func TestWrongQueryType(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "wrong_query_type",
	})

	response := &baseResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.BadRequest, response.Err.Code)
	assert.Equal(t, "Wrong query parameter 'query_type' = 'wrong_query_type'", response.Err.Event)
}

func TestWithoutQueryType(t *testing.T) {
	body := getResponseBody(t, postParams{})

	response := &baseResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.BadRequest, response.Err.Code)
	assert.Equal(t, "Wrong query parameter 'query_type' = ''", response.Err.Event)
}

func TestTooMuchParams(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "create_member",
		"some_param": "irrelevant info",
		"name":       testutils.RandomString(),
	})

	firstMemberResponse := &createMemberResponse{}
	unmarshalResponse(t, body, firstMemberResponse)

	firstMemberRef := firstMemberResponse.Reference
	assert.NotEqual(t, "", firstMemberRef)
}

func TestQueryTypeAsIntParams(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": 100,
	})

	response := &baseResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.BadRequest, response.Err.Code)
	assert.Equal(t, "Bad request", response.Err.Event)
}

func TestWrongTypeInParams(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "create_member",
		"name":       128182187,
	})

	response := &baseResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.BadRequest, response.Err.Code)
	assert.Equal(t, "Bad request", response.Err.Event)
}

// TODO: unskip test after doing errors in smart contracts
func _TestWrongReferenceInParams(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "get_balance",
		"reference":  testutils.RandomString(),
	})

	response := &baseResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.BadRequest, response.Err.Code)
	assert.Equal(t, "Bad request", response.Err.Event)
}
