// Copyright 2020 Insolar Network Ltd.
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

// +build slowtest

package preprocessor

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
	"github.com/insolar/insolar/testutils"
)

type SagasSuite struct {
	suite.Suite
}

var sagaTestContract = `
package main

import (
"fmt"
"errors"

"github.com/insolar/insolar/logicrunner/builtin/foundation"
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

func (w *SagaTestWallet) RegularMethodAfterRollback(amount int) error {
	w.Amount *= 2; // lol
}
`

// Make sure proxy doesn't contain:
// 1. Rollback method of the saga
// 2. AsImmutable-versions of Accept/Rollback methods
func (s *SagasSuite) TestSagaAdditionalMethodsAreMissingInProxy() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = goplugintestutils.WriteFile(tmpDir, testContract, sagaTestContract)
	s.NoError(err)

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
	s.NoError(err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(gen.Reference().String(), &bufProxy)
	s.NoError(err)
	proxyCode := bufProxy.String()

	s.Contains(proxyCode, "TheAcceptMethod")
	s.NotContains(proxyCode, "TheRollbackMethod")
	s.NotContains(proxyCode, "TheAcceptMethodAsImmutable")
	s.NotContains(proxyCode, "TheRollbackMethodAsImmutable")

	// Make sure that a regular method after the rollback method was processes as usual
	s.Contains(proxyCode, "RegularMethodAfterRollback")
	s.Contains(proxyCode, "RegularMethodAfterRollbackAsImmutable")
}

