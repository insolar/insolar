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
