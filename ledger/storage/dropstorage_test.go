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

package storage_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/platformpolicy"
)

type dropSuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()
	db      storage.DBContext

	dropStorage storage.DropStorage
}

func NewDropSuite() *dropSuite {
	return &dropSuite{
		Suite: suite.Suite{},
	}
}

// Init and run suite
func TestDrop(t *testing.T) {
	suite.Run(t, NewDropSuite())
}

func (s *dropSuite) BeforeTest(suiteName, testName string) {
	s.cm = &component.Manager{}
	s.ctx = inslogger.TestContext(s.T())

	db, cleaner := storagetest.TmpDB(s.ctx, s.T())

	s.db = db
	s.cleaner = cleaner
	s.dropStorage = storage.NewDropStorage(10)

	s.cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		s.db,
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

func (s *dropSuite) AfterTest(suiteName, testName string) {
	err := s.cm.Stop(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager stop failed", err)
	}
	s.cleaner()
}

// testing storage CreateDrop lock logic
//
// 1) wait update transaction start, when start CreateDrop (should lock)
// 2) transaction waits start of CreateDrop call and waits 'waittime' (200ms)
// (could be unstable in really slow environments)
// 3) wait CreateDrop and transaction finished
// 4) compare finish time of CreateDrop and transaction
// CreateDrop should happen after transaction (after 'waittime' timeout happens)
func (s *dropSuite) TestStore_DropWaitWrites() {
	// s.T().Parallel()

	var txFin time.Time
	var dropFin time.Time
	waittime := time.Millisecond * 200

	var wg sync.WaitGroup
	wg.Add(2)
	txstarted := make(chan bool)
	dropwaits := make(chan bool)
	var err error
	go func() {
		err = s.db.Update(s.ctx, func(tx *storage.TransactionManager) error {
			log.Debug("start tx")
			close(txstarted)
			<-dropwaits
			time.Sleep(waittime)
			txFin = time.Now()
			return nil
		})
		log.Debug("end tx")
		wg.Done()
	}()

	go func() {
		<-txstarted
		log.Debug("start CreateDrop")
		close(dropwaits)
		_, _, dropSize, droperr := s.dropStorage.CreateDrop(s.ctx, core.TODOJetID, 0, []byte{})
		if droperr != nil {
			panic(droperr)
		}
		require.NotEqual(s.T(), 0, dropSize)
		dropFin = time.Now()
		log.Debug("end CreateDrop")
		wg.Done()
	}()
	wg.Wait()

	log.Debug("R: tx end t:", txFin)
	log.Debug("R: drop   t:", dropFin)

	require.NoError(s.T(), err)
	assert.Conditionf(s.T(), func() bool {
		return dropFin.After(txFin)
	}, "drop should happens after transaction ending")
}
