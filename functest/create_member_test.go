//
// Copyright 2019 Insolar Technologies GmbH
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
//

// +build functest

package functest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateMember(t *testing.T) {
	result, err := signedRequest(&root, "contract.createMember", map[string]interface{}{"publicKey": "000"})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)
	require.NotEqual(t, "", ref)
}

func TestCreateMemberWrongKeyType(t *testing.T) {
	_, err := signedRequest(&root, "contract.createMember", map[string]interface{}{"publicKey": "000"})
	require.EqualError(t, err, "[ makeCall ] Error in called method: [ createMemberCall ]: [ Deserialize ]: EOF")
}

// no error
func _TestCreateMemberOneParameter(t *testing.T) {
	_, err := signedRequest(&root, "contract.createMember", map[string]interface{}{"publicKey": "000"})
	require.NoError(t, err)
}

func TestCreateMemberOneParameterOtherType(t *testing.T) {
	_, err := signedRequest(&root, "contract.createMember", map[string]interface{}{"publicKey": 123})
	require.EqualError(t, err, "[ makeCall ] Error in called method: [ createMemberCall ]: [ Deserialize ]: EOF")
}

func TestCreateMembersWithSameName(t *testing.T) {
	firstMemberRef, err := signedRequest(&root, "contract.createMember", map[string]interface{}{"publicKey": "000"})
	require.NoError(t, err)
	secondMemberRef, err := signedRequest(&root, "contract.createMember", map[string]interface{}{"publicKey": "000"})
	require.NoError(t, err)

	require.NotEqual(t, firstMemberRef, secondMemberRef)
}

func TestCreateMemberByNoRoot(t *testing.T) {
	member := createMember(t, "Member1")
	_, err := signedRequest(member, "contract.createMember", map[string]interface{}{"publicKey": "000"})
	require.NoError(t, err)
}
