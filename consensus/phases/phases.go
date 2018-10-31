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

package phases

import (
	"github.com/pkg/errors"
)

// FirstPhase is a first phase.
type FirstPhase struct {
	next Phase
}

func (fp *FirstPhase) Calculate(args ...interface{}) error {
	pulse, ok := args[0].(*PulseData)
	if !ok {
		return errors.New("failed to cast to pulse")
	}
	// do work with pulse
	result := fp.work(pulse)
	return fp.next.Calculate(result)
}

func (fp *FirstPhase) work(pulse *PulseData) *NodePulseProof {
	return NewNodePulseProof()
}

// SecondPhase is a second phase.
type SecondPhase struct {
	next Phase
}

func (sp *SecondPhase) Calculate(args ...interface{}) error {
	_, ok := args[0].(*NodePulseProof)
	if !ok {
		return errors.New("failed to cast to pulse proof")
	}
	return nil
}
