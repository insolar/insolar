// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/applicationbase/testutils/testrequest"
)

func TestSingleContract(t *testing.T) {
	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.New",
		map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)

	res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.Get",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
	require.Equal(t, float64(0), res.(float64))

	res, err = testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.Inc",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
	require.Equal(t, float64(1), res.(float64))

	res, err = testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.Get",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
	require.Equal(t, float64(1), res.(float64))

	res, err = testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.Dec",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
	require.Equal(t, float64(0), res.(float64))

	res, err = testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.Get",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
	require.Equal(t, float64(0), res.(float64))

}

func TestContractCallingContract(t *testing.T) {
	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.New",
		map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)

	res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.Hello",
		map[string]interface{}{"reference": ref, "name": "ins"})
	require.NoError(t, err)
	require.Equal(t, "Hi, ins! Two said: Hello you too, ins. 1 times!", res.(string))

	for i := 2; i <= 5; i++ {
		res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.Again",
			map[string]interface{}{"reference": ref, "name": "ins"})
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("Hi, ins! Two said: Hello you too, ins. %d times!", i), res.(string))
	}

	ref2, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.GetFriend",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
	for i := 6; i <= 9; i++ {
		res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "second.Hello",
			map[string]interface{}{"reference": ref2, "name": "Insolar"})
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("Hello you too, Insolar. %d times!", i), res.(string))
	}

	type Payload struct {
		Int int
		Str string
	}

	res, err = testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.TestPayload",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
	fmt.Println(res)
	bytes, err := json.Marshal(res)
	require.NoError(t, err)
	resPayload := Payload{}
	err = json.Unmarshal(bytes, &resPayload)
	require.NoError(t, err)
	require.Equal(t, Payload{Int: 10, Str: "HiHere"}, resPayload)

	res, err = testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.ManyTimes",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
	require.Nil(t, res)
}

func TestSagaSimpleCall(t *testing.T) {
	balance := float64(100)
	amount := float64(10)

	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.NewSaga",
		map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)

	res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.Transfer",
		map[string]interface{}{"reference": ref, "amount": int(amount)})
	require.NoError(t, err)
	ref2, ok := res.(string)
	require.True(t, ok)

	checkPassed := false

	for attempt := 0; attempt <= 10; attempt++ {
		res2, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.GetBalance",
			map[string]interface{}{"reference": ref2})
		require.NoError(t, err)
		if res2.(float64) != balance+amount {
			// money are not accepted yet
			time.Sleep(10 * time.Millisecond)
			continue
		}

		res1, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.GetBalance",
			map[string]interface{}{"reference": ref})
		require.NoError(t, err)
		require.Equal(t, balance-amount, res1.(float64))
		require.Equal(t, balance+amount, res2.(float64))

		checkPassed = true
		break
	}

	require.True(t, checkPassed)
}

func TestSagaCallFromSagaAcceptMethod(t *testing.T) {
	balance := float64(100)
	amount := float64(10)

	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.NewSaga",
		map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)

	res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.TransferWithRollback",
		map[string]interface{}{"reference": ref, "amount": int(amount)})
	require.NoError(t, err)
	ref2, ok := res.(string)
	require.True(t, ok)

	checkPassed := false

	for attempt := 0; attempt <= 10; attempt++ {
		res2, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.GetBalance",
			map[string]interface{}{"reference": ref2})
		require.NoError(t, err)
		if res2.(float64) != balance+amount {
			// money are not accepted yet
			time.Sleep(1 * time.Second)
			continue
		}

		res1, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.GetBalance",
			map[string]interface{}{"reference": ref})
		require.NoError(t, err)
		if res1.(float64) != balance-amount {
			// money are not accepted yet
			time.Sleep(1 * time.Second)
			continue
		}

		checkPassed = true
		break
	}

	require.True(t, checkPassed)
}

