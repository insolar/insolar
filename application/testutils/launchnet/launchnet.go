// Copyright 2020 Insolar Network Ltd.
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

package launchnet

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/application/api/requester"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/defaults"
)

const HOST = "http://localhost:"
const AdminPort = "19002"
const PublicPort = "19102"
const HostDebug = "http://localhost:8001"
const TestAdminRPCUrl = "/admin-api/rpc"

var AdminHostPort = HOST + AdminPort
var TestRPCUrl = HOST + AdminPort + TestAdminRPCUrl
var TestRPCUrlPublic = HOST + PublicPort + "/api/rpc"
var disableLaunchnet = false
var testRPCUrlVar = "INSOLAR_FUNC_RPC_URL"
var testRPCUrlPublicVar = "INSOLAR_FUNC_RPC_URL_PUBLIC"
var keysPathVar = "INSOLAR_FUNC_KEYS_PATH"

const insolarRootMemberKeys = "root_member_keys.json"
const insolarMigrationAdminMemberKeys = "migration_admin_member_keys.json"
const insolarFeeMemberKeys = "fee_member_keys.json"

var cmd *exec.Cmd
var cmdCompleted = make(chan error, 1)
var stdin io.WriteCloser
var stdout io.ReadCloser
var stderr io.ReadCloser

var ApplicationIncentives [application.GenesisAmountApplicationIncentivesMembers]*User
var NetworkIncentives [application.GenesisAmountNetworkIncentivesMembers]*User
var Enterprise [application.GenesisAmountEnterpriseMembers]*User
var Foundation [application.GenesisAmountFoundationMembers]*User
var Funds [application.GenesisAmountFundsMembers]*User

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

		os.Exit(2)
	}()

	pulseWatcher, config := pulseWatcherPath()

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
var FeeMember User
var MigrationDaemons [application.GenesisAmountMigrationDaemonMembers]*User

type User struct {
	Ref              string
	PrivKey          string
	PubKey           string
	MigrationAddress string
}

func launchnetPath(a ...string) (string, error) {
	keysPath := os.Getenv(keysPathVar)
	if keysPath != "" {
		p := []string{keysPath}
		p = append(p, a[len(a)-1])
		return filepath.Join(p...), nil
	}
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
		Nodes         []interface{} `yaml:"nodes"`
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

	return len(conf.DiscoverNodes) + len(conf.Nodes), nil
}

