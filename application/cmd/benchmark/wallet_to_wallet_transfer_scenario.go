// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
