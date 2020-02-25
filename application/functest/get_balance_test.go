// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"math/big"
	"testing"

	"github.com/insolar/insolar/application/testutils/launchnet"
	"github.com/insolar/insolar/insolar/gen"

	"github.com/stretchr/testify/require"
)

func TestGetBalance(t *testing.T) {
	firstMember := createMember(t)
	firstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	r := big.NewInt(0)
	require.Equal(t, r, firstBalance)
}

func TestGetBalanceWrongRef(t *testing.T) {
	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl, &launchnet.Root, "member.getBalance",
		map[string]interface{}{"reference": gen.Reference().String()})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "failed to fetch index from heavy")
}
