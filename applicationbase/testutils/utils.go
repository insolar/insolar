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

package testutils

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"
)

func GenerateNodePublicKey(t *testing.T) string {
	ks := platformpolicy.NewKeyProcessor()

	privKey, err := ks.GeneratePrivateKey()
	require.NoError(t, err)

	pubKeyStr, err := ks.ExportPublicKeyPEM(ks.ExtractPublicKey(privKey))
	require.NoError(t, err)

	return string(pubKeyStr)
}

func ExpectedError(t *testing.T, trace []string, expected string) {
	found := hasSubstring(trace, expected)
	require.True(t, found, "Expected error (%s) not found in trace: %v", expected, trace)
}

func hasSubstring(trace []string, expected string) bool {
	found := false
	for _, trace := range trace {
		found = strings.Contains(trace, expected)
		if found {
			return found
		}
	}
	return found
}

func MakeSignedRequest(URL string, user launchnet.User, method string, params interface{}) (interface{}, string, error) {
	ctx := context.TODO()
	rootCfg, err := requester.CreateUserConfig(user.GetReference(), user.GetPrivateKey(), user.GetPublicKey())
	if err != nil {
		var suffix string
		if requesterError, ok := err.(*requester.Error); ok {
			suffix = " [" + strings.Join(requesterError.Data.Trace, ": ") + "]"
		}
		fmt.Println(err.Error() + suffix)
		return nil, "", err
	}

	var caller string
	fpcs := make([]uintptr, 1)
	for i := 2; i < 10; i++ {
		if n := runtime.Callers(i, fpcs); n == 0 {
			break
		}
		caller = runtime.FuncForPC(fpcs[0] - 1).Name()
		if ok, _ := regexp.MatchString(`\.Test`, caller); ok {
			break
		}
		caller = ""
	}

	seed, err := requester.GetSeed(URL)
	if err != nil {
		return nil, "", err
	}

	res, err := requester.SendWithSeed(ctx, URL, rootCfg, &requester.Params{
		CallSite:   method,
		CallParams: params,
		PublicKey:  user.GetPublicKey(),
		Reference:  user.GetReference(),
		Test:       caller}, seed)

	if err != nil {
		return nil, "", err
	}

	resp := requester.ContractResponse{}
	err = json.Unmarshal(res, &resp)
	if err != nil {
		return nil, "", err
	}

	if resp.Error != nil {
		return nil, "", resp.Error
	}

	if resp.Result == nil {
		return nil, "", errors.New("Error and result are nil")
	}
	return resp.Result.CallResult, resp.Result.RequestReference, nil

}

func SignedRequest(t *testing.T, URL string, user launchnet.User, method string, params interface{}) (interface{}, error) {
	res, refStr, err := MakeSignedRequest(URL, user, method, params)

	if err != nil {
		var suffix string
		requesterError, ok := err.(*requester.Error)
		if ok {
			suffix = " [" + strings.Join(requesterError.Data.Trace, ": ") + "]"
		}
		t.Error("[" + method + "]" + err.Error() + suffix)
	}
	require.NotEmpty(t, refStr, "request ref is empty")
	require.NotEqual(t, insolar.NewEmptyReference().String(), refStr, "request ref is zero")

	_, err = insolar.NewReferenceFromString(refStr)
	require.Nil(t, err)

	return res, err
}

func SignedRequestWithEmptyRequestRef(t *testing.T, URL string, user launchnet.User, method string, params interface{}) (interface{}, error) {
	res, refStr, err := MakeSignedRequest(URL, user, method, params)

	require.Equal(t, "", refStr)

	return res, err
}
