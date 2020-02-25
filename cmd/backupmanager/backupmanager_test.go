// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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

// prepare
func TestNoPrepareBackupToEmptyDb(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	for i := 0; i < 3; i++ {
		output := invokeExpectFailure(t, "prepare_backup", "-d", tmpdir)
		require.Contains(t, output, "no backup start keys")
	}
}

func TestPrepareBackupToEmptyDb(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	db, err := store.NewBadgerDB(badger.DefaultOptions(tmpdir))
	require.NoError(t, err)

	var key executor.BackupStartKey
	err = db.Set(key, []byte{})
	require.NoError(t, err)

	err = db.Stop(context.Background())
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		output := invokeExpectFailure(t, "prepare_backup", "-d", tmpdir)
		require.Contains(t, output, "failed to finalizeLastPulse")
	}
}
