// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/insolar/insolar/api/sdk"
)

const (
	migrationAmount = 101

	txHashPattern = "0x89f2d6994e5d152bece9ec291f6098d236ab81f76f0d4d52fb69d0cd6b6fd70d"
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

	if traceID, err := s.insSDK.Migration(s.migrationDaemons[0], replaceLast(txHashPattern, strconv.Itoa(repetitionIndex)), big.NewInt(migrationAmount).String(), migrationMember.MigrationAddress); err != nil {
		return traceID, err
	}
	return s.insSDK.Migration(s.migrationDaemons[1], replaceLast(txHashPattern, strconv.Itoa(repetitionIndex)), big.NewInt(migrationAmount).String(), migrationMember.MigrationAddress)
}

func (s *migrationScenario) getBalanceCheckMembers() []sdk.Member {
	return s.balanceCheckMembers
}

func replaceLast(str, replace string) string {
	if len(replace) >= len(str) {
		return replace
	}
	return str[:len(str)-len(replace)] + replace
}
