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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/preprocessor"
)

type outputFlag struct {
	path   string
	writer io.Writer
}

func newOutputFlag(path string) *outputFlag {
	return &outputFlag{path: path, writer: os.Stdout}
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
		res, err = os.OpenFile(arg, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
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

type machineTypeFlag struct {
	name string
	num  insolar.MachineType
}

func newMachineTypeFlag(name string) *machineTypeFlag {
	flag := machineTypeFlag{}
	if err := flag.Set(name); err != nil {
		panic(fmt.Sprintf("unknown error: %s", err))
	}
	return &flag
}

func (r *machineTypeFlag) Set(arg string) error {
	switch arg {
	case "":
		fallthrough
	case "go":
		fallthrough
	case "golang":
		r.num = insolar.MachineTypeGoPlugin
	case "builtin":
		r.num = insolar.MachineTypeBuiltin
	default:
		return fmt.Errorf("unknown machine type: %s", arg)
	}
	r.name = arg
	return nil
}

func (r *machineTypeFlag) String() string {
	return r.name
}

func (r *machineTypeFlag) Type() string {
	return "machineType"
}

func (r *machineTypeFlag) Value() insolar.MachineType {
	return r.num
}

func main() {

	var reference, outdir string
	output := newOutputFlag("-")
	proxyOut := newOutputFlag("")
	machineType := newMachineTypeFlag("go")

	var cmdProxy = &cobra.Command{
		Use:   "proxy [flags] <file name to process>",
		Short: "Generate contract's proxy",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) != 1 {
				fmt.Println("proxy command should be followed by exactly one file name to process")
				os.Exit(1)
			}

			parsed, err := preprocessor.ParseFile(args[0], machineType.Value())
			if err != nil {
				fmt.Println(errors.Wrap(err, "couldn't parse"))
				os.Exit(1)
			}

			if proxyOut.String() == "" {
				p, err := preprocessor.GetRealApplicationDir("proxy")
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				proxyPackage, err := parsed.ProxyPackageName()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				_, err = os.Stat(path.Join(p, proxyPackage))
				if err != nil {
					err := os.Mkdir(path.Join(p, proxyPackage), 0755)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				}

				err = proxyOut.Set(path.Join(p, proxyPackage, proxyPackage+".go"))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			err = parsed.WriteProxy(reference, proxyOut.writer)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdProxy.Flags().StringVarP(&reference, "code-reference", "r", "", "reference to code of")
	cmdProxy.Flags().VarP(proxyOut, "output", "o", "output file (use - for STDOUT)")
	cmdProxy.Flags().VarP(machineType, "machine-type", "m", "machine type (one of builtin/go)")

	var cmdWrapper = &cobra.Command{
		Use:   "wrapper [flags] <file name to process>",
		Short: "Generate contract's wrapper",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("wrapper command should be followed by exactly one file name to process")
				os.Exit(1)
			}
			parsed, err := preprocessor.ParseFile(args[0], machineType.Value())
			if err != nil {
				fmt.Println(errors.Wrap(err, "couldn't parse"))
				os.Exit(1)
			}

			err = parsed.WriteWrapper(output.writer)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdWrapper.Flags().VarP(output, "output", "o", "output file (use - for STDOUT)")
	cmdWrapper.Flags().VarP(machineType, "machine-type", "m", "machine type (one of builtin/go)")

	var cmdImports = &cobra.Command{
		Use:   "imports [flags] <file name to process>",
		Short: "Rewrite imports in contract file",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("imports command should be followed by exactly one file name to process")
				os.Exit(1)
			}
			parsed, err := preprocessor.ParseFile(args[0], machineType.Value())
			if err != nil {
				fmt.Println(errors.Wrap(err, "couldn't parse"))
				os.Exit(1)
			}

			err = parsed.Write(output.writer)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdImports.Flags().VarP(output, "output", "o", "output file (use - for STDOUT)")
	cmdImports.Flags().VarP(machineType, "machine-type", "m", "machine type (one of builtin/go)")

	// PLEASE NOTE that `insgocc compile` is in fact not used for compiling contracts by insolard.
	// Instead contracts are compiled when `insolard genesis` is executed without using `insgocc`.
	keepTemp := false
	var cmdCompile = &cobra.Command{
		Use:   "compile [flags] <file name to compile>",
		Short: "Compile contract",
		Run: func(cmd *cobra.Command, args []string) {
			dir, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if len(args) != 1 {
				fmt.Println("compile command should be followed by exactly one file name to compile")
				os.Exit(1)
			}
			parsed, err := preprocessor.ParseFile(args[0], machineType.Value())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// make temporary dir
			tmpDir, err := ioutil.TempDir("", "temp-")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			defer func() {
				if keepTemp {
					fmt.Printf("Temp directory: %s\n", tmpDir)
				} else {
					os.RemoveAll(tmpDir) // nolint: errcheck
				}
			}()

			name := parsed.ContractName()

			contract, err := os.Create(filepath.Join(tmpDir, name+".go"))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer contract.Close()

			parsed.ChangePackageToMain()
			err = parsed.Write(contract)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			wrapper, err := os.Create(filepath.Join(tmpDir, name+".wrapper.go"))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer wrapper.Close()

			err = parsed.WriteWrapper(wrapper)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = os.Chdir(tmpDir)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			out, err := exec.Command("go", "build", "-buildmode=plugin", "-o", path.Join(dir, outdir, name+".so")).CombinedOutput()
			if err != nil {
				fmt.Println(errors.Wrap(err, "can't build contract: "+string(out)))
				os.Exit(1)
			}
		},
	}
	// default value for string flags is displayed automatically
	cmdCompile.Flags().StringVarP(&outdir, "output-dir", "o", ".", "output dir")
	// default value for bool flags is not displayed automatically, thus it's done manually here
	cmdCompile.Flags().BoolVarP(&keepTemp, "keep-temp", "k", false, "keep temp directory (default \"false\")")
	cmdCompile.Flags().VarP(machineType, "machine-type", "m", "machine type (one of builtin/go)")

	var rootCmd = &cobra.Command{Use: "insgocc"}
	rootCmd.AddCommand(cmdProxy, cmdWrapper, cmdImports, cmdCompile)
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
