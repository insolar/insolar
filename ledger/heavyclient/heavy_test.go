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

package heavyclient_test

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/artifactmanager"
	"github.com/insolar/insolar/ledger/pulsemanager"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/index"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/nodes"
	"github.com/insolar/insolar/ledger/storage/record"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type heavySuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()
	db      storage.DBContext

	jetStorage     storage.JetStorage
	nodeAccessor   *nodes.AccessorMock
	nodeSetter     *nodes.SetterMock
	pulseTracker   storage.PulseTracker
	replicaStorage storage.ReplicaStorage
	objectStorage  storage.ObjectStorage
	dropStorage    storage.DropStorage
	storageCleaner storage.Cleaner
}

func NewHeavySuite() *heavySuite {
	return &heavySuite{
		Suite: suite.Suite{},
	}
}

// Init and run suite
func TestHeavySuite(t *testing.T) {
	suite.Run(t, NewHeavySuite())
}

func (s *heavySuite) BeforeTest(suiteName, testName string) {
	s.cm = &component.Manager{}
	s.ctx = inslogger.TestContext(s.T())

	db, cleaner := storagetest.TmpDB(s.ctx, s.T())
	s.cleaner = cleaner
	s.db = db
	s.jetStorage = storage.NewJetStorage()
	s.nodeAccessor = nodes.NewAccessorMock(s.T())
	s.nodeSetter = nodes.NewSetterMock(s.T())
	s.pulseTracker = storage.NewPulseTracker()
	s.replicaStorage = storage.NewReplicaStorage()
	s.objectStorage = storage.NewObjectStorage()
	s.dropStorage = storage.NewDropStorage(10)
	s.storageCleaner = storage.NewCleaner()

	s.cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		s.db,
		s.jetStorage,
		s.nodeAccessor,
		s.nodeSetter,
		s.pulseTracker,
		s.replicaStorage,
		s.objectStorage,
		s.dropStorage,
		s.storageCleaner,
	)

	s.nodeSetter.SetMock.Return(nil)
	s.nodeAccessor.AllMock.Return(nil, nil)

	err := s.cm.Init(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager init failed", err)
	}
	err = s.cm.Start(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager start failed", err)
	}
}

func (s *heavySuite) AfterTest(suiteName, testName string) {
	err := s.cm.Stop(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager stop failed", err)
	}
	s.cleaner()
}

func (s *heavySuite) TestPulseManager_SendToHeavyHappyPath() {
	sendToHeavy(s, false)
}

func (s *heavySuite) TestPulseManager_SendToHeavyWithRetry() {
	sendToHeavy(s, true)
}

