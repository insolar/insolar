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
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

const defaultStdoutPath = "-"

var (
	input       string
	output      string
	concurrent  int
	repetitions int
	withInit    bool
)

func parseInputParams() {
	pflag.StringVarP(&input, "input", "i", "", "path to file with initial data for loads")
	pflag.StringVarP(&output, "output", "o", defaultStdoutPath, "output file (use - for STDOUT)")
	pflag.IntVarP(&concurrent, "concurrent", "c", 1, "concurrent users")
	pflag.IntVarP(&repetitions, "repetitions", "r", 1, "repetitions for one user")
	pflag.BoolVar(&withInit, "with_init", false, "do initialization before run load")
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

func getMembersRef(fileName string) ([]string, error) {
	var members []string

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't open file for reading")
	}
	defer file.Close() //nolint: errcheck

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		members = append(members, scanner.Text())
	}

	return members, nil
}

func runScenarios(out io.Writer, members []string, concurrent int, repetitions int) {
	firstScenario := &transferDifferentMembersScenario{
		concurrent:  concurrent,
		repetitions: repetitions,
		members:     members,
		name:        "TransferDifferentMembers",
		out:         out,
	}
	startScenario(firstScenario)
}

func startScenario(s scenario) {
	var wg sync.WaitGroup

	err := s.canBeStarted()
	check(fmt.Sprintf("Scenario %s can not be started:", s.getName()), err)

	writeToOutput(s.getOut(), fmt.Sprintf("Scenario %s: Start to transfer\n", s.getName()))

	start := time.Now()
	s.start(&wg)
	wg.Wait()
	elapsed := time.Since(start)

	writeToOutput(s.getOut(), fmt.Sprintf("Scenario %s: Transfering took %s \n", s.getName(), elapsed))
	elapsedInSeconds := float64(elapsed) / float64(time.Second)
	speed := float64(s.getOperationsNumber()) / float64(elapsedInSeconds)
	writeToOutput(s.getOut(), fmt.Sprintf("Scenario %s: Speed - %f tr/s \n", s.getName(), speed))
}

func main() {
	parseInputParams()

	out, err := chooseOutput(output)
	check("Problems with output file:", err)

	var members []string

	if withInit {
		members, err = createMembers(concurrent, repetitions)
		check("Problems with create members. One of creating request ended with error: ", err)
	}

	if input != "" {
		members, err = getMembersRef(input)
		check("Problems with parsing input:", err)
	}

	runScenarios(out, members, concurrent, repetitions)
}
