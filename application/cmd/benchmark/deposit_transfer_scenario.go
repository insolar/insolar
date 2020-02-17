// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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

	s.migrationDaemons, err = s.insSDK.GetAndActivateMigrationDaemonMembers()
	check("failed to get and activate migration daemons: ", err)

	s.balanceCheckMembers = make([]sdk.Member, len(s.members), len(s.members)+2)
	copy(s.balanceCheckMembers, s.members)
	s.balanceCheckMembers = append(s.balanceCheckMembers, s.insSDK.GetFeeMember())
	s.balanceCheckMembers = append(s.balanceCheckMembers, s.insSDK.GetMigrationAdminMember())

	for _, m := range s.members {
		mm, ok := m.(*sdk.MigrationMember)
		if !ok {
			fmt.Println("failed to cast member to migration member")
			os.Exit(1)
		}
		for i := 0; i < repetition; i++ {
			_, err := s.insSDK.FullMigration(s.migrationDaemons, replaceLast(txHashPattern, strconv.Itoa(i)), big.NewInt(migrationAmount).String(), mm.MigrationAddress)
			if err != nil && !strings.Contains(err.Error(), "migration is done for this deposit") {
				check("Error while migrating tokens: ", err)
			}
		}
	}

	// wait for hold period end
	time.Sleep(30 * time.Second)

}

func (s *depositTransferScenario) start(concurrentIndex int, repetitionIndex int) (string, error) {
	var migrationMember *sdk.MigrationMember
	migrationMember, ok := s.members[concurrentIndex].(*sdk.MigrationMember)
	if !ok {
		return "", fmt.Errorf("unexpected member type: %T", s.members[concurrentIndex])
	}

	return s.insSDK.DepositTransfer(big.NewInt(migrationAmount*10).String(), migrationMember, replaceLast(txHashPattern, strconv.Itoa(repetitionIndex)))
}

func (s *depositTransferScenario) getBalanceCheckMembers() []sdk.Member {
	return s.balanceCheckMembers
}
