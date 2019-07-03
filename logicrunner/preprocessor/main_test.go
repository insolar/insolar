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
	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
	"github.com/insolar/insolar/testutils"
)

type PreprocessorSuite struct {
	suite.Suite
}

var randomTestCode = `
package main

import (
	"fmt"
	"errors"

	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
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

	err = goplugintestutils.WriteFile(tmpDir, "main.go", randomTestCode)
	s.NoError(err)

	parsed, err := ParseFile(filepath.Join(tmpDir, "main.go"), insolar.MachineTypeGoPlugin)
	s.NoError(err)
	s.NotNil(parsed)

	s.T().Run("wrapper", func(t *testing.T) {
		t.Parallel()
		a := assert.New(t)

		buf := bytes.Buffer{}
		err := parsed.WriteWrapper(&buf, parsed.ContractName())
		a.NoError(err)

		code, err := ioutil.ReadAll(&buf)
		a.NoError(err)
		a.NotEmpty(code)
	})

	s.T().Run("proxy", func(t *testing.T) {
		t.Parallel()
		a := assert.New(t)

		buf := bytes.Buffer{}
		err := parsed.WriteProxy(testutils.RandomRef().String(), &buf)
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

	err = goplugintestutils.WriteFile(tmpDir, "code1", code)
	s.NoError(err)

	info, err := ParseFile(filepath.Join(tmpDir, "code1"), insolar.MachineTypeGoPlugin)
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

	err = goplugintestutils.WriteFile(tmpDir, "code1", code)
	s.NoError(err)

	_, err = ParseFile(filepath.Join(tmpDir, "code1"), insolar.MachineTypeGoPlugin)
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

	err = goplugintestutils.WriteFile(tmpDir, "code1", code)
	s.NoError(err)

	_, err = ParseFile(filepath.Join(tmpDir, "code1"), insolar.MachineTypeGoPlugin)
	s.Error(err)
}

func (s *PreprocessorSuite) TestCompileContractProxy() {

	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir) // nolint: errcheck

	err = os.MkdirAll(filepath.Join(tmpDir, "src/secondary"), 0777)
	s.NoError(err)

	cwd, err := os.Getwd()
	s.NoError(err)

	// XXX: dirty hack to make `dep` installed packages available in generated code
	err = os.Symlink(filepath.Join(cwd, "../../../vendor"), filepath.Join(tmpDir, "src/secondary/vendor"))
	s.NoError(err)

	proxyFh, err := os.OpenFile(filepath.Join(tmpDir, "/src/secondary/main.go"), os.O_WRONLY|os.O_CREATE, 0644)
	s.NoError(err)

	err = goplugintestutils.WriteFile(filepath.Join(tmpDir, "/contracts/secondary/"), "main.go", randomTestCode)
	s.NoError(err)

	parsed, err := ParseFile(filepath.Join(tmpDir, "/contracts/secondary/main.go"), insolar.MachineTypeGoPlugin)
	s.NoError(err)

	err = parsed.WriteProxy(testutils.RandomRef().String(), proxyFh)
	s.NoError(err)

	err = proxyFh.Close()
	s.NoError(err)

	err = goplugintestutils.WriteFile(tmpDir, "/test.go", `
package test

import (
	"github.com/insolar/insolar/insolar"
	"secondary"
)

func main() {
	ref, _ := insolar.NewReferenceFromBase58("4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.7ZQboaH24PH42sqZKUvoa7UBrpuuubRtShp6CKNuWGZa")
	_ = secondary.GetObject(*ref)
}
	`)
	s.NoError(err)

	cmd := exec.Command("go", "build", filepath.Join(tmpDir, "test.go"))
	cmd.Env = append(os.Environ(), "GOPATH="+goplugintestutils.PrependGoPath(tmpDir))
	out, err := cmd.CombinedOutput()
	s.NoError(err, string(out))
}

func (s *PreprocessorSuite) TestFailIfThereAreNoContract() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir) //nolint: errcheck

	testContract := "/test.go"
	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main
type A struct{
	ttt ppp.TTT
}
`)
	s.NoError(err)

	_, err = ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
	s.EqualError(err, "Only one smart contract must exist")
}

func (s *PreprocessorSuite) TestInitializationFunctionParamsProxy() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir) // nolint: errcheck

	testContract := "/test.go"

	err = goplugintestutils.WriteFile(tmpDir, testContract, `
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

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(testutils.RandomRef().String(), &bufProxy)
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

	err = goplugintestutils.WriteFile(tmpDir, testContract, `
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

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
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

	err = goplugintestutils.WriteFile(tmpDir, testContract, `
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

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
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

	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main
// A contains object of type foundation.BaseContract, but it must embed it
type A struct{
	tt foundation.BaseContract
}
`)
	s.NoError(err)

	_, err = ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
	s.EqualError(err, "Only one smart contract must exist")
}

