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
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/insolar/insolar/api/sdk"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/backoff"
	"github.com/insolar/insolar/insolar/defaults"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

const defaultStdoutPath = "-"
const defaultMemberFileName = "members.txt"

const backoffAttemptsCount = 20

var (
	defaultMemberFileDir      = filepath.Join(defaults.ArtifactsDir(), "bench-members")
	defaultDiscoveryNodesLogs = filepath.Join(defaults.LaunchnetDir(), "logs", "discoverynodes")

	memberFilesDir     string
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
	pflag.StringVarP(&memberFilesDir, "members-dir", "", defaultMemberFileDir, "dir for saving memebers data")
	pflag.BoolVarP(&noCheckBalance, "nocheckbalance", "b", false, "don't check balance at the end")
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

func newScenarios(out io.Writer, insSDK *sdk.SDK, members []*sdk.Member, concurrent int, repetitions int, penRetries int32) scenario {
	return &transferDifferentMembersScenario{
		concurrent:  concurrent,
		repetitions: repetitions,
		name:        "TransferDifferentMembers",
		out:         out,
		members:     members,
		insSDK:      insSDK,
		penRetries:  penRetries,
	}
}

func startScenario(ctx context.Context, s scenario) {
	err := s.canBeStarted()
	check(fmt.Sprintf("Scenario %s can not be started:", s.getName()), err)

	writeToOutput(s.getOut(), fmt.Sprintf("Scenario %s: Start to transfer\n", s.getName()))

	start := time.Now()
	logReaderCloseChan := nodesErrorLogReader(s)

	s.start(ctx)
	elapsed := time.Since(start)
	writeToOutput(s.getOut(), fmt.Sprintf("Scenario %s: Transferring took %s \n", s.getName(), elapsed))

	close(logReaderCloseChan)
	printResults(s)
}

func printResults(s scenario) {
	speed := s.getOperationPerSecond()
	writeToOutput(s.getOut(), fmt.Sprintf("Scenario %s: Speed - %f resp/s \n", s.getName(), speed))
	writeToOutput(
		s.getOut(),
		fmt.Sprintf(
			"Scenario %s: Average Request Duration - %s\n",
			s.getName(), s.getAverageOperationDuration(),
		),
	)
	s.printResult()
}

func nodesErrorLogReader(s scenario) chan struct{} {
	closeChan := make(chan struct{})
	wg := sync.WaitGroup{}

	logs, err := getLogs(discoveryNodesLogs)
	if err != nil {
		writeToOutput(s.getOut(), fmt.Sprintf("Can't find node logs: %s", err))
	}

	wg.Add(len(logs))
	for _, fileName := range logs {
		fName := fileName // be careful using loops and values in parallel code
		go readLogs(s, &wg, fName, closeChan)
	}

	wg.Wait()
	return closeChan
}

func getLogs(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == "output.log" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func readLogs(s scenario, wg *sync.WaitGroup, fileName string, closeChan chan struct{}) {
	defer wg.Done()

	file, err := os.Open(fileName)
	if err != nil {
		writeToOutput(s.getOut(), fmt.Sprintln("Can't open log file ", fileName, ", error : ", err))
	}
	_, err = file.Seek(-1, io.SeekEnd)
	if err != nil {
		writeToOutput(s.getOut(), fmt.Sprintln("Can't seek through log file ", fileName, ", error : ", err))
	}

	// for making wg.Done()
	go findErrorsInLog(s, fileName, file, closeChan)
}

func findErrorsInLog(s scenario, fName string, file io.ReadCloser, closeChan chan struct{}) {
	defer file.Close()
	reader := bufio.NewReader(file)

	ok := true
	for ok {
		select {
		case <-time.After(time.Millisecond):
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				writeToOutput(s.getOut(), fmt.Sprintln("Can't read string from ", fName, ", error: ", err))
				ok = false
			}

			if strings.Contains(line, " ERR ") {
				writeToOutput(s.getOut(), fmt.Sprintln("!!! THERE ARE ERRORS IN ERROR LOG !!! ", fName))
				ok = false
			}
		case <-closeChan:
			ok = false
		}

	}
}

func addMigrationAddresses(insSDK *sdk.SDK) int32 {
	var err error
	var retriesCount int32

	bof := backoff.Backoff{Min: 1 * time.Second, Max: 10 * time.Second}
	for bof.Attempt() < backoffAttemptsCount {
		migrationAddresses := []string{}
		for j := 0; j < concurrent*2; j++ {
			migrationAddresses = append(migrationAddresses, "fake_burn_address_"+strconv.Itoa(j))
		}
		traceID, err := insSDK.AddMigrationAddresses(migrationAddresses)
		if err == nil {
			break
		}

		if strings.Contains(err.Error(), insolar.ErrTooManyPendingRequests.Error()) {
			retriesCount++
		} else {
			fmt.Printf("Retry to add burn address. TraceID: %s Error is: %s\n", traceID, err.Error())
		}
		time.Sleep(bof.Duration())
	}
	check(fmt.Sprintf("Couldn't add burn address after retries: %d", backoffAttemptsCount), err)
	bof.Reset()

	return retriesCount
}

func createMembers(insSDK *sdk.SDK, count int) ([]*sdk.Member, int32) {
	var members []*sdk.Member
	var member *sdk.Member
	var traceID string
	var err error
	var retriesCount int32

	for i := 0; i < count; i++ {
		bof := backoff.Backoff{Min: 1 * time.Second, Max: 10 * time.Second}
		for bof.Attempt() < backoffAttemptsCount {
			member, traceID, err = insSDK.CreateMember()
			if err == nil {
				members = append(members, member)
				break
			}

			if strings.Contains(err.Error(), insolar.ErrTooManyPendingRequests.Error()) {
				retriesCount++
			} else {
				fmt.Printf("Retry to create member. TraceID: %s Error is: %s\n", traceID, err.Error())
			}
			time.Sleep(bof.Duration())
		}
		check(fmt.Sprintf("Couldn't create member after retries: %d", backoffAttemptsCount), err)
		bof.Reset()
	}
	return members, retriesCount
}

func getTotalBalance(insSDK *sdk.SDK, members []*sdk.Member) (totalBalance *big.Int, penRetires int32) {
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
			bof := backoff.Backoff{Min: 1 * time.Second, Max: 10 * time.Second}

			res := Result{num: num}
			for bof.Attempt() < backoffAttemptsCount {
				res.balance, res.err = insSDK.GetBalance(m)
				if res.err == nil {
					break
				}
				if strings.Contains(res.err.Error(), insolar.ErrTooManyPendingRequests.Error()) {
					atomic.AddInt32(&penRetires, 1)
				} else {
					// retry
					fmt.Printf("Retry to fetch balance for %v-th member: %v\n", res.num, res.err)
				}
				time.Sleep(bof.Duration())
			}
			results <- res
			wg.Done()
		}(members[i], i)
	}

	wg.Wait()
	totalBalance = big.NewInt(0)
	for i := 0; i < nmembers; i++ {
		res := <-results
		if res.err != nil {
			if !strings.Contains(res.err.Error(), insolar.ErrTooManyPendingRequests.Error()) {
				fmt.Printf("Can't get balance for %v-th member: %v\n", res.num, res.err)
			}
			continue
		}
		b := totalBalance
		totalBalance.Add(b, res.balance)
	}

	return totalBalance, penRetires
}

