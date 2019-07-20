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

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/heavy/replica/integrity"
	"github.com/insolar/insolar/ledger/heavy/sequence"
)

func TestTarget_Notify(t *testing.T) {
	var (
		ctx        = inslogger.TestContext(t)
		parent     = NewParentMock(t)
		jetKeeper  = NewJetKeeperMock(t)
		sequencer  = sequence.NewSequencerMock(t)
		validator  = integrity.NewValidatorMock(t)
		firstPulse = insolar.GenesisPulse.PulseNumber
		tar        = NewTarget(configuration.Replica{}, parent)
	)
	parent.PullMock.Return([]byte{}, 0, nil)
	jetKeeper.TopSyncPulseMock.Return(firstPulse)
	jetKeeper.UpdateMock.Return(nil)
	tar.(*target).JetKeeper = jetKeeper
	sequencer.LenMock.Return(0)
	sequencer.UpsertMock.Return(nil)
	sequencer.SliceMock.Return([]sequence.Item{})
	tar.(*target).Sequencer = sequencer
	validator.UnwrapAndValidateMock.Return([]sequence.Item{})
	tar.(*target).Validator = validator

	err := tar.(*target).Start(ctx)
	require.NoError(t, err)

	err = tar.Notify(ctx, firstPulse)
	require.NoError(t, err)
}

func TestTarget_nextPulse(t *testing.T) {
	var (
		ctx       = inslogger.TestContext(t)
		db        = store.NewMemoryMockDB()
		pulses    = pulse.NewDB(db)
		sequencer = sequence.NewSequencer(db)
		tar       = NewTarget(configuration.Replica{}, nil)
		first     = insolar.GenesisPulse.PulseNumber
		second    = first + 10
	)
	tar.(*target).Sequencer = sequencer
	err := pulses.Append(ctx, insolar.Pulse{
		PulseNumber: first,
	})
	require.NoError(t, err)
	err = pulses.Append(ctx, insolar.Pulse{
		PulseNumber: second,
	})
	require.NoError(t, err)

	actual := tar.(*target).nextPulse(first)
	require.Equal(t, second, actual)
}
