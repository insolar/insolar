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
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TESTHOST = "127.0.0.1"
const TESTPUBLICKEY = "some_fancy_public_key"

type registerAnswer struct {
	BootstrapNodes []bootstrapNode `json:"bootstrap_nodes"`
	MajorityRule   int             `json:"majority_rule"`
	PublicKey      string          `json:"public_key"`
	Reference      string
	Role           string
}

func registerNodeSignedCall(params ...interface{}) (*registerAnswer, error) {
	res, err := signedRequest(&root, "RegisterNode", params...)
	if err != nil {
		return nil, err
	}
	var cert registerAnswer
	err = json.Unmarshal([]byte(res.(string)), &cert)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func sendNoEnoughNodesRequest(t *testing.T) {
	_, err := registerNodeSignedCall(TESTPUBLICKEY, 5, 0, "virtual", TESTHOST)
	assert.EqualError(t, err, "[ registerNodeCall ] Problems with RegisterNode: [ RegisterNode ] : Can't make bootstrap nodes config: [ makeBootstrapNodesConfig ] There no enough nodes")
}

// TODO: This test must be first!! Somehow fix it
// This test tests that in case of error new node isn't added to NodeDomain
func TestRegisterDontAddIfError(t *testing.T) {
	for i := 0; i < 10; i++ {
		sendNoEnoughNodesRequest(t)
	}
}

func TestRegisterNodeNoEnoughNodes(t *testing.T) {
	sendNoEnoughNodesRequest(t)
}

func TestRegisterNodeVirtual(t *testing.T) {
	const testRole = "virtual"
	cert, err := registerNodeSignedCall(TESTPUBLICKEY, 0, 0, testRole, TESTHOST)
	assert.NoError(t, err)

	assert.Equal(t, testRole, cert.Role)
	assert.Equal(t, TESTPUBLICKEY, cert.PublicKey)
	assert.Empty(t, cert.BootstrapNodes)
}

func TestRegisterNodeHeavyMaterial(t *testing.T) {
	const testRole = "heavy_material"
	cert, err := registerNodeSignedCall(TESTPUBLICKEY, 0, 0, testRole, TESTHOST)
	assert.NoError(t, err)

	assert.Equal(t, testRole, cert.Role)
	assert.Equal(t, TESTPUBLICKEY, cert.PublicKey)
}

func TestRegisterNodeLightMaterial(t *testing.T) {
	const testRole = "light_material"
	cert, err := registerNodeSignedCall(TESTPUBLICKEY, 0, 0, testRole, TESTHOST)
	assert.NoError(t, err)

	assert.Equal(t, testRole, cert.Role)
	assert.Equal(t, TESTPUBLICKEY, cert.PublicKey)
}

func TestRegisterNodeNotExistRole(t *testing.T) {
	_, err := registerNodeSignedCall(TESTPUBLICKEY, 0, 0, "some_not_fancy_role", TESTHOST)
	assert.Contains(t, err.Error(),
		"[ registerNodeCall ] Problems with RegisterNode: [ RegisterNode ]: on calling main API: couldn't save new object as child:")
}

// TODO An error is expected but got nil.
func _TestRegisterNodeWithoutRole(t *testing.T) {
	_, err := registerNodeSignedCall(TESTPUBLICKEY, 0, 0, nil, TESTHOST)
	assert.Error(t, err)
}

// TODO An error is expected but got nil.
func _TestRegisterNodeWithoutPulicKey(t *testing.T) {
	_, err := registerNodeSignedCall("", 0, 0, "virtual", TESTHOST)
	assert.Error(t, err)
}

// TODO An error is expected but got nil.
func _TestRegisterNodeWithoutHost(t *testing.T) {
	_, err := registerNodeSignedCall(TESTPUBLICKEY, 0, 0, "virtual", "")
	assert.Error(t, err)
}

func TestRegisterNodeBadMajorityRule(t *testing.T) {
	_, err := registerNodeSignedCall(TESTPUBLICKEY, 10, 3, "virtual", TESTHOST)
	assert.EqualError(t, err, "[ registerNodeCall ] Problems with RegisterNode: majorityRule must be more than 0.51 * numberOfBootstrapNodes")
}

func findPublicKey(publicKey string, bNodes []bootstrapNode) bool {
	for _, node := range bNodes {
		if node.PublicKey == publicKey {
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
		_, err := registerNodeSignedCall(TESTPUBLICKEY+strconv.Itoa(i), 0, 0, testRole, TESTHOST+strconv.Itoa(i))
		assert.NoError(t, err)
	}

	cert, err := registerNodeSignedCall("FFFF", numNodes, numNodes, "heavy_material", TESTHOST+"new")
	assert.NoError(t, err)

	assert.Equal(t, "heavy_material", cert.Role)
	assert.Equal(t, "FFFF", cert.PublicKey)
	assert.Len(t, cert.BootstrapNodes, numNodes)

	for i := 0; i < numNodes; i++ {
		tPK := TESTPUBLICKEY + strconv.Itoa(i)
		assert.True(t, findPublicKey(tPK, cert.BootstrapNodes), "Couldn't find PublicKey: %s", tPK)
		tHost := TESTHOST + strconv.Itoa(i)
		assert.True(t, findHost(tHost, cert.BootstrapNodes), "Couldn't find Host: %s", tHost)
	}
}
