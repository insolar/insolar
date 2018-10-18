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

const TESTHOST = "127.0.0.1"
const TESTPUBLICKEY = "some_fancy_public_key"

func TestRegisterNodeVirtual(t *testing.T) {
	const testRole = "virtual"
	body := getResponseBody(t, postParams{
		"query_type": "register_node",
		"public_key": TESTPUBLICKEY,
		"host":       TESTHOST,
		"roles":      []string{testRole},
	})

	response := &registerNodeResponse{}
	unmarshalResponse(t, body, response)

	cert := response.Certificate
	assert.Len(t, cert.Roles, 1)
	assert.Equal(t, testRole, cert.Roles[0])
	assert.Equal(t, TESTPUBLICKEY, cert.PublicKey)
}

func TestRegisterNodeHeavyMaterial(t *testing.T) {
	const testRole = "heavy_material"
	body := getResponseBody(t, postParams{
		"query_type": "register_node",
		"public_key": TESTPUBLICKEY,
		"host":       TESTHOST,
		"roles":      []string{testRole},
	})

	response := &registerNodeResponse{}
	unmarshalResponse(t, body, response)

	cert := response.Certificate
	assert.Len(t, cert.Roles, 1)
	assert.Equal(t, testRole, cert.Roles[0])
	assert.Equal(t, TESTPUBLICKEY, cert.PublicKey)
}

func TestRegisterNodeLightMaterial(t *testing.T) {
	const testRole = "light_material"
	body := getResponseBody(t, postParams{
		"query_type": "register_node",
		"public_key": TESTPUBLICKEY,
		"host":       TESTHOST,
		"roles":      []string{testRole},
	})

	response := &registerNodeResponse{}
	unmarshalResponse(t, body, response)

	cert := response.Certificate
	assert.Len(t, cert.Roles, 1)
	assert.Equal(t, testRole, cert.Roles[0])
	assert.Equal(t, TESTPUBLICKEY, cert.PublicKey)
}

func _TestRegisterNodeNotExistRole(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "register_node",
		"public_key": TESTPUBLICKEY,
		"host":       TESTHOST,
		"roles":      []string{"some_not_fancy_role"},
	})

	response := &registerNodeResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.HandlerError, response.Err.Code)
	assert.Equal(t, "Error: role 'some_not_fancy_role' doesn't exist", response.Err.Message)
}

func TestRegisterNodeWithoutRole(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "register_node",
		"host":       TESTHOST,
		"public_key": TESTPUBLICKEY,
	})

	response := &registerNodeResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.HandlerError, response.Err.Code)
	assert.Equal(t, "Handler error: field 'roles' is required", response.Err.Message)
}

func TestRegisterNodeWithoutPK(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "register_node",
		"host":       TESTHOST,
		"roles":      []string{"virtual"},
	})

	response := &registerNodeResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.HandlerError, response.Err.Code)
	assert.Equal(t, "Handler error: field 'public_key' is required", response.Err.Message)
}

func TestRegisterNodeWithoutPHost(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "register_node",
		"roles":      []string{"virtual"},
		"public_key": TESTPUBLICKEY,
	})

	response := &registerNodeResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.HandlerError, response.Err.Code)
	assert.Equal(t, "Handler error: field 'host' is required", response.Err.Message)
}
