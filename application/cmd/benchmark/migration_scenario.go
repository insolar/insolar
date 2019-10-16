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

	"github.com/insolar/insolar/application/api/sdk"
)

const (
	migrationAmount = 101

	txHashPrefix = "tx_hash_"
)

type migrationScenario struct {
	insSDK           *sdk.SDK
	members          []sdk.Member
	migrationDaemons []sdk.Member

	balanceCheckMembers []sdk.Member
}

func (s *migrationScenario) canBeStarted() error {
	if len(s.members) < concurrent {
		return fmt.Errorf("not enough members for start")
	}

	if len(s.migrationDaemons) < 2 {
		return fmt.Errorf("not enough migration daemons")
	}
	return nil
}

func (s *migrationScenario) prepare(repetition int) {
	members, err := getMembers(s.insSDK, concurrent, true)
	check("Error while loading members: ", err)

	if useMembersFromFile {
		members = members[:len(members)-2]
	}

	s.members = members

	s.migrationDaemons, err = s.insSDK.GetAndActivateMigrationDaemonMembers()
	check("failed to get and activate migration daemons: ", err)

	s.balanceCheckMembers = make([]sdk.Member, len(s.members), len(s.members)+2)
	copy(s.balanceCheckMembers, s.members)
	s.balanceCheckMembers = append(s.balanceCheckMembers, s.insSDK.GetFeeMember())
	s.balanceCheckMembers = append(s.balanceCheckMembers, s.insSDK.GetMigrationAdminMember())
}

func (s *migrationScenario) start(concurrentIndex int, repetitionIndex int) (string, error) {
	var migrationMember *sdk.MigrationMember
	migrationMember, ok := s.members[concurrentIndex].(*sdk.MigrationMember)
	if !ok {
		return "", fmt.Errorf("unexpected member type: %T", s.members[concurrentIndex])
	}

	if traceID, err := s.insSDK.Migration(s.migrationDaemons[0], txHashPrefix+strconv.Itoa(repetitionIndex), big.NewInt(migrationAmount).String(), migrationMember.MigrationAddress); err != nil {
		return traceID, err
	}
	return s.insSDK.Migration(s.migrationDaemons[1], txHashPrefix+strconv.Itoa(repetitionIndex), big.NewInt(migrationAmount).String(), migrationMember.MigrationAddress)
}

func (s *migrationScenario) getBalanceCheckMembers() []sdk.Member {
	return s.balanceCheckMembers
}
