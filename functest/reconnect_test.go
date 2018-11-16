package functest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsgorundReload(t *testing.T) {
	checkAuthRequest(t)

	stopInsgorund()
	err := startInsgorund()
	require.NoError(t, err)

	checkAuthRequest(t)
}

func checkAuthRequest(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "is_auth",
	})

	isAuthResponse := &isAuthorized{}
	unmarshalResponse(t, body, isAuthResponse)

	assert.Equal(t, 1, isAuthResponse.Role)
	assert.NotEmpty(t, isAuthResponse.PublicKey)
	assert.Equal(t, true, isAuthResponse.NetCoordCheck)
}
