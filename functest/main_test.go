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
const TestURL = HOST + "/api/v1"
const insolarImportPath = "github.com/insolar/insolar"

var cmd *exec.Cmd
var stdin io.WriteCloser
var stdout io.ReadCloser
var stderr io.ReadCloser
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
	out, err := exec.Command(
		"go", "build",
		"-o", insolardPath,
		insolarImportPath+"/cmd/insolard/",
	).CombinedOutput()
	return errors.Wrapf(err, "[ buildInsolard ] could't build insolard: %s", out)
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

var insgorundPath string

func buildGinsiderCLI() (err error) {
	insgorundPath, _, err = testutil.Build()
	return errors.Wrap(err, "[ buildGinsiderCLI ] could't build ginsider CLI: ")
}

func waitForLaunch() error {
	done := make(chan bool, 1)
	timeout := 40 * time.Second

	go func() {
		scanner := bufio.NewScanner(stdout)
		fmt.Println("Insolard output: ")
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			if strings.Contains(line, "======= Host info ======") {
				done <- true
			}
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
		}
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return errors.Errorf("[ waitForLaunch ] could't wait for launch: timeout of %s was exceeded", timeout)
	}

}

func startInsolard() error {
	cmd = exec.Command(
		insolardPath,
	)
	var err error

	stdin, err = cmd.StdinPipe()
	if err != nil {
		return errors.Wrap(err, "[ startInsolard ] could't set stdin: ")
	}

	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "[ startInsolard ] could't set stdout: ")
	}

	stderr, err = cmd.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "[ startInsolard ] could't set stderr: ")
	}

	err = cmd.Start()
	if err != nil {
		return errors.Wrap(err, "[ startInsolard ] could't start insolard command: ")
	}

	err = waitForLaunch()
	if err != nil {
		return errors.Wrap(err, "[ startInsolard ] could't wait for insolard to start completely: ")
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
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	io.WriteString(stdin, "exit\n")
	err := cmd.Wait()
	if err != nil {
		fmt.Println("[ stopInsolard ] try to kill, wait done with error: ", err)
		err := cmd.Process.Kill()
		if err != nil {
			return errors.Wrap(err, "[ stopInsolard ] failed to kill process: ")
		}
	}
	return nil
}

var insgorundCleaner func()

func startInsgorund() (err error) {
	insgorundCleaner, err = testutils.StartInsgorund(insgorundPath, "127.0.0.1:18181", "127.0.0.1:18182")
	if err != nil {
		return errors.Wrap(err, "[ startInsolard ] could't wait for insolard to start completely: ")
	}
	return err
}

func stopInsgorund() error {
	if insgorundCleaner != nil {
		insgorundCleaner()
	}
	return nil
}

func setup() error {
	err := createDirForContracts()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't create dirs for test: ")
	}
	fmt.Println("[ setup ] directory for contracts cache was successfully created")

	err = buildGinsiderCLI()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't build ginsider CLI: ")
	}
	fmt.Println("[ setup ] ginsider CLI was successfully builded")

	err = buildInsolard()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't build insolard: ")
	}
	fmt.Println("[ setup ] insolard was successfully builded")

	err = startInsgorund()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't start insgorund: ")
	}
	fmt.Println("[ setup ] insgorund was successfully started")

	err = startInsolard()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't start insolard: ")
	}
	fmt.Println("[ setup ] insolard was successfully started")

	return nil
}