func GetNumShards() (int, error) {
	type bootstrapConf struct {
		PKShardCount int `yaml:"ma_shard_count"`
	}

	var conf bootstrapConf

	path, err := launchnetPath("bootstrap.yaml")
	if err != nil {
		return 0, err
	}
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, errors.Wrap(err, "[ GetNumShards ] Can't read bootstrap config")
	}

	err = yaml.Unmarshal(buff, &conf)
	if err != nil {
		return 0, errors.Wrap(err, "[ GetNumShards ] Can't parse bootstrap config")
	}

	return conf.PKShardCount, nil
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
	path, err = launchnetPath("configs", insolarFeeMemberKeys)
	if err != nil {
		return err
	}
	err = loadMemberKeys(path, &FeeMember)
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

	for i := 0; i < application.GenesisAmountApplicationIncentivesMembers; i++ {
		path, err := launchnetPath("configs", "application_incentives_"+strconv.Itoa(i)+"_member_keys.json")
		if err != nil {
			return err
		}
		var md User
		err = loadMemberKeys(path, &md)
		if err != nil {
			return err
		}
		ApplicationIncentives[i] = &md
	}

	for i := 0; i < application.GenesisAmountNetworkIncentivesMembers; i++ {
		path, err := launchnetPath("configs", "network_incentives_"+strconv.Itoa(i)+"_member_keys.json")
		if err != nil {
			return err
		}
		var md User
		err = loadMemberKeys(path, &md)
		if err != nil {
			return err
		}
		NetworkIncentives[i] = &md
	}

	for i := 0; i < application.GenesisAmountFoundationMembers; i++ {
		path, err := launchnetPath("configs", "foundation_"+strconv.Itoa(i)+"_member_keys.json")
		if err != nil {
			return err
		}
		var md User
		err = loadMemberKeys(path, &md)
		if err != nil {
			return err
		}
		Foundation[i] = &md
	}

	for i := 0; i < application.GenesisAmountFundsMembers; i++ {
		path, err := launchnetPath("configs", "funds_"+strconv.Itoa(i)+"_member_keys.json")
		if err != nil {
			return err
		}
		var md User
		err = loadMemberKeys(path, &md)
		if err != nil {
			return err
		}
		Funds[i] = &md
	}

	for i := 0; i < application.GenesisAmountEnterpriseMembers; i++ {
		path, err := launchnetPath("configs", "enterprise_"+strconv.Itoa(i)+"_member_keys.json")
		if err != nil {
			return err
		}
		var md User
		err = loadMemberKeys(path, &md)
		if err != nil {
			return err
		}
		Enterprise[i] = &md
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
	numAttempts := 270
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
	defer func() {
		_ = os.Chdir(cwd)
	}()

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

var logRotatorEnableVar = "LOGROTATOR_ENABLE"

// LogRotateEnabled checks is log rotation enabled by environment variable.
func LogRotateEnabled() bool {
	return os.Getenv(logRotatorEnableVar) == "1"
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

func RunOnlyWithLaunchnet(t *testing.T) {
	if disableLaunchnet {
		t.Skip()
	}
}

func setup() error {
	testRPCUrl := os.Getenv(testRPCUrlVar)
	testRPCUrlPublic := os.Getenv(testRPCUrlPublicVar)

	if testRPCUrl == "" || testRPCUrlPublic == "" {
		err := startNet()
		if err != nil {
			return errors.Wrap(err, "[ setup ] could't startNet")
		}
	} else {
		TestRPCUrl = testRPCUrl
		TestRPCUrlPublic = testRPCUrlPublic
		url := strings.Split(TestRPCUrlPublic, "/")
		AdminHostPort = strings.Join(url[0:len(url)-1], "/")
		disableLaunchnet = true
	}

	err := loadAllMembersKeys()
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

	return nil
}

func pulseWatcherPath() (string, string) {
	insDir := insolar.RootModuleDir()
	pulseWatcher := filepath.Join(insDir, "bin", "pulsewatcher")

	baseDir := defaults.PathWithBaseDir(defaults.LaunchnetDir(), insDir)
	config := filepath.Join(baseDir, "pulsewatcher.yaml")
	return pulseWatcher, config
}

func teardown() {
	err := stopInsolard()
	if err != nil {
		fmt.Println("[ teardown ]  failed to stop insolard:", err)
		return
	}
	fmt.Println("[ teardown ] insolard was successfully stopped")
}

// RotateLogs rotates launchnet logs, verbose flag enables printing what happens.
func RotateLogs(verbose bool) {
	launchnetDir := defaults.PathWithBaseDir(defaults.LaunchnetDir(), insolar.RootModuleDir())
	dirPattern := filepath.Join(launchnetDir, "logs/*/*/*.log")

	rmCmd := "rm -vf " + dirPattern
	cmd := exec.Command("sh", "-c", rmCmd)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal("RotateLogs: failed to execute shell command: ", rmCmd)
	}
	if verbose {
		fmt.Println("RotateLogs removed files:\n", string(out))
	}

	rotateCmd := "killall -v -SIGUSR2 inslogrotator"
	cmd = exec.Command("sh", "-c", rotateCmd)
	out, err = cmd.Output()
	if err != nil {
		if verbose {
			println("RotateLogs killall output:", string(out))
		}
		log.Fatal("RotateLogs: failed to execute shell command:", rotateCmd)
	}
}

var dumpMetricsEnabledVar = "DUMP_METRICS_ENABLE"

// LogRotateEnabled checks is log rotation enabled by environment variable.
func DumpMetricsEnabled() bool {
	return os.Getenv(dumpMetricsEnabledVar) == "1"
}

// FetchAndSaveMetrics fetches all nodes metric endpoints and saves result to files in
// logs/metrics/$iteration/<node-addr>.txt files.
func FetchAndSaveMetrics(iteration int) ([][]byte, error) {
	n, err := GetNodesCount()
	if err != nil {
		return nil, err
	}
	addrs := make([]string, n)
	for i := 0; i < n; i++ {
		addrs[i] = fmt.Sprintf(HOST+"80%02d", i+1)
	}
	results := make([][]byte, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for i, addr := range addrs {
		i := i
		addr := addr + "/metrics"
		go func() {
			defer wg.Done()

			r, err := fetchMetrics(addr)
			if err != nil {
				fetchErr := fmt.Sprintf("%v fetch failed: %v\n", addr, err.Error())
				results[i] = []byte(fetchErr)
				return
			}
			results[i] = r
		}()
	}
	wg.Wait()

	insDir := insolar.RootModuleDir()
	subDir := fmt.Sprintf("%04d", iteration)
	outDir := filepath.Join(insDir, defaults.LaunchnetDir(), "logs/metrics", subDir)
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "failed to create metrics subdirectory")
	}

	for i, b := range results {
		outFile := addrs[i][strings.Index(addrs[i], "://")+3:]
		outFile = strings.ReplaceAll(outFile, ":", "-")
		outFile = filepath.Join(outDir, outFile) + ".txt"

		err := ioutil.WriteFile(outFile, b, 0640)
		if err != nil {
			return nil, errors.Wrap(err, "write metrics failed")
		}
		fmt.Printf("Dump metrics from %v to %v\n", addrs[i], outFile)
	}
	return results, nil
}

func fetchMetrics(fetchURL string) ([]byte, error) {
	r, err := http.Get(fetchURL)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch metrics, got %v code", r.StatusCode)
	}
	return ioutil.ReadAll(r.Body)
}
