/*
 *    Copyright 2019 Insolar Technologies
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
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"go.opencensus.io/stats"
	"golang.org/x/sync/singleflight"
)

// Pool manages state of heavy sync clients (one client per jet id).
type Pool struct {
	bus            core.MessageBus
	pulseStorage   core.PulseStorage
	pulseTracker   storage.PulseTracker
	replicaStorage storage.ReplicaStorage
	dbContext      storage.DBContext

	clientDefaults Options

	sync.Mutex
	clients map[core.RecordID]*JetClient

	cleanupGroup singleflight.Group
}

// NewPool constructor of new pool.
func NewPool(
	bus core.MessageBus,
	pulseStorage core.PulseStorage,
	tracker storage.PulseTracker,
	replicaStorage storage.ReplicaStorage,
	dbContext storage.DBContext,
	clientDefaults Options,
) *Pool {
	return &Pool{
		bus:            bus,
		pulseStorage:   pulseStorage,
		pulseTracker:   tracker,
		replicaStorage: replicaStorage,
		clientDefaults: clientDefaults,
		dbContext:      dbContext,
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
		client = NewJetClient(
			scp.replicaStorage,
			scp.bus,
			scp.pulseStorage,
			scp.pulseTracker,
			scp.dbContext,
			jetID,
			scp.clientDefaults,
		)

		scp.clients[jetID] = client
	}
	scp.Unlock()

	client.addPulses(ctx, pns)

	if shouldrun {
		client.runOnce(ctx)
		if len(client.signal) == 0 {
			// send signal we have new pulse
			client.signal <- struct{}{}
		}
	}
	return client
}

func (scp *Pool) AllClients(ctx context.Context) []*JetClient {
	scp.Lock()
	defer scp.Unlock()
	clients := make([]*JetClient, 0, len(scp.clients))
	for _, c := range scp.clients {
		clients = append(clients, c)
	}
	return clients
}

func (scp *Pool) LightCleanup(ctx context.Context, untilPN core.PulseNumber, rsp recentstorage.Provider) error {
	inslog := inslogger.FromContext(ctx)
	start := time.Now()
	defer func() {
		latency := time.Since(start)
		inslog.Infof("cleanLightData db clean phase time spend=%v", latency)
		stats.Record(ctx, statCleanLatencyDB.M(latency.Nanoseconds()/1e6))

	}()

	_, err, _ := scp.cleanupGroup.Do("lightcleanup", func() (interface{}, error) {
		startCleanup := time.Now()
		defer func() {
			latency := time.Since(startCleanup)
			inslog.Infof("cleanLightData db clean phase job time spend=%v", latency)
		}()

		// This is how we can get all jets served by
		// jets, err := scp.db.GetAllSyncClientJets(ctx)

		allClients := scp.AllClients(ctx)
		var wg sync.WaitGroup
		wg.Add(len(allClients))
		for _, c := range allClients {
			jetID := c.jetID
			go func() {
				defer wg.Done()
				inslogger.FromContext(ctx).Debugf("Start light cleanup, until pulse = %v, jet = %v",
					untilPN, jetID.DebugString())
				rmStat, err := scp.dbContext.RemoveAllForJetUntilPulse(ctx, jetID, untilPN, rsp.GetStorage(ctx, jetID))
				if err != nil {
					inslogger.FromContext(ctx).Errorf("Error on light cleanup, until pulse = %v, jet = %v: %v",
						untilPN, jetID.DebugString(), err)
					return
				}
				inslogger.FromContext(ctx).Debugf("End light cleanup, rm stat=%#v (until pulse = %v, jet = %v)",
					rmStat, untilPN, jetID.DebugString())
			}()
		}
		wg.Wait()
		return nil, nil
	})
	return err
}
