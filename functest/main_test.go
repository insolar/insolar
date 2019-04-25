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

// +build functest

package functest

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/defaults"
	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
	"github.com/pkg/errors"
)

const HOST = "http://localhost:19101"
const HOST_DEBUG = "http://localhost:8001"
const TestAPIURL = HOST + "/api"
const TestRPCUrl = TestAPIURL + "/rpc"
const TestCallUrl = TestAPIURL + "/call"

const insolarRootMemberKeys = "root_member_keys.json"

var cmd *exec.Cmd
var cmdCompleted = make(chan error, 1)
var stdin io.WriteCloser
var stdout io.ReadCloser
var stderr io.ReadCloser

var (
	insolarRootMemberKeysPath = launchnetPath("configs", insolarRootMemberKeys)
	insolarGenesisConfigPath  = launchnetPath("genesis.yaml")
)

func launchnetPath(a ...string) string {
	d := defaults.LaunchnetDir()
	parts := []string{"..", d}
	if strings.HasPrefix(d, "/") {
		parts = []string{d}
	}
	parts = append(parts, a...)
	return filepath.Join(parts...)
}

var info *requester.InfoResponse
var root user

type user struct {
	ref     string
	privKey string
	pubKey  string
}

func getNumberNodes() (int, error) {
	type genesisConf struct {
		DiscoverNodes []interface{} `yaml:"discovery_nodes"`
	}

	var conf genesisConf

	buff, err := ioutil.ReadFile(insolarGenesisConfigPath)
	if err != nil {
		return 0, errors.Wrap(err, "[ getNumberNodes ] Can't read genesis conf")
	}

	err = yaml.Unmarshal(buff, &conf)
	if err != nil {
		return 0, errors.Wrap(err, "[ getNumberNodes ] Can't parse genesis conf")
	}

	return len(conf.DiscoverNodes), nil
}

func functestPath() string {
	p, err := build.Default.Import("github.com/insolar/insolar", "", build.FindOnly)
	if err != nil {
		panic(err)
	}
	return filepath.Join(p.Dir, "functest")
}

func envVarWithDefault(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value != "" {
		return value
	}
	return defaultValue
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
	var err error
	info, err = requester.Info(TestAPIURL)
	if err != nil {
		return errors.Wrap(err, "[ setInfo ] error sending request")
	}
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
var secondInsgorundCleaner func()

func startInsgorund(listenPort string, upstreamPort string) (func(), error) {
	// It starts on ports of "virtual" node
	cleaner, err := goplugintestutils.StartInsgorund(insgorundPath, "tcp", "127.0.0.1:"+listenPort, "tcp", "127.0.0.1:"+upstreamPort)
	if err != nil {
		return cleaner, errors.Wrap(err, "[ startInsgorund ] couldn't wait for insolard to start completely: ")
	}
	return cleaner, nil
}

func startAllInsgorunds() (err error) {
	insgorundCleaner, err = startInsgorund("33305", "33306")
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't start insgorund: ")
	}
	fmt.Println("[ startAllInsgorunds ] insgorund was successfully started")

	secondInsgorundCleaner, err = startInsgorund("33327", "33328")
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't start second insgorund: ")
	}
	fmt.Println("[ startAllInsgorunds ] second insgorund was successfully started")

	return nil
}

func stopAllInsgorunds() error {
	if insgorundCleaner == nil || secondInsgorundCleaner == nil {
		return errors.New("[ stopInsgorund ] cleaner func not found")
	}
	insgorundCleaner()
	secondInsgorundCleaner()
	return nil
}

func waitForNet() error {
	numAttempts := 90
	ports := []string{"19101", "19102", "19103", "19104", "19105"}
	numNodes := len(ports)
	currentOk := 0
	for i := 0; i < numAttempts; i++ {
		currentOk = 0
		for _, port := range ports {
			resp, err := requester.Status(fmt.Sprintf("http://127.0.0.1:%s/api", port))
			if err != nil {
				fmt.Println("[ waitForNet ] Problem with port " + port + ". Err: " + err.Error())
				break
			}
			if resp.NetworkState != insolar.CompleteNetworkState.String() {
				fmt.Println("[ waitForNet ] Good response from port " + port + ". Net is not ready. Response: " + resp.NetworkState)
				break
			}
			fmt.Println("[ waitForNet ] Good response from port " + port + ". Net is ready. Response: " + resp.NetworkState)
			currentOk++
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

	cmd = exec.Command("./scripts/insolard/launchnet.sh", "-ngw")
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
	timeout := 240 * time.Second

	go func() {
		scanner := bufio.NewScanner(stdout)
		fmt.Println("Insolard output: ")
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			if strings.Contains(line, "start discovery nodes ...") {
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
	err := buildGinsiderCLI()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't build ginsider CLI: ")
	}
	fmt.Println("[ setup ] ginsider CLI was successfully builded")

	err = startAllInsgorunds()
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
			fmt.Printf("[ setup ] Couldn't setInfo. Attempt %d/%d. Err: %s\n", i, numAttempts, err)
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
	var envSetting = os.Getenv("TEST_ENV")
	var err error
	fmt.Println("TEST_ENV: ", envSetting)
	if envSetting != "CI" {
		err = stopInsolard()

		if err != nil {
			fmt.Println("[ teardown ]  failed to stop insolard: ", err)
		}
		fmt.Println("[ teardown ] insolard was successfully stoped")
	}

	err = stopAllInsgorunds()
	if err != nil {
		fmt.Println("[ teardown ]  failed to stop all insgrounds: ", err)
	}
	fmt.Println("[ teardown ] insgorund was successfully stoped")

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
