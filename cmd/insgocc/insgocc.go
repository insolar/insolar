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
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/insolar/insolar/logicrunner/goplugin/preprocessor"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

func printUsage() {
	fmt.Println("usage: insgocc <command> [<args>]")
	fmt.Println("Commands: ")
	fmt.Println(" wrapper   generate contract's wrapper")
	fmt.Println(" proxy     generate contract's proxy")
	fmt.Println(" compile   compile contract")
	fmt.Println(" imports   rewrite imports")
}

type outputFlag struct {
	path   string
	writer io.Writer
}

func newOutputFlag() *outputFlag {
	return &outputFlag{path: "-", writer: os.Stdout}
}

func (r *outputFlag) String() string {
	return r.path
}
func (r *outputFlag) Set(arg string) error {
	var res io.Writer
	if arg == "-" {
		res = os.Stdout
	} else {
		var err error
		res, err = os.OpenFile(arg, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return errors.Wrap(err, "couldn't open file for writing")
		}
	}
	r.path = arg
	r.writer = res
	return nil
}
func (r *outputFlag) Type() string {
	return "file"
}

func main() {

	if len(os.Args) == 1 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "wrapper":
		fs := flag.NewFlagSet("wrapper", flag.ExitOnError)
		output := newOutputFlag()
		fs.VarP(output, "output", "o", "output file (use - for STDOUT)")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		for _, fn := range fs.Args() {
			parsed, err := preprocessor.ParseFile(fn)
			if err != nil {
				panic(errors.Wrap(err, "couldn't parse"))
			}

			err = preprocessor.GenerateContractWrapper(parsed, output.writer)
			if err != nil {
				panic(err)
			}
		}
	case "proxy":
		fs := flag.NewFlagSet("proxy", flag.ExitOnError)
		output := newOutputFlag()
		fs.VarP(output, "output", "o", "output file (use - for STDOUT)")
		var reference string
		fs.StringVarP(&reference, "code-reference", "r", "testRef", "reference to code of")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		if fs.NArg() != 1 {
			panic(errors.New("proxy command should be followed by exactly one file name to process"))
		}

		parsed, err := preprocessor.ParseFile(fs.Arg(0))
		if err != nil {
			panic(errors.Wrap(err, "couldn't parse"))
		}

		err = preprocessor.GenerateContractProxy(parsed, reference, output.writer)
		if err != nil {
			panic(err)
		}
	case "imports":
		fs := flag.NewFlagSet("imports", flag.ExitOnError)
		output := newOutputFlag()
		fs.VarP(output, "output", "o", "output file (use - for STDOUT)")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		if fs.NArg() != 1 {
			panic(errors.New("imports command should be followed by exactly one file name to process"))
		}

		parsed, err := preprocessor.ParseFile(fs.Arg(0))
		if err != nil {
			panic(errors.Wrap(err, "couldn't parse"))
		}

		err = preprocessor.CmdRewriteImports(parsed, output.writer)
		if err != nil {
			panic(err)
		}
	// use compile [-o <output-directory>] <smart-contract-file>
	case "compile":
		fs := flag.NewFlagSet("compile", flag.ExitOnError)
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		var output string
		fs.StringVarP(&output, "output", "o", dir, "output directory (default .)")
		err = fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		parsed, err := preprocessor.ParseFile(fs.Arg(0))
		if err != nil {
			panic(err)
		}

		// make temporary dir
		tmpDir, err := ioutil.TempDir("", "test-")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(tmpDir) // nolint: errcheck

		name := preprocessor.GetContractName(parsed)

		contract, err := os.Create(tmpDir + "/" + name + ".go")
		if err != nil {
			panic(err)
		}
		defer contract.Close()

		preprocessor.RewriteContractPackage(parsed, contract)

		wrapper, err := os.Create(tmpDir + "/" + name + ".wrapper.go")
		if err != nil {
			panic(err)
		}
		defer wrapper.Close()

		err = preprocessor.GenerateContractWrapper(parsed, wrapper)
		if err != nil {
			panic(err)
		}

		err = os.Chdir(tmpDir)
		if err != nil {
			panic(err)
		}

		out, err := exec.Command("go", "build", "-buildmode=plugin", "-o", path.Join(dir, output, name+".so")).CombinedOutput()
		if err != nil {
			panic(errors.Wrap(err, "can't build contract: "+string(out)))
		}

	default:
		printUsage()
		fmt.Printf("\n\n%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
}
