/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted (subject to the limitations in the disclaimer below) provided that
 * the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *  * Neither the name of Insolar Technologies nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
 * BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
 * CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING,
 * BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package storage

import (
	"github.com/dgraph-io/badger"
	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

type testBadgerKey struct {
	id    []byte
	scope Scope
}

func (k testBadgerKey) Scope() Scope {
	return k.scope
}

func (k testBadgerKey) ID() []byte {
	return k.id
}

func TestBadgerDB_Get(t *testing.T) {
	t.Parallel()

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)

	db, err := NewBadgerDB(configuration.ServiceNetwork{CacheDirectory: tmpdir})
	require.NoError(t, err)

	var (
		key           testBadgerKey
		expectedValue []byte
	)
	f := fuzz.New().NilChance(0)
	f.Fuzz(&key)
	f.Fuzz(&expectedValue)
	err = db.db.Update(func(txn *badger.Txn) error {
		return txn.Set(append(key.Scope().Bytes(), key.ID()...), expectedValue)
	})
	require.NoError(t, err)
	value, err := db.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)
}

func TestBadgerDB_Set(t *testing.T) {
	t.Parallel()

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)

	db, err := NewBadgerDB(configuration.ServiceNetwork{CacheDirectory: tmpdir})
	require.NoError(t, err)

	var (
		key           testBadgerKey
		expectedValue []byte
		value         []byte
	)
	f := fuzz.New().NilChance(0)
	f.Fuzz(&key)
	f.Fuzz(&expectedValue)
	err = db.Set(key, expectedValue)
	assert.NoError(t, err)

	err = db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(append(key.Scope().Bytes(), key.ID()...))
		require.NoError(t, err)
		value, err = item.ValueCopy(nil)
		require.NoError(t, err)
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)
}
