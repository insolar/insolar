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

package main

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/insolar/insolar/api/sdk"
)

const transferAmount = 101

type transferDifferentMembersScenario struct {
	insSDK     *sdk.SDK
	concurrent int
	members    []*sdk.Member

	totalBalanceBefore  *big.Int
	balanceCheckMembers []*sdk.Member
}

func (s *transferDifferentMembersScenario) canBeStarted() error {
	if len(s.members) < s.concurrent*2 {
		return fmt.Errorf("not enough members for scenario")
	}
	return nil
}

func (s *transferDifferentMembersScenario) prepare() {
	s.balanceCheckMembers = make([]*sdk.Member, len(s.members))

	if !noCheckBalance {
		copy(s.balanceCheckMembers, s.members)
		s.balanceCheckMembers = append(s.balanceCheckMembers, s.insSDK.GetFeeMember())
		s.totalBalanceBefore = getTotalBalance(s.insSDK, s.balanceCheckMembers)
	}
}

func (s *transferDifferentMembersScenario) scenario(index int) (string, error) {
	from := s.members[index]
	to := s.members[index+1]

	return s.insSDK.Transfer(big.NewInt(transferAmount).String(), from, to)
}

func (s *transferDifferentMembersScenario) checkResult() {
	if !noCheckBalance {
		totalBalanceAfter := big.NewInt(0)
		for nretries := 0; nretries < balanceCheckRetries; nretries++ {
			totalBalanceAfter = getTotalBalance(s.insSDK, s.balanceCheckMembers)
			if totalBalanceAfter.Cmp(s.totalBalanceBefore) == 0 {
				break
			}
			fmt.Printf("Total balance before and after don't match: %v vs %v - retrying in %s ...\n",
				s.totalBalanceBefore, totalBalanceAfter, balanceCheckDelay)
			time.Sleep(balanceCheckDelay)

		}
		fmt.Printf("Total balance before: %v and after: %v\n", s.totalBalanceBefore, totalBalanceAfter)
		if totalBalanceAfter.Cmp(s.totalBalanceBefore) != 0 {
			log.Fatal("Total balance mismatch!\n")
		}
	}
}
