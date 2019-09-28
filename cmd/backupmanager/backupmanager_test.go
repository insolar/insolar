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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/executor"
)

const BackupManagerBinary = "backupmanager"

var backupManagerBinaryPath string

func findProjectRoot() (string, error) {
	workdir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if filepath.Clean(workdir) == "/" {
			return "", errors.New("no directory contains .gir")
		}
		if fileInfo, err := os.Stat(filepath.Join(workdir, ".git")); err == nil && fileInfo.IsDir() {
			return filepath.Clean(workdir), nil
		}
		workdir = filepath.Join(workdir, "..")
	}
}

func init() {
	rootDir, err := findProjectRoot()
	if err != nil {
		panic("failed to find project root: " + err.Error())
	}
	fmt.Println("Found project root:", rootDir)

	cmd := exec.Command("make", BackupManagerBinary)
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to build backupmanager:", err.Error())
		for _, line := range strings.Split(string(output), "\n") {
			fmt.Println(line)
		}
		runtime.Goexit()
	}
	fmt.Println("Successfully build backupmanager")

	backupManagerBinaryPath = filepath.Join(rootDir, "bin", BackupManagerBinary)
}

type BadgerLogger struct {
	insolar.Logger
}

func (b BadgerLogger) Warningf(fmt string, args ...interface{}) {
	b.Warnf(fmt, args...)
}

func defaultBadgerOpts(ctx context.Context, dbPath string) badger.Options {
	_, logger := inslogger.WithField(ctx, "component", "badger")
	badgerLogger := BadgerLogger{Logger: logger}

	opts := badger.DefaultOptions(dbPath)
	opts.Logger = badgerLogger
	return opts
}

func logOutput(t testing.TB, text string) {
	t.Log("Stdout+Stderr of backup manager invocation:")
	for _, line := range strings.Split(text, "\n") {
		t.Log(line)
	}
}

func invoke(args ...string) (string, error) {
	cmd := exec.Command(backupManagerBinaryPath, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func invokeExpectFailure(t testing.TB, args ...string) string {
	output, err := invoke(args...)
	if !assert.IsType(t, (*exec.ExitError)(nil), err) {
		logOutput(t, output)
		t.FailNow()
	}
	return output
}

func prepareTemporaryDir(t *testing.T, init bool) (string, func(*testing.T)) {
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	if err != nil {
		t.Fatal("Failed to create testing directory: ", err.Error())
	}

	cleanup := func(t *testing.T) {
		err := os.RemoveAll(tmpdir)
		if err != nil {
			t.Log("Failed to cleanup (remove directory): ", err.Error())
		}
	}

	invokeExpectSuccess := func(t testing.TB, args ...string) bool {
		output, err := invoke(args...)
		success := assert.NoError(t, err)
		if !success {
			logOutput(t, output)
		}
		return success
	}

	if init && !invokeExpectSuccess(t, "create", "-d", tmpdir) {
		cleanup(t)
		t.FailNow()
	}

	return tmpdir, cleanup
}

// create
func TestNoCreateToExistingDir(t *testing.T) {
	t.Parallel()

	tmpdir, cleanup := prepareTemporaryDir(t, true)
	defer cleanup(t)

	for i := 0; i < 3; i++ {
		output := invokeExpectFailure(t, "create", "-d", tmpdir)
		require.Contains(t, output, "database must be empty")
	}
}

func TestCreateHappyPath(t *testing.T) {
	t.Parallel()

	tmpdir, cleanup := prepareTemporaryDir(t, true)
	defer cleanup(t)

	{
		ctx := context.Background()

		db, err := store.NewBadgerDB(defaultBadgerOpts(ctx, tmpdir))
		require.NoError(t, err)

		var key executor.DBInitializedKey
		val, err := db.Get(key)
		require.NoError(t, err)

		timeValue := time.Time{}
		err = timeValue.UnmarshalBinary(val)
		require.NoError(t, err, "failed to parse time")
		require.False(t, timeValue.IsZero())
	}
}

// merge
func TestFailToMergeBadBackupFile(t *testing.T) {
	t.Parallel()

	tmpdir, cleanup := prepareTemporaryDir(t, true)
	defer cleanup(t)

	{
		bkpFile := tmpdir + "/incr.bkp"
		err := ioutil.WriteFile(bkpFile, []byte("test Data"), 0600)
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			invokeExpectFailure(t, "merge", "-t", tmpdir, "-n", bkpFile)
		}
	}
}

func TestNoMergeToEmptyDb(t *testing.T) {
	t.Parallel()

	tmpdir, cleanup := prepareTemporaryDir(t, false)
	defer cleanup(t)

	for i := 0; i < 3; i++ {
		output := invokeExpectFailure(t, "merge", "-t", tmpdir, "-n", "TEST")
		require.Contains(t, output, "database must not be empty")
	}
}

func TestMergeNoBackupFile(t *testing.T) {
	t.Parallel()

	tmpdir, cleanup := prepareTemporaryDir(t, true)
	defer cleanup(t)

	for i := 0; i < 3; i++ {
		output := invokeExpectFailure(t, "merge", "-t", tmpdir, "-n", "TEST")
		require.Contains(t, output, "open TEST: no such file or directory")
	}
}

// prepare
func TestNoPrepareBackupToEmptyDb(t *testing.T) {
	t.Parallel()

	tmpdir, cleanup := prepareTemporaryDir(t, false)
	defer cleanup(t)

	for i := 0; i < 3; i++ {
		output := invokeExpectFailure(t, "prepare_backup", "-d", tmpdir, "-l", "TEST")
		require.Contains(t, output, "no backup start keys")
	}
}

func TestPrepareBackupToEmptyDb(t *testing.T) {
	t.Parallel()

	tmpdir, cleanup := prepareTemporaryDir(t, true)
	defer cleanup(t)

	{
		ctx := context.Background()

		db, err := store.NewBadgerDB(defaultBadgerOpts(ctx, tmpdir))
		require.NoError(t, err)

		var key executor.BackupStartKey
		err = db.Set(key, []byte{})
		require.NoError(t, err)

		err = db.Stop(ctx)
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			output := invokeExpectFailure(t, "prepare_backup", "-d", tmpdir, "-l", "TEST")
			require.Contains(t, output, "failed to finalizeLastPulse")
		}
	}
}