func (s *PreprocessorSuite) TestOnlyOneSmartContractMustExist() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir) //nolint: errcheck

	testContract := "/test.go"

	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main

type A struct{
	foundation.BaseContract
}

type B struct{
	foundation.BaseContract
}
`)
	s.NoError(err)

	_, err = ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
	s.EqualError(err, ": more than one contract in a file")
}

func (s *PreprocessorSuite) TestImportsFromContract() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main
import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
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

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(testutils.RandomRef().String(), &bufProxy)
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
	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main
import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
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

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(testutils.RandomRef().String(), &bufProxy)
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
	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main
import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
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

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(testutils.RandomRef().String(), &bufProxy)
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
	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main
import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
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

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(testutils.RandomRef().String(), &bufProxy)
	s.NoError(err)
	s.Contains(bufProxy.String(), `"some/test/import/path"`)
	code, err := ioutil.ReadAll(&bufProxy)
	s.NoError(err)
	s.NotEqual(len(code), 0)

	var bufWrapper bytes.Buffer
	err = parsed.WriteWrapper(&bufWrapper, parsed.ContractName())
	s.NoError(err)
	s.NotContains(bufWrapper.String(), `"some/test/import/path"`)
}

func (s *PreprocessorSuite) TestNotMatchFileNameForProxy() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test_not_go_file.test"
	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main

type A struct{
	foundation.BaseContract
}
`)
	s.NoError(err)

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(testutils.RandomRef().String(), &bufProxy)
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
			t.Parallel()
			a, r := assert.New(t), require.New(t)

			parsed, err := ParseFile(path.Join(contractDir, contract, contract+".go"), insolar.MachineTypeGoPlugin)
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

var sagaTestContract = `
package main

import (
"fmt"
"errors"

"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type SagaTestWallet struct {
	foundation.BaseContract
	Amount int
}

//ins:saga(TheRollbackMethod)
func (w *SagaTestWallet) TheAcceptMethod(amount int) error {
	w.Amount += amount
}

func (w *SagaTestWallet) TheRollbackMethod(amount int) error {
	w.Amount -= amount
}
`

// Make sure proxy doesn't contain:
// 1. Rollback method of the saga
// 2. AsImmutable-versions of Accept/Rollback methods
// 3. NoWait-versions of Accept/Rollback methods
func (s *PreprocessorSuite) TestSagaAdditionalMethodsAreMissingInProxy() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = goplugintestutils.WriteFile(tmpDir, testContract, sagaTestContract)
	s.NoError(err)

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(testutils.RandomRef().String(), &bufProxy)
	s.NoError(err)
	proxyCode := bufProxy.String()
	s.Contains(proxyCode, "TheAcceptMethod")
	s.NotContains(proxyCode, "TheRollbackMethod")
	s.NotContains(proxyCode, "TheAcceptMethodNoWait")
	s.NotContains(proxyCode, "TheRollbackMethodNoWait")
	s.NotContains(proxyCode, "TheAcceptMethodAsImmutable")
	s.NotContains(proxyCode, "TheRollbackMethodAsImmutable")
}

// Make sure wrapper contains meta information about saga
func (s *PreprocessorSuite) TestSagaMetaInfoIsPresentInProxy() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = goplugintestutils.WriteFile(tmpDir, testContract, sagaTestContract)
	s.NoError(err)

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteWrapper(&bufProxy, parsed.ContractName())
	s.NoError(err)
	proxyCode := bufProxy.String()
	s.Contains(proxyCode, "INSMETHOD_TheAcceptMethod")
	s.Contains(proxyCode, "INSMETHOD_TheRollbackMethod")
	s.Contains(proxyCode, `
func INS_META_INFO() []map[string]string {
	result := make([]map[string]string, 0)

	{
		info := make(map[string]string, 3)
		info["Type"] = "SagaInfo"
		info["MethodName"] = "TheAcceptMethod"
		info["RollbackMethodName"] = "TheRollbackMethod"
		result = append(result, info)
	}

	return result
}
`)
}

