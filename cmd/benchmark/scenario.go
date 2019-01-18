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
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/insolar/insolar/api/sdk"
	"github.com/pkg/errors"
)

type scenario interface {
	canBeStarted() error
	start()
	getOperationsNumber() int
	getAverageOperationDuration() time.Duration
	getName() string
	getOut() io.Writer
	printResult()
}

type transferDifferentMembersScenario struct {
	name        string
	concurrent  int
	repetitions int
	out         io.Writer
	totalTime   int64
	successes   uint32
	errors      uint32
	timeouts    uint32
	members     []*sdk.Member
	insSDK      *sdk.SDK
}

func (s *transferDifferentMembersScenario) getOperationsNumber() int {
	return s.concurrent * s.repetitions
}

func (s *transferDifferentMembersScenario) getAverageOperationDuration() time.Duration {
	return time.Duration(s.totalTime / int64(s.getOperationsNumber()))
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
		from := s.members[index]
		to := s.members[index+1]

		start := time.Now()
		traceID, err := s.insSDK.Transfer(1, from, to)
		atomic.AddInt64(&s.totalTime, int64(time.Since(start)))

		if err == nil {
			atomic.AddUint32(&s.successes, 1)
		} else if netErr, ok := errors.Cause(err).(net.Error); ok && netErr.Timeout() {
			atomic.AddUint32(&s.timeouts, 1)
			writeToOutput(s.out, fmt.Sprintf("[Member №%d] Transfer error with traceID: %s. Timeout.\n", index, traceID))
		} else {
			atomic.AddUint32(&s.errors, 1)
			if strings.Contains(err.Error(), "Incorrect message pulse") {
				writeToOutput(s.out, fmt.Sprintf("[Member №%d] Incorrect message pulse [its ok]: %s. Response: %s.\n", index, traceID, err.Error()))
			} else {
				writeToOutput(s.out, fmt.Sprintf("[Member №%d] Transfer error with traceID: %s. Response: %s.\n", index, traceID, err.Error()))
			}
		}
	}
}

func (s *transferDifferentMembersScenario) printResult() {
	writeToOutput(s.out, fmt.Sprintf("Scenario result:\n\tSuccesses: %d\n\tErrors: %d\n\tTimeouts: %d\n", s.successes, s.errors, s.timeouts))
}
