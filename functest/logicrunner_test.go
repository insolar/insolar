///
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
///

// +build functest

package functest

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
)

func TestSingleContract(t *testing.T) {
	var contractCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type One struct {
	foundation.BaseContract
	Number int
}

func New() (*One, error) {
	return &One{}, nil
}

func (c *One) Inc() (int, error) {
	c.Number++
	return c.Number, nil
}

func (c *One) Get() (int, error) {
	return c.Number, nil
}

func (c *One) Dec() (int, error) {
	c.Number--
	return c.Number, nil
}
`
	objectRef := callConstructor(t, uploadContractOnce(t, "test", contractCode), "New")

	// be careful - jsonUnmarshal convert json numbers to float64
	result := callMethod(t, objectRef, "Get")
	require.Empty(t, result.Error)
	require.Equal(t, float64(0), result.ExtractedReply)

	result = callMethod(t, objectRef, "Inc")
	require.Empty(t, result.Error)
	require.Equal(t, float64(1), result.ExtractedReply)

	result = callMethod(t, objectRef, "Get")
	require.Empty(t, result.Error)
	require.Equal(t, float64(1), result.ExtractedReply)

	result = callMethod(t, objectRef, "Dec")
	require.Empty(t, result.Error)
	require.Equal(t, float64(0), result.ExtractedReply)

	result = callMethod(t, objectRef, "Get")
	require.Empty(t, result.Error)
	require.Equal(t, float64(0), result.ExtractedReply)
}

func TestContractCallingContract(t *testing.T) {
	var contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "github.com/insolar/insolar/application/proxy/two"
import "github.com/insolar/insolar/insolar"
import "errors"

type One struct {
	foundation.BaseContract
	Friend insolar.Reference
}

func New() (*One, error) {
	return &One{}, nil
}

func (r *One) Hello(s string) (string, error) {
	holder := two.New()
	friend, err := holder.AsChild(r.GetReference())
	if err != nil {
		return "1", err
	}

	res, err := friend.Hello(s)
	if err != nil {
		return "2", err
	}

	r.Friend = friend.GetReference()
	return "Hi, " + s + "! Two said: " + res, nil
}

func (r *One) Again(s string) (string, error) {
	res, err := two.GetObject(r.Friend).Hello(s)
	if err != nil {
		return "", err
	}

	return "Hi, " + s + "! Two said: " + res, nil
}

func (r *One)GetFriend() (string, error) {
	return r.Friend.String(), nil
}

func (r *One) TestPayload() (two.Payload, error) {
	f := two.GetObject(r.Friend)
	err := f.SetPayload(two.Payload{Int: 10, Str: "HiHere"})
	if err != nil { return two.Payload{}, err }

	p, err := f.GetPayload()
	if err != nil { return two.Payload{}, err }

	str, err := f.GetPayloadString()
	if err != nil { return two.Payload{}, err }

	if p.Str != str { return two.Payload{}, errors.New("Oops") }

	return p, nil

}

`

	var contractTwoCode = `
package main

import (
	"fmt"

	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Two struct {
	foundation.BaseContract
	X int
	P Payload
}

type Payload struct {
	Int int
	Str string
}

func New() (*Two, error) {
	return &Two{X:0}, nil
}

func (r *Two) Hello(s string) (string, error) {
	r.X ++
	return fmt.Sprintf("Hello you too, %s. %d times!", s, r.X), nil
}

func (r *Two) GetPayload() (Payload, error) {
	return r.P, nil
}

func (r *Two) SetPayload(P Payload) (error) {
	r.P = P
	return nil
}

func (r *Two) GetPayloadString() (string, error) {
	return r.P.Str, nil
}
`

	uploadContractOnce(t, "two", contractTwoCode)
	objectRef := callConstructor(t, uploadContractOnce(t, "one", contractOneCode), "New")

	resp := callMethod(t, objectRef, "Hello", "ins")
	require.Empty(t, resp.Error)
	require.Equal(t, "Hi, ins! Two said: Hello you too, ins. 1 times!", resp.ExtractedReply)

	for i := 2; i <= 5; i++ {
		resp = callMethod(t, objectRef, "Again", "ins")
		require.Empty(t, resp.Error)
		require.Equal(t, fmt.Sprintf("Hi, ins! Two said: Hello you too, ins. %d times!", i), resp.ExtractedReply)
	}

	resp = callMethod(t, objectRef, "GetFriend")
	require.Empty(t, resp.Error)

	two, err2 := insolar.NewReferenceFromBase58(resp.ExtractedReply.(string))
	require.NoError(t, err2)

	for i := 6; i <= 9; i++ {
		resp = callMethod(t, two, "Hello", "Insolar")
		require.Empty(t, resp.Error)
		require.Equal(t, fmt.Sprintf("Hello you too, Insolar. %d times!", i), resp.ExtractedReply)
	}

	type Payload struct {
		Int int
		Str string
	}

	expected := []interface{}{Payload{Int: 10, Str: "HiHere"}, nil}

	resp = callMethod(t, objectRef, "TestPayload")
	require.Equal(
		t,
		goplugintestutils.CBORMarshal(t, expected),
		resp.Reply.Result,
	)
}

