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

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/insolar/insolar/api/sdk"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

const defaultStdoutPath = "-"
const defaultMemberFileDir = "scripts/insolard/benchmark"
const defaultMemberFileName = "members.txt"

var (
	output             string
	concurrent         int
	repetitions        int
	rootMemberKeys     string
	apiURLs            []string
	logLevel           string
	saveMembersToFile  bool
	useMembersFromFile bool
)

func parseInputParams() {
	pflag.StringVarP(&output, "output", "o", defaultStdoutPath, "output file (use - for STDOUT)")
	pflag.IntVarP(&concurrent, "concurrent", "c", 1, "concurrent users")
	pflag.IntVarP(&repetitions, "repetitions", "r", 1, "repetitions for one user")
	pflag.StringVarP(&rootMemberKeys, "rootmemberkeys", "k", "", "path to file with RootMember keys")
	pflag.StringArrayVarP(&apiURLs, "apiurl", "u", []string{"http://localhost:19191/api"}, "url to api")
	pflag.StringVarP(&logLevel, "loglevel", "l", "info", "log level for benchmark")
	pflag.BoolVarP(&saveMembersToFile, "savemembers", "s", false, "save members to file")
	pflag.BoolVarP(&useMembersFromFile, "usemembers", "m", false, "use members from file")
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

func runScenarios(out io.Writer, insSDK *sdk.SDK, members []*sdk.Member, concurrent int, repetitions int) {
	transferDifferentMembers := &transferDifferentMembersScenario{
		concurrent:  concurrent,
		repetitions: repetitions,
		name:        "TransferDifferentMembers",
		out:         out,
		members:     members,
		insSDK:      insSDK,
	}
	startScenario(transferDifferentMembers)
}

func startScenario(s scenario) {
	err := s.canBeStarted()
	check(fmt.Sprintf("Scenario %s can not be started:", s.getName()), err)

	writeToOutput(s.getOut(), fmt.Sprintf("Scenario %s: Start to transfer\n", s.getName()))

	start := time.Now()
	s.start()
	elapsed := time.Since(start)

	writeToOutput(s.getOut(), fmt.Sprintf("Scenario %s: Transferring took %s \n", s.getName(), elapsed))
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

var numRetries = 3

func createMembers(insSDK *sdk.SDK, count int) []*sdk.Member {
	var members []*sdk.Member
	var member *sdk.Member
	var traceID string
	var err error

	for i := 0; i < count; i++ {
		for j := 0; j < numRetries; j++ {
			member, traceID, err = insSDK.CreateMember()
			if err == nil {
				members = append(members, member)
				break
			}

			fmt.Printf("Retry to create member. TraceID: %s Error is: %s\n", traceID, err.Error())
			time.Sleep(time.Second)
		}
		check(fmt.Sprintf("Couldn't create member after retries: %d", numRetries), err)
	}
	return members
}

func getTotalBalance(insSDK *sdk.SDK, members []*sdk.Member) uint64 {
	type Result struct {
		num     int
		balance uint64
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
			for attempt := 0; attempt < 5; attempt++ {
				res.balance, res.err = insSDK.GetBalance(m)
				if res.err == nil {
					break
				}
				// retry
				time.Sleep(1 * time.Second)
			}
			results <- res
			wg.Done()
		}(members[i], i)
	}

	wg.Wait()
	totalBalance := uint64(0)
	for i := 0; i < nmembers; i++ {
		res := <-results
		if res.err != nil {
			fmt.Printf("Can't get balance for %v-th member: %v\n", res.num, res.err)
			continue
		}
		totalBalance += res.balance
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
		members = createMembers(insSDK, concurrent*2)
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
	_, err = file.Write([]byte(result))
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

func main() {
	parseInputParams()

	// Start benchmark time
	t := time.Now()
	fmt.Printf("Start: %s\n\n", t.String())

	err := log.SetLevel(logLevel)
	check(fmt.Sprintf("Can't set '%s' level on logger:", logLevel), err)

	out, err := chooseOutput(output)
	check("Problems with output file:", err)

	insSDK, err := sdk.NewSDK(apiURLs, rootMemberKeys)
	check("SDK is not initialized: ", err)

	members, err := getMembers(insSDK)
	check("Error while loading members: ", err)
	totalBalanceBefore := getTotalBalance(insSDK, members)

	runScenarios(out, insSDK, members, concurrent, repetitions)

	// Finish benchmark time
	t = time.Now()
	fmt.Printf("\nFinish: %s\n\n", t.String())

	totalBalanceAfter := uint64(0)
	for nretries := 0; nretries < 5; nretries++ {
		totalBalanceAfter = getTotalBalance(insSDK, members)
		if totalBalanceAfter == totalBalanceBefore {
			break
		}
		fmt.Printf("Total balance before and after don't match: %v vs %v - retrying in 3 seconds...\n",
			totalBalanceBefore, totalBalanceAfter)
		time.Sleep(3 * time.Second)

	}
	fmt.Printf("Total balance before: %v and after: %v\n", totalBalanceBefore, totalBalanceAfter)
	if totalBalanceBefore != totalBalanceAfter {
		panic("Total balance mismatch!\n")
	}
}