func teardown() {
	err := stopInsolard()
	if err != nil {
		fmt.Println("[ teardown ] failed to stop insolard: ", err)
	}
	fmt.Println("[ teardown ] insolard was successfully stoped")

	err = stopInsgorund()
	if err != nil {
		fmt.Println("[ teardown ] failed to stop insgorund: ", err)
	}
	fmt.Println("[ teardown ] insgorund was successfully stoped")

	err = deleteDirForData()
	if err != nil {
		fmt.Println("[ teardown ] failed to remove data directory for func tests: ", err)
	}
	fmt.Println("[ teardown ] data directory was successfully deleted")

	err = deleteDirForContracts()
	if err != nil {
		fmt.Println("[ teardown ] failed to remove directory for contracts cache for func tests: ", err)
	}
	fmt.Println("[ teardown ] directory for contracts cache was successfully deleted")
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

func createMember(t *testing.T) string {
	body := getResponseBody(t, postParams{
		"query_type": "create_member",
		"name":       testutils.RandomString(),
	})

	firstMemberResponse := &createMemberResponse{}
	unmarshalResponse(t, body, firstMemberResponse)

	return firstMemberResponse.Reference
}

func getBalance(t *testing.T, reference string) int {
	body := getResponseBody(t, postParams{
		"query_type": "get_balance",
		"reference":  reference,
	})

	firstBalanceResponse := &getBalanceResponse{}
	unmarshalResponse(t, body, firstBalanceResponse)

	return int(firstBalanceResponse.Amount)
}

func getResponseBody(t *testing.T, postParams map[string]interface{}) []byte {
	jsonValue, _ := json.Marshal(postParams)
	postResp, err := http.Post(TestURL, "application/json", bytes.NewBuffer(jsonValue))
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
	firstMemberRef := createMember(t)
	secondMemberRef := createMember(t)
	oldFirstBalance := getBalance(t, firstMemberRef)
	oldSecondBalance := getBalance(t, secondMemberRef)

	amount := 111

	// Transfer money from one member to another
	body := getResponseBody(t, postParams{
		"query_type": "send_money",
		"from":       secondMemberRef,
		"to":         firstMemberRef,
		"amount":     amount,
	})

	transferResponse := &sendMoneyResponse{}
	unmarshalResponse(t, body, transferResponse)

	assert.Equal(t, true, transferResponse.Success)

	newFirstBalance := getBalance(t, firstMemberRef)
	newSecondBalance := getBalance(t, secondMemberRef)

	assert.Equal(t, oldFirstBalance+amount, newFirstBalance)
	assert.Equal(t, oldSecondBalance-amount, newSecondBalance)
}

func TestWrongUrl(t *testing.T) {
	jsonValue, _ := json.Marshal(postParams{
		"query_type": "dump_all_users",
	})
	testURL := HOST + "/not_api/v1"
	postResp, err := http.Post(testURL, "application/json", bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, postResp.StatusCode)
}

func TestGetRequest(t *testing.T) {
	postResp, err := http.Get(TestURL)
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
	postResp, err := http.Post(TestURL, "application/json", bytes.NewBuffer([]byte("some not json value")))
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

func TestTransferNegativeAmount(t *testing.T) {
	firstMemberRef := createMember(t)
	secondMemberRef := createMember(t)
	oldFirstBalance := getBalance(t, firstMemberRef)
	oldSecondBalance := getBalance(t, secondMemberRef)

	body := getResponseBody(t, postParams{
		"query_type": "send_money",
		"from":       secondMemberRef,
		"to":         firstMemberRef,
		"amount":     -111,
	})

	transferResponse := &sendMoneyResponse{}
	unmarshalResponseWithError(t, body, transferResponse)

	assert.Equal(t, api.BadRequest, transferResponse.Err.Code)
	assert.Equal(t, "Bad request", transferResponse.Err.Event)

	newFirstBalance := getBalance(t, firstMemberRef)
	newSecondBalance := getBalance(t, secondMemberRef)

	assert.Equal(t, oldFirstBalance, newFirstBalance)
	assert.Equal(t, oldSecondBalance, newSecondBalance)

}

func TestTransferAllAmount(t *testing.T) {
	firstMemberRef := createMember(t)
	secondMemberRef := createMember(t)
	oldFirstBalance := getBalance(t, firstMemberRef)
	oldSecondBalance := getBalance(t, secondMemberRef)

	body := getResponseBody(t, postParams{
		"query_type": "send_money",
		"from":       secondMemberRef,
		"to":         firstMemberRef,
		"amount":     oldSecondBalance,
	})

	transferResponse := &sendMoneyResponse{}
	unmarshalResponse(t, body, transferResponse)

	assert.Equal(t, true, transferResponse.Success)

	newFirstBalance := getBalance(t, firstMemberRef)
	newSecondBalance := getBalance(t, secondMemberRef)

	assert.Equal(t, oldFirstBalance+oldSecondBalance, newFirstBalance)
	assert.Equal(t, 0, newSecondBalance)

}

func _TestTransferMoreThanAvailableAmount(t *testing.T) {
	firstMemberRef := createMember(t)
	secondMemberRef := createMember(t)
	oldFirstBalance := getBalance(t, firstMemberRef)
	oldSecondBalance := getBalance(t, secondMemberRef)

	body := getResponseBody(t, postParams{
		"query_type": "send_money",
		"from":       secondMemberRef,
		"to":         firstMemberRef,
		"amount":     10000000000,
	})

	transferResponse := &sendMoneyResponse{}
	unmarshalResponse(t, body, transferResponse)

	// Add checking than contract gives specific error

	newFirstBalance := getBalance(t, firstMemberRef)
	newSecondBalance := getBalance(t, secondMemberRef)

	assert.Equal(t, oldFirstBalance, newFirstBalance)
	assert.Equal(t, oldSecondBalance, newSecondBalance)
}

func _TestTransferToMyself(t *testing.T) {
	memberRef := createMember(t)
	oldBalance := getBalance(t, memberRef)

	body := getResponseBody(t, postParams{
		"query_type": "send_money",
		"from":       memberRef,
		"to":         memberRef,
		"amount":     oldBalance - 1,
	})

	transferResponse := &sendMoneyResponse{}
	unmarshalResponse(t, body, transferResponse)

	assert.Equal(t, true, transferResponse.Success)

	newBalance := getBalance(t, memberRef)

	assert.Equal(t, oldBalance, newBalance)
}

// TODO: test to check overflow of balance
// TODO: check transfer zero amount

func TestTransferTwoTimes(t *testing.T) {
	firstMemberRef := createMember(t)
	secondMemberRef := createMember(t)
	oldFirstBalance := getBalance(t, firstMemberRef)
	oldSecondBalance := getBalance(t, secondMemberRef)

	firstBody := getResponseBody(t, postParams{
		"query_type": "send_money",
		"from":       secondMemberRef,
		"to":         firstMemberRef,
		"amount":     100,
	})
	transferResponse := &sendMoneyResponse{}
	unmarshalResponse(t, firstBody, transferResponse)
	assert.Equal(t, true, transferResponse.Success)

	secondBody := getResponseBody(t, postParams{
		"query_type": "send_money",
		"from":       secondMemberRef,
		"to":         firstMemberRef,
		"amount":     100,
	})
	unmarshalResponse(t, secondBody, transferResponse)
	assert.Equal(t, true, transferResponse.Success)

	newFirstBalance := getBalance(t, firstMemberRef)
	newSecondBalance := getBalance(t, secondMemberRef)

	assert.Equal(t, oldFirstBalance+200, newFirstBalance)
	assert.Equal(t, oldSecondBalance-200, newSecondBalance)
}

func TestCreateMembersWithSameName(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "create_member",
		"name":       "NameForTestCreateMembersWithSameName",
	})

	memberResponse := &createMemberResponse{}
	unmarshalResponse(t, body, memberResponse)

	firstMemberRef := memberResponse.Reference
	assert.NotEqual(t, "", firstMemberRef)

	body = getResponseBody(t, postParams{
		"query_type": "create_member",
		"name":       "NameForTestCreateMembersWithSameName",
	})

	unmarshalResponse(t, body, memberResponse)

	secondMemberRef := memberResponse.Reference
	assert.NotEqual(t, "", secondMemberRef)

	assert.NotEqual(t, firstMemberRef, secondMemberRef)
}
