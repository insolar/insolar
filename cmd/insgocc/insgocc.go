// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/insolar/insolar/applicationbase/genesisrefs"
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

func getAppropriateContractDir(machineType insolar.MachineType, dir string) string {
	if machineType == insolar.MachineTypeBuiltin {
		return path.Join(dir, "proxy")
	}
	panic(fmt.Sprintf("unknown machine type %v", machineType))
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
	parsed *preprocessor.ParsedFile,
	dir string) error {

	p := getAppropriateContractDir(machineType, dir)

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

func openDefaultInitializationPath(output *outputFlag, initPath string) error {
	err := output.SetJoin(initPath, "initialization.go")
	return err
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
	var reference string
	output := newOutputFlag("-")
	proxyOut := newOutputFlag("")
	machineType := newMachineTypeFlag("builtin")
	var panicIsLogicalError bool

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
				err = openDefaultProxyPath(proxyOut, machineType.Value(), parsed, "")
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
			if panicIsLogicalError {
				parsed.SetPanicIsLogicalError()
			}

			err = parsed.WriteWrapper(output.writer, parsed.ContractName())
			checkError(err)
		},
	}
	cmdWrapper.Flags().VarP(output, "output", "o", "output file (use - for STDOUT)")
	cmdWrapper.Flags().VarP(machineType, "machine-type", "m", "machine type (one of builtin/go)")
	cmdWrapper.Flags().BoolVarP(&panicIsLogicalError, "panic-logical", "p", false, "panics are logical errors (turned off by default)")

	var (
		importPath    string
		contractsPath string
	)
	var cmdGenerateBuiltins = &cobra.Command{
		Use:   "regen-builtin",
		Short: "Build builtin proxy, wrappers and initializator",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if !path.IsAbs(contractsPath) {
				dir, err := os.Getwd()
				checkError(err)
				contractsPath = path.Join(dir, contractsPath)
			}

			buildInPath := path.Join(contractsPath, "..")

			fileList, err := ioutil.ReadDir(contractsPath)
			checkError(err)

			contractList := make(preprocessor.ContractList, 0)

			// find all contracts in the folder
			for _, file := range fileList {
				if file.IsDir() {
					contractDirPath := path.Join(contractsPath, file.Name())

					contractPath := findContractPath(contractDirPath)
					if contractPath != nil {
						parsedFile, err := preprocessor.ParseFile(*contractPath, insolar.MachineTypeBuiltin)
						checkError(err)

						contract := preprocessor.ContractListEntry{
							Name:       file.Name(),
							Path:       *contractPath,
							Parsed:     parsedFile,
							ImportPath: path.Join(importPath, file.Name()),
						}
						contractList = append(contractList, contract)
					}
				}
			}

			for _, contract := range contractList {
				/* write proxy */
				output := newOutputFlag("")
				err := openDefaultProxyPath(output, insolar.MachineTypeBuiltin, contract.Parsed, buildInPath)
				checkError(err)
				reference := genesisrefs.GenerateProtoReferenceFromContractID(preprocessor.PrototypeType, contract.Name, contract.Version)
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
			err = openDefaultInitializationPath(initializeOutput, buildInPath)
			checkError(err)

			err = preprocessor.GenerateInitializationList(initializeOutput.writer, contractList)
			checkError(err)
		},
	}
	cmdGenerateBuiltins.Flags().StringVarP(
		&importPath, "importPath", "i", "", "import path for builtin contracts packages, example: github.com/insolar/insolar/application/builtin/contract")
	cmdGenerateBuiltins.Flags().StringVarP(
		&contractsPath, "contractsPath", "c", "", "dir path to builtin contracts, example: application/builtin/contract")

	var rootCmd = &cobra.Command{Use: "insgocc"}
	rootCmd.AddCommand(
		cmdProxy, cmdWrapper, cmdGenerateBuiltins)
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
