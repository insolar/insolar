// Copyright 2020 Insolar Network Ltd.
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

// +build slowtest

package integration_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/applicationbase/genesis"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
)

func Test_FinalizePulse(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultHeavyConfig()
	defer os.RemoveAll(cfg.Ledger.Storage.DataDirectory)
	heavyConfig := genesis.HeavyConfig{}
	s, err := NewServer(ctx, cfg, heavyConfig, nil)
	assert.NoError(t, err)
	defer s.Stop()

	s.SetPulse(ctx)
	s.SetPulse(ctx)

	targetPulse := s.Pulse() - PulseStep

	_, done := s.Send(ctx, &payload.GotHotConfirmation{
		JetID: insolar.ZeroJetID,
		Pulse: targetPulse,
		Split: false,
	})
	done()

	require.Equal(t, insolar.GenesisPulse.PulseNumber, s.JetKeeper.TopSyncPulse())

	d := drop.Drop{
		Pulse: targetPulse,
		JetID: insolar.ZeroJetID,
		Split: false,
	}

	_, done = s.Send(ctx, &payload.Replication{
		JetID: insolar.ZeroJetID,
		Pulse: targetPulse,
		Drop:  d,
	})
	done()

	numIterations := 20
	for s.JetKeeper.TopSyncPulse() == insolar.GenesisPulse.PulseNumber && numIterations > 0 {
		time.Sleep(500 * time.Millisecond)
		numIterations--
	}
	require.Equal(t, targetPulse, s.JetKeeper.TopSyncPulse())
}
