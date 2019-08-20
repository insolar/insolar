/*
 * Copyright 2019 Insolar Technologies GmbH
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package executor

import (
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/object"
)

var (
	pending insolar.PulseNumber = 780
	topSync insolar.PulseNumber = 800
	current insolar.PulseNumber = 850
)

func indexesFixture() []record.Index {
	return []record.Index{
		{
			Lifeline: record.Lifeline{
				EarliestOpenRequest: &pending,
			},
		},
		{
			Lifeline: record.Lifeline{
				EarliestOpenRequest: &pending,
			},
		},
		{
			Lifeline: record.Lifeline{
				EarliestOpenRequest: &pending,
			},
		},
		{
			Lifeline: record.Lifeline{
				EarliestOpenRequest: &pending,
			},
		},
		{
			Lifeline: record.Lifeline{
				EarliestOpenRequest: &pending,
			},
		},
	}
}

func TestInitialStateKeeper_GetOnStart(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	ctx := inslogger.TestContext(t)

	topSyncPulse := insolar.Pulse{PulseNumber: 800}
	jetKeeper := NewJetKeeperMock(mc)
	jetKeeper.TopSyncPulseMock.Return(topSyncPulse.PulseNumber)

	jetAccessor := jet.NewAccessorMock(mc)
	jetAccessor.AllMock.Expect(ctx, topSyncPulse.PulseNumber).Return(nil)

	indexAccessor := object.NewIndexAccessorMock(mc)
	indexAccessor.ForPulseMock.Expect(ctx, topSyncPulse.PulseNumber).Return(nil, nil)

}
