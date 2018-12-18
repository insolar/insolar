/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/insolar/insolar/instrumentation/inslogger"
)

type scenario interface {
	canBeStarted() error
	start()
	getOperationsNumber() int
	getName() string
	getOut() io.Writer
}

type transferDifferentMembersScenario struct {
	name        string
	concurrent  int
	repetitions int
	members     []memberInfo
	out         io.Writer
}

func (s *transferDifferentMembersScenario) getOperationsNumber() int {
	return s.concurrent * s.repetitions
}

func (s *transferDifferentMembersScenario) getName() string {
	return s.name
}

func (s *transferDifferentMembersScenario) getOut() io.Writer {
	return s.out
}

func (s *transferDifferentMembersScenario) canBeStarted() error {
	writeToOutput(s.getOut(), fmt.Sprint("canBeStarted\n"))
	if len(s.members) < s.concurrent*2 {
		return fmt.Errorf("not enough members for scenario %s", s.getName())
	}
	return nil
}

func (s *transferDifferentMembersScenario) start() {
	var wg sync.WaitGroup
	for i := 0; i < s.concurrent*2; i = i + 2 {
		wg.Add(1)
		go s.startMember(i, &wg)
	}
	wg.Wait()
}

func (s *transferDifferentMembersScenario) startMember(index int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := 0; j < s.repetitions; j = j + 1 {
		ctx := inslogger.ContextWithTrace(context.Background(), fmt.Sprintf("transferFromMemberNumber%d", index))
		from := s.members[index]
		to := s.members[index+1]
		response := transfer(ctx, 1, from, to)
		if response != "success" {
			writeToOutput(s.out, fmt.Sprintf("[Member №%d] Transfer from %s to %s. Response: %s.\n", index, from.ref, to.ref, response))
		}
	}
}
