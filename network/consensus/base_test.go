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
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
)

func TestBaseConsensus_exchangeDataWithOtherParticipants(t *testing.T) {
	list2 := []core.Node{newActiveNode(12, 0)}
	list5 := []core.Node{newActiveNode(15, 0), newActiveNode(25, 0)}

	self := NewParticipant(1, nil)

	participant2 := NewParticipant(2, list2)
	participant3 := NewParticipant(3, nil)
	participant4 := NewParticipant(4, nil)
	participant5 := NewParticipant(5, list5)

	c := baseConsensus{self: self,
		allParticipants: []Participant{participant2, self, participant3, participant5, participant4},
		communicator:    &testCommunicator{self: self},
		holder:          &mockUnsyncHolder{},
		results:         newExchangeResults(5),
	}

	ctx := context.Background()
	c.exchangeDataWithOtherParticipants(ctx)

	assert.Equal(t, 1, len(c.results.data[participant2.GetID()]))
	assert.Equal(t, newActiveNode(12, 0), c.results.data[participant2.GetID()][0])

	assert.Equal(t, 2, len(c.results.data[participant5.GetID()]))
	assert.Equal(t, newActiveNode(15, 0), c.results.data[participant5.GetID()][0])
	assert.Equal(t, newActiveNode(25, 0), c.results.data[participant5.GetID()][1])
}
