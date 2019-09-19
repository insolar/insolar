// /
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
// /

package main

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/insolar/insolar/api/sdk"
)

const migrationAmount = 101

type MigrationScenario struct {
	insSDK           *sdk.SDK
	members          []sdk.Member
	migrationDaemons []sdk.Member

	totalBalanceBefore  *big.Int
	balanceCheckMembers []sdk.Member
}

func (s *MigrationScenario) canBeStarted() error {
	if len(s.members) < concurrent {
		return fmt.Errorf("not enough members for start")
	}

	if len(s.migrationDaemons) < 3 {
		return fmt.Errorf("not enough migration daemons")
	}
	return nil
}

func (s *MigrationScenario) prepare() {
	members, err := getMembers(s.insSDK, concurrent, true)
	check("Error while loading members: ", err)
	s.members = members

	s.migrationDaemons = s.insSDK.GetMigrationDaemonMembers()

	if !noCheckBalance {
		s.balanceCheckMembers = make([]sdk.Member, len(s.members))
		copy(s.balanceCheckMembers, s.members)
		s.balanceCheckMembers = append(s.balanceCheckMembers, s.insSDK.GetMigrationAdminMember())
	}
}

func (s *MigrationScenario) start(concurrentIndex int, repetitionIndex int) (string, error) {
	var migrationMember *sdk.MigrationMember
	migrationMember, ok := s.members[concurrentIndex].(*sdk.MigrationMember)
	if !ok {
		return "", fmt.Errorf("unexpected member type: %T", s.members[concurrentIndex])
	}

	if traceId, err := s.insSDK.Migration(s.migrationDaemons[0], "tx_hash_"+strconv.Itoa(repetitionIndex), big.NewInt(migrationAmount).String(), migrationMember.MigrationAddress); err != nil {
		return traceId, err
	}
	return s.insSDK.Migration(s.migrationDaemons[1], "tx_hash_"+strconv.Itoa(repetitionIndex), big.NewInt(migrationAmount).String(), migrationMember.MigrationAddress)
}

func (s *MigrationScenario) getBalanceCheckMembers() []sdk.Member {
	return s.balanceCheckMembers
}
