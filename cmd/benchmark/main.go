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
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

const defaultStdoutPath = "-"

var (
	input          string
	output         string
	concurrent     int
	repetitions    int
	rootmemberkeys string

	rootMember memberInfo
)

func parseInputParams() {
	pflag.StringVarP(&input, "input", "i", "", "path to file with initial data for benchmark")
	pflag.StringVarP(&output, "output", "o", defaultStdoutPath, "output file (use - for STDOUT)")
	pflag.IntVarP(&concurrent, "concurrent", "c", 1, "concurrent users")
	pflag.IntVarP(&repetitions, "repetitions", "r", 1, "repetitions for one user")
	pflag.StringVarP(&rootmemberkeys, "rqootmemberkeys", "k", "", "path to file with RootMember keys")
	pflag.Parse()
}

// TODO FIXME use buffered output
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

// TODO FIXME out is not thread safe! Also it should be flushed.
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

type memberInfo struct {
	ref        string
	privateKey string
}

const memberInfoFieldsNumber = 2

func getMembersInfo(fileName string) ([]memberInfo, error) {
	var members []memberInfo

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't open file for reading")
	}
	defer file.Close() //nolint: errcheck

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		info := strings.Fields(scanner.Text())
		if len(info) != memberInfoFieldsNumber {
			check("problem with getting member info", errors.New("not enough info for single member"))
		}
		members = append(members, memberInfo{ref: info[0], privateKey: info[1]})
	}

	return members, nil
}

type memberKeys struct {
	Private string `json:"private_key"`
	Public  string `json:"public_key"`
}

func getRootMemberRef() string {
	infoResp := info()
	return infoResp.RootMember
}

func getRootMemberInfo(fileName string) memberInfo {

	rawConf, err := ioutil.ReadFile(fileName)
	check("problem with reading root member keys file", err)

	keys := memberKeys{}
	err = json.Unmarshal(rawConf, &keys)
	check("problem with unmarshaling root member keys", err)

	return memberInfo{getRootMemberRef(), keys.Private}
}

func runScenarios(out io.Writer, members []memberInfo, concurrent int, repetitions int) {
	transferDifferentMembers := &transferDifferentMembersScenario{
		concurrent:  concurrent,
		repetitions: repetitions,
		members:     members,
		name:        "TransferDifferentMembers",
		out:         out,
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
}

func main() {
	log.Printf("Starting benchmark. Calling parseInputParams()...\n")
	parseInputParams()

	log.Printf("Calling chooseOutput()...\n")
	out, err := chooseOutput(output)
	check("Problems with output file:", err)

	var members []memberInfo

	log.Printf("Calling getRootMemberInfo()...\n")
	rootMember = getRootMemberInfo(rootmemberkeys)

	if input != "" {
		log.Printf("Calling getMembersInfo()...\n")
		members, err = getMembersInfo(input)
		check("Problems with parsing input:", err)
	} else {
		log.Printf("Calling createMembers(), this may take some time...\n") // TODO: THIS DOES TAKE TIME
		members, err = createMembers(concurrent, repetitions)
		check("Problems with create members. One of creating request ended with error: ", err)
	}

	log.Printf("Calling runScenarios()...\n")
	runScenarios(out, members, concurrent, repetitions)
}
