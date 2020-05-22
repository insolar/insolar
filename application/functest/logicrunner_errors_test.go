// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest_error

package functest

import (
	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/applicationbase/testutils/testrequest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPanicIsSystemError(t *testing.T) {
	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.New",
		map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)

	_, err = testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, &Root, "first.Panic",
		map[string]interface{}{"reference": ref})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "CallMethod returns error")
	require.Contains(t, data.Trace, "AAAAAAAA!")

	_, err = testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, &Root, "first.NewPanic",
		map[string]interface{}{})
	data = checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "CallMethod returns error")
	require.Contains(t, data.Trace, "BBBBBBBB!")
}

func TestPanicIsLogicalError(t *testing.T) {
	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "panicAsLogicalError.New",
		map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)

	_, err = testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, &Root, "panicAsLogicalError.Panic",
		map[string]interface{}{"reference": ref})
	data := checkConvertRequesterError(t, err).Data
	require.NotContains(t, data.Trace, "CallMethod returns error")
	require.Contains(t, data.Trace, "AAAAAAAA!")

	_, err = testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, &Root, "panicAsLogicalError.NewPanic",
		map[string]interface{}{})
	data = checkConvertRequesterError(t, err).Data
	require.NotContains(t, data.Trace, "CallMethod returns error")
	require.Contains(t, data.Trace, "BBBBBBBB!")
}

func TestRecursiveCallError(t *testing.T) {
	obj := callConstructor(t, "first", "New")
	_, _, err := testrequest.MakeSignedRequest(
		launchnet.TestRPCUrlPublic,
		&Root,
		"first.Recursive",
		map[string]interface{}{"reference": obj})
	data := checkConvertRequesterError(t, err).Data

	require.Contains(t, data.Trace[len(data.Trace)-1], "loop detected")
}

func TestPrototypeMismatch(t *testing.T) {
	objSecond := callConstructor(t, "third", "New")
	objTest := callConstructor(t, "first", "New")

	_, _, err := testrequest.MakeSignedRequest(
		launchnet.TestRPCUrlPublic,
		&Root,
		"first.Test",
		map[string]interface{}{
			"reference": objTest,
			"firstRef":  objSecond,
		})
	data := checkConvertRequesterError(t, err).Data

	require.Contains(t, data.Trace, "try to call method of prototype as method of another prototype")
}

func TestContractWithEmbeddedConstructor(t *testing.T) {
	_ = callConstructor(t, "first", "NewZero")
	_, _, err := testrequest.MakeSignedRequest(
		launchnet.TestRPCUrlPublic,
		&Root,
		"second.NewWithOne",
		map[string]interface{}{
			"oneNumber": "10",
		})
	data := checkConvertRequesterError(t, err).Data

	require.Contains(t, data.Trace, "object is not activated")
}

// TestDeactivation

func TestErrorInterface(t *testing.T) {
	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.New",
		map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)

	_, _, err = testrequest.MakeSignedRequest(
		launchnet.TestRPCUrlPublic,
		&Root,
		"first.AnError",
		map[string]interface{}{"reference": ref})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "an error")

	_, err = testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.NoError",
		map[string]interface{}{"reference": ref})
	require.NoError(t, err)
}

// If a contract constructor returns `nil, nil` it's considered a logical error,
// which is returned to the calling contract and/or API.
func TestConstructorReturnNil(t *testing.T) {
	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.New",
		map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)

	_, _, err = testrequest.MakeSignedRequest(
		launchnet.TestRPCUrlPublic,
		&Root,
		"first.ConstructorReturnNil",
		map[string]interface{}{"reference": ref})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "constructor returned nil")

}

// If a contract constructor fails it's considered a logical error,
// which is returned to the calling contract and/or API.
func TestConstructorReturnError(t *testing.T) {
	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.New",
		map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)

	_, _, err = testrequest.MakeSignedRequest(
		launchnet.TestRPCUrlPublic,
		&Root,
		"first.ConstructorReturnError",
		map[string]interface{}{"reference": ref})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "Epic fail in NewWithErr()")
}
