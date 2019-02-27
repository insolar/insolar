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

package jet

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/instrumentation/inslogger"
)

func TestJetStorage_GetJetTree(t *testing.T) {
	ctx := inslogger.TestContext(t)
	js := NewJetStorage()

	tree := js.(*jetStorage).getJetTree(ctx, 100)
	require.Equal(t, "root (level=0 actual=false)\n", tree.String())
}

func TestJetStorage_UpdateJetTree(t *testing.T) {
	ctx := inslogger.TestContext(t)
	js := NewJetStorage()

	js.UpdateJetTree(ctx, 100, true, core.RecordID(*storage.NewID(0, nil)))

	tree := js.(*jetStorage).getJetTree(ctx, 100)
	require.Equal(t, "root (level=0 actual=true)\n", tree.String())
}

func TestJetStorage_SplitJetTree(t *testing.T) {
	ctx := inslogger.TestContext(t)
	js := NewJetStorage()

	l, r, err := js.SplitJetTree(ctx, 100, core.RecordID(*storage.NewID(0, nil)))
	require.NoError(t, err)
	require.Equal(t, "[JET 1 0]", l.DebugString())
	require.Equal(t, "[JET 1 1]", r.DebugString())

	tree := js.(*jetStorage).getJetTree(ctx, 100)
	require.Equal(t, "root (level=0 actual=false)\n 0 (level=1 actual=false)\n 1 (level=1 actual=false)\n", tree.String())
}

func TestJetStorage_CloneJetTree(t *testing.T) {
	ctx := inslogger.TestContext(t)
	js := NewJetStorage()

	js.UpdateJetTree(ctx, 100, true, core.RecordID(*storage.NewID(0, nil)))

	tree := js.(*jetStorage).getJetTree(ctx, 100)
	require.Equal(t, "root (level=0 actual=true)\n", tree.String())

	tree = js.CloneJetTree(ctx, 100, 101)
	require.Equal(t, "root (level=0 actual=false)\n", tree.String())

	tree = js.(*jetStorage).getJetTree(ctx, 101)
	require.Equal(t, "root (level=0 actual=false)\n", tree.String())

	tree = js.(*jetStorage).getJetTree(ctx, 100)
	require.Equal(t, "root (level=0 actual=true)\n", tree.String())
}

func TestJetStorage_DeleteJetTree(t *testing.T) {
	ctx := inslogger.TestContext(t)
	js := NewJetStorage()

	_, _, err := js.SplitJetTree(ctx, 100, core.RecordID(*storage.NewID(0, nil)))
	require.NoError(t, err)

	js.DeleteJetTree(ctx, 100)

	tree := js.(*jetStorage).getJetTree(ctx, 100)
	require.NoError(t, err)
	require.Equal(t, "root (level=0 actual=false)\n", tree.String())
}
