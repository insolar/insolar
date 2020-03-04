///
// Copyright 2020 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
///

package testutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/platformpolicy"
)

func GenerateNodePublicKey(t *testing.T) string {
	ks := platformpolicy.NewKeyProcessor()

	privKey, err := ks.GeneratePrivateKey()
	require.NoError(t, err)

	pubKeyStr, err := ks.ExportPublicKeyPEM(ks.ExtractPublicKey(privKey))
	require.NoError(t, err)

	return string(pubKeyStr)
}

func ExpectedError(t *testing.T, trace []string, expected string) {
	found := hasSubstring(trace, expected)
	require.True(t, found, "Expected error (%s) not found in trace: %v", expected, trace)
}

func hasSubstring(trace []string, expected string) bool {
	found := false
	for _, trace := range trace {
		found = strings.Contains(trace, expected)
		if found {
			return found
		}
	}
	return found
}
