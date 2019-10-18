package functest

import (
	"github.com/insolar/insolar/application/testutils/launchnet"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFoundationMemberCreate(t *testing.T) {
	for _, m := range launchnet.Foundation {
		err := verifyFundsMembersAndDeposits(t, m)
		if err != nil {
			require.NoError(t, err)
		}
	}
}

func TestEnterpriseMemberCreate(t *testing.T) {
	for _, m := range launchnet.Enterprise {
		err := verifyFundsMembersAndDeposits(t, m)
		if err != nil {
			require.NoError(t, err)
		}
	}
}

func TestNetworkIncentivesMemberCreate(t *testing.T) {
	for _, m := range launchnet.NetworkIncentives {
		err := verifyFundsMembersAndDeposits(t, m)
		if err != nil {
			require.NoError(t, err)
		}
	}
}

func TestApplicationIncentivesMemberCreate(t *testing.T) {
	for _, m := range launchnet.ApplicationIncentives {
		err := verifyFundsMembersAndDeposits(t, m)
		if err != nil {
			require.NoError(t, err)
		}
	}
}

func TestFundsMemberCreate(t *testing.T) {
	for _, m := range launchnet.Funds {
		err := verifyFundsMembersAndDeposits(t, m)
		if err != nil {
			require.NoError(t, err)
		}
	}
}
