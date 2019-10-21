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
	"testing"

	"github.com/insolar/insolar/application/testutils/launchnet"
	"github.com/stretchr/testify/require"
)

func TestFoundationMemberCreate(t *testing.T) {
	for _, m := range launchnet.Foundation {
		err := verifyFundsMembersAndDeposits(t, m)
		if err != nil {
			require.NoError(t, err)
		}
	}
}

func TestEnterpriseMemberCreate(t *testing.T) {
	for _, m := range launchnet.Enterprise {
		err := verifyFundsMembersAndDeposits(t, m)
		if err != nil {
			require.NoError(t, err)
		}
	}
}

func TestNetworkIncentivesMemberCreate(t *testing.T) {
	for _, m := range launchnet.NetworkIncentives {
		err := verifyFundsMembersAndDeposits(t, m)
		if err != nil {
			require.NoError(t, err)
		}
	}
}

func TestApplicationIncentivesMemberCreate(t *testing.T) {
	for _, m := range launchnet.ApplicationIncentives {
		err := verifyFundsMembersAndDeposits(t, m)
		if err != nil {
			require.NoError(t, err)
		}
	}
}

func TestFundsMemberCreate(t *testing.T) {
	for _, m := range launchnet.Funds {
		err := verifyFundsMembersAndDeposits(t, m)
		if err != nil {
			require.NoError(t, err)
		}
	}
}