// Make sure saga doesn't compile when saga's rollback method doesn't exist
func (s *PreprocessorSuite) TestSagaDoesntCompileWhenRollbackIsMissing() {
	var testSaga = `
package main

import (
"fmt"
"errors"

"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type SagaTestWallet struct {
	foundation.BaseContract
	Amount int
}

//ins:saga(TheRollbackMethod)
func (w *SagaTestWallet) TheAcceptMethod(amount int) error {
	w.Amount += amount
    return nil
}
`
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = goplugintestutils.WriteFile(tmpDir, testContract, testSaga)
	s.NoError(err)

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteWrapper(&bufProxy, parsed.ContractName())
	s.Error(err)
	s.Equal("Semantic error: 'TheAcceptMethod' is a saga with rollback method 'TheRollbackMethod', "+
		"but 'TheRollbackMethod' is not declared. Maybe a typo?", err.Error())

	err = parsed.WriteProxy(testutils.RandomRef().String(), &bufProxy)
	s.Error(err)
	s.Equal("Semantic error: 'TheAcceptMethod' is a saga with rollback method 'TheRollbackMethod', "+
		"but 'TheRollbackMethod' is not declared. Maybe a typo?", err.Error())
}

// Make sure saga doesn't compile if the accept method has more then one argument
func (s *PreprocessorSuite) TestSagaDoesntCompileWhenAcceptHasMultipleArguments() {
	var testSaga = `
package main

import (
"fmt"
"errors"

"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type SagaTestWallet struct {
	foundation.BaseContract
	Amount int
}

//ins:saga(TheRollbackMethod)
func (w *SagaTestWallet) TheAcceptMethod(arg1 int, arg2 string) error {
    return nil
}

func (w *SagaTestWallet) TheRollbackMethod(arg1 int, arg2 string) error {
    return nil
}
`
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = goplugintestutils.WriteFile(tmpDir, testContract, testSaga)
	s.NoError(err)

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteWrapper(&bufProxy, parsed.ContractName())
	s.Error(err)
	s.Equal("Semantic error: 'TheAcceptMethod' is a saga with 2 arguments. Currently only one argument is allowed.",
		err.Error())

	err = parsed.WriteProxy(testutils.RandomRef().String(), &bufProxy)
	s.Error(err)
	s.Equal("Semantic error: 'TheAcceptMethod' is a saga with 2 arguments. Currently only one argument is allowed.",
		err.Error())
}

// Make sure saga doesn't compile when saga's rollback method has arguments that don't match
func (s *PreprocessorSuite) TestSagaDoesntCompileWhenRollbackArgumentsDontMatch() {
	var testSaga = `
package main

import (
"fmt"
"errors"

"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type SagaTestWallet struct {
	foundation.BaseContract
	Amount int
}

//ins:saga(TheRollbackMethod)
func (w *SagaTestWallet) TheAcceptMethod(amount int) error {
	w.Amount += amount
    return nil
}

func (w *SagaTestWallet) TheRollbackMethod(amount string) error {
	return nil
}
`
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = goplugintestutils.WriteFile(tmpDir, testContract, testSaga)
	s.NoError(err)

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteWrapper(&bufProxy, parsed.ContractName())
	s.Error(err)
	s.Equal("Semantic error: 'TheAcceptMethod' is a saga with arguments 'amount int' and rollback method "+
		"'TheRollbackMethod', but 'TheRollbackMethod' arguments 'amount string' dont't match. "+
		"They should be exactly the same.", err.Error())

	err = parsed.WriteProxy(testutils.RandomRef().String(), &bufProxy)
	s.Error(err)
	s.Equal("Semantic error: 'TheAcceptMethod' is a saga with arguments 'amount int' and rollback method "+
		"'TheRollbackMethod', but 'TheRollbackMethod' arguments 'amount string' dont't match. "+
		"They should be exactly the same.", err.Error())
}

// Low-level tests for extractSagaInfoFromComment procedure
func (s *PreprocessorSuite) TestExtractSagaInfoFromComment() {
	info := &SagaInfo{}
	res := extractSagaInfoFromComment("", info)
	s.Require().False(res)
	s.Require().False(info.IsSaga)

	res = extractSagaInfoFromComment("ololo", info)
	s.Require().False(res)
	s.Require().False(info.IsSaga)

	res = extractSagaInfoFromComment("//ins:saga()", info)
	s.Require().False(res)
	s.Require().False(info.IsSaga)

	res = extractSagaInfoFromComment("//ins:saga(SomeRollbackMethodName) ", info)
	s.Require().True(res)
	s.Require().True(info.IsSaga)
	s.Require().Equal(info.RollbackMethodName, "SomeRollbackMethodName")
}

func TestPreprocessor(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(PreprocessorSuite))
}
