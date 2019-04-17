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
	"time"

	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/require"
)

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

func TestToHeavySyncer_addToNotSentPayloads_BackoffConstructedProperly(t *testing.T) {
	var payload message.HeavyPayload
	fuzz.New().NilChance(0).Fuzz(&payload)
	syncer := toHeavySyncer{
		conf: configuration.LightToHeavySync{
			Backoff: configuration.Backoff{
				MaxAttempts: 10,
				Max:         10,
				Factor:      1,
				Jitter:      false,
				Min:         0,
			},
		},
	}
	syncer.addToNotSentPayloads(&payload)

	require.Equal(t, 1, len(syncer.notSentPayloads))
	require.Equal(t, 0, syncer.notSentPayloads[0].backoff.Attempt())
}

func TestToHeavySyncer_extractNotSentPayload(t *testing.T) {
	var fPayload message.HeavyPayload
	var sPayload message.HeavyPayload
	var tPayload message.HeavyPayload
	fuzzer := fuzz.New().NilChance(0)
	fuzzer.Fuzz(&fPayload)
	fuzzer.Fuzz(&sPayload)
	fuzzer.Fuzz(&tPayload)

	syncer := toHeavySyncer{
		conf: configuration.LightToHeavySync{
			Backoff: configuration.Backoff{
				Max: 2,
				Min: 1,
			},
		},
	}

	syncer.addToNotSentPayloads(&fPayload)
	syncer.addToNotSentPayloads(&sPayload)
	syncer.addToNotSentPayloads(&tPayload)

	p, ok := syncer.extractNotSentPayload()
	require.Equal(t, &fPayload, p.msg)
	require.Equal(t, true, ok)
	p, ok = syncer.extractNotSentPayload()
	require.Equal(t, &sPayload, p.msg)
	require.Equal(t, true, ok)
	p, ok = syncer.extractNotSentPayload()
	require.Equal(t, &tPayload, p.msg)
	require.Equal(t, true, ok)
	p, ok = syncer.extractNotSentPayload()
	require.Equal(t, false, ok)

	require.Equal(t, 0, len(syncer.notSentPayloads))
}

func TestToHeavySyncer_extractNotSentPayload_LongTimeout(t *testing.T) {
	var fPayload message.HeavyPayload
	fuzzer := fuzz.New().NilChance(0)
	fuzzer.Fuzz(&fPayload)

	syncer := toHeavySyncer{
		conf: configuration.LightToHeavySync{
			Backoff: configuration.Backoff{
				Max: 0, // because of the Backoff realisation
				Min: 0,
			},
		},
	}

	syncer.addToNotSentPayloads(&fPayload)

	_, ok := syncer.extractNotSentPayload()
	require.Equal(t, false, ok)

	require.Equal(t, 1, len(syncer.notSentPayloads))
	require.Equal(t, &fPayload, syncer.notSentPayloads[0].msg)
}

func TestToHeavySyncer_reAddToNotSentPayloads(t *testing.T) {
	var fPayload message.HeavyPayload
	fuzz.New().NilChance(0).Fuzz(&fPayload)

	syncer := toHeavySyncer{
		conf: configuration.LightToHeavySync{
			Backoff: configuration.Backoff{
				Max:         2,
				Min:         1,
				MaxAttempts: 3,
			},
		},
	}

	syncer.addToNotSentPayloads(&fPayload)
	time.Sleep(1 * time.Millisecond)
	slot, ok := syncer.extractNotSentPayload()
	require.Equal(t, true, ok)
	require.Equal(t, &fPayload, slot.msg)

	syncer.reAddToNotSentPayloads(inslogger.TestContext(t), slot)
	slot, ok = syncer.extractNotSentPayload()
	require.Equal(t, true, ok)
	require.Equal(t, &fPayload, slot.msg)

	t.Run("Max attempts count", func(t *testing.T) {
		// Because of MaxAttempts:3
		slot.backoff.Duration()
		slot.backoff.Duration()
		slot.backoff.Duration()
		slot.backoff.Duration()

		syncer.reAddToNotSentPayloads(inslogger.TestContext(t), slot)
		require.Equal(t, 0, len(syncer.notSentPayloads))
	})

}
