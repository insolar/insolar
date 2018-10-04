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
	"github.com/stretchr/testify/assert"
)

func TestRegisterNodeVirtual(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "register_node",
		"public_key": "some_fancy_public_key",
		"role":       "virtual",
	})

	response := &registerNodeResponse{}
	unmarshalResponse(t, body, response)

	nodeRef := response.Reference
	assert.NotEqual(t, "", nodeRef)
}

func TestRegisterNodeHeavyMaterial(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "register_node",
		"public_key": "some_fancy_public_key",
		"role":       "heavy_material",
	})

	response := &registerNodeResponse{}
	unmarshalResponse(t, body, response)

	nodeRef := response.Reference
	assert.NotEqual(t, "", nodeRef)
}

func TestRegisterNodeLightMaterial(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "register_node",
		"public_key": "some_fancy_public_key",
		"role":       "light_material",
	})

	response := &registerNodeResponse{}
	unmarshalResponse(t, body, response)

	nodeRef := response.Reference
	assert.NotEqual(t, "", nodeRef)
}

func TestRegisterNodeWithoutRole(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "register_node",
		"public_key": "some_fancy_public_key",
	})

	response := &registerNodeResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.HandlerError, response.Err.Code)
	assert.Equal(t, "Handler error: field 'role' is required", response.Err.Message)
}

func TestRegisterNodeWithoutPK(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "register_node",
		"role":       "virtual",
	})

	response := &registerNodeResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.HandlerError, response.Err.Code)
	assert.Equal(t, "Handler error: field 'public_key' is required", response.Err.Message)
}
