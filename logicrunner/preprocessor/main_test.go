// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build slowtest

package preprocessor

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/testutils"
)

const useLeakTest = false

type PreprocessorSuite struct {
	suite.Suite
}

var randomTestCode = `
package main

import (
	"fmt"
	"errors"

	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

type HelloWorlder struct {
	foundation.BaseContract
	Greeted int
}

type FullName struct {
	First string
	Last  string
}

type PersonalGreeting struct {
	Name    FullName
	Message string
}

func (hw *HelloWorlder) Hello() (string, error) {
	hw.Greeted++
	return "Hello world 2", nil
}

func (hw *HelloWorlder) Fail() (string, error) {
	hw.Greeted++
	return "", errors.New("We failed 2")
}

func (hw *HelloWorlder) Echo(s string) (string, error) {
	hw.Greeted++
	return s, nil
}

func (hw *HelloWorlder) HelloHuman(Name FullName) (PersonalGreeting, error) {
	hw.Greeted++
	return PersonalGreeting{
		Name:    Name,
		Message: fmt.Sprintf("Dear %s %s, we specially say hello to you", Name.First, Name.Last),
	}, nil
}

func (hw *HelloWorlder) HelloHumanPointer(Name FullName) (*PersonalGreeting, error) {
	hw.Greeted++
	return &PersonalGreeting{
		Name:    Name,
		Message: fmt.Sprintf("Dear %s %s, we specially say hello to you", Name.First, Name.Last),
	}, nil
}

func (hw *HelloWorlder) MultiArgs(Name FullName, s string, i int) (*PersonalGreeting, error) {
	hw.Greeted++
	return &PersonalGreeting{
		Name:    Name,
		Message: fmt.Sprintf("Dear %s %s, we specially say hello to you", Name.First, Name.Last),
	}, nil
}

func (hw HelloWorlder) ConstEcho(s string) (string, error) {
	return s, nil
}

func JustExportedStaticFunction(int, int) error { return nil }
`

func (s *PreprocessorSuite) TestBasicGeneration() {
	tmpDir, err := ioutil.TempDir("", "test_")
	s.NoError(err)
	defer os.RemoveAll(tmpDir) // nolint: errcheck

	err = WriteFile(tmpDir, "main.go", randomTestCode)
	s.NoError(err)

	parsed, err := ParseFile(filepath.Join(tmpDir, "main.go"), insolar.MachineTypeBuiltin)
	s.NoError(err)
	s.NotNil(parsed)

	s.T().Run("wrapper", func(t *testing.T) {
		if useLeakTest {
			defer testutils.LeakTester(t)
		} else {
			t.Parallel()
		}
		a := assert.New(t)

		buf := bytes.Buffer{}
		err := parsed.WriteWrapper(&buf, parsed.ContractName())
		a.NoError(err)

		code, err := ioutil.ReadAll(&buf)
		a.NoError(err)
		a.NotEmpty(code)
	})

	s.T().Run("proxy", func(t *testing.T) {
		if useLeakTest {
			defer testutils.LeakTester(t)
		} else {
			t.Parallel()
		}
		a := assert.New(t)

		buf := bytes.Buffer{}
		err := parsed.WriteProxy(gen.Reference().String(), &buf)
		a.NoError(err)

		code, err := ioutil.ReadAll(&buf)
		a.NoError(err)
		a.NotEmpty(code)
	})
}

