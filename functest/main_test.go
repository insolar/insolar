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

	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
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

func waitForLaunch(stdout io.ReadCloser) {
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if strings.Contains(line, "======= Host info ======") {
			break
		}
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
	waitForLaunch(stdout)
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

type respAPI struct {
	Qid       string                 `json:"qid"`
	Reference string                 `json:"reference"`
	Success   bool                   `json:"success"`
	Amount    uint                   `json:"amount"`
	Currency  string                 `json:"currency"`
	Err       map[string]interface{} `json:"error"`
}

func TestInsolardResponseNotErr(t *testing.T) {
	var resp respAPI
	postParams := map[string]interface{}{
		"query_type": "dump_all_users",
	}
	jsonValue, _ := json.Marshal(postParams)
	postResp, err := http.Post(TestUrl, "application/json", bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal([]byte(body), &resp)
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}(nil), resp.Err)
}

func TestTransferMoney(t *testing.T) {
	var resp respAPI
	// Create member which balance will increase
	postParams := map[string]interface{}{
		"query_type": "create_member",
		"name":       "First",
	}
	jsonValue, _ := json.Marshal(postParams)
	postResp, err := http.Post(TestUrl, "application/json", bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal([]byte(body), &resp)
	assert.NoError(t, err)
	assert.NotEqual(t, "", resp.Reference)
	firstMemberRef := resp.Reference

	// Create member which balance will decrease
	postParams = map[string]interface{}{
		"query_type": "create_member",
		"name":       "Second",
	}
	jsonValue, _ = json.Marshal(postParams)
	postResp, err = http.Post(TestUrl, "application/json", bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err = ioutil.ReadAll(postResp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal([]byte(body), &resp)
	assert.NoError(t, err)
	assert.NotEqual(t, "", resp.Reference)
	secondMemberRef := resp.Reference

	// Transfer money from one member to another
	postParams = map[string]interface{}{
		"query_type": "send_money",
		"from":       secondMemberRef,
		"to":         firstMemberRef,
		"amount":     111,
	}
	jsonValue, _ = json.Marshal(postParams)
	postResp, err = http.Post(TestUrl, "application/json", bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err = ioutil.ReadAll(postResp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal([]byte(body), &resp)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Success)

	// Check balance of first member
	postParams = map[string]interface{}{
		"query_type": "get_balance",
		"reference":  firstMemberRef,
	}
	jsonValue, _ = json.Marshal(postParams)
	postResp, err = http.Post(TestUrl, "application/json", bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err = ioutil.ReadAll(postResp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal([]byte(body), &resp)
	assert.NoError(t, err)
	assert.Equal(t, uint(1111), resp.Amount)
	assert.Equal(t, "RUB", resp.Currency)

	// Check balance of second member
	postParams = map[string]interface{}{
		"query_type": "get_balance",
		"reference":  secondMemberRef,
	}
	jsonValue, _ = json.Marshal(postParams)
	postResp, err = http.Post(TestUrl, "application/json", bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err = ioutil.ReadAll(postResp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal([]byte(body), &resp)
	assert.NoError(t, err)
	assert.Equal(t, uint(889), resp.Amount)
	assert.Equal(t, "RUB", resp.Currency)
}