// Make sure a contract can make a saga call to another contract
func TestSagaSimpleCall(t *testing.T) {
	balance := float64(100)
	amount := float64(10)
	var contractCode = `
package main

import (
"github.com/insolar/insolar/insolar"
"github.com/insolar/insolar/logicrunner/goplugin/foundation"
"github.com/insolar/insolar/application/proxy/test_saga_simple_contract"
)

type TestSagaSimpleCallContract struct {
	foundation.BaseContract
	Friend insolar.Reference
	Amount int
}

func New() (*TestSagaSimpleCallContract, error) {
	return &TestSagaSimpleCallContract{Amount: 100}, nil
}

func (r *TestSagaSimpleCallContract) Transfer(n int) (string, error) {
	second := test_saga_simple_contract.New()
	w2, err := second.AsChild(r.GetReference())
	if err != nil {
		return "1", err
	}

	r.Amount -= n

	err = w2.Accept(n)
	if err != nil {
		return "2", err
	}
	return w2.GetReference().String(), nil
}

func (w *TestSagaSimpleCallContract) GetBalance() (int, error) {
	return w.Amount, nil
}

//ins:saga(Rollback)
func (w *TestSagaSimpleCallContract) Accept(amount int) error {
	w.Amount += amount
	return nil
}

func (w *TestSagaSimpleCallContract) Rollback(amount int) error {
	w.Amount -= amount
	return nil
}
`
	prototype := uploadContractOnce(t, "test_saga_simple_contract", contractCode)
	firstWalletRef := callConstructor(t, prototype, "New")
	resp := callMethod(t, firstWalletRef, "Transfer", int(amount))
	require.Empty(t, resp.Error)

	secondWalletRef, err := insolar.NewReferenceFromBase58(resp.ExtractedReply.(string))
	require.NoError(t, err)

	checkPassed := false

	for attempt := 0; attempt <= 10; attempt++ {
		bal2 := callMethod(t, secondWalletRef, "GetBalance")
		require.Empty(t, bal2.Error)
		if bal2.ExtractedReply.(float64) != balance+amount {
			// money are not accepted yet
			time.Sleep(10 * time.Millisecond)
			continue
		}

		bal1 := callMethod(t, firstWalletRef, "GetBalance")
		require.Empty(t, bal1.Error)
		require.Equal(t, balance-amount, bal1.ExtractedReply.(float64))
		require.Equal(t, balance+amount, bal2.ExtractedReply.(float64))

		checkPassed = true
		break
	}

	require.True(t, checkPassed)
}

