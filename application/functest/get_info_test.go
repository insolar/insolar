// +build functest

package functest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetInfo(t *testing.T) {
	info := getInfo(t)
	require.NotNil(t, info)
	require.NotEqual(t, "", info.RootDomain)
	require.NotEqual(t, "", info.RootMember)
	require.NotEqual(t, "", info.NodeDomain)
}
