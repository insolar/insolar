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
	"testing"

	"github.com/insolar/insolar/log"
	"github.com/stretchr/testify/assert"
)

func initNodes(count int) []TestNode {
	participants := make([]Participant, count)
	for i := range participants {
		participants[i] = NewParticipant(byte(i), nil)

	}

	nodes := make([]TestNode, count)
	for i := range nodes {
		c := NewConsensus(&testCommunicator{})

		nodes[i] = TestNode{self: participants[i],
			allParticipants: participants,
			consensus:       c,
			ctx:             context.Background(),
		}
		log.Infof("Node %d has id %s", i, nodes[i].self.GetActiveNode().ID().String())
	}

	return nodes
}

func TestNodeConsensus_DoConsensus(t *testing.T) {
	//t.Skip("not redy")
	nodes := initNodes(5)

	wg := &sync.WaitGroup{}

	for _, n := range nodes {
		wg.Add(1)
		go func(p TestNode, wg *sync.WaitGroup) {
			defer wg.Done()
			log.Info("Do consensus for ", p.self.GetActiveNode().ID().String())
			result, err := p.consensus.DoConsensus(p.ctx, &mockUnsyncHolder{}, p.self, p.allParticipants)
			//consensusResult, err := p.consensus.DoConsensus(p.ctx, p.self, p.allParticipants)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(result))
			//log.Infof("%s consensus result %b", p.self.GetActiveNode().ID().String(), consensusResult)
		}(n, wg)
	}

	wg.Wait()
}