// Make sure a contract can make a saga call to itself
func TestSagaSelfCall(t *testing.T) {
	var contractCode = `
package main

import (
"github.com/insolar/insolar/logicrunner/goplugin/foundation"
"github.com/insolar/insolar/application/proxy/test_saga_self_contract"
)

type TestSagaSelfCallContract struct {
	foundation.BaseContract
	SagaCallsNum int
}

func New() (*TestSagaSelfCallContract, error) {
	return &TestSagaSelfCallContract{SagaCallsNum: 0}, nil
}

func (c *TestSagaSelfCallContract) Transfer(delta int) error {
	proxy := test_saga_self_contract.GetObject(c.GetReference())
	err := proxy.Accept(delta)
	if err != nil {
		return err
	}
	return nil
}

func (c *TestSagaSelfCallContract) GetSagaCallsNum() (int, error) {
	return c.SagaCallsNum, nil
}

//ins:saga(Rollback)
func (c *TestSagaSelfCallContract) Accept(delta int) error {
	c.SagaCallsNum += delta
	return nil
}

func (c *TestSagaSelfCallContract) Rollback(delta int) error {
	c.SagaCallsNum -= delta
	return nil
}
`
	prototype := uploadContractOnce(t, "test_saga_self_contract", contractCode)
	contractRef := callConstructor(t, prototype, "New")

	res1 := callMethod(t, contractRef, "GetSagaCallsNum")
	require.Empty(t, res1.Error)
	require.Equal(t, float64(0), res1.ExtractedReply.(float64))

	resp := callMethod(t, contractRef, "Transfer", 1)
	require.Empty(t, resp.Error)

	checkPassed := false

	for attempt := 0; attempt <= 10; attempt++ {
		res2 := callMethod(t, contractRef, "GetSagaCallsNum")
		require.Empty(t, res2.Error)
		if res2.ExtractedReply.(float64) != float64(1) {
			// saga is not accepted yet
			time.Sleep(10 * time.Millisecond)
			continue
		}
		checkPassed = true
		break
	}

	require.True(t, checkPassed)
}

