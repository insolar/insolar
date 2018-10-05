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

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGetBalance(t *testing.T) {
	memberRef := createMember(t)

	body := getResponseBody(t, postParams{
		"query_type": "get_balance",
		"reference":  memberRef,
	})

	response := &getBalanceResponse{}
	unmarshalResponse(t, body, response)

	assert.Equal(t, 1000, int(response.Amount))
	assert.Equal(t, "RUB", response.Currency)
}

func TestGetBalanceWrongRef(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "get_balance",
		"reference":  core.RandomRef(),
	})

	response := &getBalanceResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.BadRequest, response.Err.Code)
	assert.Equal(t, "Bad request", response.Err.Message)
}

// TODO: unskip test after doing errors in smart contracts
func _TestWrongReferenceInParams(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "get_balance",
		"reference":  testutils.RandomString(),
	})

	response := &baseResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.BadRequest, response.Err.Code)
	assert.Equal(t, "Bad request", response.Err.Message)
}
