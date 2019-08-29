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

package jetcoordinator

import (
	"context"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/pulse"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
)

type jetCoordinatorSuite struct {
	suite.Suite

	cm  *component.Manager
	ctx context.Context

	pulseAccessor insolarPulse.Accessor
	pulseAppender insolarPulse.Appender

	jetStorage  jet.Storage
	nodeStorage *node.AccessorMock
	coordinator *Coordinator
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

	ps := insolarPulse.NewStorageMem()
	s.pulseAppender = ps
	s.pulseAccessor = ps
	storage := jet.NewStore()
	s.jetStorage = storage
	s.nodeStorage = node.NewAccessorMock(s.T())
	s.coordinator = NewJetCoordinator(5)
	s.coordinator.OriginProvider = network.NewOriginProviderMock(s.T())

	s.cm.Inject(
		testutils.NewPlatformCryptographyScheme(),
		ps,
		storage,
		s.nodeStorage,
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
}

func (s *jetCoordinatorSuite) TestJetCoordinator_QueryRole() {
	err := s.pulseAppender.Append(s.ctx, insolar.Pulse{PulseNumber: 0, Entropy: insolar.Entropy{1, 2, 3}})
	require.NoError(s.T(), err)
	var nds []insolar.Node
	var nodeRefs []insolar.Reference
	for i := 0; i < 100; i++ {
		ref := *insolar.NewReference(*insolar.NewID(100, []byte{byte(i)}))
		nds = append(nds, insolar.Node{ID: ref, Role: insolar.StaticRoleLightMaterial})
		nodeRefs = append(nodeRefs, ref)
	}
	require.NoError(s.T(), err)

	s.nodeStorage.InRoleMock.Return(nds, nil)

	objID := insolar.NewID(100, []byte{1, 42, 123})
	err = s.jetStorage.Update(s.ctx, 0, true, *insolar.NewJetID(50, []byte{1, 42, 123}))
	require.NoError(s.T(), err)

	selected, err := s.coordinator.QueryRole(s.ctx, insolar.DynamicRoleLightValidator, *objID, 0)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 3, len(selected))

	// Indexes are hard-coded from previously calculated values.
	assert.Equal(s.T(), []insolar.Reference{nodeRefs[16], nodeRefs[21], nodeRefs[78]}, selected)
}

func TestJetCoordinator_Me(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedID := gen.Reference()
	nodeNet := network.NewNodeNetworkMock(t)
	node := network.NewNetworkNodeMock(t)
	nodeNet.GetOriginMock.Return(node)
	node.IDMock.Return(expectedID)
	jc := NewJetCoordinator(1)
	jc.OriginProvider = nodeNet

	// Act
	resultID := jc.Me()

	// Assert
	require.Equal(t, expectedID, resultID)
}

func TestNewJetCoordinator(t *testing.T) {
	t.Parallel()
	// Act
	calc := NewJetCoordinator(12)

	// Assert
	require.NotNil(t, calc)
	require.Equal(t, 12, calc.lightChainLimit)
}

func TestJetCoordinator_IsBeyondLimit_ProblemsWithTracker(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)
	pulseCalculator := insolarPulse.NewCalculatorMock(t)
	pulseCalculator.BackwardsMock.Return(insolar.Pulse{}, errors.New("it's expected"))
	pulseAccessor := insolarPulse.NewAccessorMock(t)
	pulseAccessor.LatestMock.Return(insolar.Pulse{PulseNumber: pulse.MinTimePulse + 2}, nil)
	calc := NewJetCoordinator(12)
	calc.PulseCalculator = pulseCalculator
	calc.PulseAccessor = pulseAccessor

	// Act
	res, err := calc.IsBeyondLimit(ctx, pulse.MinTimePulse+1)

	// Assert
	require.NotNil(t, err)
	require.Equal(t, false, res)
}

func TestJetCoordinator_IsBeyondLimit_OutsideOfLightChainLimit(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)

	coord := NewJetCoordinator(25)
	pulseCalculator := insolarPulse.NewCalculatorMock(t)
	pulseCalculator.BackwardsMock.Expect(ctx, pulse.MinTimePulse, 25).Return(insolar.Pulse{PulseNumber: 34}, nil)
	pulseAccessor := insolarPulse.NewAccessorMock(t)
	pulseAccessor.LatestMock.Return(insolar.Pulse{PulseNumber: pulse.MinTimePulse}, nil)
	coord.PulseCalculator = pulseCalculator
	coord.PulseAccessor = pulseAccessor

	// Act
	res, err := coord.IsBeyondLimit(ctx, 10)

	// Assert
	require.Nil(t, err)
	require.Equal(t, true, res)
}

func TestJetCoordinator_IsBeyondLimit_PulseNotFoundIsNotBeyondLimit(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	coord := NewJetCoordinator(25)
	pulseCalculator := insolarPulse.NewCalculatorMock(mc)
	pulseCalculator.BackwardsMock.Expect(ctx, pulse.MinTimePulse+2, 1).Return(insolar.Pulse{}, insolarPulse.ErrNotFound)
	pulseAccessor := insolarPulse.NewAccessorMock(mc)
	pulseAccessor.LatestMock.Return(insolar.Pulse{PulseNumber: pulse.MinTimePulse + 2}, nil)
	coord.PulseCalculator = pulseCalculator
	coord.PulseAccessor = pulseAccessor

	// Act
	res, err := coord.IsBeyondLimit(ctx, pulse.MinTimePulse+1)

	// Assert
	require.Nil(t, err)
	require.Equal(t, true, res)
}