func (s *PreprocessorSuite) TestConstructorsParsing() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir) // nolint: errcheck

	code := `
package main

type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

func NewFromString(s string) (*One, error) {
	return &One{}, nil
}
`

	err = WriteFile(tmpDir, "code1", code)
	s.NoError(err)

	info, err := ParseFile(filepath.Join(tmpDir, "code1"), insolar.MachineTypeBuiltin)
	s.NoError(err)

	s.Equal(1, len(info.constructors))
	s.Equal(2, len(info.constructors["One"]))
	s.Equal("New", info.constructors["One"][0].Name.Name)
	s.Equal("NewFromString", info.constructors["One"][1].Name.Name)

	code = `
package main

type One struct {
	foundation.BaseContract
}

func New() {
	return
}
`

	err = WriteFile(tmpDir, "code1", code)
	s.NoError(err)

	_, err = ParseFile(filepath.Join(tmpDir, "code1"), insolar.MachineTypeBuiltin)
	s.Error(err)

	code = `
package main

type One struct {
	foundation.BaseContract
}

func New() *One {
	return &One{}
}
`

	err = WriteFile(tmpDir, "code1", code)
	s.NoError(err)

	_, err = ParseFile(filepath.Join(tmpDir, "code1"), insolar.MachineTypeBuiltin)
	s.Error(err)
}

func (s *PreprocessorSuite) TestFailIfThereAreNoContract() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir) //nolint: errcheck

	testContract := "/test.go"
	err = WriteFile(tmpDir, testContract, `
package main
type A struct{
	ttt ppp.TTT
}
`)
	s.NoError(err)

	_, err = ParseFile(tmpDir+testContract, insolar.MachineTypeBuiltin)
	s.EqualError(err, "Only one smart contract must exist")
}

func (s *PreprocessorSuite) TestInitializationFunctionParamsProxy() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir) // nolint: errcheck

	testContract := "/test.go"

	err = WriteFile(tmpDir, testContract, `
package main

type A struct{
	foundation.BaseContract
}

func ( a *A ) Get(
	a int, b bool, c string, d foundation.Reference,
) (
	int, bool, string, foundation.Reference, error,
) {
	return nil
}
`)

	s.NoError(err)

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeBuiltin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(gen.Reference().String(), &bufProxy)
	s.NoError(err)
	s.Contains(bufProxy.String(), "var ret0 int")
	s.Contains(bufProxy.String(), "ret[0] = &ret0")

	s.Contains(bufProxy.String(), "var ret1 bool")
	s.Contains(bufProxy.String(), "ret[1] = &ret1")

	s.Contains(bufProxy.String(), "var ret2 string")
	s.Contains(bufProxy.String(), "ret[2] = &ret2")

	s.Contains(bufProxy.String(), "var ret3 foundation.Reference")
	s.Contains(bufProxy.String(), "ret[3] = &ret3")
}

func (s *PreprocessorSuite) TestInitializationFunctionParamsWrapper() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir) //nolint: errcheck

	testContract := "/test.go"

	err = WriteFile(tmpDir, testContract, `
package main

type A struct{
	foundation.BaseContract
}

func (a *A) Get(
	a int, b bool, c string, d foundation.Reference,
) (
	int, bool, string, foundation.Reference, error,
) {
	return
}
`)
	s.NoError(err)

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeBuiltin)
	s.NoError(err)

	var bufWrapper bytes.Buffer
	err = parsed.WriteWrapper(&bufWrapper, parsed.ContractName())
	s.NoError(err)
	s.Contains(bufWrapper.String(), "var args0 int")
	s.Contains(bufWrapper.String(), "args[0] = &args0")

	s.Contains(bufWrapper.String(), "var args1 bool")
	s.Contains(bufWrapper.String(), "args[1] = &args1")

	s.Contains(bufWrapper.String(), "var args2 string")
	s.Contains(bufWrapper.String(), "args[2] = &args2")

	s.Contains(bufWrapper.String(), "var args3 foundation.Reference")
	s.Contains(bufWrapper.String(), "args[3] = &args3")
}

