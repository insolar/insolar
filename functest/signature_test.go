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
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api/requester"
)

//type signCases struct {
//	input string
//	expectedErr string
//}

func IncorrectSign(t *testing.T, signature string) ([]byte, error) {
	testMember := createMember(t)
	seed, err := requester.GetSeed(TestAPIURL)
	require.NoError(t, err)
	body, err := requester.GetResponseBodyContract(
		TestCallUrl,
		requester.Request{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "api.call",
			Params:  requester.Params{Seed: seed, Reference: testMember.ref, PublicKey: testMember.pubKey, CallSite: "wallet.getBalance", CallParams: map[string]interface{}{"reference": testMember.ref}},
		},
		signature)
	return body, nil
}

//func TestTableSignData(t *testing.T){
//	data := []signCases{
//		{"MEQCIAvgBR42vSccBKynBIC7gb5GffqtW8q2XWRP+DlJ0IeUAiAeKCxZNSSRSsYcz2d49CT6KlSLpr5L7VlOokOiI9dsvQ==", "invalid signature"},
//		{"","invalid signature"},
//	}
//
//
//}

func TestIncorrectSign(t *testing.T) {
	body, err := IncorrectSign(t, "MEQCIAvgBR42vSccBKynBIC7gb5GffqtW8q2XWRP+DlJ0IeUAiAeKCxZNSSRSsYcz2d49CT6KlSLpr5L7VlOokOiI9dsvQ==")
	require.NoError(t, err)
	var res requester.ContractAnswer
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)
	require.Contains(t, res.Error.Message, "invalid signature")
}

func TestEmptyStrSign(t *testing.T) {
	body, err := IncorrectSign(t, "")
	require.NoError(t, err)
	var res requester.ContractAnswer
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)
	require.Contains(t, res.Error.Message, "invalid signature")
}

func TestSignWithSpaces(t *testing.T) {
	body, err := IncorrectSign(t, "  MEQCIAvgBR42vSccBKynBIC7gb5GffqtW8q2XWRP+DlJ0IeUAiAeKCxZNSSRSsYcz2d49CT6KlSLpr5L7VlOokOiI9dsvQ==   ")
	require.NoError(t, err)
	var res requester.ContractAnswer
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)
	require.Contains(t, res.Error.Message, "invalid signature")
}

func TestSignOtherMemberWithPubKeyWords(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	body, err := IncorrectSign(t, member.pubKey)
	require.NoError(t, err)
	var res requester.ContractAnswer
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)
	require.Contains(t, res.Error.Message, "invalid signature")
}

func TestSignOtherMember(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	s := strings.Split(member.pubKey, "-----")
	sign := s[2]
	body, err := IncorrectSign(t, sign)
	require.NoError(t, err)
	var res requester.ContractAnswer
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)
	require.Contains(t, res.Error.Message, "invalid signature")
}
