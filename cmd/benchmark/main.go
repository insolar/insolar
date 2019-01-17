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
	"fmt"
	"io"
	"os"
	"time"

	"github.com/insolar/insolar/api/sdk"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

const defaultStdoutPath = "-"

var (
	output         string
	concurrent     int
	repetitions    int
	rootMemberKeys string
	apiURLs        []string
	logLevel       string
)

func parseInputParams() {
	pflag.StringVarP(&output, "output", "o", defaultStdoutPath, "output file (use - for STDOUT)")
	pflag.IntVarP(&concurrent, "concurrent", "c", 1, "concurrent users")
	pflag.IntVarP(&repetitions, "repetitions", "r", 1, "repetitions for one user")
	pflag.StringVarP(&rootMemberKeys, "rootmemberkeys", "k", "", "path to file with RootMember keys")
	pflag.StringArrayVarP(&apiURLs, "apiurl", "u", []string{"http://localhost:19191/api"}, "url to api")
	pflag.StringVarP(&logLevel, "loglevel", "l", "info", "log level for benchmark")
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
	elapsedInSeconds := float64(elapsed) / float64(time.Second)
	speed := float64(s.getOperationsNumber()) / float64(elapsedInSeconds)
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
	var err error
	for i := 0; i < count; i++ {

		for j := 0; j < numRetries; j++ {
			member, _, err = insSDK.CreateMember()
			if err == nil {
				members = append(members, member)
				break
			}

			fmt.Println("Retry to create member. Error is: ", err.Error())
			time.Sleep(time.Second)
		}
		check(fmt.Sprintf("Couldn't create member after retries: %d", numRetries), err)
	}
	return members
}

func main() {
	parseInputParams()

	err := log.SetLevel(logLevel)
	check(fmt.Sprintf("Can't set '%s' level on logger:", logLevel), err)

	out, err := chooseOutput(output)
	check("Problems with output file:", err)

	insSDK, err := sdk.NewSDK(apiURLs, rootMemberKeys)
	check("SDK is not initialized: ", err)

	members := createMembers(insSDK, concurrent*2)

	runScenarios(out, insSDK, members, concurrent, repetitions)
}
