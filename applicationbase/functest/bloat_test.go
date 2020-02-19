///
// Copyright 2020 Insolar Technologies GmbH
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
///

// +build functest
// +build bloattest

package functest

import (
	"strings"
	"testing"

	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/stretchr/testify/require"
)

// Make sure that panic() in a contract causes a system error and that this error
// is returned by API.
func TestPanicIsSystemError(t *testing.T) {
	launchnet.RunOnlyWithLaunchnet(t)
	var panicContractCode = `
package main

import "github.com/insolar/insolar/logicrunner/builtin/foundation"

type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

var INSATTR_Panic_API = true
func (r *One) Panic() error {
	panic("AAAAAAAA!")
	return nil
}

func NewPanic() (*One, error) {
	panic("BBBBBBBB!")
}
`
	prototype := uploadContractOnceExt(t, "panicAsSystemError", panicContractCode, false)
	obj := callConstructor(t, prototype, "New")

	resp1 := callMethodNoChecks(t, obj, "Panic")
	require.Contains(t, resp1.Error.Message, "AAAAAAAA!")

	resp2 := callConstructorNoChecks(t, prototype, "NewPanic")
	require.Contains(t, resp2.Error.Message, "BBBBBBBB!")
}

// Make sure that panic() in a contract causes a system error and that this error
// is returned by API.
func TestPanicIsLogicalError(t *testing.T) {
	launchnet.RunOnlyWithLaunchnet(t)
	var panicContractCode = `
package main

import "github.com/insolar/insolar/logicrunner/builtin/foundation"

type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

var INSATTR_Panic_API = true
func (r *One) Panic() error {
	panic("AAAAAAAA!")
	return nil
}

func NewPanic() (*One, error) {
	panic("BBBBBBBB!")
}
`
	prototype := uploadContractOnceExt(t, "panicAsLogicalError", panicContractCode, true)
	obj := callConstructor(t, prototype, "New")

	resp1 := callMethodNoChecks(t, obj, "Panic")
	require.Contains(t, resp1.Result.ExtractedError, "AAAAAAAA!")

	resp2 := callConstructorNoChecks(t, prototype, "NewPanic")
	require.Contains(t, resp2.Result.Error.S, "BBBBBBBB!")

}

func TestRecursiveCallError(t *testing.T) {
	launchnet.RunOnlyWithLaunchnet(t)
	var contractOneCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	recursive "github.com/insolar/insolar/applicationbase/proxy/recursive_call_one"
)
type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

var INSATTR_Recursive_API = true
func (r *One) Recursive() (error) {
	remoteSelf := recursive.GetObject(r.GetReference())
	err := remoteSelf.Recursive()
	return err
}

`
	protoRef := uploadContractOnce(t, "recursive_call_one", contractOneCode)

	// for now Recursive calls may cause timeouts. Dont remove retries until we make new loop detection algorithm
	var err string
	for i := 0; i <= 5; i++ {
		obj := callConstructor(t, protoRef, "New")
		resp := callMethodNoChecks(t, obj, "Recursive")

		err = resp.Error.Error()
		if !strings.Contains(err, "timeout") {
			// system error is not timeout, loop detected is in response
			err = resp.Result.ExtractedError
			break
		}
	}

	require.NotEmpty(t, err)
	require.Contains(t, err, "loop detected")
}

func TestPrototypeMismatch(t *testing.T) {
	launchnet.RunOnlyWithLaunchnet(t)
	testContract := `
package main

import (
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	first "github.com/insolar/insolar/applicationbase/proxy/prototype_mismatch_first"
	"github.com/insolar/insolar/insolar"
)

func New() (*Contract, error) {
	return &Contract{}, nil
}

type Contract struct {
	foundation.BaseContract
}

var INSATTR_Test_API = true
func (c *Contract) Test(firstRef *insolar.Reference) (string, error) {
	return first.GetObject(*firstRef).GetName()
}
`

	// right contract
	firstContract := `
package main

import (
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

type First struct {
	foundation.BaseContract
}

var INSATTR_GetName_API = true
func (c *First) GetName() (string, error) {
	return "first", nil
}
`

	// malicious contract with same method signature and another behaviour
	secondContract := `
package main

import (
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

type First struct {
	foundation.BaseContract
}

func New() (*First, error) {
	return &First{}, nil
}

var INSATTR_GetName_API = true
func (c *First) GetName() (string, error) {
	return "YOU ARE ROBBED!", nil
}
`

	uploadContractOnce(t, "prototype_mismatch_first", firstContract)
	secondObj := callConstructor(t, uploadContractOnce(t, "prototype_mismatch_second", secondContract), "New")
	testObj := callConstructor(t, uploadContractOnce(t, "prototype_mismatch_test", testContract), "New")

	resp := callMethodNoChecks(t, testObj, "Test", *secondObj)
	require.Empty(t, resp.Error)
	require.Contains(t, resp.Result.Error.S, "try to call method of prototype as method of another prototype")
}

func TestContractWithEmbeddedConstructor(t *testing.T) {
	launchnet.RunOnlyWithLaunchnet(t)
	var contractOneCode = `
package main

import ("github.com/insolar/insolar/logicrunner/builtin/foundation")

type One struct {
	foundation.BaseContract
	Number int
}

func New() (*One, error) {
	return &One{Number: 0}, nil
}

func NewWithNumber(num int) (*One, error) {
	return &One{Number: num}, nil
}

var INSATTR_Get_API = true

func (c *One) Get() (int, error) {
	return c.Number, nil
}

var INSATTR_DoNothing_API = true

func (r *One) DoNothing() (error) {
	return nil
}
`

	var contractTwoCode = `
package main

import (
    "github.com/insolar/insolar/application/appfoundation"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	one "github.com/insolar/insolar/applicationbase/proxy/parent_one"
)

type Two struct {
	foundation.BaseContract
	Number int
	OneRef insolar.Reference
}


func New() (*Two, error) {
	return &Two{Number: 10, OneRef: *insolar.NewEmptyReference()}, nil
}

func NewWithOne(oneNumber int) (*Two, error) {
	holder := one.NewWithNumber(oneNumber)

	objOne, err := holder.AsChild(appfoundation.GetRootDomain())

	if err != nil {
		return nil, err
	}

	return &Two{Number: oneNumber, OneRef: objOne.GetReference() }, nil
}

var INSATTR_DoNothing_API = true

func (r *Two) DoNothing() (error) {
	return nil
}

var INSATTR_Get_API = true

func (c * Two) Get() (int, error) {
	return c.Number, nil
}
`
	codeOneRef := uploadContractOnce(t, "parent_one", contractOneCode)
	codeTwoRef := uploadContractOnce(t, "parent_two", contractTwoCode)

	_ = callConstructor(t, codeOneRef, "New")
	errorMsg := callConstructorExpectSystemError(t, codeTwoRef, "NewWithOne", 10)
	require.Contains(t, errorMsg, "object is not activated")
}
