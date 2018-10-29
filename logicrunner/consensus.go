/*
 *
 *  *    Copyright 2018 Insolar
 *  *
 *  *    Licensed under the Apache License, Version 2.0 (the "License");
 *  *    you may not use this file except in compliance with the License.
 *  *    You may obtain a copy of the License at
 *  *
 *  *        http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  *    Unless required by applicable law or agreed to in writing, software
 *  *    distributed under the License is distributed on an "AS IS" BASIS,
 *  *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  *    See the License for the specific language governing permissions and
 *  *    limitations under the License.
 *
 */

package logicrunner

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/log"
)

type ConsensusRecord struct {
	Steps int
	Error error
}

// Consensus is an object for one validation process where all validated results will be compared.
type Consensus struct {
	lr          *LogicRunner
	Have        int
	Need        int
	Total       int
	Results     map[Ref]ConsensusRecord
	CaseRecords []core.CaseRecord
}

func newConsensus(lr *LogicRunner, refs []Ref) *Consensus {
	c := &Consensus{
		lr:      lr,
		Results: make(map[Ref]ConsensusRecord),
	}
	for _, r := range refs {
		c.Results[r] = ConsensusRecord{}
	}
	c.Total = len(c.Results)
	c.Need = c.Total/2 + 1
	return c
}

// AddValidated adds results from validators
func (c *Consensus) AddValidated(m *message.ValidationResults, node Ref, sig []byte) {
	caller := *m.GetCaller()
	if _, ok := c.Results[caller]; !ok {
		// why ??
	} else {
		c.Results[caller] = ConsensusRecord{
			Steps: m.PassedStepsCount,
			Error: m.Error,
		}
	}
	c.Have++
	c.CheckReady()
}

func (c *Consensus) AddExecutor(m *message.ExecutorResults, node Ref, sig []byte) {
	c.CaseRecords = m.CaseRecords
	c.CheckReady()
}

func (c *Consensus) CheckReady() {
	if c.CaseRecords == nil {
		return
	} else if c.Have < c.Need {
		return
	}
	steps := make(map[int]int)
	maxSame := 0
	stepsSame := 0
	for _, r := range c.Results {
		steps[r.Steps]++
		if maxSame < steps[r.Steps] {
			maxSame = steps[r.Steps]
			stepsSame = r.Steps
		}
	}
	if maxSame < c.Need && c.Total == c.Have {
		log.Debugf("Contract failed, agrred for %d steps by %d nodes", stepsSame, maxSame)
	}
	log.Debugf("Contract checking validation")
}
