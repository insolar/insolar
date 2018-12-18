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
	"context"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/pkg/errors"
)

type ConsensusRecord struct {
	Steps   int
	Error   string
	Message core.Parcel
}

// Consensus is an object for one validation process where all validated results will be compared.
type Consensus struct {
	sync.Mutex
	lr       *LogicRunner
	ready    bool
	Have     int
	Need     int
	Total    int
	Results  map[Ref]ConsensusRecord
	CaseBind CaseBind
	Message  core.Parcel
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
func (c *Consensus) AddValidated(ctx context.Context, sm core.Parcel, msg *message.ValidationResults) error {
	source := sm.GetSender()
	c.Lock()
	defer c.Unlock()
	if _, ok := c.Results[source]; !ok {
		return errors.Errorf("Validation packet from non validation node for %#v", sm)
	} else {
		c.Results[source] = ConsensusRecord{
			Steps: msg.PassedStepsCount,
			Error: msg.Error,
		}
	}
	c.Have++
	c.CheckReady(ctx)
	return nil
}

func (c *Consensus) AddExecutor(ctx context.Context, sm core.Parcel, msg *message.ExecutorResults) {
	c.Lock()
	defer c.Unlock()
	c.CaseBind = *NewCaseBindFromExecutorResultsMessage(msg)
	c.Message = sm
	c.CheckReady(ctx)
}

func (c *Consensus) CheckReady(ctx context.Context) {
	if c.Have < c.Need {
		return
	}
	steps := make(map[int]int)
	maxSame := 0   // count of nodes with same result
	stepsSame := 0 // steps agreed by maximum nodes
	for _, r := range c.Results {
		steps[r.Steps]++
		if maxSame < steps[r.Steps] {
			maxSame = steps[r.Steps]
			stepsSame = r.Steps
		}
	}
	var err error
	if maxSame < c.Need && c.Total == c.Have {
		c.ready = true
		err = c.lr.ArtifactManager.RegisterValidation(ctx, c.GetReference(), *c.FindRequestBefore(stepsSame), false, c.GetValidatorSignatures())
	} else if maxSame >= c.Need && stepsSame == len(c.CaseBind.Requests) {
		c.ready = true
		err = c.lr.ArtifactManager.RegisterValidation(ctx, c.GetReference(), *c.FindRequestBefore(stepsSame), true, c.GetValidatorSignatures())
	}
	if err != nil {
		panic(err)
	}
}

func (c *Consensus) GetReference() Ref {
	return c.Message.Message().(*message.ExecutorResults).RecordRef
}

//
func (c *Consensus) GetValidatorSignatures() (messages []core.Message) {
	for _, x := range c.Results {
		messages = append(messages, x.Message)
	}
	return messages
}

// FindRequestBefore returns request placed before step (last valid request)
func (c *Consensus) FindRequestBefore(steps int) *core.RecordID {
	// TODO: resurrect this part
	return nil
}
