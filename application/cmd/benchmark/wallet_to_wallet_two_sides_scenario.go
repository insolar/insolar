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