func (s *PreprocessorSuite) TestConstructorsWrapper() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir) //nolint: errcheck

	testContract := "/test.go"

	err = WriteFile(tmpDir, testContract, `
package main

type A struct{
	foundation.BaseContract
}

func New() (*A, error) {
    return &A{}, nil
}

func NewWithNumber(i int) (*A, error) {
    return &A{}, nil
}
`)
	s.NoError(err)

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeBuiltin)
	s.NoError(err)

	var bufWrapper bytes.Buffer
	err = parsed.WriteWrapper(&bufWrapper, parsed.ContractName())
	s.NoError(err)

	str := bufWrapper.String()
	s.Contains(str, "INSCONSTRUCTOR_New(")
	s.Contains(str, "INSCONSTRUCTOR_NewWithNumber(")
}

func (s *PreprocessorSuite) TestContractOnlyIfEmbedBaseContract() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir) //nolint: errcheck

	testContract := "/test.go"

	err = WriteFile(tmpDir, testContract, `
package main
// A contains object of type foundation.BaseContract, but it must embed it
type A struct{
	tt foundation.BaseContract
}
`)
	s.NoError(err)

	_, err = ParseFile(tmpDir+testContract, insolar.MachineTypeBuiltin)
	s.EqualError(err, "Only one smart contract must exist")
}

func (s *PreprocessorSuite) TestOnlyOneSmartContractMustExist() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir) //nolint: errcheck

	testContract := "/test.go"

	err = WriteFile(tmpDir, testContract, `
package main

type A struct{
	foundation.BaseContract
}

type B struct{
	foundation.BaseContract
}
`)
	s.NoError(err)

	_, err = ParseFile(tmpDir+testContract, insolar.MachineTypeBuiltin)
	s.EqualError(err, ": more than one contract in a file")
}

func (s *PreprocessorSuite) TestImportsFromContract() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = WriteFile(tmpDir, testContract, `
package main
import (
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"some/test/import/path"
	"some/test/import/pointerPath"
)

type A struct{
	foundation.BaseContract
}

func ( A ) Get(i path.SomeType) error {
	return nil
}

func ( A ) GetPointer(i *pointerPath.SomeType) error {
	return nil
}
`)
	s.NoError(err)

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeBuiltin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(gen.Reference().String(), &bufProxy)
	s.NoError(err)
	s.Contains(bufProxy.String(), `"some/test/import/path"`)
	s.Contains(bufProxy.String(), `"some/test/import/pointerPath"`)
	s.Contains(bufProxy.String(), `"github.com/insolar/insolar/logicrunner/common"`)
	code, err := ioutil.ReadAll(&bufProxy)
	s.NoError(err)
	s.NotEqual(len(code), 0)

	var bufWrapper bytes.Buffer
	err = parsed.WriteWrapper(&bufWrapper, parsed.ContractName())
	s.NoError(err)
	s.Contains(bufWrapper.String(), `"some/test/import/path"`)
	s.Contains(bufWrapper.String(), `"some/test/import/pointerPath"`)
}

func (s *PreprocessorSuite) TestAliasImportsFromContract() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = WriteFile(tmpDir, testContract, `
package main
import (
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	someAlias "some/test/import/path"
)

type A struct{
	foundation.BaseContract
}

func ( A ) Get(i someAlias.SomeType) error {
	return nil
}
`)
	s.NoError(err)

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeBuiltin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(gen.Reference().String(), &bufProxy)
	s.NoError(err)
	s.Contains(bufProxy.String(), `someAlias "some/test/import/path"`)
	s.Contains(bufProxy.String(), `"github.com/insolar/insolar/logicrunner/common"`)
	code, err := ioutil.ReadAll(&bufProxy)
	s.NoError(err)
	s.NotEqual(len(code), 0)

	var bufWrapper bytes.Buffer
	err = parsed.WriteWrapper(&bufWrapper, parsed.ContractName())
	s.NoError(err)
	s.Contains(bufWrapper.String(), `someAlias "some/test/import/path"`)
	s.NotContains(bufProxy.String(), `"github.com/insolar/insolar/logicrunner/common"`)
}

