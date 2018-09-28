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
	"fmt"
	"testing"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
)

func TestDumpAllUsers(t *testing.T) {
	createMember(t)

	body := getResponseBody(t, postParams{
		"query_type": "dump_all_users",
	})

	response := &dumpAllUsersResponse{}
	unmarshalResponse(t, body, response)

	assert.NotEqual(t, []userInfo{}, response.DumpInfo)
}

func TestDumpUser(t *testing.T) {
	memberRef := createMember(t)

	body := getResponseBody(t, postParams{
		"query_type": "dump_user_info",
		"reference":  memberRef,
	})

	response := &dumpUserInfoResponse{}
	unmarshalResponse(t, body, response)
	fmt.Println(response)

	assert.NotEmpty(t, response.DumpInfo.Member)
	assert.Equal(t, getBalance(t, memberRef), int(response.DumpInfo.Wallet))
}

func TestDumpUserWrongRef(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "dump_user_info",
		"reference":  core.RandomRef(),
	})

	response := &dumpUserInfoResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.BadRequest, response.Err.Code)
	assert.Equal(t, "Bad request", response.Err.Message)
}
