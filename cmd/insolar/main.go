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
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/insolar/insolar/configuration"
	"github.com/pkg/errors"
)

var (
	output string
	cmd    string
)

func chooseOutput(path string) (io.Writer, error) {
	var res io.Writer
	if path == "-" {
		res = os.Stdout
	} else {
		var err error
		res, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't open file for writing")
		}
	}
	return res, nil
}

func parseInputParams() {
	flag.StringVar(&output, "output", "-", "output file (use - for STDOUT)")
	flag.StringVar(&cmd, "cmd", "default_config", "type of cmd")

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	flag.Parse()
}

func printDefaultConfig(out io.Writer) {
	cfgHolder := configuration.NewHolder()

	out.Write([]byte(configuration.ToString(cfgHolder.Configuration)))
}

func main() {
	parseInputParams()
	out, err := chooseOutput(output)
	if err != nil {
		fmt.Println("Problems with parsing input:", err)
		os.Exit(1)
	}

	switch cmd {
	case "default_config":
		printDefaultConfig(out)
	}
}
