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
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
	"github.com/pkg/errors"
)

const HOST = "http://localhost:19191"
const TestURL = HOST + "/api/v1"
const insolarImportPath = "github.com/insolar/insolar"

//const insolarNodeKeys = "bootstrap_keys.json"
const insolarRootMemberKeys = "root_member_keys.json"

//const insolarCertificate = "certificate.json"

var cmd *exec.Cmd
var cmdCompleted = make(chan error, 1)
var stdin io.WriteCloser
var stdout io.ReadCloser
var stderr io.ReadCloser
var insolarPath = filepath.Join(testdataPath(), "insolar")
var insolardPath = filepath.Join(testdataPath(), "insolard")

//var insolarNodeKeysPath = filepath.Join(testdataPath(), insolarNodeKeys)
var insolarRootMemberKeysPath = filepath.Join(testdataPath(), insolarRootMemberKeys)
var insolarNodesKeysPath = filepath.Join(testdataPath(), "discovery_node_")

//var insolarCertificatePath = filepath.Join(testdataPath(), insolarCertificate)

var info infoResponse
var root user

type user struct {
	ref     string
	privKey string
	pubKey  string
}

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

func buildInsolar() error {
	out, err := exec.Command(
		"go", "build",
		"-o", insolarPath,
		insolarImportPath+"/cmd/insolar/",
	).CombinedOutput()
	return errors.Wrapf(err, "[ buildInsolar ] could't build insolar: %s", out)
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

/*func generateNodeKeys() error {
	out, err := exec.Command(
		insolarPath, "-c", "gen_keys",
		"-o", insolarNodeKeysPath).CombinedOutput()
	return errors.Wrapf(err, "[ generateNodeKeys ] could't generate node keys: %s", out)
}*/

func generateRootMemberKeys() error {
	out, err := exec.Command(
		insolarPath, "-c", "gen_keys",
		"-o", insolarRootMemberKeysPath).CombinedOutput()
	return errors.Wrapf(err, "[ generateRootMemberKeys ] could't generate root member keys: %s", out)
}

func generateDiscoveryNodesKeys() error {
	for i := 0; i < 5; i++ {
		out, err := exec.Command(
			insolarPath, "-c", "gen_keys",
			"-o", insolarNodesKeysPath+strconv.Itoa(i+1)+".json").CombinedOutput()
		if err != nil {
			return errors.Wrapf(err, "[ generateDiscoveryNodesKeys ] could't generate discovery node keys: %s", out)
		}
	}
	return nil
}

/*func generateCertificate() error {
	out, err := exec.Command(
		insolarPath, "-c", "gen_certificate", "-g", insolarNodeKeysPath,
		"-o", insolarCertificatePath).CombinedOutput()
	return errors.Wrapf(err, "[ generateCertificate ] could't generate certificate: %s", out)
}*/

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
	resp, err := http.Get(TestURL + "/info")
	if err != nil {
		return errors.Wrapf(err, "[ setInfo ] couldn't request %s", TestURL+"/info")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "[ setInfo ] couldn't read answer")
	}
	err = json.Unmarshal(body, &info)
	if err != nil {
		return errors.Wrapf(err, "[ setInfo ] couldn't unmarshall answer")
	}
	return nil
}

var insgorundPath string

func buildGinsiderCLI() (err error) {
	insgorundPath, _, err = goplugintestutils.Build()
	return errors.Wrap(err, "[ buildGinsiderCLI ] could't build ginsider CLI: ")
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

func startInsolard() error {
	cmd = exec.Command(
		insolardPath, "--genesis", filepath.Join(functestPath(), "genesis.yaml"),
		"--keyout", testdataPath(),
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
	err := <-cmdCompleted
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
	insgorundCleaner, err = goplugintestutils.StartInsgorund(insgorundPath, "tcp", "127.0.0.1:18181", "tcp", "127.0.0.1:18182")
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

	err := deleteDirForData()
	if err != nil {
		fmt.Println("[ setup ] failed to remove data directory for func tests: ", err)
	}

	err = createDirForContracts()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't create dirs for test: ")
	}
	fmt.Println("[ setup ] directory for contracts cache was successfully created")

	err = buildGinsiderCLI()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't build ginsider CLI: ")
	}
	fmt.Println("[ setup ] ginsider CLI was successfully builded")

	err = buildInsolar()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't build insolar: ")
	}
	fmt.Println("[ setup ] insolar was successfully builded")

	/*err = generateNodeKeys()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't generate node keys: ")
	}
	fmt.Println("[ setup ] node keys successfully generated")

	err = generateCertificate()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't generate certificate: ")
	}
	fmt.Println("[ setup ] certificate successfully generated")*/

	err = generateRootMemberKeys()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't generate root member keys: ")
	}
	fmt.Println("[ setup ] root member keys successfully generated")

	err = loadRootKeys()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't load root keys: ")
	}
	fmt.Println("[ setup ] root keys successfully loaded")

	err = generateDiscoveryNodesKeys()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't generate discovery node keys: ")
	}
	fmt.Println("[ setup ] discovery nodes keys successfully generated")

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
	err = setInfo()
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
