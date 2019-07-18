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

package integration_test

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BootstrapCalls(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)

	t.Run("message before pulse received returns error", func(t *testing.T) {
		p, _ := setCode(ctx, t, s)
		_, ok := p.(*payload.Error)
		assert.True(t, ok)
	})

	// First pulse goes in storage then interrupts.
	s.Pulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.Pulse(ctx)

	t.Run("messages after two pulses return result", func(t *testing.T) {
		p, _ := setCode(ctx, t, s)
		requirePayloadNotError(t, p)
	})
}

func Test_BasicOperations(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)

	// First pulse goes in storage then interrupts.
	s.Pulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.Pulse(ctx)

	runner := func(t *testing.T) {
		// Creating root reason request.
		var reasonID insolar.ID
		{
			p, _ := setIncomingRequest(ctx, t, s, gen.ID(), gen.ID(), record.CTSaveAsChild)
			requirePayloadNotError(t, p)
			reasonID = p.(*payload.RequestInfo).RequestID
		}
		// Save and check code.
		{
			p, sent := setCode(ctx, t, s)
			requirePayloadNotError(t, p)

			p = getCode(ctx, t, s, p.(*payload.ID).ID)
			requirePayloadNotError(t, p)
			material := record.Material{}
			err := material.Unmarshal(p.(*payload.Code).Record)
			require.NoError(t, err)
			require.Equal(t, &sent, material.Virtual)
		}
		var objectID insolar.ID
		// Set, get request.
		{
			p, sent := setIncomingRequest(ctx, t, s, gen.ID(), reasonID, record.CTSaveAsChild)
			requirePayloadNotError(t, p)

			p = getRequest(ctx, t, s, p.(*payload.RequestInfo).RequestID)
			requirePayloadNotError(t, p)
			require.Equal(t, sent, p.(*payload.Request).Request)
			objectID = p.(*payload.Request).RequestID
		}
		// Activate and check object.
		{
			p, state := activateObject(ctx, t, s, objectID)
			requirePayloadNotError(t, p)
			_, material := requireGetObject(ctx, t, s, objectID)
			require.Equal(t, &state, material.Virtual)
		}
		// Amend and check object.
		{
			p, _ := setIncomingRequest(ctx, t, s, objectID, reasonID, record.CTMethod)
			requirePayloadNotError(t, p)
			p, state := amendObject(ctx, t, s, objectID, p.(*payload.RequestInfo).RequestID)
			requirePayloadNotError(t, p)
			_, material := requireGetObject(ctx, t, s, objectID)
			require.Equal(t, &state, material.Virtual)
		}
		// Deactivate and check object.
		{
			p, _ := setIncomingRequest(ctx, t, s, objectID, reasonID, record.CTMethod)
			requirePayloadNotError(t, p)
			deactivateObject(ctx, t, s, objectID, p.(*payload.RequestInfo).RequestID)

			lifeline, _ := getObject(ctx, t, s, objectID)
			_, ok := lifeline.(*payload.Error)
			assert.True(t, ok)
		}
	}

	t.Run("happy basic", runner)

	t.Run("happy concurrent", func(t *testing.T) {
		count := 1000
		pulseAt := rand.Intn(count)
		var wg sync.WaitGroup
		wg.Add(count)
		for i := 0; i < count; i++ {
			if i == pulseAt {
				// FIXME: find out why it hangs.
				// s.Pulse(ctx)
			}
			i := i
			go func() {
				t.Run(fmt.Sprintf("iter %d", i), runner)
				wg.Done()
			}()
		}

		wg.Wait()
	})
}

func requirePayloadNotError(t *testing.T, pl payload.Payload) {
	if err, ok := pl.(*payload.Error); ok {
		t.Fatal(err)
	}
}

func requireGetObject(ctx context.Context, t *testing.T, s *Server, objectID insolar.ID) (record.Lifeline, record.Material) {
	lifelinePL, statePL := getObject(ctx, t, s, objectID)
	requirePayloadNotError(t, lifelinePL)
	requirePayloadNotError(t, statePL)

	lifeline := record.Lifeline{}
	err := lifeline.Unmarshal(lifelinePL.(*payload.Index).Index)
	require.NoError(t, err)

	state := record.Material{}
	err = state.Unmarshal(statePL.(*payload.State).Record)
	require.NoError(t, err)

	return lifeline, state
}
