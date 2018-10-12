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

package consensus

import (
	"context"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
)

// exchangeResults is thread safe results struct
type exchangeResults struct {
	mutex *sync.RWMutex
	data  map[core.RecordRef][]*core.ActiveNode
	hash  []*NodeUnsyncHash
}

func (r *exchangeResults) writeResultData(id core.RecordRef, data []*core.ActiveNode) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.data[id] = data
}

func (r *exchangeResults) calculateResultHash() []*NodeUnsyncHash {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for id, x := range r.data {
		d, err := CalculateNodeUnsyncHash(id, x)
		if err != nil {
			log.Error(err)
			continue
		}
		r.hash = append(r.hash, d)
	}
	return r.hash
}

func newExchangeResults(participantsCount int) *exchangeResults {
	return &exchangeResults{
		mutex: &sync.RWMutex{},
		data:  make(map[core.RecordRef][]*core.ActiveNode, participantsCount),
		hash:  make([]*NodeUnsyncHash, participantsCount),
	}
}

type baseConsensus struct {
	self            Participant
	allParticipants []Participant
	communicator    Communicator
	holder          UnsyncHolder
	results         *exchangeResults
}

// DoConsensus implements consensus interface
func (c *baseConsensus) DoConsensus(ctx context.Context, holder UnsyncHolder, self Participant, allParticipants []Participant) ([]*core.ActiveNode, error) {
	c.self = self
	c.allParticipants = allParticipants
	c.holder = holder
	c.results = newExchangeResults(len(c.allParticipants))

	c.exchangeDataWithOtherParticipants(ctx)
	c.exchangeHashWithOtherParticipants(ctx)
	return c.analyzeResults()
}

func (c *baseConsensus) exchangeDataWithOtherParticipants(ctx context.Context) {

	wg := &sync.WaitGroup{}
	for _, p := range c.allParticipants {
		wg.Add(1)

		go func(wg *sync.WaitGroup, participant Participant) {
			defer wg.Done()
			if participant.GetActiveNode().NodeID != c.self.GetActiveNode().NodeID {
				log.Infof("data exchage with %s", participant.GetActiveNode().NodeID.String())

				data, err := c.communicator.ExchangeData(ctx, c.holder.GetPulse(), participant, c.holder.GetUnsync())
				if err != nil {
					log.Errorln(err.Error())
				}
				c.results.writeResultData(participant.GetActiveNode().NodeID, data)
			} else {
				c.results.writeResultData(participant.GetActiveNode().NodeID, c.holder.GetUnsync())
			}
		}(wg, p)
	}
	wg.Wait()
}

func (c *baseConsensus) exchangeHashWithOtherParticipants(ctx context.Context) {

	hash := c.results.calculateResultHash()
	c.holder.SetHash(hash)

	wg := &sync.WaitGroup{}
	for _, p := range c.allParticipants {
		wg.Add(1)

		go func(wg *sync.WaitGroup, participant Participant) {
			defer wg.Done()
			if participant.GetActiveNode().NodeID != c.self.GetActiveNode().NodeID {
				log.Infof("data exchage with %s", participant.GetActiveNode().NodeID.String())

				_, err := c.communicator.ExchangeHash(ctx, c.holder.GetPulse(), participant, hash)
				if err != nil {
					log.Errorln(err.Error())
				}
			}
		}(wg, p)
	}
	wg.Wait()
}

func (c *baseConsensus) analyzeResults() ([]*core.ActiveNode, error) {

	return c.holder.GetUnsync(), nil
}
