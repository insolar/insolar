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
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/insolar/insolar/logicrunner/goplugin/preprocessor"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

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

	var reference, outdir string
	output := newOutputFlag()

	var cmdProxy = &cobra.Command{
		Use:   "proxy [flags] <file name to process>",
		Short: "Generate contract's proxy",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				panic(errors.New("proxy command should be followed by exactly one file name to process"))
			}
			parsed, err := preprocessor.ParseFile(args[0])
			if err != nil {
				panic(errors.Wrap(err, "couldn't parse"))
			}

			err = preprocessor.GenerateContractProxy(parsed, reference, output.writer)
			if err != nil {
				panic(err)
			}
		},
	}
	cmdProxy.Flags().StringVarP(&reference, "code-reference", "r", "", "reference to code of")
	cmdProxy.Flags().VarP(output, "output", "o", "output file (use - for STDOUT)")

	var cmdWrapper = &cobra.Command{
		Use:   "wrapper [flags] <file name to process>",
		Short: "Generate contract's wrapper",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				panic(errors.New("wrapper command should be followed by exactly one file name to process"))
			}
			parsed, err := preprocessor.ParseFile(args[0])
			if err != nil {
				panic(errors.Wrap(err, "couldn't parse"))
			}

			err = preprocessor.GenerateContractWrapper(parsed, output.writer)
			if err != nil {
				panic(err)
			}
		},
	}
	cmdWrapper.Flags().VarP(output, "output", "o", "output file (use - for STDOUT)")

	var cmdImports = &cobra.Command{
		Use:   "imports [flags] <file name to process>",
		Short: "Rewrite imports in contract file",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				panic(errors.New("imports command should be followed by exactly one file name to process"))
			}
			parsed, err := preprocessor.ParseFile(args[0])
			if err != nil {
				panic(errors.Wrap(err, "couldn't parse"))
			}

			err = preprocessor.CmdRewriteImports(parsed, output.writer)
			if err != nil {
				panic(err)
			}
		},
	}
	cmdImports.Flags().VarP(output, "output", "o", "output file (use - for STDOUT)")

	var cmdCompile = &cobra.Command{
		Use:   "compile [flags] <file name to compile>",
		Short: "Compile contract",
		Run: func(cmd *cobra.Command, args []string) {
			dir, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			if len(args) != 1 {
				panic(errors.New("compile command should be followed by exactly one file name to compile"))
			}
			parsed, err := preprocessor.ParseFile(args[0])
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

			out, err := exec.Command("go", "build", "-buildmode=plugin", "-o", path.Join(dir, outdir, name+".so")).CombinedOutput()
			if err != nil {
				panic(errors.Wrap(err, "can't build contract: "+string(out)))
			}
		},
	}
	cmdCompile.Flags().StringVarP(&outdir, "output-dir", "o", ".", "output dir (default .)")

	var rootCmd = &cobra.Command{Use: "insgocc"}
	rootCmd.AddCommand(cmdProxy, cmdWrapper, cmdImports, cmdCompile)
	rootCmd.Execute()
}
