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

package heavyclient

import (
	"context"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage"
)

// Pool manages state of heavy sync clients (one client per jet id).
type Pool struct {
	Bus          core.MessageBus
	PulseStorage core.PulseStorage
	db           *storage.DB

	ClientDefaults Options

	sync.Mutex
	clients map[core.RecordID]*JetClient
}

// NewPool constructor of new pool.
func NewPool(db *storage.DB, clientDefaults Options) *Pool {
	return &Pool{
		db:             db,
		ClientDefaults: clientDefaults,
		clients:        map[core.RecordID]*JetClient{},
	}
}

// Stop send stop signals to all managed heavy clients and waits when until all of them will stop.
func (scp *Pool) Stop(ctx context.Context) {
	scp.Lock()
	defer scp.Unlock()

	var wg sync.WaitGroup
	wg.Add(len(scp.clients))
	for _, c := range scp.clients {
		c := c
		go func() {
			c.Stop(ctx)
			wg.Done()
		}()
	}
	wg.Wait()
}

// AddPulsesToSyncClient add pulse numbers to the end of jet's heavy client queue.
//
// Bool flag 'shouldrun' controls should heavy client be started (if not already) or not.
func (scp *Pool) AddPulsesToSyncClient(
	ctx context.Context,
	jetID core.RecordID,
	shouldrun bool,
	pns ...core.PulseNumber,
) *JetClient {
	scp.Lock()
	client, ok := scp.clients[jetID]
	if !ok {
		client = NewJetClient(jetID, scp.ClientDefaults)
		client.db = scp.db
		client.Bus = scp.Bus
		client.PulseStorage = scp.PulseStorage

		scp.clients[jetID] = client
	}
	scp.Unlock()

	client.addPulses(pns)

	if shouldrun {
		client.runOnce(ctx)
		if len(client.signal) == 0 {
			// send signal we have new pulse
			client.signal <- struct{}{}
		}
	}
	return client
}
