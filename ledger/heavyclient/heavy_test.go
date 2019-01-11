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

package heavyclient_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/pulsemanager"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/index"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/record"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
)

func TestPulseManager_SendToHeavyHappyPath(t *testing.T) {
	sendToHeavy(t, false)
}

func TestPulseManager_SendToHeavyWithRetry(t *testing.T) {
	sendToHeavy(t, true)
}

func sendToHeavy(t *testing.T, withretry bool) {
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	// TODO: test should work with any JetID (add new test?) - 14.Dec.2018 @nordicdyno
	jetID := jet.ZeroJetID
	// Mock N1: LR mock do nothing
	lrMock := testutils.NewLogicRunnerMock(t)
	lrMock.OnPulseMock.Return(nil)

	// Mock N2: we are light material
	nodeMock := network.NewNodeMock(t)
	nodeMock.RoleMock.Return(core.StaticRoleLightMaterial)
	nodeMock.IDMock.Return(core.RecordRef{})

	// Mock N3: nodenet returns mocked node (above)
	// and add stub for GetActiveNodes
	nodenetMock := network.NewNodeNetworkMock(t)
	nodenetMock.GetActiveNodesMock.Return(nil)
	nodenetMock.GetOriginMock.Return(nodeMock)

	// Mock N4: message bus for Send method
	busMock := testutils.NewMessageBusMock(t)
	busMock.OnPulseFunc = func(context.Context, core.Pulse) error {
		return nil
	}

	// Mock5: RecentStorageMock
	recentMock := recentstorage.NewRecentStorageMock(t)
	recentMock.ClearZeroTTLObjectsMock.Return()
	recentMock.GetObjectsMock.Return(nil)
	recentMock.GetRequestsMock.Return(nil)
	recentMock.ClearObjectsMock.Return()

	// Mock6: JetCoordinatorMock
	jcMock := testutils.NewJetCoordinatorMock(t)
	jcMock.LightExecutorForJetMock.Return(&core.RecordRef{}, nil)
	jcMock.MeMock.Return(core.RecordRef{})

	// Mock N7: GIL mock
	gilMock := testutils.NewGlobalInsolarLockMock(t)
	gilMock.AcquireFunc = func(context.Context) {}
	gilMock.ReleaseFunc = func(context.Context) {}

	// Mock N8: Active List Swapper mock
	alsMock := testutils.NewActiveListSwapperMock(t)
	alsMock.MoveSyncToActiveFunc = func() {}

	// Mock N8: Crypto things mock
	cryptoServiceMock := testutils.NewCryptographyServiceMock(t)
	cryptoServiceMock.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		signature := core.SignatureFromBytes(nil)
		return &signature, nil
	}
	cryptoScheme := testutils.NewPlatformCryptographyScheme()

	// mock bus.Mock method, store synced records, and calls count with HeavyRecord
	var synckeys []key
	var syncsended int
	type messageStat struct {
		size int
		keys []key
	}
	syncmessagesPerMessage := map[int]*messageStat{}
	var bussendfailed int32
	busMock.SendFunc = func(ctx context.Context, msg core.Message, ops *core.MessageSendOptions) (core.Reply, error) {
		// fmt.Printf("got msg: %T (%s)\n", msg, msg.Type())
		heavymsg, ok := msg.(*message.HeavyPayload)
		if ok {
			if withretry && atomic.AddInt32(&bussendfailed, 1) < 2 {
				return &reply.HeavyError{
					SubType: reply.ErrHeavySyncInProgress,
					Message: "retryable error",
				}, nil
			}

			syncsended++
			var size int
			var keys []key

			for _, rec := range heavymsg.Records {
				keys = append(keys, rec.K)
				size += len(rec.K) + len(rec.V)
			}
			synckeys = append(synckeys, keys...)
			syncmessagesPerMessage[syncsended] = &messageStat{
				size: size,
				keys: keys,
			}
		}
		return nil, nil
	}

	// build PulseManager
	minretry := 20 * time.Millisecond
	kb := 1 << 10
	pmconf := configuration.PulseManager{
		HeavySyncEnabled:      true,
		HeavySyncMessageLimit: 2 * kb,
		HeavyBackoff: configuration.Backoff{
			Jitter: true,
			Min:    minretry,
			Max:    minretry * 2,
			Factor: 2,
		},
		SplitThreshold: 10 * 1000 * 1000,
	}
	pm := pulsemanager.NewPulseManager(
		db,
		configuration.Ledger{
			PulseManager:    pmconf,
			LightChainLimit: 10,
		},
	)
	pm.LR = lrMock
	pm.NodeNet = nodenetMock
	pm.Bus = busMock
	pm.JetCoordinator = jcMock
	pm.GIL = gilMock
	pm.PulseStorage = storage.NewPulseStorage(db)

	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetStorageFunc = func(p core.RecordID) (r recentstorage.RecentStorage) {
		return recentMock
	}
	pm.RecentStorageProvider = provideMock

	pm.ActiveListSwapper = alsMock
	pm.CryptographyService = cryptoServiceMock
	pm.PlatformCryptographyScheme = cryptoScheme

	// Actial test logic
	// start PulseManager
	err := pm.Start(ctx)
	assert.NoError(t, err)

	// store last pulse as light material and set next one
	lastpulse := core.FirstPulseNumber + 1
	err = setpulse(ctx, pm, lastpulse)
	require.NoError(t, err)

	for i := 0; i < 2; i++ {
		// fmt.Printf("%v: call addRecords for pulse %v\n", t.Name(), lastpulse)
		addRecords(ctx, t, db, jetID, core.PulseNumber(lastpulse+i))
	}

	fmt.Println("Case1: sync after db fill and with new received pulses")
	for i := 0; i < 2; i++ {
		lastpulse++
		err = setpulse(ctx, pm, lastpulse)
		require.NoError(t, err)
	}

	fmt.Println("Case2: sync during db fill")
	for i := 0; i < 2; i++ {
		// fill DB with records, indexes (TODO: add blobs)
		addRecords(ctx, t, db, jetID, core.PulseNumber(lastpulse))

		lastpulse++
		err = setpulse(ctx, pm, lastpulse)
		require.NoError(t, err)
	}
	// set last pulse
	lastpulse++
	err = setpulse(ctx, pm, lastpulse)
	require.NoError(t, err)

	// give sync chance to complete and start sync loop again
	time.Sleep(2 * minretry)

	err = pm.Stop(ctx)
	assert.NoError(t, err)

	synckeys = uniqkeys(sortkeys(synckeys))

	recs := getallkeys(db.GetBadgerDB())
	recs = filterkeys(recs, func(k key) bool {
		return k.pulse() != 0
	})

	// fmt.Println("synckeys")
	// printkeys(synckeys, "  ")
	// fmt.Println("getallkeys")
	// printkeys(recs, "  ")
	require.Equal(t, len(recs), len(synckeys), "synced keys count are the same as records count in storage")
	assert.Equal(t, recs, synckeys, "synced keys are the same as records in storage")
}

