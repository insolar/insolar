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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/insolar/insolar/api/sdk"
	"github.com/insolar/insolar/insolar/defaults"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/testutils"
)

const (
	defaultStdoutPath   = "-"
	createMemberRetries = 5
	balanceCheckRetries = 10
	balanceCheckDelay   = 5 * time.Second
)

var (
	defaultMemberFile         = filepath.Join(defaults.ArtifactsDir(), "bench-members", "members.txt")
	defaultDiscoveryNodesLogs = defaults.LaunchnetDiscoveryNodesLogsDir()

	memberFile         string
	output             string
	concurrent         int
	repetitions        int
	memberKeys         string
	adminAPIURLs       []string
	publicAPIURLs      []string
	logLevel           string
	logLevelServer     string
	saveMembersToFile  bool
	useMembersFromFile bool
	noCheckBalance     bool
	scenarioName       string
	discoveryNodesLogs string
)

func parseInputParams() {
	pflag.StringVarP(&output, "output", "o", defaultStdoutPath, "output file (use - for STDOUT)")
	pflag.IntVarP(&concurrent, "concurrent", "c", 1, "concurrent users")
	pflag.IntVarP(&repetitions, "repetitions", "r", 1, "repetitions for one user")
	pflag.StringVarP(&memberKeys, "memberkeys", "k", "", "path to dir with members keys")
	pflag.StringArrayVarP(&adminAPIURLs, "adminurls", "a", []string{"http://localhost:19001/admin-api/rpc"}, "url to admin api")
	pflag.StringArrayVarP(&publicAPIURLs, "publicurls", "p", []string{"http://localhost:19101/api/rpc"}, "url to public api")
	pflag.StringVarP(&logLevel, "loglevel", "l", "info", "log level for benchmark")
	pflag.StringVarP(&logLevelServer, "loglevelserver", "L", "", "server log level")
	pflag.BoolVarP(&saveMembersToFile, "savemembers", "s", false, "save members to file")
	pflag.BoolVarP(&useMembersFromFile, "usemembers", "m", false, "use members from file")
	pflag.StringVarP(&memberFile, "members-file", "", defaultMemberFile, "dir for saving members data")
	pflag.BoolVarP(&noCheckBalance, "nocheckbalance", "b", false, "don't check balance at the end")
	pflag.StringVarP(&scenarioName, "scenarioname", "t", "", "name of scenario")
	pflag.StringVarP(&discoveryNodesLogs, "discovery-nodes-logs-dir", "", defaultDiscoveryNodesLogs, "launchnet logs dir for checking errors")
	pflag.Parse()
}

func chooseOutput(path string) (io.Writer, error) {
	var res io.Writer
	if path == defaultStdoutPath {
		res = os.Stdout
	} else {
		var err error
		res, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't open file for writing")
		}
	}
	return res, nil
}

func writeToOutput(out io.Writer, data string) {
	_, err := out.Write([]byte(data))
	check("Can't write data to output", err)
}

