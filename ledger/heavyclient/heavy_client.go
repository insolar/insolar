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
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/utils/backoff"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

// Options contains heavy client configuration params.
type Options struct {
	SyncMessageLimit int
	PulsesDeltaLimit int
	BackoffConf      configuration.Backoff
}

// JetClient heavy replication client. Replicates records for one jet.
type JetClient struct {
	Bus          core.MessageBus
	PulseStorage core.PulseStorage

	db   *storage.DB
	opts Options

	// life cycle control
	//
	startOnce sync.Once
	cancel    context.CancelFunc
	signal    chan struct{}
	// syncdone closes when syncloop is gracefully finished
	syncdone chan struct{}

	// state:
	jetID       core.RecordID
	muPulses    sync.Mutex
	leftPulses  []core.PulseNumber
	syncbackoff *backoff.Backoff
}

// NewJetClient heavy replication client constructor.
//
// First argument defines what jet it serve.
func NewJetClient(jetID core.RecordID, opts Options) *JetClient {
	jsc := &JetClient{
		jetID:       jetID,
		syncbackoff: backoffFromConfig(opts.BackoffConf),
		signal:      make(chan struct{}),
		syncdone:    make(chan struct{}),
		opts:        opts,
	}
	return jsc
}

// should be called from protected by mutex code
func (c *JetClient) updateLeftPulsesMetrics(ctx context.Context) {
	// instrumentation
	var pn core.PulseNumber
	if len(c.leftPulses) > 0 {
		pn = c.leftPulses[0]
	}
	ctx = insmetrics.InsertTag(ctx, tagJet, c.jetID.DebugString())
	stats.Record(ctx,
		statUnsyncedPulsesCount.M(int64(len(c.leftPulses))),
		statFirstUnsyncedPulse.M(int64(pn)),
	)
}

// addPulses add pulse numbers for syncing.
func (c *JetClient) addPulses(ctx context.Context, pns []core.PulseNumber) {
	c.muPulses.Lock()
	c.leftPulses = append(c.leftPulses, pns...)

	if err := c.db.SetSyncClientJetPulses(ctx, c.jetID, c.leftPulses); err != nil {
		inslogger.FromContext(ctx).Errorf(
			"attempt to persist jet sync state failed: jetID=%v: %v", c.jetID, err.Error())
	}

	c.updateLeftPulsesMetrics(ctx)
	c.muPulses.Unlock()
}

func (c *JetClient) pulsesLeft() int {
	c.muPulses.Lock()
	defer c.muPulses.Unlock()
	return len(c.leftPulses)
}

// unshiftPulse removes and returns pulse number from head of processing queue.
func (c *JetClient) unshiftPulse(ctx context.Context) *core.PulseNumber {
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

	if err := c.db.SetSyncClientJetPulses(ctx, c.jetID, c.leftPulses); err != nil {
		inslogger.FromContext(ctx).Errorf(
			"attempt to persist jet sync state failed: jetID=%v: %v", c.jetID, err.Error())
	}

	c.updateLeftPulsesMetrics(ctx)
	return &result
}

func (c *JetClient) nextPulseNumber() (core.PulseNumber, bool) {
	c.muPulses.Lock()
	defer c.muPulses.Unlock()

	if len(c.leftPulses) == 0 {
		return 0, false
	}
	return c.leftPulses[0], true
}

func (c *JetClient) runOnce(ctx context.Context) {
	// retrydelay = m.syncbackoff.ForAttempt(attempt)
	c.startOnce.Do(func() {
		// TODO: reset TraceID from context, or just don't use context?
		// (TraceID not meaningful in async sync loop)
		ctx, cancel := context.WithCancel(context.Background())
		c.cancel = cancel
		go c.syncloop(ctx)
	})
}

func (c *JetClient) syncloop(ctx context.Context) {
	inslog := inslogger.FromContext(ctx)
	defer close(c.syncdone)

	var (
		syncPN     core.PulseNumber
		hasNext    bool
		retrydelay time.Duration
	)

	finishpulse := func() {
		_ = c.unshiftPulse(ctx)
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
			syncPN, hasNext = c.nextPulseNumber()
			if hasNext {
				inslog.Debugf("synchronization next sync pulse num: %v (left=%v)", syncPN, c.leftPulses)
				break
			}

			inslog.Debug("synchronization waiting signal what new pulse happens")
			_, ok := <-c.signal
			if !ok {
				inslog.Debug("stop is called, so we are should just stop syncronization loop")
				return
			}
		}

		if isPulseNumberOutdated(ctx, c.db, c.PulseStorage, syncPN, c.opts.PulsesDeltaLimit) {
			inslog.Infof("pulse %v on jet %v is outdated, skip it", syncPN, c.jetID)
			finishpulse()
			continue
		}

		inslog.Infof("start synchronization to heavy for pulse %v", syncPN)

		shouldretry := false
		isretry := c.syncbackoff.Attempt() > 0

		syncerr := c.HeavySync(ctx, syncPN, isretry)
		if syncerr != nil {
			if heavyerr, ok := syncerr.(*reply.HeavyError); ok {
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
		} else {
			ctx = insmetrics.InsertTag(ctx, tagJet, c.jetID.DebugString())
			stats.Record(ctx,
				statSyncedPulsesCount.M(1),
			)
		}

		finishpulse()
	}

}

// Stop stops heavy client replication
func (c *JetClient) Stop(ctx context.Context) {
	// cancel should be set if client has started
	if c.cancel != nil {
		// two signals for sync loop to stop
		c.cancel()
		close(c.signal)
		// waits sync loop to stop
		<-c.syncdone
	}
}

func backoffFromConfig(bconf configuration.Backoff) *backoff.Backoff {
	return &backoff.Backoff{
		Jitter: bconf.Jitter,
		Min:    bconf.Min,
		Max:    bconf.Max,
		Factor: bconf.Factor,
	}
}

func isPulseNumberOutdated(ctx context.Context, db *storage.DB, pstore core.PulseStorage, pn core.PulseNumber, delta int) bool {
	current, err := pstore.Current(ctx)
	if err != nil {
		panic(err)
	}

	currentPulse, err := db.GetPulse(ctx, current.PulseNumber)
	if err != nil {
		panic(err)
	}

	pnPulse, err := db.GetPulse(ctx, pn)
	if err != nil {
		inslogger.FromContext(ctx).Errorf("Can't get pulse by pulse number: %v", pn)
		return true
	}

	return currentPulse.SerialNumber-delta > pnPulse.SerialNumber
}
