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

	"github.com/stretchr/testify/assert"
)

func TestIsAuthorized(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "is_auth",
	})

	isAuthResponse := &isAuthorized{}
	unmarshalResponse(t, body, isAuthResponse)

	assert.Equal(t, 1, isAuthResponse.Role)
	assert.NotEmpty(t, isAuthResponse.PublicKey)
	assert.Equal(t, true, isAuthResponse.NetCoordCheck)
}