func TestJetCoordinator_IsBeyondLimit_InsideOfLightChainLimit(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)
	coord := NewJetCoordinator(25)
	pulseCalculator := insolarPulse.NewCalculatorMock(t)
	pulseCalculator.BackwardsMock.Expect(ctx, pulse.MinTimePulse+1, 25).Return(insolar.Pulse{PulseNumber: 15}, nil)
	pulseAccessor := insolarPulse.NewAccessorMock(t)
	pulseAccessor.LatestMock.Return(insolar.Pulse{PulseNumber: pulse.MinTimePulse + 1}, nil)
	coord.PulseCalculator = pulseCalculator
	coord.PulseAccessor = pulseAccessor

	// Act
	res, err := coord.IsBeyondLimit(ctx, pulse.MinTimePulse+2)

	// Assert
	require.Nil(t, err)
	require.Equal(t, false, res)
}

func TestJetCoordinator_NodeForJet_CheckLimitFailed(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)
	pulseCalculator := insolarPulse.NewCalculatorMock(t)
	pulseCalculator.BackwardsMock.Return(insolar.Pulse{}, errors.New("it's expected"))
	pulseAccessor := insolarPulse.NewAccessorMock(t)
	pulseAccessor.LatestMock.Return(insolar.Pulse{PulseNumber: pulse.MinTimePulse + 2}, nil)
	calc := NewJetCoordinator(12)
	calc.PulseCalculator = pulseCalculator
	calc.PulseAccessor = pulseAccessor

	// Act
	res, err := calc.NodeForJet(ctx, insolar.ID(gen.JetID()), pulse.MinTimePulse+1)

	// Assert
	require.NotNil(t, err)
	require.Nil(t, res)
}

func TestJetCoordinator_NodeForJet_GoToHeavy(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)
	pulseCalculator := insolarPulse.NewCalculatorMock(t)
	pulseCalculator.BackwardsMock.Return(insolar.Pulse{PulseNumber: 11}, nil)
	pulseAccessor := insolarPulse.NewAccessorMock(t)
	generator := entropygenerator.StandardEntropyGenerator{}
	pulseAccessor.LatestMock.Return(insolar.Pulse{PulseNumber: pulse.MinTimePulse, Entropy: generator.GenerateEntropy()}, nil)

	expectedID := insolar.NewReference(gen.ID())
	activeNodesStorageMock := node.NewAccessorMock(t)
	activeNodesStorageMock.InRoleMock.Set(func(p insolar.PulseNumber, p1 insolar.StaticRole) (r []insolar.Node, r1 error) {
		require.Equal(t, pulse.MinTimePulse, int(p))
		require.Equal(t, insolar.StaticRoleHeavyMaterial, p1)

		return []insolar.Node{{ID: *expectedID}}, nil
	})

	coord := NewJetCoordinator(25)
	coord.PulseCalculator = pulseCalculator
	coord.Nodes = activeNodesStorageMock
	coord.PlatformCryptographyScheme = platformpolicy.NewPlatformCryptographyScheme()
	coord.PulseAccessor = pulseAccessor

	// Act
	resNode, err := coord.NodeForJet(ctx, gen.ID(), 10)

	// Assert
	require.Nil(t, err)
	require.Equal(t, expectedID, resNode)
}

func TestJetCoordinator_NodeForJet_GoToLight(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)

	pulseCalculator := insolarPulse.NewCalculatorMock(t)
	pulseCalculator.BackwardsMock.Return(insolar.Pulse{PulseNumber: pulse.MinTimePulse - 100}, nil)
	pulseAccessor := insolarPulse.NewAccessorMock(t)
	generator := entropygenerator.StandardEntropyGenerator{}
	pulseAccessor.LatestMock.Return(insolar.Pulse{PulseNumber: insolar.PulseNumber(pulse.MinTimePulse + 1), Entropy: generator.GenerateEntropy()}, nil)

	expectedID := insolar.NewReference(gen.ID())
	activeNodesStorageMock := node.NewAccessorMock(t)
	activeNodesStorageMock.InRoleMock.Set(func(p insolar.PulseNumber, p1 insolar.StaticRole) (r []insolar.Node, r1 error) {
		require.Equal(t, pulse.MinTimePulse+1, int(p))
		require.Equal(t, insolar.StaticRoleLightMaterial, p1)

		return []insolar.Node{{ID: *expectedID}}, nil
	})

	coord := NewJetCoordinator(25)
	coord.PulseAccessor = pulseAccessor
	coord.PulseCalculator = pulseCalculator
	coord.Nodes = activeNodesStorageMock
	coord.PlatformCryptographyScheme = platformpolicy.NewPlatformCryptographyScheme()

	// Act
	resNode, err := coord.NodeForJet(ctx, insolar.ID(gen.JetID()), pulse.MinTimePulse+1)

	// Assert
	require.Nil(t, err)
	require.Equal(t, expectedID, resNode)
}