func setpulse(ctx context.Context, pm core.PulseManager, pulsenum int) error {
	// fmt.Printf("CALL setpulse %v\n", pulsenum)
	return pm.Set(ctx, core.Pulse{PulseNumber: core.PulseNumber(pulsenum)}, true)
}

func addRecords(
	ctx context.Context,
	t *testing.T,
	db *storage.DB,
	jetID core.RecordID,
	pn core.PulseNumber,
) {
	// fmt.Printf("CALL addRecords for pulse %v\n", pn)
	// set record
	parentID, err := db.SetRecord(
		ctx,
		jetID,
		pn,
		&record.ObjectActivateRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain: testutils.RandomRef(),
			},
		},
	)
	require.NoError(t, err)

	_, err = db.SetBlob(ctx, jetID, pn, []byte("100500"))
	require.NoError(t, err)

	// set index of record
	err = db.SetObjectIndex(ctx, jetID, parentID, &index.ObjectLifeline{
		LatestState: parentID,
	})
	require.NoError(t, err)
	return
}

var (
	scopeIDLifeline = byte(1)
	scopeIDRecord   = byte(2)
	scopeIDJetDrop  = byte(3)
	scopeIDBlob     = byte(7)
)

type key []byte

func getallkeys(db *badger.DB) (records []key) {
	txn := db.NewTransaction(true)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		k := item.KeyCopy(nil)
		if key(k).pulse() == 0 {
			continue
		}
		switch k[0] {
		case
			scopeIDRecord,
			scopeIDJetDrop,
			scopeIDLifeline,
			scopeIDBlob:
			records = append(records, k)
		}
	}
	return
}

func (b key) pulse() core.PulseNumber {
	pulseStartsAt := 1
	pulseEndsAt := 1 + core.PulseNumberSize
	// if jet defined for record type
	switch b[0] {
	case
		scopeIDRecord,
		scopeIDJetDrop,
		scopeIDLifeline,
		scopeIDBlob:

		pulseStartsAt += core.RecordIDSize
		pulseEndsAt += core.RecordIDSize
	}
	return core.NewPulseNumber(b[pulseStartsAt:pulseEndsAt])
}

func (b key) String() string {
	return hex.EncodeToString(b)
}

func printkeys(keys []key, prefix string) {
	for _, k := range keys {
		fmt.Printf("%v%v (%v)\n", prefix, k, k.pulse())
	}
}

func filterkeys(keys []key, check func(key) bool) (keyout []key) {
	for _, k := range keys {
		if check(k) {
			keyout = append(keyout, k)
		}
	}
	return
}

func uniqkeys(keys []key) (keyout []key) {
	uniq := map[string]bool{}
	for _, k := range keys {
		if uniq[string(k)] {
			continue
		}
		uniq[string(k)] = true
		keyout = append(keyout, k)
	}
	return
}

func sortkeys(keys []key) []key {
	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i], keys[j]) < 0
	})
	return keys
}