func (s *PreprocessorSuite) TestImportsFromContractUseInsideFunc() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = WriteFile(tmpDir, testContract, `
package main
import (
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"some/test/import/path"
)

type A struct{
	foundation.BaseContract
}

func ( A ) Get() error {
	path.SomeMethod()
	return nil
}
`)
	s.NoError(err)

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeBuiltin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(gen.Reference().String(), &bufProxy)
	s.NoError(err)
	s.NotContains(bufProxy.String(), `"some/test/import/path"`)
	code, err := ioutil.ReadAll(&bufProxy)
	s.NoError(err)
	s.NotEqual(len(code), 0)

	var bufWrapper bytes.Buffer
	err = parsed.WriteWrapper(&bufWrapper, parsed.ContractName())
	s.NoError(err)
	s.NotContains(bufWrapper.String(), `"some/test/import/path"`)
}

func (s *PreprocessorSuite) TestImportsFromContractUseForReturnValue() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = WriteFile(tmpDir, testContract, `
package main
import (
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"some/test/import/path"
)
type A struct{
	foundation.BaseContract
}
func ( A ) Get() (path.SomeValue, error) {
	f := path.SomeMethod()
	return f, nil
}
`)
	s.NoError(err)

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeBuiltin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(gen.Reference().String(), &bufProxy)
	s.NoError(err)
	s.Contains(bufProxy.String(), `"some/test/import/path"`)
	code, err := ioutil.ReadAll(&bufProxy)
	s.NoError(err)
	s.NotEqual(len(code), 0)

	var bufWrapper bytes.Buffer
	err = parsed.WriteWrapper(&bufWrapper, parsed.ContractName())
	s.NoError(err)
	s.Contains(bufWrapper.String(), `"some/test/import/path"`)
}

func (s *PreprocessorSuite) TestNotMatchFileNameForProxy() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test_not_go_file.test"
	err = WriteFile(tmpDir, testContract, `
package main

type A struct{
	foundation.BaseContract
}
`)
	s.NoError(err)

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeBuiltin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(gen.Reference().String(), &bufProxy)
	s.EqualError(err, "couldn't match filename without extension and path")
}

func (s *PreprocessorSuite) TestProxyGeneration() {
	s.T().Skip()

	contracts, err := GetRealContractsNames()
	s.Require().NoError(err)

	contractDir, err := GetRealApplicationDir("contract")
	s.Require().NoError(err)

	for _, contract := range contracts {
		// Make a copy for proper work of closure inside goroutine
		contract := contract

		s.T().Run(contract, func(t *testing.T) {
			if useLeakTest {
				defer testutils.LeakTester(t)
			} else {
				t.Parallel()
			}
			a, r := assert.New(t), require.New(t)

			parsed, err := ParseFile(path.Join(contractDir, contract, contract+".go"), insolar.MachineTypeBuiltin)
			a.NotNil(parsed, "have parsed object")
			a.NoError(err)

			proxyPath, err := GetRealApplicationDir("proxy")
			a.NoError(err)

			name, err := parsed.ProxyPackageName()
			a.NoError(err)

			proxy := path.Join(proxyPath, name, name+".go")
			_, err = os.Stat(proxy)
			a.NoError(err)

			buff := bytes.NewBufferString("")
			err = parsed.WriteProxy("", buff)
			r.NoError(err)

			cmd := exec.Command("diff", "-u", proxy, "-")
			cmd.Stdin = buff
			out, err := cmd.CombinedOutput()
			a.NoError(err, string(out))
		})
	}
}

func TestPreprocessor(t *testing.T) {
	if useLeakTest {
		defer testutils.LeakTester(t)
	} else {
		t.Parallel()
	}
	suite.Run(t, new(PreprocessorSuite))
}

// WriteFile dumps `text` into file named `name` into directory `dir`.
// Creates directory if needed as well as file
func WriteFile(dir string, name string, text string) error {
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(dir, name), []byte(text), 0644)
}
