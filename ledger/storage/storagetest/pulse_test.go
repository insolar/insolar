/*
 *    Copyright 2019 Insolar
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
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDB_AddPulse_IncrementsSerialNumber(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	db, cleaner := TmpDB(ctx, t)
	defer cleaner()

	err := db.AddPulse(ctx, core.Pulse{})
	require.NoError(t, err)
	pulse, err := db.GetLatestPulse(ctx)
	assert.Equal(t, 2, pulse.SerialNumber)

	err = db.AddPulse(ctx, core.Pulse{})
	require.NoError(t, err)
	pulse, err = db.GetLatestPulse(ctx)
	assert.Equal(t, 3, pulse.SerialNumber)

	err = db.AddPulse(ctx, core.Pulse{})
	require.NoError(t, err)
	pulse, err = db.GetLatestPulse(ctx)
	assert.Equal(t, 4, pulse.SerialNumber)
}
