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

/*
//go:generate stringer -type=step
type step int

const (
	stepUndefined step = iota + 1
	stepExchageData
	stepExchageHash
)
*/
type exchangeResults map[core.RecordRef][]*core.ActiveNode

type baseConsensus struct {
	self            Participant
	allParticipants []Participant
	communicator    Communicator
	holder          UnsyncHolder
	results         exchangeResults
	resultsHash     []*NodeUnsyncHash
}

// DoConsensus implements consensus interface
func (c *baseConsensus) DoConsensus(ctx context.Context, holder UnsyncHolder, self Participant, allParticipants []Participant) ([]*core.ActiveNode, error) {
	c.self = self
	c.allParticipants = allParticipants
	c.holder = holder
	c.results = make(exchangeResults, len(c.allParticipants))

	c.exchangeDataWithOtherParticipants(ctx)
	c.exchangeHashWithOtherParticipants(ctx)

	return c.holder.GetUnsync(), nil
}

func (c *baseConsensus) exchangeDataWithOtherParticipants(ctx context.Context) {

	wg := &sync.WaitGroup{}
	for _, p := range c.allParticipants {
		wg.Add(1)

		go func(wg *sync.WaitGroup, participant Participant) {
			defer wg.Done()
			if p.GetActiveNode().NodeID != c.self.GetActiveNode().NodeID {
				log.Infof("data exchage with %s", participant.GetActiveNode().NodeID.String())

				// goroutine
				data, err := c.communicator.ExchangeData(ctx, c.holder.GetPulse(), participant, c.holder.GetUnsync())
				if err != nil {
					log.Errorln(err.Error())
				}
				c.results[participant.GetActiveNode().NodeID] = data
			} else {
				c.results[participant.GetActiveNode().NodeID] = c.holder.GetUnsync()
			}
		}(wg, p)
	}
	wg.Wait()

	hashes := make([]*NodeUnsyncHash, len(c.results))
	for id, x := range c.results {
		r, _ := CalculateHash(id, x)
		hashes = append(hashes, r)
	}

	c.resultsHash = hashes
	c.holder.SetHash(hashes)
}

func (c *baseConsensus) exchangeHashWithOtherParticipants(ctx context.Context) {
	wg := &sync.WaitGroup{}
	for _, p := range c.allParticipants {
		wg.Add(1)

		go func(wg *sync.WaitGroup, participant Participant) {
			defer wg.Done()
			if p.GetActiveNode().NodeID != c.self.GetActiveNode().NodeID {
				log.Infof("data exchage with %s", participant.GetActiveNode().NodeID.String())

				// goroutine
				_, err := c.communicator.ExchangeHash(ctx, c.holder.GetPulse(), participant, c.resultsHash)
				if err != nil {
					log.Errorln(err.Error())
				}
			}
		}(wg, p)
	}
	wg.Wait()
}