func getMembers(insSDK *sdk.SDK) ([]*sdk.Member, int32, error) {
	var members []*sdk.Member
	var err error
	var retriesCount int32

	if useMembersFromFile {
		members, err = loadMembers(concurrent * 2)
		if err != nil {
			return nil, 0, errors.Wrap(err, "error while loading members: ")
		}
	} else {
		start := time.Now()
		members, retriesCount = createMembers(insSDK, concurrent*2)
		creationTime := time.Since(start)
		fmt.Printf("Members were created in %s\n", creationTime)
		fmt.Printf("Average creation of member time - %s\n", time.Duration(int64(creationTime)/int64(concurrent*2)))
	}

	if saveMembersToFile {
		err = saveMembers(members)
		if err != nil {
			return nil, 0, errors.Wrap(err, "save member done with error: ")
		}
	}
	return members, retriesCount, nil
}

func saveMembers(members []*sdk.Member) error {
	err := os.MkdirAll(defaultMemberFileDir, 0777)
	if err != nil {
		return errors.Wrap(err, "couldn't create dir for file")
	}
	file, err := os.Create(filepath.Join(defaultMemberFileDir, defaultMemberFileName))
	if err != nil {
		return errors.Wrap(err, "couldn't create file")
	}
	defer file.Close() //nolint: errcheck

	result, err := json.MarshalIndent(members, "", "    ")
	if err != nil {
		return errors.Wrap(err, "couldn't marshal members in json")
	}
	_, err = file.Write(result)
	return errors.Wrap(err, "couldn't save members in file")
}

func loadMembers(count int) ([]*sdk.Member, error) {
	var members []*sdk.Member

	rawMembers, err := ioutil.ReadFile(filepath.Join(defaultMemberFileDir, defaultMemberFileName))
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

func calcFee(amount int64) int64 {
	return transferFee
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

	crBaPenBefore := addMigrationAddresses(insSDK)
	check("Error while adding burn addresses: ", err)

	members, crMemPenBefore, err := getMembers(insSDK)
	check("Error while loading members: ", err)

	var totalBalanceBefore *big.Int
	var balancePenRetries int32
	if !noCheckBalance {
		totalBalanceBefore, balancePenRetries = getTotalBalance(insSDK, members)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP)

	s := newScenarios(out, insSDK, members, concurrent, repetitions, crBaPenBefore+crMemPenBefore+balancePenRetries)
	go func() {
		stopGracefully := true
		for {
			sig := <-sigChan

			switch sig {
			case syscall.SIGHUP:
				printResults(s)
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

	startScenario(ctx, s)

	// Finish benchmark time
	t = time.Now()
	fmt.Printf("\nFinish: %s\n\n", t.String())

	if !noCheckBalance {
		totalBalanceAfter := big.NewInt(0)
		totalBalanceAfterWithFee := big.NewInt(0)
		for nretries := 0; nretries < 3; nretries++ {
			totalBalanceAfter, _ = getTotalBalance(insSDK, members)
			totalBalanceAfterWithFee = new(big.Int).Add(totalBalanceAfter, big.NewInt(calcFee(transferAmount)*int64(repetitions*concurrent)))
			if totalBalanceAfterWithFee.Cmp(totalBalanceBefore) == 0 {
				break
			}
			fmt.Printf("Total balance before and after don't match: %v vs %v - retrying in 3 seconds...\n",
				totalBalanceBefore, totalBalanceAfterWithFee)
			time.Sleep(3 * time.Second)

		}
		fmt.Printf("Total balance before: %v and after: %v\n", totalBalanceBefore, totalBalanceAfterWithFee)
		if totalBalanceAfterWithFee.Cmp(totalBalanceBefore) != 0 {
			log.Fatal("Total balance mismatch!\n")
		}
	}
}
