// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/application/testutils/launchnet"
)

func TestActivateDaemonDoubleCall(t *testing.T) {
	activeDaemons := activateDaemons(t, countTwoActiveDaemon)
	for _, daemon := range activeDaemons {
		_, _, err := makeSignedRequest(launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.activateDaemon",
			map[string]interface{}{"reference": daemon.Ref})
		data := checkConvertRequesterError(t, err).Data
		require.Contains(t, data.Trace, "daemon member already activated")
	}
}

func TestActivateDeactivateDaemon(t *testing.T) {
	activeDaemons := activateDaemons(t, countTwoActiveDaemon)
	for _, daemon := range activeDaemons {
		_, err := signedRequest(t, launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.deactivateDaemon",
			map[string]interface{}{"reference": daemon.Ref})
		require.NoError(t, err)
	}

	for _, daemon := range activeDaemons {
		res, _, err := makeSignedRequest(launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.checkDaemon",
			map[string]interface{}{"reference": daemon.Ref})
		require.NoError(t, err)
		status := res.(map[string]interface{})["status"].(string)
		require.Equal(t, status, "inactive")
	}

	for _, daemon := range activeDaemons {
		_, err := signedRequest(t, launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.activateDaemon",
			map[string]interface{}{"reference": daemon.Ref})
		require.NoError(t, err)
	}

	for _, daemon := range activeDaemons {
		res, _, err := makeSignedRequest(launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.checkDaemon",
			map[string]interface{}{"reference": daemon.Ref})
		require.NoError(t, err)
		status := res.(map[string]interface{})["status"].(string)
		require.Equal(t, status, "active")
	}
}
func TestDeactivateDaemonDoubleCall(t *testing.T) {
	activeDaemons := activateDaemons(t, countTwoActiveDaemon)
	for _, daemon := range activeDaemons {
		_, _, err := makeSignedRequest(launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.deactivateDaemon",
			map[string]interface{}{"reference": daemon.Ref})
		require.NoError(t, err)
	}
	for _, daemon := range activeDaemons {
		_, _, err := makeSignedRequest(launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.deactivateDaemon",
			map[string]interface{}{"reference": daemon.Ref})
		data := checkConvertRequesterError(t, err).Data
		require.Contains(t, data.Trace, "daemon member already deactivated")
	}
}
func TestActivateAccess(t *testing.T) {

	member := createMigrationMemberForMA(t)
	_, _, err := makeSignedRequest(launchnet.TestRPCUrl, member, "migration.activateDaemon",
		map[string]interface{}{"reference": launchnet.MigrationDaemons[0].Ref})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "only migration admin can activate migration demons")
}

func TestDeactivateAccess(t *testing.T) {

	member := createMigrationMemberForMA(t)
	_, _, err := makeSignedRequest(launchnet.TestRPCUrl, member, "migration.deactivateDaemon",
		map[string]interface{}{"reference": launchnet.MigrationDaemons[0].Ref})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "only migration admin can deactivate migration demons")
}

func TestCheckDaemonAccess(t *testing.T) {

	member := createMigrationMemberForMA(t)
	_, _, err := makeSignedRequest(launchnet.TestRPCUrl, member, "migration.checkDaemon",
		map[string]interface{}{"reference": launchnet.MigrationDaemons[0].Ref})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "permission denied to information about migration daemons")
}
