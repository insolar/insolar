// +build functest

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
	"syscall"
	"testing"
	"time"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
	"github.com/pkg/errors"
)

const HOST = "http://localhost:19191"
const TestAPIURL = HOST + "/api"
const TestRPCUrl = TestAPIURL + "/rpc"
const TestCallUrl = TestAPIURL + "/call"

const insolarRootMemberKeys = "root_member_keys.json"

var cmd *exec.Cmd
var cmdCompleted = make(chan error, 1)
var stdin io.WriteCloser
var stdout io.ReadCloser
var stderr io.ReadCloser

var insolarRootMemberKeysPath = filepath.Join("../scripts/insolard/configs", insolarRootMemberKeys)

var info infoResponse
var root user

type user struct {
	ref     string
	privKey string
	pubKey  string
}

func functestPath() string {
	p, err := build.Default.Import("github.com/insolar/insolar", "", build.FindOnly)
	if err != nil {
		panic(err)
	}
	return filepath.Join(p.Dir, "functest")
}

func createDirForContracts() error {
	return os.MkdirAll(filepath.Join(functestPath(), "contractstorage"), 0777)
}

func deleteDirForContracts() error {
	return os.RemoveAll(filepath.Join(functestPath(), "contractstorage"))
}

func loadRootKeys() error {
	text, err := ioutil.ReadFile(insolarRootMemberKeysPath)
	if err != nil {
		return errors.Wrapf(err, "[ loadRootKeys ] could't load root keys")
	}
	var data map[string]string
	err = json.Unmarshal(text, &data)
	if err != nil {
		return errors.Wrapf(err, "[ loadRootKeys ] could't unmarshal root keys")
	}
	if data["private_key"] == "" || data["public_key"] == "" {
		return errors.New("[ loadRootKeys ] could't find any keys")
	}
	root.privKey = data["private_key"]
	root.pubKey = data["public_key"]

	return nil
}

func setInfo() error {
	jsonValue, err := json.Marshal(postParams{
		"jsonrpc": "2.0",
		"method":  "info.Get",
		"id":      "",
	})
	if err != nil {
		return errors.Wrap(err, "[ setInfo ] couldn't marshal post params")
	}
	postResp, err := http.Post(TestRPCUrl, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return errors.Wrapf(err, "[ setInfo ] couldn't send request to %s", TestRPCUrl)
	}
	body, err := ioutil.ReadAll(postResp.Body)
	if err != nil {
		return errors.Wrapf(err, "[ setInfo ] couldn't read answer")
	}
	infoResp := &rpcInfoResponse{}
	err = json.Unmarshal(body, infoResp)
	if err != nil {
		return errors.Wrapf(err, "[ setInfo ] couldn't unmarshall answer")
	}
	info = infoResp.Result
	return nil
}

var insgorundPath string

func buildGinsiderCLI() (err error) {
	insgorundPath, _, err = goplugintestutils.Build()
	return errors.Wrap(err, "[ buildGinsiderCLI ] could't build ginsider CLI: ")
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

	err := cmd.Process.Signal(syscall.SIGHUP)
	if err != nil {
		return errors.Wrap(err, "[ stopInsolard ] failed to kill process:")
	}

	pState, err := cmd.Process.Wait()
	if err != nil {
		return errors.Wrap(err, "[ stopInsolard ] failed to wait process:")
	}

	fmt.Println("[ stopInsolard ] State: ", pState.String())

	return nil
}

var insgorundCleaner func()

func startInsgorund() (err error) {
	// It starts on ports of "virtual" node
	insgorundCleaner, err = goplugintestutils.StartInsgorund(insgorundPath, "tcp", "127.0.0.1:18181", "tcp", "127.0.0.1:18182")
	if err != nil {
		return errors.Wrap(err, "[ startInsgorund ] couldn't wait for insolard to start completely: ")
	}
	return nil
}

