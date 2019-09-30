///
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
///

package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/insolar/insolar/api/sdk"
)

const ethTxHash = "ethTxHash_test"

type depositTransferScenario struct {
	insSDK           *sdk.SDK
	members          []sdk.Member
	migrationDaemons []sdk.Member

	balanceCheckMembers []sdk.Member
}

func (s *depositTransferScenario) canBeStarted() error {
	if len(s.members) < concurrent {
		return fmt.Errorf("not enough members for start")
	}

	return nil
}

func (s *depositTransferScenario) prepare(repetition int) {
	members, err := getMembers(s.insSDK, concurrent, true)
	check("Error while loading members: ", err)

	if useMembersFromFile {
		members = members[:len(members)-2]
	}

	s.members = members

	s.migrationDaemons, err = s.insSDK.GetActivateMigrationDaemonMembers()
	check("failed to get migration daemons: ", err)

	for _, md := range s.migrationDaemons {
		_, err := s.insSDK.ActivateDaemon(md.GetReference())
		if err != nil && !strings.Contains(err.Error(), "[daemon member already activated]") {
			check("Error while activating daemons: ", err)
		}
	}

	s.balanceCheckMembers = make([]sdk.Member, len(s.members), len(s.members)+2)
	copy(s.balanceCheckMembers, s.members)
	s.balanceCheckMembers = append(s.balanceCheckMembers, s.insSDK.GetFeeMember())

	for _, m := range s.members {
		mm, ok := m.(*sdk.MigrationMember)
		if !ok {
			fmt.Println("failed to cast member to migration member")
			os.Exit(1)
		}
		for i := 0; i < repetition; i++ {
			_, err := s.insSDK.FullMigration(s.migrationDaemons, txHashPrefix+strconv.Itoa(i), big.NewInt(migrationAmount).String(), mm.MigrationAddress)
			check("Error while migrating tokens: ", err)
		}
	}

	time.Sleep(30 * time.Second)

}

func (s *depositTransferScenario) start(concurrentIndex int, repetitionIndex int) (string, error) {
	var migrationMember *sdk.MigrationMember
	migrationMember, ok := s.members[concurrentIndex].(*sdk.MigrationMember)
	if !ok {
		return "", fmt.Errorf("unexpected member type: %T", s.members[concurrentIndex])
	}

	return s.insSDK.DepositTransfer(big.NewInt(migrationAmount).String(), migrationMember, txHashPrefix+strconv.Itoa(repetitionIndex))
}

func (s *depositTransferScenario) getBalanceCheckMembers() []sdk.Member {
	return s.balanceCheckMembers
}
