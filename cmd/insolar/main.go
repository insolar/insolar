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
	"github.com/insolar/insolar/core"
	ecdsa_helper "github.com/insolar/insolar/crypto_helpers/ecdsa"
	"github.com/insolar/insolar/version"
	"github.com/pkg/errors"
)

var (
	output string
	cmd    string
)

const defaultStdoutPath = "-"

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

func parseInputParams() {
	flag.StringVar(&output, "output", defaultStdoutPath, "output file (use - for STDOUT)")
	flag.StringVar(&cmd, "cmd", "default_config", "available commands: default_config | random_ref | version | gen_keys")

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	flag.Parse()
}

func writeToOutput(out io.Writer, data string) {
	_, err := out.Write([]byte(data))
	if err != nil {
		fmt.Println("Can't write data to output", err)
		os.Exit(1)
	}
}

func printDefaultConfig(out io.Writer) {
	cfgHolder := configuration.NewHolder()

	writeToOutput(out, configuration.ToString(cfgHolder.Configuration))
}

func randomRef(out io.Writer) {
	ref := core.RandomRef()

	writeToOutput(out, ref.String()+"\n")
}

func generateKeysPair(out io.Writer) {
	privKey, err := ecdsa_helper.GeneratePrivateKey()
	if err != nil {
		fmt.Println("Problems with generating of private key:", err)
		os.Exit(1)
	}

	privKeyStr, err := ecdsa_helper.ExportPrivateKey(privKey)
	if err != nil {
		fmt.Println("Problems with serialization of private key:", err)
		os.Exit(1)
	}

	pubKeyStr, err := ecdsa_helper.ExportPublicKey(&privKey.PublicKey)
	if err != nil {
		fmt.Println("Problems with serialization of public key:", err)
		os.Exit(1)
	}

	result := fmt.Sprintf("Public key:\n %s\n", pubKeyStr)
	result += fmt.Sprintf("Private key:\n %s", privKeyStr)

	writeToOutput(out, result)
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
	case "random_ref":
		randomRef(out)
	case "version":
		fmt.Println(version.GetFullVersion())
	case "gen_keys":
		generateKeysPair(out)
	}
}