func stopInsgorund() error {
	if insgorundCleaner != nil {
		insgorundCleaner()
	}
	return nil
}

func waitForNet() error {
	numAttempts := 90
	ports := []string{"19191", "19192", "19193"}
	numNodes := len(ports)
	currentOk := 0
	for i := 0; i < numAttempts; i++ {
		currentOk = 0
		for _, port := range ports {
			resp, err := requester.Status(fmt.Sprintf("http://127.0.0.1:%s/api", port))
			if err != nil {
				fmt.Println("[ waitForNet ] Problem with port " + port + ". Err: " + err.Error())
				break
			} else {
				fmt.Println("[ waitForNet ] Good response from port " + port + ". Response: " + resp.NetworkState)
				currentOk++
			}
		}
		if currentOk == numNodes {
			fmt.Printf("[ waitForNet ] All %d nodes have started\n", numNodes)
			break
		}

		time.Sleep(time.Second)
		fmt.Printf("[ waitForNet ] Waiting for net: attempt %d/%d\n", i, numAttempts)
	}

	if currentOk != numNodes {
		return errors.New("[ waitForNet ] Can't Start net: No attempts left")
	}

	return nil
}

func startNet() error {
	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "[ startNet ] Can't get current working directory")
	}
	defer os.Chdir(cwd)

	err = os.Chdir("../")
	if err != nil {
		return errors.Wrap(err, "[ startNet  ] Can't change dir")
	}

	cmd = exec.Command("./scripts/insolard/launchnet.sh", "-ng")
	stdout, _ = cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "[ startNet ] could't set stdout: ")
	}

	stderr, err = cmd.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "[ startNet] could't set stderr: ")
	}

	err = cmd.Start()
	if err != nil {
		return errors.Wrap(err, "[ startNet ] Can't run cmd")
	}

	err = waitForLaunch()
	if err != nil {
		return errors.Wrap(err, "[ startNet ] couldn't waitForLaunch more")
	}

	err = waitForNet()
	if err != nil {
		return errors.Wrap(err, "[ startNet ] couldn't waitForNet more")
	}

	return nil

}

func waitForLaunch() error {
	done := make(chan bool, 1)
	timeout := 120 * time.Second

	go func() {
		scanner := bufio.NewScanner(stdout)
		fmt.Println("Insolard output: ")
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			if strings.Contains(line, "start nodes ...") {
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

	go func() { cmdCompleted <- cmd.Wait() }()
	select {
	case err := <-cmdCompleted:
		cmdCompleted <- nil
		return errors.New("[ waitForLaunch ] insolard finished unexpectedly: " + err.Error())
	case <-done:
		return nil
	case <-time.After(timeout):
		return errors.Errorf("[ waitForLaunch ] could't wait for launch: timeout of %s was exceeded", timeout)
	}
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

	err = startInsgorund()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't start insgorund: ")
	}
	fmt.Println("[ setup ] insgorund was successfully started")

	err = startNet()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't startNet")
	}

	err = loadRootKeys()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't load root keys: ")
	}
	fmt.Println("[ setup ] root keys successfully loaded")

	numAttempts := 60
	for i := 0; i < numAttempts; i++ {
		err = setInfo()
		if err != nil {
			fmt.Printf("[ setup ] Couldn't setInfo. Attempt %d/%d. Err: %s", i, numAttempts, err)
		} else {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't receive root reference ")
	}

	fmt.Println("[ setup ] root reference successfully received")
	root.ref = info.RootMember

	return nil
}

func teardown() {
	err := stopInsolard()
	if err != nil {
		fmt.Println("[ teardown ]  failed to stop insolard: ", err)
	}
	fmt.Println("[ teardown ] insolard was successfully stoped")

	err = stopInsgorund()
	if err != nil {
		fmt.Println("[ teardown ] failed to stop insgorund: ", err)
	}
	fmt.Println("[ teardown ] insgorund was successfully stoped")

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
