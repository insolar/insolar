//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package replica

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/replica/integrity"
	"github.com/insolar/insolar/ledger/heavy/sequence"
)

func TestParent_Subscribe(t *testing.T) {
	var (
		ctx       = inslogger.TestContext(t)
		syncPulse = insolar.GenesisPulse.PulseNumber
		page      = Page{Pulse: syncPulse}
		sequencer = sequence.NewSequencerMock(t)
		jetKeeper = NewJetKeeperMock(t)
		provider  = integrity.NewProviderMock(t)
		target    = NewTargetMock(t)
		p         = NewParent()
	)
	p.(*parent).Sequencer = sequencer
	jetKeeper.TopSyncPulseMock.Return(syncPulse)
	p.(*parent).JetKeeper = jetKeeper
	p.(*parent).Provider = provider
	target.NotifyMock.Return(nil)

	err := p.Subscribe(ctx, target, page)
	require.NoError(t, err)
}

func TestParent_Pull(t *testing.T) {
	var (
		ctx       = inslogger.TestContext(t)
		syncPulse = insolar.GenesisPulse.PulseNumber
		page      = Page{Pulse: syncPulse}
		reply     = []byte{1, 2, 3}
		total     = uint32(0)
		sequencer = sequence.NewSequencerMock(t)
		jetKeeper = NewJetKeeperMock(t)
		provider  = integrity.NewProviderMock(t)
		p         = NewParent()
	)
	sequencer.SliceMock.Return([]sequence.Item{})
	sequencer.LenMock.Return(0)
	p.(*parent).Sequencer = sequencer
	jetKeeper.TopSyncPulseMock.Return(syncPulse)
	p.(*parent).JetKeeper = jetKeeper
	provider.WrapMock.Return(reply)
	p.(*parent).Provider = provider

	actualReply, actualTotal, err := p.Pull(ctx, page)
	require.NoError(t, err)
	require.Equal(t, reply, actualReply)
	require.Equal(t, total, actualTotal)
}
