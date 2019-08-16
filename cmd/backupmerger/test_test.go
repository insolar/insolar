///
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
///

package main_test

import (
	"context"
	"encoding/binary"
	"io/ioutil"
	"os"
	"testing"

	"github.com/insolar/insolar/insolar/store"
	"github.com/stretchr/testify/require"
)

type testKey struct {
	id int64
}

func (t *testKey) ID() []byte {
	bs := make([]byte, 8)
	binary.PutVarint(bs, t.id)
	return bs
}

func (t *testKey) Scope() store.Scope {
	return store.ScopeJetDrop
}

func TestMakeBackupFile(t *testing.T) {

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	db, err := store.NewBadgerDB(tmpdir)
	require.NoError(t, err)
	defer db.Stop(context.Background())

	for i := int64(0); i < 20; i++ {
		err = db.Set(&testKey{id: i}, []byte{})
		require.NoError(t, err)
	}

	bkpKile, err := os.OpenFile("./incr.bkp", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	require.NoError(t, err)
	_, err = db.Backup(bkpKile, 0)
	require.NoError(t, err)
}