// Make sure wrapper contains meta information about saga
func (s *SagasSuite) TestSagaMetaInfoIsPresentInWrapper() {
	tmpDir, err := ioutil.TempDir("", "test-")
	s.NoError(err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = goplugintestutils.WriteFile(tmpDir, testContract, sagaTestContract)
	s.NoError(err)

	parsed, err := ParseFile(tmpDir+testContract, insolar.MachineTypeGoPlugin)
	s.NoError(err)

	var bufWrapper bytes.Buffer
	err = parsed.WriteWrapper(&bufWrapper, parsed.ContractName())
	s.NoError(err)
	wrapperCode := bufWrapper.String()
	s.Contains(wrapperCode, "INSMETHOD_TheAcceptMethod")
	s.Contains(wrapperCode, "INSMETHOD_TheRollbackMethod")
	s.Contains(wrapperCode, `
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
func (s *SagasSuite) TestSagaDoesntCompileWhenRollbackIsMissing() {
	var testSaga = `
package main

import (
"fmt"
"errors"

"github.com/insolar/insolar/logicrunner/builtin/foundation"
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

	err = parsed.WriteProxy(gen.Reference().String(), &bufProxy)
	s.Error(err)
	s.Equal("Semantic error: 'TheAcceptMethod' is a saga with rollback method 'TheRollbackMethod', "+
		"but 'TheRollbackMethod' is not declared. Maybe a typo?", err.Error())
}

// Make sure saga's accept method can't be marked as immutable
func (s *SagasSuite) TestSagaAcceptMethodCantBeImmutable() {
	var testSaga = `
package main

import (
"fmt"
"errors"

"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

type SagaTestWallet struct {
	foundation.BaseContract
	Amount int
}

//ins:immutable
//ins:saga(TheRollbackMethod)
func (w *SagaTestWallet) TheAcceptMethod(amount int) error {
	w.Amount += amount
}

func (w *SagaTestWallet) TheRollbackMethod(amount int) error {
	w.Amount -= amount
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
	s.Equal("semantic error: 'TheAcceptMethod' can't be a saga because it's immutable", err.Error())

	err = parsed.WriteProxy(gen.Reference().String(), &bufProxy)
	s.Error(err)
	s.Equal("semantic error: 'TheAcceptMethod' can't be a saga because it's immutable", err.Error())
}

// Make sure saga's rollback method can't be marked as immutable
func (s *SagasSuite) TestSagaRollbackMethodCantBeImmutable() {
	var testSaga = `
package main

import (
"fmt"
"errors"

"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

type SagaTestWallet struct {
	foundation.BaseContract
	Amount int
}

//ins:saga(TheRollbackMethod)
func (w *SagaTestWallet) TheAcceptMethod(amount int) error {
	w.Amount += amount
}

//ins:immutable
func (w *SagaTestWallet) TheRollbackMethod(amount int) error {
	w.Amount -= amount
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
	s.Equal("semantic error: 'TheRollbackMethod' can't be saga's rollback method because it's immutable", err.Error())

	err = parsed.WriteProxy(gen.Reference().String(), &bufProxy)
	s.Error(err)
	s.Equal("semantic error: 'TheRollbackMethod' can't be saga's rollback method because it's immutable", err.Error())
}

// Make sure saga's accept method can't return anything but single `error` value
func (s *SagasSuite) TestSagaAcceptMethodReturnsOnlySingleErrorValue() {
	var testSaga = `
package main

import (
"fmt"
"errors"

"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

type SagaTestWallet struct {
	foundation.BaseContract
	Amount int
}

//ins:saga(TheRollbackMethod)
func (w *SagaTestWallet) TheAcceptMethod(amount int) (string, error) {
	w.Amount += amount
	return "lol", nil
}

func (w *SagaTestWallet) TheRollbackMethod(amount int) error {
	w.Amount -= amount
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
	s.Equal("semantic error: 'TheAcceptMethod' is a saga accept method and thus should return only a single `error` value", err.Error())

	err = parsed.WriteProxy(gen.Reference().String(), &bufProxy)
	s.Error(err)
	s.Equal("semantic error: 'TheAcceptMethod' is a saga accept method and thus should return only a single `error` value", err.Error())
}

// Make sure saga's rollback method can't return anything but single `error` value
func (s *SagasSuite) TestSagaRollbackMethodReturnsOnlySingleErrorValue() {
	var testSaga = `
package main

import (
"fmt"
"errors"

"github.com/insolar/insolar/logicrunner/builtin/foundation"
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

func (w *SagaTestWallet) TheRollbackMethod(amount int) (string, error) {
	w.Amount -= amount
	return "lol", nil
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
	s.Equal("semantic error: 'TheRollbackMethod' is a saga rollback method and thus should return only a single `error` value", err.Error())

	err = parsed.WriteProxy(gen.Reference().String(), &bufProxy)
	s.Error(err)
	s.Equal("semantic error: 'TheRollbackMethod' is a saga rollback method and thus should return only a single `error` value", err.Error())
}

// Make sure saga compiles without a rollback method if special flag is used
func (s *SagasSuite) TestSagaCompilesWhenRollbackIsMissingSpecialFlag() {
	var testSaga = `
package main

import (
"fmt"
"errors"

"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

type SagaTestWallet struct {
	foundation.BaseContract
	Amount int
}

//ins:saga(INS_FLAG_NO_ROLLBACK_METHOD)
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

	var bufWrapper bytes.Buffer
	err = parsed.WriteWrapper(&bufWrapper, parsed.ContractName())
	s.NoError(err)
	wrapperCode := bufWrapper.String()

	s.Contains(wrapperCode, "TheAcceptMethod")
	s.NotContains(wrapperCode, "TheAcceptMethodAsImmutable")

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(gen.Reference().String(), &bufProxy)
	s.NoError(err)
	proxyCode := bufProxy.String()

	s.Contains(proxyCode, "TheAcceptMethod")
	s.NotContains(proxyCode, "TheAcceptMethodAsImmutable")
}

// Make sure saga doesn't compile if the accept method has more then one argument
func (s *SagasSuite) TestSagaDoesntCompileWhenAcceptHasMultipleArguments() {
	var testSaga = `
package main

import (
"fmt"
"errors"

"github.com/insolar/insolar/logicrunner/builtin/foundation"
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
	s.Equal("Semantic error: 'TheAcceptMethod' is a saga with 2 arguments. Currently only one argument is allowed (hint: use a structure).",
		err.Error())

	err = parsed.WriteProxy(gen.Reference().String(), &bufProxy)
	s.Error(err)
	s.Equal("Semantic error: 'TheAcceptMethod' is a saga with 2 arguments. Currently only one argument is allowed (hint: use a structure).",
		err.Error())
}

// Make sure saga doesn't compile when saga's rollback method has arguments that don't match
func (s *SagasSuite) TestSagaDoesntCompileWhenRollbackArgumentsDontMatch() {
	var testSaga = `
package main

import (
"fmt"
"errors"

"github.com/insolar/insolar/logicrunner/builtin/foundation"
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

	err = parsed.WriteProxy(gen.Reference().String(), &bufProxy)
	s.Error(err)
	s.Equal("Semantic error: 'TheAcceptMethod' is a saga with arguments 'amount int' and rollback method "+
		"'TheRollbackMethod', but 'TheRollbackMethod' arguments 'amount string' dont't match. "+
		"They should be exactly the same.", err.Error())
}

// Make sure it's impossible to make a saga constructor.
// This case is implicitly forbidden by the specification.
func (s *SagasSuite) TestItsImpossibleToMakeASagaConstructor() {
	var testSaga = `
package main

import (
"github.com/insolar/insolar/insolar"
"github.com/insolar/insolar/logicrunner/builtin/foundation"
"github.com/insolar/insolar/application/proxy/test_saga_simple_contract"
)

type TestSagaSimpleCallContract struct {
	foundation.BaseContract
	Friend insolar.Reference
	Amount int
}

//ins:saga(GetBalance)
func New() (*TestSagaSimpleCallContract, error) {
	return &TestSagaSimpleCallContract{Amount: 100}, nil
}

func (w *TestSagaSimpleCallContract) GetBalance() (int, error) {
	return w.Amount, nil
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
	s.Equal("semantic error: 'New' can't be a saga because it's a constructor", err.Error())

	err = parsed.WriteProxy(gen.Reference().String(), &bufProxy)
	s.Error(err)
	s.Equal("semantic error: 'New' can't be a saga because it's a constructor", err.Error())
}

// Low-level tests for extractSagaInfoFromComment procedure
func (s *SagasSuite) TestExtractSagaInfoFromComment() {
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

	// Spaces after '//' - IDE of a contractor author may force this formatting
	res = extractSagaInfoFromComment("// ins:saga(SomeRollbackMethodName) ", info)
	s.Require().True(res)
	s.Require().True(info.IsSaga)
	s.Require().Equal(info.RollbackMethodName, "SomeRollbackMethodName")
}

func TestSagas(t *testing.T) {
	if useLeakTest {
		defer testutils.LeakTester(t)
	} else {
		t.Parallel()
	}
	suite.Run(t, new(SagasSuite))
}
