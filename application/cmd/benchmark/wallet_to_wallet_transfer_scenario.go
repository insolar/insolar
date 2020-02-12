// Copyright 2020 Insolar Network Ltd.
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

package main

import (
	"fmt"
	"math/big"

	"github.com/insolar/insolar/application/api/sdk"
)

const transferAmount = 101

type walletToWalletTransferScenario struct {
	insSDK  *sdk.SDK
	members []sdk.Member

	balanceCheckMembers []sdk.Member
}

func (s *walletToWalletTransferScenario) canBeStarted() error {
	if len(s.members) < concurrent*2 {
		return fmt.Errorf("not enough members for start")
	}
	return nil
}

func (s *walletToWalletTransferScenario) prepare(repetition int) {
	members, err := getMembers(s.insSDK, concurrent*2, false)
	check("Error while loading members: ", err)

	if useMembersFromFile {
		members = members[:len(members)-2]
	}

	s.members = members

	s.balanceCheckMembers = make([]sdk.Member, len(s.members), len(s.members)+2)
	copy(s.balanceCheckMembers, s.members)
	s.balanceCheckMembers = append(s.balanceCheckMembers, s.insSDK.GetFeeMember())
	s.balanceCheckMembers = append(s.balanceCheckMembers, s.insSDK.GetMigrationAdminMember())
}

func (s *walletToWalletTransferScenario) start(concurrentIndex int, repetitionIndex int) (string, error) {
	from := s.members[concurrentIndex*2]
	to := s.members[concurrentIndex*2+1]

	return s.insSDK.Transfer(big.NewInt(transferAmount).String(), from, to)
}

func (s *walletToWalletTransferScenario) getBalanceCheckMembers() []sdk.Member {
	return s.balanceCheckMembers
}
