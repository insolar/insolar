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

package shutdown

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
)

func TestFlag(t *testing.T) {
	ctx := inslogger.TestContext(t)
	syncT := testutils.SyncT{T: t}

	flag := NewFlag()

	require.False(&syncT, flag.IsStopped())

	syncOut := make(chan struct{}, 2)

	sync1 := make(chan struct{})
	go func() {
		waitFinish := flag.Stop(ctx)
		sync1 <- struct{}{}
		require.NotNil(&syncT, waitFinish)
		waitFinish()
		syncOut <- struct{}{}
	}()

	<-sync1
	require.True(&syncT, flag.IsStopped())

	sync2 := make(chan struct{})
	go func() {
		waitFinish := flag.Stop(ctx)
		sync2 <- struct{}{}
		require.NotNil(&syncT, waitFinish)
		waitFinish()
		syncOut <- struct{}{}
	}()

	<-sync2
	require.True(&syncT, flag.IsStopped())

	// done - without stopping
	flag.Done(ctx, func() bool { return false })

	// done - with stopping
	flag.Done(ctx, func() bool { return true })

	// done - with stopping (duplicated call)
	flag.Done(ctx, func() bool { return true })

	for i := 0; i < 2; i++ {
		select {
		case <-syncOut:
		case <-time.After(1 * time.Minute):
			require.Fail(&syncT, "failed to wait waitFinish to exit")
		}
	}
}
