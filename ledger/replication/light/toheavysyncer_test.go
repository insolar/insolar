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

package light

import (
	"testing"

	"github.com/google/gofuzz"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/message"
	"github.com/stretchr/testify/require"
)

func TestToHeavySyncer_addToWaitingPulses(t *testing.T) {
	syncer := toHeavySyncer{}

	syncer.addToWaitingPulses(gen.PulseNumber())
	syncer.addToWaitingPulses(gen.PulseNumber())

	require.Equal(t, 2, len(syncer.syncWaitingPulses))
}

func TestToHeavySyncer_extractWaitingPulse(t *testing.T) {
	fPn := gen.PulseNumber()
	sPn := gen.PulseNumber()
	tPn := gen.PulseNumber()

	syncer := toHeavySyncer{}
	syncer.addToWaitingPulses(fPn)
	syncer.addToWaitingPulses(sPn)
	syncer.addToWaitingPulses(tPn)

	pn, ok := syncer.extractWaitingPulse()
	require.Equal(t, fPn, pn)
	require.Equal(t, true, ok)
	pn, ok = syncer.extractWaitingPulse()
	require.Equal(t, sPn, pn)
	require.Equal(t, true, ok)
	pn, ok = syncer.extractWaitingPulse()
	require.Equal(t, tPn, pn)
	require.Equal(t, true, ok)
	pn, ok = syncer.extractWaitingPulse()
	require.Equal(t, false, ok)
}

func TestToHeavySyncer_addToNotSentPayloads(t *testing.T) {
	var fPayload message.HeavyPayload
	var sPayload message.HeavyPayload
	var tPayload message.HeavyPayload
	fuzzer := fuzz.New().NilChance(0)
	fuzzer.Fuzz(&fPayload)
	fuzzer.Fuzz(&sPayload)
	fuzzer.Fuzz(&tPayload)

	syncer := toHeavySyncer{}

	syncer.addToNotSentPayloads(&fPayload)
	syncer.addToNotSentPayloads(&sPayload)
	syncer.addToNotSentPayloads(&tPayload)

	require.Equal(t, 3, len(syncer.notSentPayloads))
	require.Equal(t, &fPayload, syncer.notSentPayloads[0].msg)
	require.Equal(t, &sPayload, syncer.notSentPayloads[1].msg)
	require.Equal(t, &tPayload, syncer.notSentPayloads[2].msg)
}
