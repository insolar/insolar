// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

//
// import (
// 	"testing"
//
// 	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
// 	"github.com/stretchr/testify/require"
// )
//
// func TestCallUploadedContract(t *testing.T) {
// 	launchnet.RunOnlyWithLaunchnet(t)
// 	contractCode := `
// 		package main
// 		import "github.com/insolar/insolar/logicrunner/builtin/foundation"
// 		type One struct {
// 			foundation.BaseContract
// 		}
// 		func New() (*One, error){
// 			return &One{}, nil}
//
// 		var INSATTR_Hello_API = true
// 		func (r *One) Hello(str string) (string, error) {
// 			return str, nil
// 		}`
//
// 	prototypeRef := uploadContractOnce(t, "CallUploadedContract", contractCode)
// 	objectRef := callConstructor(t, prototypeRef, "New")
//
// 	testParam := "test"
// 	methodResult := callMethod(t, objectRef, "Hello", testParam)
// 	require.Empty(t, methodResult.Error)
// 	require.Equal(t, testParam, methodResult.ExtractedReply)
// }
