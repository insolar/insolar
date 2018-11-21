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
	"strings"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
)

// exchangeResults is thread safe results struct
type exchangeResults struct {
	mutex *sync.Mutex
	data  map[core.RecordRef][]core.Node
	hash  []*network.NodeUnsyncHash
}

func (r *exchangeResults) writeResultData(id core.RecordRef, data []core.Node) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.data[id] = data
}

func (r *exchangeResults) calculateResultHash() []*network.NodeUnsyncHash {
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

func (r *exchangeResults) getAllCollectedNodes() []core.Node {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	result := make([]core.Node, 0)
	for _, nodes := range r.data {
		result = append(result, nodes...)
	}
	return result
}

func newExchangeResults(participantsCount int) *exchangeResults {
	return &exchangeResults{
		mutex: &sync.Mutex{},
		data:  make(map[core.RecordRef][]core.Node, participantsCount),
		hash:  make([]*network.NodeUnsyncHash, 0),
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
func (c *baseConsensus) DoConsensus(ctx context.Context, holder UnsyncHolder, self Participant, allParticipants []Participant) ([]core.Node, error) {
	log.Infof("Start consensus between %d participants about %d unsyncs", len(allParticipants), len(holder.GetUnsync()))
	c.self = self
	c.allParticipants = allParticipants
	c.holder = holder
	c.results = newExchangeResults(len(c.allParticipants))

	c.exchangeDataWithOtherParticipants(ctx)
	c.exchangeHashWithOtherParticipants(ctx)
	return c.analyzeResults()
}

func (c *baseConsensus) exchangeDataWithOtherParticipants(ctx context.Context) {
	log.Debugln("Start exchange data between consensus participants")
	wg := &sync.WaitGroup{}
	wg.Add(len(c.allParticipants))
	for _, p := range c.allParticipants {

		go func(wg *sync.WaitGroup, participant Participant) {
			defer wg.Done()
			if participant.GetActiveNode().ID() != c.self.GetActiveNode().ID() {
				log.Infof("data exchage with %s", participant.GetActiveNode().ID().String())

				data, err := c.communicator.ExchangeData(ctx, c.holder.GetPulse(), participant, c.holder.GetUnsync())
				receivedNodes := make([]string, 0)
				for _, node := range data {
					receivedNodes = append(receivedNodes, node.ID().String())
				}
				log.Debugf("received from %s unsync list: %s", participant.GetActiveNode().ID(), strings.Join(receivedNodes, ", "))
				if err != nil {
					log.Errorln(err.Error())
				}
				c.results.writeResultData(participant.GetActiveNode().ID(), data)
			} else {
				c.results.writeResultData(participant.GetActiveNode().ID(), c.holder.GetUnsync())
			}
		}(wg, p)
	}
	wg.Wait()
	log.Debugln("End exchange data between consensus participants")
}

func (c *baseConsensus) exchangeHashWithOtherParticipants(ctx context.Context) {
	log.Debugln("Start exchange hashes between consensus participants")

	hash := c.results.calculateResultHash()
	c.holder.SetHash(hash)

	wg := &sync.WaitGroup{}
	for _, p := range c.allParticipants {
		wg.Add(1)

		go func(wg *sync.WaitGroup, participant Participant) {
			defer wg.Done()
			if participant.GetActiveNode().ID() != c.self.GetActiveNode().ID() {
				log.Infof("data exchage with %s", participant.GetActiveNode().ID().String())

				_, err := c.communicator.ExchangeHash(ctx, c.holder.GetPulse(), participant, hash)
				if err != nil {
					log.Errorln(err.Error())
				}
			}
		}(wg, p)
	}
	wg.Wait()
	log.Debugln("End exchange hashes between consensus participants")
}

func (c *baseConsensus) analyzeResults() ([]core.Node, error) {
	return c.results.getAllCollectedNodes(), nil
}
