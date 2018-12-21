/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

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
	conf, err := ReadRequestConfigFromFile("testdata/requestConfig.json")
	require.NoError(t, err)
	require.Equal(t, "CreateMember", conf.Method)

	require.Len(t, conf.Params, 2)
	require.Equal(t, float64(200), conf.Params[0])
	require.Equal(t, "Test", conf.Params[1])
}

func TestReadUserConfigFromFile(t *testing.T) {
	conf, err := ReadUserConfigFromFile("testdata/userConfig.json")
	require.NoError(t, err)
	require.Contains(t, conf.PrivateKey, "MHcCAQEEIPOsF3ujjM7jnb7V")
	require.Equal(t, "4FFB8zfQoGznSmzDxwv4njX1aR9ioL8GHSH17QXH2AFa.4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu", conf.Caller)
}
