// +build functest

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

package functest

import (
	"testing"

	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestGetBalance(t *testing.T) {
	firstMember := createMember(t, "Member1")
	firstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	require.Equal(t, 1000, firstBalance)
}

func TestGetBalanceWrongRef(t *testing.T) {
	_, err := getBalance(&root, testutils.RandomRef().String())
	require.Contains(t, err.Error(), "[ getBalance ] : [ GetDelegate ] on calling main API")
}
