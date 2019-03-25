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

package storagetest

import (
	"context"
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type pulseSuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()

	pulseTracker storage.PulseTracker
}

func NewPulseSuite() *pulseSuite {
	return &pulseSuite{
		Suite: suite.Suite{},
	}
}

// Init and run suite
func TestPulse(t *testing.T) {
	suite.Run(t, NewPulseSuite())
}

func (s *pulseSuite) BeforeTest(suiteName, testName string) {
	s.cm = &component.Manager{}
	s.ctx = inslogger.TestContext(s.T())

	db, cleaner := TmpDB(s.ctx, nil, s.T())
	s.pulseTracker = storage.NewPulseTracker()

	s.cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		db,
		s.pulseTracker,
	)

	err := s.cm.Init(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager init failed", err)
	}
	err = s.cm.Start(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager start failed", err)
	}

	s.cleaner = cleaner
}

func (s *pulseSuite) AfterTest(suiteName, testName string) {
	err := s.cm.Stop(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager stop failed", err)
	}
	s.cleaner()
}

func (s *pulseSuite) TestDB_AddPulse_IncrementsSerialNumber() {
	err := s.pulseTracker.AddPulse(s.ctx, insolar.Pulse{PulseNumber: 1})
	require.NoError(s.T(), err)
	pulse, err := s.pulseTracker.GetPulse(s.ctx, 1)
	assert.Equal(s.T(), 2, pulse.SerialNumber)

	err = s.pulseTracker.AddPulse(s.ctx, insolar.Pulse{PulseNumber: 2})
	require.NoError(s.T(), err)
	pulse, err = s.pulseTracker.GetPulse(s.ctx, 2)
	assert.Equal(s.T(), 3, pulse.SerialNumber)

	err = s.pulseTracker.AddPulse(s.ctx, insolar.Pulse{PulseNumber: 3})
	require.NoError(s.T(), err)
	pulse, err = s.pulseTracker.GetPulse(s.ctx, 3)
	assert.Equal(s.T(), 4, pulse.SerialNumber)
}
