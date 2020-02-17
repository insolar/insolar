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

type walletToWalletTwoSidesScenario struct {
	insSDK  *sdk.SDK
	members []sdk.Member

	balanceCheckMembers []sdk.Member
}

type transferResult struct {
	err     error
	traceID string
}

func (s *walletToWalletTwoSidesScenario) canBeStarted() error {
	if len(s.members) < concurrent*2 {
		return fmt.Errorf("not enough members for start")
	}
	return nil
}

func (s *walletToWalletTwoSidesScenario) prepare(repetition int) {
	members, err := getMembers(s.insSDK, concurrent*2, false)
	check("Error while loading members: ", err)

	if useMembersFromFile {
		members = members[:len(members)-2]
	}

	s.members = members

	s.balanceCheckMembers = make([]sdk.Member, len(s.members))
	copy(s.balanceCheckMembers, s.members)
	s.balanceCheckMembers = append(s.balanceCheckMembers, s.insSDK.GetFeeMember())
	s.balanceCheckMembers = append(s.balanceCheckMembers, s.insSDK.GetMigrationAdminMember())
}

func (s *walletToWalletTwoSidesScenario) start(concurrentIndex int, repetitionIndex int) (string, error) {
	c := make(chan *transferResult)

	first := s.members[concurrentIndex*2]
	second := s.members[concurrentIndex*2+1]

	transfer := func(from sdk.Member, to sdk.Member) {
		traceID, err := s.insSDK.Transfer(big.NewInt(transferAmount).String(), from, to)
		c <- &transferResult{err: err, traceID: traceID}
	}

	go transfer(first, second)
	go transfer(second, first)

	firstRes := <-c
	secondRes := <-c

	if firstRes.err != nil || secondRes.err != nil {
		fmt.Println("One of double transfer has error in it:")
		fmt.Println("Error: ", firstRes.err, ". Trace: ", firstRes.traceID)
		fmt.Println("Error: ", secondRes.err, ". Trace: ", secondRes.traceID)
	}

	if firstRes.err != nil {
		return firstRes.traceID, firstRes.err
	}

	if secondRes.err != nil {
		return secondRes.traceID, secondRes.err
	}

	return "", nil
}

func (s *walletToWalletTwoSidesScenario) getBalanceCheckMembers() []sdk.Member {
	return s.balanceCheckMembers
}
