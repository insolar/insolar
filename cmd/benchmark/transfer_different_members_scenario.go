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

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"math/big"
// 	"time"
//
// 	"github.com/insolar/insolar/api/sdk"
// )
//
// const transferAmount = 101
//
// type MemberWithBalance struct {
// 	Reference  string
// 	PrivateKey string
// 	PublicKey  string
// 	Balance    *big.Int
// }
//
// type transferDifferentMembersScenario struct {
// 	insSDK     *sdk.SDK
// 	concurrent int
// 	members    []*sdk.Member
//
// 	totalBalanceBefore  *big.Int
// 	balanceCheckMembers []*sdk.Member
// }
//
// func (s *transferDifferentMembersScenario) canBeStarted() error {
// 	if len(s.members) < s.concurrent*2 {
// 		return fmt.Errorf("not enough members for scenario")
// 	}
// 	return nil
// }
//
// func (s *transferDifferentMembersScenario) prepare() {
// 	s.balanceCheckMembers = make([]*sdk.Member, len(s.members))
//
// 	if !noCheckBalance {
// 		var membersWithBalanceMap map[string]*big.Int
// 		copy(s.balanceCheckMembers, s.members)
// 		s.balanceCheckMembers = append(s.balanceCheckMembers, s.insSDK.GetFeeMember())
// 		s.totalBalanceBefore, membersWithBalanceMap = getTotalBalance(s.insSDK, s.balanceCheckMembers)
//
// 		if useMembersFromFile {
// 			var membersWithBalance []*MemberWithBalance
// 			rawMembers, err := ioutil.ReadFile(memberFile)
// 			check("Error while read members with balance from file: ", err)
// 			err = json.Unmarshal(rawMembers, &membersWithBalance)
//
// 			check("Error while unmarshal members with balance: ", err)
// 			for _, m := range membersWithBalance {
// 				if membersWithBalanceMap[m.Reference] != nil && m.Balance != nil {
// 					if m.Balance.Cmp(membersWithBalanceMap[m.Reference]) == 0 {
// 						log.Fatalf("Balance mismatch: member with ref %s, balance at file - %s, balance at system - %s \n", m.Reference, m.Balance, membersWithBalanceMap[m.Reference])
// 					}
// 					log.Fatalf("Balance match: member with ref %s, balance at file - %s, balance at system - %s \n", m.Reference, m.Balance, membersWithBalanceMap[m.Reference])
// 				}
// 				fmt.Printf("NOOOOOOOOOO")
// 			}
// 		}
//
// 	}
// }
//
// func (s *transferDifferentMembersScenario) scenario(index int) (string, error) {
// 	from := s.members[index]
// 	to := s.members[index+1]
//
// 	return s.insSDK.Transfer(big.NewInt(transferAmount).String(), from, to)
// }
//
// func (s *transferDifferentMembersScenario) checkResult() {
// 	if !noCheckBalance {
// 		totalBalanceAfter := big.NewInt(0)
// 		var membersWithBalanceMap map[string]*big.Int
//
// 		for nretries := 0; nretries < balanceCheckRetries; nretries++ {
// 			totalBalanceAfter, membersWithBalanceMap = getTotalBalance(s.insSDK, s.balanceCheckMembers)
// 			if totalBalanceAfter.Cmp(s.totalBalanceBefore) == 0 {
// 				break
// 			}
// 			fmt.Printf("Total balance before and after don't match: %v vs %v - retrying in %s ...\n",
// 				s.totalBalanceBefore, totalBalanceAfter, balanceCheckDelay)
// 			time.Sleep(balanceCheckDelay)
//
// 		}
// 		fmt.Printf("Total balance before: %v and after: %v\n", s.totalBalanceBefore, totalBalanceAfter)
// 		if totalBalanceAfter.Cmp(s.totalBalanceBefore) != 0 {
// 			log.Fatal("Total balance mismatch!\n")
// 		}
//
// 		var membersWithBalance []*MemberWithBalance
// 		rawMembers, err := ioutil.ReadFile(memberFile)
// 		check("Error while read members with balance from file: ", err)
//
// 		err = json.Unmarshal(rawMembers, &membersWithBalance)
// 		check("Error while unmarshal members with balance: ", err)
//
// 		for _, m := range membersWithBalance {
// 			m.Balance = membersWithBalanceMap[m.Reference]
// 		}
//
// 		if saveMembersToFile {
// 			err = saveMembers(membersWithBalance)
// 			check("Error while saving members with balance to file: ", err)
// 		}
// 	}
// }
