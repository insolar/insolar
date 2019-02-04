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

package jetcoordinator

import (
	"context"
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type jetCoordinatorSuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()

	pulseStorage *storage.PulseStorage
	pulseTracker storage.PulseTracker
	jetStorage   storage.JetStorage
	nodeStorages storage.NodeStorage
	coordinator  *JetCoordinator
}

func NewJetCoordinatorSuite() *jetCoordinatorSuite {
	return &jetCoordinatorSuite{
		Suite: suite.Suite{},
	}
}

func TestCoordinator(t *testing.T) {
	suite.Run(t, NewJetCoordinatorSuite())
}

func (s *jetCoordinatorSuite) BeforeTest(suiteName, testName string) {
	s.cm = &component.Manager{}
	s.ctx = inslogger.TestContext(s.T())

	db, cleaner := storagetest.TmpDB(s.ctx, s.T())

	s.cleaner = cleaner
	s.pulseTracker = storage.NewPulseTracker()
	s.pulseStorage = storage.NewPulseStorage()
	s.jetStorage = storage.NewJetStorage()
	s.nodeStorages = storage.NewNodeStorage()
	s.coordinator = NewJetCoordinator()
	s.coordinator.NodeNet = network.NewNodeNetworkMock(s.T())

	s.cm.Inject(
		testutils.NewPlatformCryptographyScheme(),
		db,
		s.pulseTracker,
		s.pulseStorage,
		s.jetStorage,
		s.nodeStorages,
		s.coordinator,
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

func (s *jetCoordinatorSuite) AfterTest(suiteName, testName string) {
	err := s.cm.Stop(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager stop failed", err)
	}
	s.cleaner()
}

func (s *jetCoordinatorSuite) TestJetCoordinator_QueryRole() {
	err := s.pulseTracker.AddPulse(s.ctx, core.Pulse{PulseNumber: 0, Entropy: core.Entropy{1, 2, 3}})
	require.NoError(s.T(), err)
	var nodes []core.Node
	var nodeRefs []core.RecordRef
	for i := 0; i < 100; i++ {
		ref := *core.NewRecordRef(core.DomainID, *core.NewRecordID(0, []byte{byte(i)}))
		nodes = append(nodes, storage.Node{FID: ref, FRole: core.StaticRoleLightMaterial})
		nodeRefs = append(nodeRefs, ref)
	}
	err = s.nodeStorages.SetActiveNodes(0, nodes)
	require.NoError(s.T(), err)

	objID := core.NewRecordID(0, []byte{1, 42, 123})
	err = s.jetStorage.UpdateJetTree(s.ctx, 0, true, *jet.NewID(50, []byte{1, 42, 123}))
	require.NoError(s.T(), err)

	selected, err := s.coordinator.QueryRole(s.ctx, core.DynamicRoleLightValidator, *objID, 0)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 3, len(selected))

	// Indexes are hard-coded from previously calculated values.
	assert.Equal(s.T(), []core.RecordRef{nodeRefs[16], nodeRefs[21], nodeRefs[78]}, selected)
}
