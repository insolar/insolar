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
	"go/build"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"

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

func (r *outputFlag) SetJoin(pathParts ...string) error {
	return r.Set(path.Join(pathParts...))
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

var (
	rootProjectDir   string
	rootProjectError error
	rootProjectOnce  sync.Once
)

func getRootProjectDir() (string, error) {
	rootProjectOnce.Do(func() {
		gopath := build.Default.GOPATH
		if gopath == "" {
			rootProjectDir, rootProjectError = "", errors.New("GOPATH is not set")
			return
		}
		contractsPath := ""
		for _, p := range strings.Split(gopath, ":") {
			contractsPath = path.Join(p, "src/github.com/insolar/insolar/")
			_, err := os.Stat(contractsPath)
			if err == nil {
				rootProjectDir, rootProjectError = contractsPath, nil
				return
			}
		}

		rootProjectDir, rootProjectError = "", errors.New("Not found github.com/insolar/insolar in GOPATH")
	})
	return rootProjectDir, rootProjectError
}

func getBuiltinContractDir(dir string) (string, error) {
	projectRoot, err := getRootProjectDir()
	if err != nil {
		return "", err
	}
	return path.Join(projectRoot, "logicrunner", "builtin", dir), nil
}

func getApplicationContractDir(dir string) (string, error) {
	projectRoot, err := getRootProjectDir()
	if err != nil {
		return "", err
	}
	return path.Join(projectRoot, "application", dir), nil
}

func getAppropriateContractDir(machineType insolar.MachineType, dir string) (string, error) {
	if machineType == insolar.MachineTypeBuiltin {
		return getBuiltinContractDir(dir)
	} else if machineType == insolar.MachineTypeGoPlugin {
		return getApplicationContractDir(dir)
	}
	panic("unreachable")
}

func mkdirIfNotExists(pathParts ...string) (string, error) {
	newPath := path.Join(pathParts...)
	stat, err := os.Stat(newPath)
	if err == nil {
		if stat.IsDir() {
			return newPath, nil
		}
		return "", fmt.Errorf("failed to mkdir '%s': already exists and is not dir", newPath)
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(newPath, 0755); err != nil {
			return "", errors.Wrap(err, "failed to mkdir "+newPath)
		}
		return newPath, nil
	}
	return "", errors.Wrap(err, "failed to mkdir "+newPath)
}

func openDefaultProxyPath(proxyOut *outputFlag,
	machineType insolar.MachineType,
	parsed *preprocessor.ParsedFile) error {

	p, err := getAppropriateContractDir(machineType, "proxy")
	if err != nil {
		return err
	}

	proxyPackage, err := parsed.ProxyPackageName()
	if err != nil {
		return err
	}

	proxyPath, err := mkdirIfNotExists(p, proxyPackage)
	if err != nil {
		return err
	}

	err = proxyOut.SetJoin(proxyPath, proxyPackage+".go")
	if err != nil {
		return err
	}

	return nil
}

func openDefaultInitializationPath(output *outputFlag) error {
	initPath, err := getBuiltinContractDir("")
	if err != nil {
		return err
	}

	err = output.SetJoin(initPath, "initialization.go")
	if err != nil {
		return err
	}

	return nil
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func findContractPath(contractDirPath string) *string {
	contractName := path.Base(contractDirPath)
	for _, contractFileName := range []string{"main.go", contractName + ".go"} {
		contractPath := path.Join(contractDirPath, contractFileName)
		if stat, err := os.Stat(contractPath); err == nil && !stat.IsDir() {
			return &contractPath
		}
	}
	return nil
}

func main() {
	var reference, outdir string
	output := newOutputFlag("-")
	proxyOut := newOutputFlag("")
	machineType := newMachineTypeFlag("go")

	var cmdProxy = &cobra.Command{
		Use:   "proxy [flags] <file name to process>",
		Short: "Generate contract's proxy",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			parsed, err := preprocessor.ParseFile(args[0], machineType.Value())
			if err != nil {
				fmt.Println(errors.Wrap(err, "couldn't parse"))
				os.Exit(1)
			}

			if proxyOut.String() == "" {
				err = openDefaultProxyPath(proxyOut, machineType.Value(), parsed)
				checkError(err)
			}

			err = parsed.WriteProxy(reference, proxyOut.writer)
			checkError(err)
		},
	}
	cmdProxy.Flags().StringVarP(&reference, "code-reference", "r", "", "reference to code of")
	cmdProxy.Flags().VarP(proxyOut, "output", "o", "output file (use - for STDOUT)")
	cmdProxy.Flags().VarP(machineType, "machine-type", "m", "machine type (one of builtin/go)")

	var cmdWrapper = &cobra.Command{
		Use:   "wrapper [flags] <file name to process>",
		Short: "Generate contract's wrapper",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			parsed, err := preprocessor.ParseFile(args[0], machineType.Value())
			if err != nil {
				fmt.Println(errors.Wrap(err, "couldn't parse"))
				os.Exit(1)
			}

			err = parsed.WriteWrapper(output.writer, "main")
			checkError(err)
		},
	}
	cmdWrapper.Flags().VarP(output, "output", "o", "output file (use - for STDOUT)")
	cmdWrapper.Flags().VarP(machineType, "machine-type", "m", "machine type (one of builtin/go)")

	var cmdImports = &cobra.Command{
		Use:   "imports [flags] <file name to process>",
		Short: "Rewrite imports in contract file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			parsed, err := preprocessor.ParseFile(args[0], machineType.Value())
			if err != nil {
				fmt.Println(errors.Wrap(err, "couldn't parse"))
				os.Exit(1)
			}

			err = parsed.Write(output.writer)
			checkError(err)
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
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dir, err := os.Getwd()
			checkError(err)

			parsed, err := preprocessor.ParseFile(args[0], machineType.Value())
			checkError(err)

			// make temporary dir
			tmpDir, err := ioutil.TempDir("", "temp-")
			checkError(err)

			defer func() {
				if keepTemp {
					fmt.Printf("Temp directory: %s\n", tmpDir)
				} else {
					os.RemoveAll(tmpDir) // nolint: errcheck
				}
			}()

			name := parsed.ContractName()

			contract, err := os.Create(filepath.Join(tmpDir, name+".go"))
			checkError(err)
			defer contract.Close()

			parsed.ChangePackageToMain()
			err = parsed.Write(contract)
			checkError(err)

			wrapper, err := os.Create(filepath.Join(tmpDir, name+".wrapper.go"))
			checkError(err)
			defer wrapper.Close()

			err = parsed.WriteWrapper(wrapper, "main")
			checkError(err)

			err = os.Chdir(tmpDir)
			checkError(err)

			contractPath := path.Join(dir, outdir, name+".so")
			cmdArgs := []string{"build", "-buildmode=plugin", "-o", contractPath}
			out, err := exec.Command("go", cmdArgs...).CombinedOutput()
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

	var cmdGenerateBuiltins = &cobra.Command{
		Use:   "regen-builtin [flags] <dir path to builtin contracts>",
		Short: "Build builtin proxy, wrappers and initializator",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			contractPath, err := getBuiltinContractDir("contract")
			checkError(err)

			fileList, err := ioutil.ReadDir(contractPath)
			checkError(err)

			contractList := make(preprocessor.ContractList, 0)

			rootProjectDir, err := getRootProjectDir()
			checkError(err)

			// find all contracts in the folder
			for _, file := range fileList {
				if file.IsDir() {
					contractDirPath := path.Join(contractPath, file.Name())

					contractPath := findContractPath(contractDirPath)
					if contractPath != nil {
						parsedFile, err := preprocessor.ParseFile(*contractPath, insolar.MachineTypeBuiltin)
						checkError(err)

						contract := preprocessor.ContractListEntry{
							Name:       file.Name(),
							Path:       *contractPath,
							Parsed:     parsedFile,
							ImportPath: "github.com/insolar/insolar/" + contractDirPath[len(rootProjectDir)+1:],
						}
						contractList = append(contractList, contract)
					}
				}
			}

			for _, contract := range contractList {
				/* write proxy */
				output := newOutputFlag("")
				err := openDefaultProxyPath(output, insolar.MachineTypeBuiltin, contract.Parsed)
				checkError(err)

				reference := contract.GenerateReference(preprocessor.PrototypeType)
				err = contract.Parsed.WriteProxy(reference.String(), output.writer)
				checkError(err)

				/* write wrappers */
				err = output.SetJoin(path.Dir(contract.Path), contract.Name+".wrapper.go")
				checkError(err)

				err = contract.Parsed.WriteWrapper(output.writer, contract.Parsed.ContractName())
				checkError(err)
			}

			// write include contract + write initialization function
			initializeOutput := newOutputFlag("")
			err = openDefaultInitializationPath(initializeOutput)
			checkError(err)

			err = preprocessor.GenerateInitializationList(initializeOutput.writer, contractList)
			checkError(err)
		},
	}

	var rootCmd = &cobra.Command{Use: "insgocc"}
	rootCmd.AddCommand(cmdProxy, cmdWrapper, cmdImports, cmdCompile, cmdGenerateBuiltins)
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