func TestNoWaitCall(t *testing.T) {
	var contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "github.com/insolar/insolar/insolar"
import two "github.com/insolar/insolar/application/proxy/basic_notification_call_two"

type One struct {
	foundation.BaseContract
	Friend insolar.Reference
}

func New() (*One, error) {
	return &One{}, nil
}

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

func (r *One) Value() (int, error) {
	return two.GetObject(r.Friend).GetValue()
}
`

	var contractTwoCode = `
package main

import (
	"fmt"

	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Two struct {
	foundation.BaseContract
	X int
}

func New() (*Two, error) {
	return &Two{X:322}, nil
}

func (r *Two) Multiply() (string, error) {
	r.X *= 2
	return fmt.Sprintf("Hello %d times!", r.X), nil
}

func (r *Two) GetValue() (int, error) {
	return r.X, nil
}
`
	uploadContractOnce(t, "basic_notification_call_two", contractTwoCode)
	obj := callConstructor(t, uploadContractOnce(t, "basic_notification_call_one", contractOneCode), "New")

	resp := callMethod(t, obj, "Hello")
	require.Empty(t, resp.Error)

	for i := 0; i < 25; i++ {
		resp = callMethod(t, obj, "Value")
		require.Empty(t, resp.Error)

		if float64(322) != resp.ExtractedReply {
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}

	require.Equal(t, float64(644), resp.ExtractedReply)
}

func TestContextPassing(t *testing.T) {
	var contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

func (r *One) Hello() (string, error) {
	return r.GetPrototype().String(), nil
}
`
	prototype := uploadContractOnce(t, "context_passing", contractOneCode)
	obj := callConstructor(t, prototype, "New")

	resp := callMethod(t, obj, "Hello")
	require.Empty(t, resp.Error)
	require.Equal(t, prototype.String(), resp.ExtractedReply)
}

func TestDeactivation(t *testing.T) {
	var contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

func (r *One) Kill() error {
	r.SelfDestruct()
	return nil
}
`

	obj := callConstructor(t, uploadContractOnce(t, "deactivation", contractOneCode), "New")

	resp := callMethod(t, obj, "Kill")
	require.Empty(t, resp.Error)
}

// Make sure that panic() in a contract causes a system error and that this error
// is returned by API.
func TestPanic(t *testing.T) {
	var panicContractCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

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

func TestErrorInterface(t *testing.T) {
	var contractOneCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	two "github.com/insolar/insolar/application/proxy/error_interface_two"
)

type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

func (r *One) AnError() error {
	holder := two.New()
	friend, err := holder.AsChild(r.GetReference())
	if err != nil {
		return err
	}

	return friend.AnError()
}

func (r *One) NoError() error {
	holder := two.New()
	friend, err := holder.AsChild(r.GetReference())
	if err != nil {
		return err
	}

	return friend.NoError()
}
`

	var contractTwoCode = `
package main

import (
	"errors"

	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Two struct {
	foundation.BaseContract
}
func New() (*Two, error) {
	return &Two{}, nil
}
func (r *Two) AnError() error {
	return errors.New("an error")
}
func (r *Two) NoError() error {
	return nil
}
`
	uploadContractOnce(t, "error_interface_two", contractTwoCode)
	obj := callConstructor(t, uploadContractOnce(t, "error_interface_one", contractOneCode), "New")

	resp := callMethod(t, obj, "AnError")
	require.Equal(t, "an error", resp.ExtractedError)

	resp = callMethod(t, obj, "NoError")
	require.Nil(t, resp.Error)
}

func TestNilResult(t *testing.T) {
	var contractOneCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	two "github.com/insolar/insolar/application/proxy/nil_result_two"
)

type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

func (r *One) Hello() (*string, error) {
	holder := two.New()
	friend, err := holder.AsChild(r.GetReference())
	if err != nil {
		return nil, err
	}

	return friend.Hello()
}
`

	var contractTwoCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Two struct {
	foundation.BaseContract
}
func New() (*Two, error) {
	return &Two{}, nil
}
func (r *Two) Hello() (*string, error) {
	return nil, nil
}
`

	uploadContractOnce(t, "nil_result_two", contractTwoCode)
	obj := callConstructor(t, uploadContractOnce(t, "nil_result_one", contractOneCode), "New")

	resp := callMethod(t, obj, "Hello")
	require.Empty(t, resp.Error)
	require.Nil(t, resp.ExtractedReply)
}

func TestConstructorReturnNil(t *testing.T) {
	var contractOneCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	two "github.com/insolar/insolar/application/proxy/constructor_return_nil_two"
)

type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

func (r *One) Hello() (*string, error) {
	holder := two.New()
	_, err := holder.AsChild(r.GetReference())
	if err != nil {
		return nil, err
	}
	ok := "all was well"
	return &ok, nil
}
`

	var contractTwoCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Two struct {
	foundation.BaseContract
}
func New() (*Two, error) {
	return nil, nil
}
`
	uploadContractOnce(t, "constructor_return_nil_two", contractTwoCode)
	obj := callConstructor(t, uploadContractOnce(t, "constructor_return_nil_one", contractOneCode), "New")

	resp := callMethodExpectError(t, obj, "Hello")
	require.Empty(t, resp.Reply.Result)
}

func TestRecursiveCallError(t *testing.T) {
	var contractOneCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	recursive "github.com/insolar/insolar/application/proxy/recursive_call_one"
)
type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

func (r *One) Recursive() (error) {
	remoteSelf := recursive.GetObject(r.GetReference())
	err := remoteSelf.Recursive()
	return err
}

