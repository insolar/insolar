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

package functest

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api/requester"
)

type signCases struct {
	input       string
	expectedErr string
}

func IncorrectSignTableTests(t *testing.T, signature string, error string) {
	testMember := createMember(t)
	seed, err := requester.GetSeed(TestAPIURL)
	require.NoError(t, err)

	//debug
	fmt.Println("input : ")
	fmt.Println(signature)
	fmt.Println("expectedErr : ")
	fmt.Println(error)
	//debug end

	body, err := requester.GetResponseBodyContract(
		TestCallUrl,
		requester.Request{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "api.call",
			Params:  requester.Params{Seed: seed, Reference: testMember.ref, PublicKey: testMember.pubKey, CallSite: "wallet.getBalance", CallParams: map[string]interface{}{"reference": testMember.ref}},
		},
		signature)
	require.NoError(t, err)
	var res requester.ContractAnswer
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)
	require.Contains(t, res.Error.Message, error)
}

func TestTableSignData(t *testing.T) {
	data := []signCases{
		{"MEQCIAvgBR42vSccBKynBIC7gb5GffqtW8q2XWRP+DlJ0IeUAiAeKCxZNSSRSsYcz2d49CT6KlSLpr5L7VlOokOiI9dsvQ==", "invalid signature"},
		{"", "empty signature"},
		{"  MEQCIAvgBR42vSccBKynBIC7gb5GffqtW8q2XWRP+DlJ0IeUAiAeKCxZNSSRSsYcz2d49CT6KlSLpr5L7VlOokOiI9dsvQ==   ",
			"[ makeCall ] Error in called method: error while verify signature: cant decode signature illegal base64 data at input byte 0"},
	}
	for _, tc := range data {
		IncorrectSignTableTests(t, tc.input, tc.expectedErr)
	}

}

//bug https://insolar.atlassian.net/browse/INS-3115
func TestTableSignData1(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)

	member2, err := newUserWithKeys()
	require.NoError(t, err)
	sign := strings.ReplaceAll(
		strings.Split(member2.pubKey, "-----")[2], "\n", "")

	data := []signCases{
		{member.pubKey, "empty signature"},
		{sign, "invalid signature"},
	}
	for _, tc := range data {
		IncorrectSignTableTests(t, tc.input, tc.expectedErr)
	}
}