func sendToHeavy(s *heavySuite, withretry bool) {
	// TODO: test should work with any JetID (add new test?) - 14.Dec.2018 @nordicdyno
	jetID := jet.ZeroJetID
	// Mock N1: LR mock do nothing
	lrMock := testutils.NewLogicRunnerMock(s.T())
	lrMock.OnPulseMock.Return(nil)

	// Mock N2: we are light material
	nodeMock := network.NewNodeMock(s.T())
	nodeMock.RoleMock.Return(core.StaticRoleLightMaterial)
	nodeMock.IDMock.Return(core.RecordRef{})

	// Mock N3: nodenet returns mocked node (above)
	// and add stub for GetActiveNodes
	nodenetMock := network.NewNodeNetworkMock(s.T())
	nodenetMock.GetWorkingNodesMock.Return(nil)
	nodenetMock.GetOriginMock.Return(nodeMock)

	// Mock N4: message bus for Send method
	busMock := testutils.NewMessageBusMock(s.T())
	busMock.OnPulseFunc = func(context.Context, core.Pulse) error {
		return nil
	}

	// Mock5: RecentIndexStorageMock and PendingStorageMock
	recentMock := recentstorage.NewRecentIndexStorageMock(s.T())
	recentMock.GetObjectsMock.Return(nil)
	recentMock.AddObjectMock.Return()
	recentMock.DecreaseIndexTTLMock.Return([]core.RecordID{})
	recentMock.FilterNotExistWithLockMock.Return()

	pendingStorageMock := recentstorage.NewPendingStorageMock(s.T())
	pendingStorageMock.GetRequestsMock.Return(map[core.RecordID]recentstorage.PendingObjectContext{})

	// Mock6: JetCoordinatorMock
	jcMock := testutils.NewJetCoordinatorMock(s.T())
	jcMock.LightExecutorForJetMock.Return(&core.RecordRef{}, nil)
	jcMock.MeMock.Return(core.RecordRef{})

	// Mock N7: GIL mock
	gilMock := testutils.NewGlobalInsolarLockMock(s.T())
	gilMock.AcquireFunc = func(context.Context) {}
	gilMock.ReleaseFunc = func(context.Context) {}

	// Mock N8: Active List Swapper mock
	alsMock := testutils.NewActiveListSwapperMock(s.T())
	alsMock.MoveSyncToActiveFunc = func(context.Context) error { return nil }

	// Mock N9: Crypto things mock
	cryptoServiceMock := testutils.NewCryptographyServiceMock(s.T())
	cryptoServiceMock.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		signature := core.SignatureFromBytes(nil)
		return &signature, nil
	}
	cryptoScheme := testutils.NewPlatformCryptographyScheme()

	// Mock N10: ArtifactManagerMessageHandler
	artifactManagerMessageHandlerMock := testutils.NewArtifactManagerMessageHandlerMock(s.T())
	artifactManagerMessageHandlerMock.ResetEarlyRequestCircuitBreakerMock.Return()
	artifactManagerMessageHandlerMock.CloseEarlyRequestCircuitBreakerForJetMock.Return()

	// mock bus.Mock method, store synced records, and calls count with HeavyRecord
	var statMutex sync.Mutex
	var synckeys []key
	var syncsended int32
	type messageStat struct {
		size int
		keys []key
	}
	syncmessagesPerMessage := map[int32]*messageStat{}
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

			syncsendedNewVal := atomic.AddInt32(&syncsended, 1)
			var size int
			var keys []key

			for _, rec := range heavymsg.Records {
				keys = append(keys, rec.K)
				size += len(rec.K) + len(rec.V)
			}

			statMutex.Lock()
			synckeys = append(synckeys, keys...)
			syncmessagesPerMessage[syncsendedNewVal] = &messageStat{
				size: size,
				keys: keys,
			}
			statMutex.Unlock()
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
	pm.JetStorage = s.jetStorage
	pm.Nodes = s.nodeAccessor
	pm.NodeSetter = s.nodeSetter
	pm.DBContext = s.db
	pm.PulseTracker = s.pulseTracker
	pm.ReplicaStorage = s.replicaStorage
	pm.StorageCleaner = s.storageCleaner
	pm.ObjectStorage = s.objectStorage
	pm.DropStorage = s.dropStorage

	ps := storage.NewPulseStorage()
	ps.PulseTracker = s.pulseTracker
	pm.PulseStorage = ps

	pm.HotDataWaiter = artifactmanager.NewHotDataWaiterConcrete()

	providerMock := recentstorage.NewProviderMock(s.T())
	providerMock.GetIndexStorageMock.Return(recentMock)
	providerMock.GetPendingStorageMock.Return(pendingStorageMock)
	providerMock.CloneIndexStorageMock.Return()
	providerMock.ClonePendingStorageMock.Return()
	providerMock.RemovePendingStorageMock.Return()
	providerMock.DecreaseIndexesTTLMock.Return(map[core.RecordID][]core.RecordID{})
	pm.RecentStorageProvider = providerMock

	pm.ActiveListSwapper = alsMock
	pm.CryptographyService = cryptoServiceMock
	pm.PlatformCryptographyScheme = cryptoScheme

	// Actial test logic
	// start PulseManager
	err := pm.Start(s.ctx)
	assert.NoError(s.T(), err)

	// store last pulse as light material and set next one
	lastpulse := core.FirstPulseNumber + 1
	err = setpulse(s.ctx, pm, lastpulse)
	require.NoError(s.T(), err)

	for i := 0; i < 2; i++ {
		// fmt.Printf("%v: call addRecords for pulse %v\n", t.Name(), lastpulse)
		addRecords(s.ctx, s.T(), s.objectStorage, jetID, core.PulseNumber(lastpulse+i))
	}

	fmt.Println("Case1: sync after db fill and with new received pulses")
	for i := 0; i < 2; i++ {
		lastpulse++
		err = setpulse(s.ctx, pm, lastpulse)
		require.NoError(s.T(), err)
	}

	fmt.Println("Case2: sync during db fill")
	for i := 0; i < 2; i++ {
		// fill DB with records, indexes (TODO: add blobs)
		addRecords(s.ctx, s.T(), s.objectStorage, jetID, core.PulseNumber(lastpulse))

		lastpulse++
		err = setpulse(s.ctx, pm, lastpulse)
		require.NoError(s.T(), err)
	}
	// set last pulse
	lastpulse++
	err = setpulse(s.ctx, pm, lastpulse)
	require.NoError(s.T(), err)

	// give sync chance to complete and start sync loop again
	time.Sleep(2 * minretry)

	err = pm.Stop(s.ctx)
	assert.NoError(s.T(), err)

	synckeys = uniqkeys(sortkeys(synckeys))

	recs := getallkeys(s.db.GetBadgerDB())
	recs = filterkeys(recs, func(k key) bool {
		return storage.Key(k).PulseNumber() != 0
	})

	require.Equal(s.T(), len(recs), len(synckeys), "synced keys count are the same as records count in storage")
	assert.Equal(s.T(), recs, synckeys, "synced keys are the same as records in storage")
}

func setpulse(ctx context.Context, pm core.PulseManager, pulsenum int) error {
	return pm.Set(ctx, core.Pulse{PulseNumber: core.PulseNumber(pulsenum)}, true)
}

func addRecords(
	ctx context.Context,
	t *testing.T,
	objectStorage storage.ObjectStorage,
	jetID core.RecordID,
	pn core.PulseNumber,
) {
	// set record
	parentID, err := objectStorage.SetRecord(
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

	_, err = objectStorage.SetBlob(ctx, jetID, pn, []byte("100500"))
	require.NoError(t, err)

	// set index of record
	err = objectStorage.SetObjectIndex(ctx, jetID, parentID, &index.ObjectLifeline{
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

func (k key) String() string {
	return storage.Key(k).String()
}

func getallkeys(db *badger.DB) (records []key) {
	txn := db.NewTransaction(true)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		k := item.KeyCopy(nil)
		if storage.Key(k).PulseNumber() == 0 {
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

func printkeys(keys []key, prefix string) {
	for _, k := range keys {
		sk := storage.Key(k)
		fmt.Printf("%v%v (%v)\n", prefix, sk, sk.PulseNumber())
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
