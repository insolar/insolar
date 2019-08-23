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
	"math/big"
	"testing"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/testutils/launchnet"

	"github.com/stretchr/testify/require"
)

func TestGetBalance(t *testing.T) {
	firstMember := createMember(t)
	firstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	r := big.NewInt(10000000000)
	require.Equal(t, r, firstBalance)
}

func TestGetBalanceWrongRef(t *testing.T) {
	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl, &launchnet.Root, "member.getBalance",
		map[string]interface{}{"reference": gen.Reference().String()})
	require.Error(t, err)
	require.Contains(t, err.Error(), "index not found")
}
