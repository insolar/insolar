// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"github.com/insolar/insolar/api/sdk"
)

type createMemberScenario struct {
	insSDK *sdk.SDK
}

func (s *createMemberScenario) canBeStarted() error {
	return nil
}

func (s *createMemberScenario) prepare(repetition int) {}

func (s *createMemberScenario) start(concurrentIndex int, repetitionIndex int) (string, error) {
	_, traceID, err := s.insSDK.CreateMember()
	return traceID, err
}

func (s *createMemberScenario) getBalanceCheckMembers() []sdk.Member {
	return nil
}
