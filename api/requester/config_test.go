// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package requester

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadFile_BadFile(t *testing.T) {
	err := readFile("zzz", nil)
	require.EqualError(t, err, "[ readFile ] Problem with reading config: open zzz: no such file or directory")
}

func TestReadFile_NotJson(t *testing.T) {
	err := readFile("testdata/bad_json.json", nil)
	require.EqualError(t, err, "[ readFile ] Problem with unmarshaling config: invalid character ']' after object key")
}

func TestReadRequestConfigFromFile(t *testing.T) {
	params, err := ReadRequestParamsFromFile("testdata/requestConfig.json")
	require.NoError(t, err)

	require.Equal(t, "member.create", params.CallSite)
}

func TestReadUserConfigFromFile(t *testing.T) {
	conf, err := ReadUserConfigFromFile("testdata/userConfig.json")
	require.NoError(t, err)
	require.Contains(t, conf.PrivateKey, "MHcCAQEEIPOsF3ujjM7jnb7V")
	require.Equal(t, "4FFB8zfQoGznSmzDxwv4njX1aR9ioL8GHSH17QXH2AFa.4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu", conf.Caller)
}