`
	protoRef := uploadContractOnce(t, "recursive_call_one", contractOneCode)

	obj := callConstructor(t, protoRef, "New")
	resp := callMethodNoChecks(t, obj, "Recursive")

	errstr := resp.Error.Error()
	require.NotEmpty(t, errstr)
	// if you get a timeout here add a retry loop to the test as it was before
	require.Contains(t, errstr, "loop detected")
}

func TestGetParent(t *testing.T) {
	var contractOneCode = `
 package main

 import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
 	"github.com/insolar/insolar/insolar"
	two "github.com/insolar/insolar/application/proxy/get_parent_two"
 )

 type One struct {
	foundation.BaseContract
 }


func New() (*One, error) {
	return &One{}, nil
}

func (r *One) AddChildAndReturnMyselfAsParent() (string, error) {
	holder := two.New()
	friend, err := holder.AsChild(r.GetReference())
	if err != nil {
		return insolar.Reference{}.String(), err
	}

 	return friend.GetParent()
}
`
	var contractTwoCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Two struct {
	foundation.BaseContract
}

func New() (*Two, error) {
	return &Two{}, nil
}

func (r *Two) GetParent() (string, error) {
	return r.GetContext().Parent.String(), nil
}
`

	uploadContractOnce(t, "get_parent_two", contractTwoCode)
	obj := callConstructor(t, uploadContractOnce(t, "get_parent_one", contractOneCode), "New")

	resp := callMethod(t, obj, "AddChildAndReturnMyselfAsParent")
	require.Empty(t, resp.Error)
	require.Equal(t, obj.String(), resp.ExtractedReply)
}

// TODO need to move it into jepsen tests
func TestGinsiderMustDieAfterInsolardError(t *testing.T) {
	// can't kill LR in launch.sh from functest
}

func TestGetRemoteData(t *testing.T) {
	var contractOneCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	two "github.com/insolar/insolar/application/proxy/get_remote_data_two"
	"github.com/insolar/insolar/insolar"
)

type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

func (r *One) GetChildPrototype() (string, error) {
	holder := two.New()
	child, err := holder.AsChild(r.GetReference())
	if err != nil {
		return insolar.Reference{}.String(), err
	}

	ref, err := child.GetPrototype()
 	return ref.String(), err
}
`
	var contractTwoCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)
 
type Two struct {
	foundation.BaseContract
}

func New() (*Two, error) {
	return &Two{}, nil
}
`
	codeTwoRef := uploadContractOnce(t, "get_remote_data_two", contractTwoCode)
	obj := callConstructor(t, uploadContractOnce(t, "get_remote_data_one", contractOneCode), "New")

	resp := callMethod(t, obj, "GetChildPrototype")
	require.Empty(t, resp.Error)
	require.Equal(t, codeTwoRef.String(), resp.ExtractedReply.(string))
}

func TestNoLoopsWhileNotificationCall(t *testing.T) {
	var contractOneCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	two "github.com/insolar/insolar/application/proxy/no_loops_while_notification_call_two"
)

type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

func (r *One) IncrementBy100() (int, error) {
	holder := two.New()
	child, err := holder.AsChild(r.GetReference())
	if err != nil {
		return 0, err
	}

	for i := 0; i < 100; i++ {
		child.IncreaseNoWait()
	}

 	return child.GetCounter()
}
`
	var contractTwoCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Two struct {
	foundation.BaseContract
	Counter int
}
func New() (*Two, error) {
	return &Two{}, nil
}

func (r *Two) Increase() error {
 	r.Counter++
	return nil
}

func (r *Two) GetCounter() (int, error) {
	return r.Counter, nil
}
`
	uploadContractOnce(t, "no_loops_while_notification_call_two", contractTwoCode)
	obj := callConstructor(t, uploadContractOnce(t, "no_loops_while_notification_call_one", contractOneCode), "New")

	resp := callMethod(t, obj, "IncrementBy100")
	require.Empty(t, resp.Error)
}

