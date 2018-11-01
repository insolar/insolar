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
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

func TestDumpAllUsers(t *testing.T) {
	_ = createMember(t, "Member")

	result, err := signedRequest(&root, "DumpAllUsers")
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestDumpUser(t *testing.T) {
	member := createMember(t, "Member")

	resp, err := signedRequest(&root, "DumpUserInfo", member.ref)
	assert.NoError(t, err)

	data, err := base64.StdEncoding.DecodeString(resp.(string))
	assert.NoError(t, err)

	result := struct {
		Member string
		Wallet int
	}{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)
	assert.Equal(t, "Member", result.Member)
	assert.Equal(t, 1000, result.Wallet)
}

func TestDumpUserWrongRef(t *testing.T) {
	_, err := signedRequest(&root, "DumpUserInfo", testutils.RandomRef())
	assert.EqualError(t, err, "[ DumpUserInfo ] Problem with making request: [ getUserInfoMap ] Can't get implementation: on calling main API: failed to fetch object index: storage object not found")
}

func TestDumpAllUsersNoRoot(t *testing.T) {
	member := createMember(t, "Member")

	_, err := signedRequest(member, "DumpAllUsers")
	assert.EqualError(t, err, "[ DumpUserInfo ] Only root can call this method")
}

// todo fix this deadlock
func _TestDumpUserYourself(t *testing.T) {
	member := createMember(t, "Member")

	_, err := signedRequest(member, "DumpUserInfo", member.ref)
	assert.NoError(t, err)
}

func TestDumpUserOther(t *testing.T) {
	member1 := createMember(t, "Member1")
	member2 := createMember(t, "Member2")

	_, err := signedRequest(member1, "DumpUserInfo", member2.ref)
	assert.EqualError(t, err, "[ DumpUserInfo ] You can dump only yourself")
}