func TestSagaMultipleCalls(t *testing.T) {
	balance := float64(100)
	amount := float64(10)

	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.NewSaga",
		map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)

	res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.TransferTwice",
		map[string]interface{}{"reference": ref, "amount": int(amount)})
	require.NoError(t, err)
	ref2, ok := res.(string)
	require.True(t, ok)

	checkPassed := false

	for attempt := 0; attempt <= 10; attempt++ {
		res2, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.GetBalance",
			map[string]interface{}{"reference": ref2})
		require.NoError(t, err)
		if res2.(float64) != balance+amount {
			// money are not accepted yet
			time.Sleep(1 * time.Millisecond)
			continue
		}

		res1, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.GetBalance",
			map[string]interface{}{"reference": ref})
		require.NoError(t, err)
		require.Equal(t, balance-amount, res1.(float64))
		require.Equal(t, balance+amount, res2.(float64))

		checkPassed = true
		break
	}

	require.True(t, checkPassed)
}

func TestSagaCallBetweenContractsWithoutRollback(t *testing.T) {
	balance := float64(100)
	amount := float64(10)

	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.NewSaga",
		map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)

	res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.TransferToAnotherContract",
		map[string]interface{}{"reference": ref, "amount": int(amount)})
	require.NoError(t, err)
	ref2, ok := res.(string)
	require.True(t, ok)

	checkPassed := false

	for attempt := 0; attempt <= 10; attempt++ {
		res2, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "second.GetBalance",
			map[string]interface{}{"reference": ref2})
		require.NoError(t, err)
		if res2.(float64) != balance+amount {
			// money are not accepted yet
			time.Sleep(1 * time.Millisecond)
			continue
		}

		res1, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.GetBalance",
			map[string]interface{}{"reference": ref})
		require.NoError(t, err)
		require.Equal(t, balance-amount, res1.(float64))
		require.Equal(t, balance+amount, res2.(float64))

		checkPassed = true
		break
	}

	require.True(t, checkPassed)
}

func TestSagaSelfCall(t *testing.T) {
	amount := float64(1)

	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "third.New",
		map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)

	num, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "third.GetSagaCallsNum",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
	require.Equal(t, float64(0), num.(float64))

	_, err = testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "third.Transfer",
		map[string]interface{}{"reference": ref, "amount": int(amount)})
	require.NoError(t, err)

	checkPassed := false

	for attempt := 0; attempt <= 10; attempt++ {
		res2, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "third.GetSagaCallsNum",
			map[string]interface{}{"reference": ref})
		require.NoError(t, err)
		if res2.(float64) != amount {
			// saga are not accepted yet
			time.Sleep(1 * time.Millisecond)
			continue
		}

		checkPassed = true
		break
	}

	require.True(t, checkPassed)
}

func TestContextPassing(t *testing.T) {
	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.New",
		map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)

	prototype, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.SelfRef",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
	require.Equal(t, ref, prototype.(string))
}

func TestNilResult(t *testing.T) {
	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.New",
		map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)

	res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.ReturnNil",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
	require.Nil(t, res)
}

func TestGetRemoteData(t *testing.T) {
	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.New",
		map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)

	res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.GetChildPrototype",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
	require.NotEmpty(t, res.(string))
}

func TestImmutableAnnotation(t *testing.T) {
	ref := callConstructor(t, "first", "New")

	res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.ExternalImmutableCall",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
	require.Equal(t, float64(10), res.(float64))

	res, err = testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.ExternalImmutableCallMakesExternalCall",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
}

func TestMultipleConstructorsCall(t *testing.T) {
	ref := callConstructor(t, "first", "NewZero")

	res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.Get",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
	require.Equal(t, float64(0), res.(float64))

	ref = callConstructorWithParameters(t, "first", "NewWithNumber", map[string]interface{}{"amount": 12})

	res, err = testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.Get",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
	require.Equal(t, float64(12), res.(float64))
}

func TestGetParent(t *testing.T) {
	ref := callConstructor(t, "first", "New")
	res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.AddChildAndReturnMyselfAsParent",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
	require.Equal(t, ref, res.(string))
}

func TestDeactivation(t *testing.T) {
	ref := callConstructor(t, "first", "New")
	_, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.Kill",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
}
