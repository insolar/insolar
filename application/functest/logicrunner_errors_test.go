// +build functest_endless_abandon

package functest

import (
	"testing"

	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/applicationbase/testutils/testrequest"
	"github.com/stretchr/testify/require"
)

func TestPanicIsSystemError(t *testing.T) {
	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.New",
		map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)

	_, err = testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, &Root, "first.Panic",
		map[string]interface{}{"reference": ref})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "CallMethod returns error")
	require.Contains(t, data.Trace, "AAAAAAAA!")

	_, err = testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, &Root, "first.NewPanic",
		map[string]interface{}{})
	data = checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "CallMethod returns error")
	require.Contains(t, data.Trace, "BBBBBBBB!")
}

func TestContractWithEmbeddedConstructor(t *testing.T) {
	_ = callConstructor(t, "first", "NewZero")
	_, _, err := testrequest.MakeSignedRequest(
		launchnet.TestRPCUrlPublic,
		&Root,
		"second.NewWithOne",
		map[string]interface{}{
			"oneNumber": "10",
		})
	data := checkConvertRequesterError(t, err).Data

	require.Contains(t, data.Trace, "object is not activated")
}
