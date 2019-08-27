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

package launchnet

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/defaults"
	"github.com/pkg/errors"
)

const HOST = "http://localhost:"
const AdminPort = "19002"
const PublicPort = "19102"
const HostDebug = "http://localhost:8001"
const TestAdminRPCUrl = "/admin-api/rpc"
const TestRPCUrl = HOST + AdminPort + TestAdminRPCUrl
const TestRPCUrlPublic = HOST + PublicPort + "/api/rpc"

const insolarRootMemberKeys = "root_member_keys.json"
const insolarMigrationAdminMemberKeys = "migration_admin_member_keys.json"

var cmd *exec.Cmd
var cmdCompleted = make(chan error, 1)
var stdin io.WriteCloser
var stdout io.ReadCloser
var stderr io.ReadCloser

// Method starts launchnet before execution of callback function (cb) and stops launchnet after.
// Returns exit code as a result from calling callback function.
func Run(cb func() int) int {
	err := setup()
	defer teardown()
	if err != nil {
		fmt.Println("error while setup, skip tests: ", err)
		return 1
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	go func() {
		sig := <-c
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		teardown()
	}()

	pulseWatcher, config, err := pulseWatcherPath()
	if err != nil {
		fmt.Println("PulseWatcher not found: ", err)
		return 1
	}

	code := cb()

	if code != 0 {
		out, err := exec.Command(pulseWatcher, "-c", config, "-s").CombinedOutput()
		if err != nil {
			fmt.Println("PulseWatcher execution error: ", err)
			return 1
		}
		fmt.Println(string(out))
	}
	return code
}

var info *requester.InfoResponse
var Root User
var MigrationAdmin User
var MigrationDaemons [insolar.GenesisAmountMigrationDaemonMembers]*User

type User struct {
	Ref     string
	PrivKey string
	PubKey  string
}

func launchnetPath(a ...string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "[ startNet ] Can't get current working directory")
	}
	cwdList := strings.Split(cwd, "/")
	var count int
	for i := len(cwdList); i >= 0; i-- {
		if cwdList[i-1] == "insolar" && cwdList[i-2] == "insolar" {
			break
		}
		count++
	}
	var dirUp []string
	for i := 0; i < count; i++ {
		dirUp = append(dirUp, "..")
	}

	d := defaults.LaunchnetDir()
	parts := append(dirUp, d)
	if strings.HasPrefix(d, "/") {
		parts = []string{d}
	}
	parts = append(parts, a...)
	return filepath.Join(parts...), nil
}

func GetNodesCount() (int, error) {
	type nodesConf struct {
		DiscoverNodes []interface{} `yaml:"discovery_nodes"`
	}

	var conf nodesConf

	path, err := launchnetPath("bootstrap.yaml")
	if err != nil {
		return 0, err
	}
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, errors.Wrap(err, "[ getNumberNodes ] Can't read bootstrap config")
	}

	err = yaml.Unmarshal(buff, &conf)
	if err != nil {
		return 0, errors.Wrap(err, "[ getNumberNodes ] Can't parse bootstrap config")
	}

	return len(conf.DiscoverNodes), nil
}

func loadMemberKeys(keysPath string, member *User) error {
	text, err := ioutil.ReadFile(keysPath)
	if err != nil {
		return errors.Wrapf(err, "[ loadMemberKeys ] could't load member keys")
	}
	var data map[string]string
	err = json.Unmarshal(text, &data)
	if err != nil {
		return errors.Wrapf(err, "[ loadMemberKeys ] could't unmarshal member keys")
	}
	if data["private_key"] == "" || data["public_key"] == "" {
		return errors.New("[ loadMemberKeys ] could't find any keys")
	}
	member.PrivKey = data["private_key"]
	member.PubKey = data["public_key"]

	return nil
}

func loadAllMembersKeys() error {
	path, err := launchnetPath("configs", insolarRootMemberKeys)
	if err != nil {
		return err
	}
	err = loadMemberKeys(path, &Root)
	if err != nil {
		return err
	}
	path, err = launchnetPath("configs", insolarMigrationAdminMemberKeys)
	if err != nil {
		return err
	}
	err = loadMemberKeys(path, &MigrationAdmin)
	if err != nil {
		return err
	}
	for i := range MigrationDaemons {
		path, err := launchnetPath("configs", "migration_daemon_"+strconv.Itoa(i)+"_member_keys.json")
		if err != nil {
			return err
		}
		var md User
		err = loadMemberKeys(path, &md)
		if err != nil {
			return err
		}
		MigrationDaemons[i] = &md
	}

	return nil
}

func setInfo() error {
	var err error
	info, err = requester.Info(TestRPCUrl)
	if err != nil {
		return errors.Wrap(err, "[ setInfo ] error sending request")
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

func waitForNet() error {
	numAttempts := 90
	// TODO: read ports from bootstrap config
	ports := []string{
		"19001",
		"19002",
		"19003",
		"19004",
		"19005",
		// "19106",
		// "19107",
		// "19108",
		// "19109",
		// "19110",
		// "19111",
	}
	numNodes := len(ports)
	currentOk := 0
	for i := 0; i < numAttempts; i++ {
		currentOk = 0
		for _, port := range ports {
			resp, err := requester.Status(fmt.Sprintf("%s%s%s", HOST, port, TestAdminRPCUrl))
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

	for cwd[len(cwd)-15:] != "insolar/insolar" {
		err = os.Chdir("../")
		if err != nil {
			return errors.Wrap(err, "[ startNet  ] Can't change dir")
		}
		cwd, err = os.Getwd()
		if err != nil {
			return errors.Wrap(err, "[ startNet ] Can't get current working directory")
		}
	}

	// If you want to add -n flag here please make sure that insgorund will
	// be eventually started with --log-level=debug. Otherwise someone will spent
	// a lot of time trying to figure out why insgorund debug logs are missing
	// during execution of functests.
	cmd = exec.Command("./scripts/insolard/launchnet.sh", "-gw")
	stdout, _ = cmd.StdoutPipe()

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
	err := startNet()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't startNet")
	}

	err = loadAllMembersKeys()
	if err != nil {
		return errors.Wrap(err, "[ setup ] could't load keys: ")
	}
	fmt.Println("[ setup ] all keys successfully loaded")

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
		return errors.Wrap(err, "[ setup ] could't receive Root reference ")
	}

	fmt.Println("[ setup ] references successfully received")
	Root.Ref = info.RootMember
	MigrationAdmin.Ref = info.MigrationAdminMember

	//Contracts = make(map[string]*contractInfo)

	return nil
}

func pulseWatcherPath() (string, string, error) {
	p, err := build.Default.Import("github.com/insolar/insolar", "", build.FindOnly)
	if err != nil {
		return "", "", errors.Wrap(err, "Couldn't receive path to github.com/insolar/insolar")
	}
	pulseWatcher := filepath.Join(p.Dir, "bin", "pulsewatcher")
	config := filepath.Join(p.Dir, ".artifacts", "launchnet", "pulsewatcher.yaml")
	return pulseWatcher, config, nil
}

func teardown() {
	var err error

	err = stopInsolard()
	if err != nil {
		fmt.Println("[ teardown ]  failed to stop insolard: ", err)
	}
	fmt.Println("[ teardown ] insolard was successfully stoped")
}
