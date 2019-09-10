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

// +build functest
// +build bloattest

package functest

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoWaitCall(t *testing.T) {
	var contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/builtin/foundation"
import "github.com/insolar/insolar/insolar"
import two "github.com/insolar/insolar/application/proxy/basic_notification_call_two"

type One struct {
	foundation.BaseContract
	Friend insolar.Reference
}

func New() (*One, error) {
	return &One{}, nil
}

var INSATTR_Hello_API = true
func (r *One) Hello() error {
	holder := two.New()

	friend, err := holder.AsChild(r.GetReference())
	if err != nil {
		return err
	}

	r.Friend = friend.GetReference()

	err = friend.MultiplyNoWait()
	if err != nil {
		return err
	}

	return nil
}

var INSATTR_Value_API = true
func (r *One) Value() (int, error) {
	return two.GetObject(r.Friend).GetValue()
}
`

	var contractTwoCode = `
package main

import (
	"fmt"

	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

type Two struct {
	foundation.BaseContract
	X int
}

func New() (*Two, error) {
	return &Two{X:322}, nil
}

var INSATTR_Multiply_API = true
func (r *Two) Multiply() (string, error) {
	r.X *= 2
	return fmt.Sprintf("Hello %d times!", r.X), nil
}

var INSATTR_GetValue_API = true
func (r *Two) GetValue() (int, error) {
	return r.X, nil
}
`
	uploadContractOnce(t, "basic_notification_call_two", contractTwoCode)
	obj := callConstructor(t, uploadContractOnce(t, "basic_notification_call_one", contractOneCode), "New")

	resp := callMethodNoChecks(t, obj, "Hello")
	require.NotEmpty(t, resp.Error)
	require.Contains(t, resp.Error.Error(), "reason request is not closed for a detached call")
}

// Make sure that panic() in a contract causes a system error and that this error
// is returned by API.
func TestPanic(t *testing.T) {
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
`
	prototype := uploadContractOnce(t, "panic", panicContractCode)
	obj := callConstructor(t, prototype, "New")

	resp := callMethodNoChecks(t, obj, "Panic")
	require.Contains(t, resp.Error.Message, "executor error: problem with API call: AAAAAAAA!")
}

func TestRecursiveCallError(t *testing.T) {
	var contractOneCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	recursive "github.com/insolar/insolar/application/proxy/recursive_call_one"
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
	testContract := `
package main

import (
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	first "github.com/insolar/insolar/application/proxy/prototype_mismatch_first"
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
