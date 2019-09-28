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
// +build slowtest

package main_test

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/ledger/heavy/executor"
)

var binaryPath string

func init() {
	wd, err := os.Getwd()
	binaryPath = filepath.Join(wd, "..", "..", "bin")

	if err != nil {
		panic(err.Error())
	}

	// Always rebuild backupmanager
	bashCmd := "cd " + binaryPath + " && (rm backupmanager || true) && go build ../cmd/backupmanager"
	cmd := exec.Command("bash", "-c", bashCmd)
	err = cmd.Run()
	if err != nil {
		panic(err.Error())
	}
}

func logOutput(t testing.TB, text string) {
	t.Log("Stdout+Stderr of backup manager invocation:")
	for _, line := range strings.Split(text, "\n") {
		t.Log(line)
	}
}

func invoke(args ...string) (string, error) {
	cmd := exec.Command(binaryPath+"/backupmanager", args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func invokeExpectSuccess(t testing.TB, args ...string) string {
	output, err := invoke(args...)
	if !assert.NoError(t, err) {
		logOutput(t, output)
		t.FailNow()
	}
	return output
}

func invokeExpectFailure(t testing.TB, args ...string) string {
	output, err := invoke(args...)
	if !assert.IsType(t, (*exec.ExitError)(nil), err) {
		logOutput(t, output)
		t.FailNow()
	}
	return output
}

// create
func TestNoCreateToExistingDir(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	invokeExpectSuccess(t, "create", "-d", tmpdir)

	for i := 0; i < 3; i++ {
		output, err := invoke("create", "-d", tmpdir)
		assert.IsType(t, err, (*exec.ExitError)(nil))

		require.Contains(t, output, "DB must be empty")
	}
}

func TestCreateHappyPath(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	invokeExpectSuccess(t, "create", "-d", tmpdir)

	db, err := store.NewBadgerDB(badger.DefaultOptions(tmpdir))
	require.NoError(t, err)

	var key executor.DBInitializedKey
	val, err := db.Get(key)
	require.NoError(t, err)

	timeValue := time.Time{}
	err = timeValue.UnmarshalBinary(val)
	require.NoError(t, err, "failed to parse time")
	require.False(t, timeValue.IsZero())
}

// merge
func TestFailToMergeBadBackupFile(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	invokeExpectSuccess(t, "create", "-d", tmpdir)

	bkpFile := tmpdir + "/incr.bkp"
	err = ioutil.WriteFile(bkpFile, []byte("test Data"), 0600)
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		invokeExpectFailure(t, "merge", "-t", tmpdir, "-n", bkpFile)
	}
}

func TestNoMergeToEmptyDb(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	for i := 0; i < 3; i++ {
		output := invokeExpectFailure(t, "merge", "-t", tmpdir, "-n", "TEST")
		require.Contains(t, output, "db must not be empty")
	}
}

func TestMergeNoBackupFile(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	invokeExpectSuccess(t, "create", "-d", tmpdir)

	for i := 0; i < 3; i++ {
		output := invokeExpectFailure(t, "merge", "-t", tmpdir, "-n", "TEST")
		require.Contains(t, output, "open TEST: no such file or directory")
	}
}

// prepare
func TestNoPrepareBackupToEmptyDb(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	for i := 0; i < 3; i++ {
		output := invokeExpectFailure(t, "prepare_backup", "-d", tmpdir, "-l", "TEST")
		require.Contains(t, output, "no backup start keys")
	}
}

func TestPrepareBackupToEmptyDb(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	invokeExpectSuccess(t, "create", "-d", tmpdir)

	db, err := store.NewBadgerDB(badger.DefaultOptions(tmpdir))
	require.NoError(t, err)

	var key executor.BackupStartKey
	err = db.Set(key, []byte{})
	require.NoError(t, err)

	err = db.Stop(context.Background())
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		output := invokeExpectFailure(t, "prepare_backup", "-d", tmpdir, "-l", "TEST")
		require.Contains(t, output, "failed to finalizeLastPulse")
	}
}
