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

package artifactmanager

import (
	"context"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/nodes"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type heavySuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()
	db      storage.DBContext

	scheme        core.PlatformCryptographyScheme
	pulseTracker  storage.PulseTracker
	nodeStorage   nodes.Accessor
	objectStorage storage.ObjectStorage
	jetStorage    storage.JetStorage
	dropStorage   storage.DropStorage
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
	s.scheme = platformpolicy.NewPlatformCryptographyScheme()
	s.jetStorage = storage.NewJetStorage()
	s.nodeStorage = nodes.NewStorage()
	s.pulseTracker = storage.NewPulseTracker()
	s.objectStorage = storage.NewObjectStorage()
	s.dropStorage = storage.NewDropStorage(10)

	s.cm.Inject(
		s.scheme,
		s.db,
		s.jetStorage,
		s.nodeStorage,
		s.pulseTracker,
		s.objectStorage,
		s.dropStorage,
	)

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

func (s *heavySuite) TestLedgerArtifactManager_handleHeavy() {
	jetID := testutils.RandomJet()

	// prepare mock
	heavysync := testutils.NewHeavySyncMock(s.T())
	heavysync.StartMock.Return(nil)
	heavysync.StoreMock.Set(func(ctx context.Context, jetID core.RecordID, pn core.PulseNumber, kvs []core.KV) error {
		return s.db.StoreKeyValues(ctx, kvs)
	})
	heavysync.StopMock.Return(nil)

	recentIndexMock := recentstorage.NewRecentIndexStorageMock(s.T())
	recentIndexMock.AddObjectMock.Return()
	pendingMock := recentstorage.NewPendingStorageMock(s.T())
	pendingMock.RemovePendingRequestMock.Return()
	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetIndexStorageMock.Return(recentIndexMock)
	provideMock.GetPendingStorageMock.Return(pendingMock)

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleHeavyMaterial)

	// message hanler with mok
	mh := NewMessageHandler(nil, certificate)
	mh.JetStorage = s.jetStorage
	mh.Nodes = s.nodeStorage
	mh.DBContext = s.db
	mh.PulseTracker = s.pulseTracker
	mh.ObjectStorage = s.objectStorage
	mh.RecentStorageProvider = provideMock

	mh.HeavySync = heavysync

	payload := []core.KV{
		{K: []byte("ABC"), V: []byte("CDE")},
		{K: []byte("ABC"), V: []byte("CDE")},
		{K: []byte("CDE"), V: []byte("ABC")},
	}

	parcel := &message.Parcel{
		Msg: &message.HeavyPayload{
			JetID:   jetID,
			Records: payload,
		},
	}

	var err error
	_, err = mh.handleHeavyPayload(s.ctx, parcel)
	require.NoError(s.T(), err)

	badgerdb := s.db.GetBadgerDB()
	err = badgerdb.View(func(tx *badger.Txn) error {
		for _, kv := range payload {
			item, err := tx.Get(kv.K)
			if !assert.NoError(s.T(), err) {
				continue
			}
			value, err := item.Value()
			if !assert.NoError(s.T(), err) {
				continue
			}
			assert.Equal(s.T(), kv.V, value)
		}
		return nil
	})
	require.NoError(s.T(), err)
}
