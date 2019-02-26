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

package storagetest

import (
	"context"
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/suite"
)

type jetSuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()

	dropStorage storage.DropStorage

	jetID core.RecordID
}

func NewJetSuite() *jetSuite {
	return &jetSuite{
		Suite: suite.Suite{},
	}
}

// Init and run suite
func TestJet(t *testing.T) {
	suite.Run(t, NewJetSuite())
}

func (s *jetSuite) BeforeTest(suiteName, testName string) {
	s.cm = &component.Manager{}
	s.ctx = inslogger.TestContext(s.T())
	s.jetID = core.TODOJetID

	db, cleaner := TmpDB(s.ctx, s.T())

	s.cleaner = cleaner
	s.dropStorage = storage.NewDropStorage(10)

	s.cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		db,
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

func (s *jetSuite) AfterTest(suiteName, testName string) {
	err := s.cm.Stop(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager stop failed", err)
	}
	s.cleaner()
}