func TestPrototypeMismatch(t *testing.T) {
	testContract := `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	first "github.com/insolar/insolar/application/proxy/prototype_mismatch_first"
	"github.com/insolar/insolar/insolar"
)

func New() (*Contract, error) {
	return &Contract{}, nil
}

type Contract struct {
	foundation.BaseContract
}

func (c *Contract) Test(firstRef *insolar.Reference) (string, error) {
	return first.GetObject(*firstRef).GetName()
}
`

	// right contract
	firstContract := `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type First struct {
	foundation.BaseContract
}

func (c *First) GetName() (string, error) {
	return "first", nil
}
`

	// malicious contract with same method signature and another behaviour
	secondContract := `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type First struct {
	foundation.BaseContract
}

func New() (*First, error) {
	return &First{}, nil
}

func (c *First) GetName() (string, error) {
	return "YOU ARE ROBBED!", nil
}
`

	uploadContractOnce(t, "prototype_mismatch_first", firstContract)
	secondObj := callConstructor(t, uploadContractOnce(t, "prototype_mismatch_second", secondContract), "New")
	testObj := callConstructor(t, uploadContractOnce(t, "prototype_mismatch_test", testContract), "New")

	callMethodExpectError(t, testObj, "Test", *secondObj)
}

func TestImmutableAnnotation(t *testing.T) {
	var contractOneCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	two "github.com/insolar/insolar/application/proxy/immutable_annotation_two"
)

type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

func (r *One) ExternalImmutableCall() (int, error) {
	holder := two.New()
	objTwo, err := holder.AsChild(r.GetReference())
	if err != nil {
		return 0, err
	}
	return objTwo.ReturnNumberAsImmutable()
}

func (r *One) ExternalImmutableCallMakesExternalCall() (error) {
	holder := two.New()
	objTwo, err := holder.AsChild(r.GetReference())
	if err != nil {
		return err
	}
	return objTwo.Immutable()
}
`

	var contractTwoCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	three "github.com/insolar/insolar/application/proxy/immutable_annotation_three"
)

type Two struct {
	foundation.BaseContract
}

func New() (*Two, error) {
	return &Two{}, nil
}

func (r *Two) ReturnNumber() (int, error) {
	return 42, nil
}

//ins:immutable
func (r *Two) Immutable() (error) {
	holder := three.New()
	objThree, err := holder.AsChild(r.GetReference())
	if err != nil {
		return err
	}
	return objThree.DoNothing()
}

`

	var contractThreeCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type Three struct {
	foundation.BaseContract
}

func New() (*Three, error) {
	return &Three{}, nil
}

func (r *Three) DoNothing() (error) {
	return nil
}

`

	uploadContractOnce(t, "immutable_annotation_three", contractThreeCode)
	uploadContractOnce(t, "immutable_annotation_two", contractTwoCode)
	obj := callConstructor(t, uploadContractOnce(t, "immutable_annotation_one", contractOneCode), "New")

	resp := callMethod(t, obj, "ExternalImmutableCall")
	require.Empty(t, resp.Error)
	require.Equal(t, float64(42), resp.ExtractedReply)

	resp = callMethodExpectError(t, obj, "ExternalImmutableCallMakesExternalCall")
}

func TestMultipleConstructorsCall(t *testing.T) {
	var contractCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

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

func (c *One) Get() (int, error) {
	return c.Number, nil
}
`

	prototypeRef := uploadContractOnce(t, "test_multiple_constructor", contractCode)

	objRef := callConstructor(t, prototypeRef, "New")

	// be careful - jsonUnmarshal convert json numbers to float64
	result := callMethod(t, objRef, "Get")
	require.Empty(t, result.Error)
	require.Equal(t, float64(0), result.ExtractedReply)

	objRef = callConstructor(t, prototypeRef, "NewWithNumber", 12)

	// be careful - jsonUnmarshal convert json numbers to float64
	result = callMethod(t, objRef, "Get")
	require.Empty(t, result.Error)
	require.Equal(t, float64(12), result.ExtractedReply)
}