func check(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

func newTransferDifferentMemberScenarios(out io.Writer, insSDK *sdk.SDK, members []*sdk.Member, concurrent int, repetitions int) benchmark {
	return benchmark{
		scenario: &transferDifferentMembersScenario{
			insSDK:  insSDK,
			members: members,
		},
		concurrent:  concurrent,
		repetitions: repetitions,
		name:        "TransferDifferentMembers",
		out:         out,
	}
}

func newCreateMemberScenarios(out io.Writer, insSDK *sdk.SDK, members []*sdk.Member, concurrent int, repetitions int) benchmark {
	return benchmark{
		scenario: &createMembersScenario{
			insSDK:  insSDK,
			members: members,
		},
		concurrent:  concurrent,
		repetitions: repetitions,
		name:        "CreateMember",
		out:         out,
	}
}

func startScenario(ctx context.Context, b benchmark) {
	err := b.scenario.canBeStarted()
	check(fmt.Sprintf("Scenario %s can not be started:", b.getName()), err)

	writeToOutput(b.getOut(), fmt.Sprintf("Scenario %s started: \n", b.getName()))

	start := time.Now()
	logReaderCloseChan := testutils.NodesErrorLogReader(discoveryNodesLogs, b.getOut())

	b.start(ctx)
	elapsed := time.Since(start)
	writeToOutput(b.getOut(), fmt.Sprintf("Scenario %s took: %s \n", b.getName(), elapsed))

	close(logReaderCloseChan)
	printResults(b)
}

func printResults(b benchmark) {
	speed := b.getOperationPerSecond()
	writeToOutput(b.getOut(), fmt.Sprintf("Scenario %s: Speed - %f resp/s \n", b.getName(), speed))
	writeToOutput(
		b.getOut(),
		fmt.Sprintf(
			"Scenario %s: Average Request Duration - %s\n",
			b.getName(), b.getAverageOperationDuration(),
		),
	)
	b.printResult()
}

func createMembers(insSDK *sdk.SDK, count int) []*sdk.Member {
	var members []*sdk.Member
	var member *sdk.Member
	var traceID string
	var err error

	for i := 0; i < count; i++ {
		retries := createMemberRetries
		for retries > 0 {
			member, traceID, err = insSDK.CreateMember("")
			if err == nil {
				members = append(members, member)
				break
			}
			fmt.Printf("Retry to create member. TraceID: %s Error is: %s\n", traceID, err.Error())
			retries--
		}
		check(fmt.Sprintf("Couldn't create member after retries: %d", createMemberRetries), err)
	}
	return members
}

func getTotalBalance(insSDK *sdk.SDK, members []*sdk.Member) (totalBalance *big.Int) {
	type Result struct {
		num     int
		balance *big.Int
		err     error
	}

	nmembers := len(members)
	var wg sync.WaitGroup
	wg.Add(nmembers)
	results := make(chan Result, nmembers)

	// execute all queries in parallel
	for i := 0; i < nmembers; i++ {
		go func(m *sdk.Member, num int) {
			res := Result{num: num}
			res.balance, res.err = insSDK.GetBalance(m)
			results <- res
			wg.Done()
		}(members[i], i)
	}

	wg.Wait()
	totalBalance = big.NewInt(0)
	for i := 0; i < nmembers; i++ {
		res := <-results
		if res.err != nil {
			fmt.Printf("Can't get balance for %v-th member: %v\n", res.num, res.err)
			continue
		}
		b := totalBalance
		totalBalance.Add(b, res.balance)
	}

	return totalBalance
}

func getMembers(insSDK *sdk.SDK) ([]*sdk.Member, error) {
	var members []*sdk.Member
	var err error

	if useMembersFromFile {
		members, err = loadMembers(concurrent * 2)
		if err != nil {
			return nil, errors.Wrap(err, "error while loading members: ")
		}
	} else {
		start := time.Now()
		members = createMembers(insSDK, concurrent*2)
		creationTime := time.Since(start)
		fmt.Printf("Members were created in %s\n", creationTime)
		fmt.Printf("Average creation of member time - %s\n", time.Duration(int64(creationTime)/int64(concurrent*2)))
	}

	if saveMembersToFile {
		err = saveMembers(members)
		if err != nil {
			return nil, errors.Wrap(err, "save member done with error: ")
		}
	}
	return members, nil
}

func saveMembers(members []*sdk.Member) error {
	dir, _ := path.Split(memberFile)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return errors.Wrap(err, "couldn't create dir for file")
	}
	file, err := os.Create(memberFile)
	if err != nil {
		return errors.Wrap(err, "couldn't create file")
	}
	defer file.Close() // nolint: errcheck

	result, err := json.MarshalIndent(members, "", "    ")
	if err != nil {
		return errors.Wrap(err, "couldn't marshal members in json")
	}
	_, err = file.Write(result)
	return errors.Wrap(err, "couldn't save members in file")
}

func loadMembers(count int) ([]*sdk.Member, error) {
	var members []*sdk.Member

	rawMembers, err := ioutil.ReadFile(memberFile)
	if err != nil {
		return nil, errors.Wrap(err, "can't read members from file")
	}

	err = json.Unmarshal(rawMembers, &members)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal members from file")
	}

	if count > len(members) {
		return nil, errors.Errorf("Not enough members in file: got %d, needs %d", len(members), count)
	}
	return members, nil
}

func main() {
	parseInputParams()

	// Start benchmark time
	t := time.Now()
	fmt.Printf("Start: %s\n\n", t.String())

	err := log.SetLevel(logLevel)
	check(fmt.Sprintf("Can't set '%s' level on logger:", logLevel), err)

	out, err := chooseOutput(output)
	check("Problems with output file:", err)

	insSDK, err := sdk.NewSDK(adminAPIURLs, publicAPIURLs, memberKeys)
	check("SDK is not initialized: ", err)

	err = insSDK.SetLogLevel(logLevelServer)
	check("Failed to parse log level: ", err)

	members, err := getMembers(insSDK)
	check("Error while loading members: ", err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP)

	var b benchmark
	switch scenarioName {
	case "createMember":
		b = newCreateMemberScenarios(out, insSDK, members, concurrent, repetitions)
	default:
		b = newTransferDifferentMemberScenarios(out, insSDK, members, concurrent, repetitions)
	}

	go func() {
		stopGracefully := true
		for {
			sig := <-sigChan

			switch sig {
			case syscall.SIGHUP:
				printResults(b)
			case syscall.SIGINT:
				if !stopGracefully {
					log.Fatal("Force quiting.")
				} else {
					log.Info("Gracefully finishing benchmark. Press Ctrl+C again to force quit.")
				}

				stopGracefully = false
				cancel()
			}
		}
	}()

	b.scenario.prepare()

	startScenario(ctx, b)

	// Finish benchmark time
	t = time.Now()
	fmt.Printf("\nFinish: %s\n\n", t.String())

	b.scenario.checkResult()
}
