package functest

import (
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
	"testing"
)

const insolarMigrationMember1 = "migration_daemon_0_member_keys.json"

func TestMigrationTokenNotInTheList(t *testing.T) {
	ba := testutils.RandomString()
	_, err := signedRequest(&migrationAdmin,
		"deposit.migration",
		map[string]interface{}{"amount": "1000", "ethTxHash": "TxHash", "migrationAddress": ba})
	require.Contains(t, err.Error(), "this migration daemon is not in the list")
}

func TestMigrationToken(t *testing.T) {
	ba := testutils.RandomString()
	err := createMemberWithMigrationAddress(ba)
	require.NoError(t, err)

	migrationMember, err := getMigrationDaemon(insolarMigrationMember1, 0)
	require.NoError(t, err)

	result, err := signedRequest(
		migrationMember,
		"deposit.migration",
		map[string]interface{}{"amount": "1000", "ethTxHash": "TxHash", "migrationAddress": ba})
	require.NoError(t, err)
	require.Nil(t, result)
}

func TestMigrationTokenZeroAmount(t *testing.T) {
	ba := testutils.RandomString()
	err := createMemberWithMigrationAddress(ba)
	require.NoError(t, err)

	migrationMember, err := getMigrationDaemon(insolarMigrationMember1, 0)
	require.NoError(t, err)

	result, err := signedRequest(
		migrationMember,
		"deposit.migration",
		map[string]interface{}{"amount": "0", "ethTxHash": "TxHash", "migrationAddress": ba})

	require.Contains(t, err.Error(), "amount must be greater than zero")
	require.Nil(t, result)

}

func TestMigrationTokenMistakeField(t *testing.T) {
	ba := testutils.RandomString()
	err := createMemberWithMigrationAddress(ba)
	require.NoError(t, err)

	migrationMember, err := getMigrationDaemon(insolarMigrationMember1, 0)
	require.NoError(t, err)

	result, err := signedRequest(
		migrationMember,
		"deposit.migration",
		map[string]interface{}{"amount1": "0", "ethTxHash": "TxHash", "migrationAddress": ba})
	require.Contains(t, err.Error(), " incorect input: failed to get 'amount' param")
	require.Nil(t, result)
}

func TestMigrationTokenNilValue(t *testing.T) {
	ba := testutils.RandomString()
	err := createMemberWithMigrationAddress(ba)
	require.NoError(t, err)

	migrationMember, err := getMigrationDaemon(insolarMigrationMember1, 0)
	require.NoError(t, err)

	result, err := signedRequest(migrationMember, "deposit.migration", map[string]interface{}{"amount": "20", "ethTxHash": nil, "migrationAddress": ba})
	require.Contains(t, err.Error(), "failed to get 'ethTxHash' param")
	require.Nil(t, result)

}

func TestMigrationTokenMaxAmount(t *testing.T) {
	ba := testutils.RandomString()
	err := createMemberWithMigrationAddress(ba)
	require.NoError(t, err)

	migrationMember, err := getMigrationDaemon(insolarMigrationMember1, 0)
	require.NoError(t, err)

	result, err := signedRequest(
		migrationMember,
		"deposit.migration",
		map[string]interface{}{"amount": "500000000000000000", "ethTxHash": "ethTxHash", "migrationAddress": ba})
	require.NoError(t, err)
	require.Nil(t, result)
}

func TestMigrationDoubleMigration(t *testing.T) {
	ba := testutils.RandomString()
	err := createMemberWithMigrationAddress(ba)
	require.NoError(t, err)

	migrationMember, err := getMigrationDaemon(insolarMigrationMember1, 0)
	require.NoError(t, err)

	resultMigr1, err := signedRequest(
		migrationMember, "deposit.migration", map[string]interface{}{"amount": "20", "ethTxHash": "ethTxHash", "migrationAddress": ba})
	require.NoError(t, err)
	require.Nil(t, resultMigr1)

	_, err = signedRequest(
		migrationMember,
		"deposit.migration",
		map[string]interface{}{"amount": "20", "ethTxHash": "ethTxHash", "migrationAddress": ba})
	require.Contains(t, err.Error(), "confirmed failed: confirm from the migration daemon")

}

func createMemberWithMigrationAddress(migrationAddress string) error {
	member, err := newUserWithKeys()
	if err != nil {
		return err
	}

	member.ref = root.ref
	_, err = signedRequest(&migrationAdmin, "migration.addBurnAddresses", map[string]interface{}{"burnAddresses": []string{migrationAddress}})
	if err != nil {
		return err
	}
	_, err = retryableMemberMigrationCreate(member, true)
	if err != nil {
		return err
	}
	return nil
}
