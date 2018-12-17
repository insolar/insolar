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

package pulsemanager

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/utils/backoff"
	"github.com/pkg/errors"
)

type clientOptions struct {
	syncMessageLimit int
	pulsesDeltaLimit core.PulseNumber
	backoffConf      configuration.Backoff
}

type syncClientsPool struct {
	PulseManager *PulseManager

	clientDefaults clientOptions

	// syncdone closes when all syncs is over
	// syncdone chan struct{}

	sync.Mutex
	clients map[core.RecordID]*jetSyncClient
}

func newSyncClientsPool(pm *PulseManager, clientDefaults clientOptions) *syncClientsPool {
	return &syncClientsPool{
		PulseManager:   pm,
		clientDefaults: clientDefaults,
		clients:        map[core.RecordID]*jetSyncClient{},
	}
}

func (scp *syncClientsPool) Stop(ctx context.Context) {
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

func (scp *syncClientsPool) AddPulsesToSyncClient(
	ctx context.Context,
	PulseManager *PulseManager,
	jetID core.RecordID,
	pns ...core.PulseNumber,
) *jetSyncClient {
	scp.Lock()
	client, ok := scp.clients[jetID]
	if !ok {
		client = newJetSyncClient(jetID, scp.clientDefaults)
		client.PulseManager = scp.PulseManager
		scp.clients[jetID] = client
	}
	scp.Unlock()

	client.RunOnce(ctx)

	client.AddPulses(pns)
	if len(client.signal) == 0 {
		// send signal we have new pulse
		client.signal <- struct{}{}
	}
	return client
}

type jetSyncClient struct {
	// TODO: use own db, messagebus - @nordicdyno 14/Dec/2018
	PulseManager *PulseManager

	// life cycle control
	startOnce sync.Once
	cancel    context.CancelFunc
	signal    chan struct{}
	// syncdone closes when syncloop is gracefully finished
	syncdone chan struct{}

	// state:
	jetID    core.RecordID
	muPulses sync.Mutex
	// currentInSync *core.PulseNumber
	leftPulses       []core.PulseNumber
	syncbackoff      *backoff.Backoff
	syncMessageLimit int
	pulsesDeltaLimit core.PulseNumber
	// MAYBE: we can miss next pules sync if retry the same pulse to long
	// so it probably could be better to abandon failed pulses earlier their outdate
	//  maxattempt  int
}

func newJetSyncClient(jetID core.RecordID, conf clientOptions) *jetSyncClient {
	jsc := &jetSyncClient{
		jetID:            jetID,
		syncbackoff:      backoffFromConfig(conf.backoffConf),
		syncMessageLimit: conf.syncMessageLimit,
		pulsesDeltaLimit: conf.pulsesDeltaLimit,
		signal:           make(chan struct{}),
		syncdone:         make(chan struct{}),
	}
	return jsc
}

func (c *jetSyncClient) AddPulses(pns []core.PulseNumber) {
	c.muPulses.Lock()
	c.leftPulses = append(c.leftPulses, pns...)
	c.muPulses.Unlock()
}

func (c *jetSyncClient) pulsesLeft() int {
	c.muPulses.Lock()
	defer c.muPulses.Unlock()
	left := len(c.leftPulses)
	return left
}

func (c *jetSyncClient) UnshiftPulse() *core.PulseNumber {
	c.muPulses.Lock()
	defer c.muPulses.Unlock()

	if len(c.leftPulses) == 0 {
		return nil
	}
	result := c.leftPulses[0]

	// shift array elements on one position to left
	shifted := c.leftPulses[:len(c.leftPulses)-1]
	copy(shifted, c.leftPulses[1:])
	c.leftPulses = shifted

	return &result
}

func (c *jetSyncClient) NextPulseNumber() (core.PulseNumber, bool) {
	c.muPulses.Lock()
	defer c.muPulses.Unlock()

	if len(c.leftPulses) == 0 {
		return 0, false
	}
	return c.leftPulses[0], true
}

func (c *jetSyncClient) RunOnce(ctx context.Context) {
	// retrydelay = m.syncbackoff.ForAttempt(attempt)
	c.startOnce.Do(func() {
		// TODO: reset TraceID from context, or just don't use context?
		// (TraceID not meaningful in async sync loop)
		ctx, cancel := context.WithCancel(ctx)
		c.cancel = cancel
		fmt.Printf("*START* client.syncloop for jet %v\n", c.jetID)
		go c.syncloop(ctx)
	})
}

func (c *jetSyncClient) syncloop(ctx context.Context) {
	inslog := inslogger.FromContext(ctx)
	defer close(c.syncdone)

	// TODO: use own db instance (untie from PulseManager)
	db := c.PulseManager.db

	var (
		syncPN     core.PulseNumber
		hasNext    bool
		retrydelay time.Duration
	)

	finishpulse := func() {
		_ = c.UnshiftPulse()
		c.syncbackoff.Reset()
		retrydelay = 0
	}

	for {
		select {
		case <-time.After(retrydelay):
			// for first try delay should be zero
		case <-ctx.Done():
			if c.pulsesLeft() == 0 {
				// got cancel signal and have nothing to do
				return
			}
			// client in canceled state signal but has smth to do
		}

		for {
			// if we have pulses to sync, process it
			syncPN, hasNext = c.NextPulseNumber()
			if hasNext {
				break
			}

			inslog.Debug("syncronization waiting signal what new pulses happens")
			_, ok := <-c.signal
			if !ok {
				inslog.Debug("stop is called, so we are should just stop syncronization loop")
				return
			}
			// get latest RP
			syncPN, hasNext = c.NextPulseNumber()
			if hasNext {
				// nothing to do
				continue
			}
			inslog.Debugf("syncronization next sync pulse num: %v (left=%v)", syncPN, c.leftPulses)
			break
		}

		if pulseIsOutdated(ctx, db, syncPN, c.pulsesDeltaLimit) {
			inslog.Infof("pulse %v on jet %v is outdated, skip it", syncPN, c.jetID)
			finishpulse()
			continue
		}
		inslog.Infof("start syncronization to heavy for pulse %v", syncPN)

		shouldretry := false
		isretry := c.syncbackoff.Attempt() > 0

		syncerr := c.HeavySync(ctx, syncPN, isretry)
		if syncerr != nil {
			if heavyerr, ok := syncerr.(HeavyErr); ok {
				shouldretry = heavyerr.IsRetryable()
			}

			syncerr = errors.Wrap(syncerr, "HeavySync failed")
			inslog.Errorf("%v (on attempt=%v, shouldretry=%v)",
				syncerr.Error(), c.syncbackoff.Attempt(), shouldretry)

			if shouldretry {
				retrydelay = c.syncbackoff.Duration()
				continue
			}
			// TODO: write some info to dust - 14.Dec.2018 @nordicdyno
		}

		err := db.SetReplicatedPulse(ctx, c.jetID, syncPN)
		if err != nil {
			err = errors.Wrap(err, "SetReplicatedPulse failed")
			inslog.Error(err)
			panic(err)
		}

		finishpulse()
	}

}

func pulseIsOutdated(ctx context.Context, db *storage.DB, pn core.PulseNumber, limit core.PulseNumber) bool {
	current, err := db.GetLatestPulse(ctx)
	if err != nil {
		panic(err)
	}
	return current.Pulse.PulseNumber-pn > limit
}

func (c *jetSyncClient) Stop(ctx context.Context) {
	// cancel should be set if client has started
	if c.cancel != nil {
		c.cancel()
		close(c.signal)
		<-c.syncdone
	}
}
