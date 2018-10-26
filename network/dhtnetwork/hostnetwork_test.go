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

package dhtnetwork

import (
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/cascade"
	"github.com/insolar/insolar/network/nodekeeper"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

func TestConfiguration_NewHostNetwork(t *testing.T) {

	tests := map[string]struct {
		cfg           configuration.Configuration
		expectedError bool
	}{
		// negative
		// "InvalidAddress":   {addressCfg("invalid"), true},
		// "InvalidTransport": {transportCfg("invalid"), true},

		// positive
		"DefaultConfiguration": {configuration.NewConfiguration(), false},
		/*
			"UseStun":              {stunCfg(true), false},
			"NotUseStun":           {stunCfg(false), false},
		*/
		// todo: bootstrap
	}

	cascade1 := &cascade.Cascade{}
	cfg := configuration.NewConfiguration()
	nodeID := core.NewRefFromBase58(cfg.Node.Node.ID)
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			network, err := NewHostNetwork(test.cfg, cascade1, nil, func(core.Pulse) {})
			network.SetNodeKeeper(nodekeeper.NewNodeKeeper(testutils.TestNode(nodeID)))
			if test.expectedError {
				assert.Error(t, err)
				assert.Nil(t, network)
			} else {
				// assert.NoError(t, err)
				// assert.NotNil(t, network)
				// network.Disconnect()
			}
		})
	}
}
