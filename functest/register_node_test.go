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
	"strconv"
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
	assert.Empty(t, cert.BootstrapNodes)
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

func TestRegisterNodeNotExistRole(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "register_node",
		"public_key": TESTPUBLICKEY,
		"host":       TESTHOST,
		"roles":      []string{"some_not_fancy_role"},
	})

	response := &registerNodeResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.HandlerError, response.Err.Code)
	assert.Contains(t, response.Err.Message, "Role is not supported: some_not_fancy_role")
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

func TestRegisterNodeWithoutHost(t *testing.T) {
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

func TestRegisterNodeNoEnoughNodes(t *testing.T) {
	for i := 0; i < 5; i++ {
		body := getResponseBody(t, postParams{
			"query_type":          "register_node",
			"roles":               []string{"virtual"},
			"host":                TESTHOST,
			"public_key":          TESTPUBLICKEY,
			"bootstrap_nodes_num": 2,
		})

		response := &registerNodeResponse{}
		unmarshalResponseWithError(t, body, response)

		assert.Equal(t, api.HandlerError, response.Err.Code)
		assert.Contains(t, response.Err.Message, "There no enough nodes")
	}
}

func TestRegisterNodeBadMajorityRule(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type":          "register_node",
		"roles":               []string{"virtual"},
		"host":                TESTHOST,
		"public_key":          TESTPUBLICKEY,
		"majority_rule":       3,
		"bootstrap_nodes_num": 10,
	})

	response := &registerNodeResponse{}
	unmarshalResponseWithError(t, body, response)
	assert.Contains(t, response.Err.Message, "majorityRule must be more than 0.51 * numberOfBootstrapNodes")
}

func findPublicKey(pk string, bNodes []bootstrapNode) bool {
	for _, node := range bNodes {
		if node.PublicKey == pk {
			return true
		}
	}

	return false
}

func findHost(host string, bNodes []bootstrapNode) bool {
	for _, node := range bNodes {
		if node.Host == host {
			return true
		}
	}

	return false
}

func TestRegisterNodeWithBootstrapNodes(t *testing.T) {
	const testRole = "virtual"
	const numNodes = 5
	// Adding nodes
	for i := 0; i < numNodes; i++ {
		body := getResponseBody(t, postParams{
			"query_type": "register_node",
			"public_key": TESTPUBLICKEY + strconv.Itoa(i),
			"host":       TESTHOST + strconv.Itoa(i),
			"roles":      []string{testRole},
		})

		response := &registerNodeResponse{}
		unmarshalResponse(t, body, response)
	}

	body := getResponseBody(t, postParams{
		"query_type":          "register_node",
		"public_key":          "FFFF",
		"host":                TESTHOST + "new",
		"roles":               []string{"heavy_material"},
		"majority_rule":       numNodes,
		"bootstrap_nodes_num": numNodes,
	})

	response := &registerNodeResponse{}
	unmarshalResponse(t, body, response)

	cert := response.Certificate
	assert.Len(t, cert.Roles, 1)
	assert.Equal(t, "heavy_material", cert.Roles[0])
	assert.Equal(t, "FFFF", cert.PublicKey)
	assert.Len(t, cert.BootstrapNodes, numNodes)

	for i := 0; i < numNodes; i++ {
		tPK := TESTPUBLICKEY + strconv.Itoa(i)
		assert.True(t, findPublicKey(tPK, cert.BootstrapNodes), "Couldn't find PublicKey: %s", tPK)
		tHost := TESTHOST + strconv.Itoa(i)
		assert.True(t, findHost(tHost, cert.BootstrapNodes), "Couldn't find Host: %s", tPK)
	}
}
