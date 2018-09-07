/*
 *    Copyright 2018 INS Ecosystem
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

package hostnetwork

import (
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/stretchr/testify/assert"
)

func addressCfg(address string) configuration.HostNetwork {
	cfg := configuration.NewConfiguration().Host
	cfg.Transport.Address = address
	return cfg
}

func stunCfg(useStun bool) configuration.HostNetwork {
	cfg := configuration.NewConfiguration().Host
	cfg.Transport.BehindNAT = useStun
	return cfg
}

func transportCfg(tr string) configuration.HostNetwork {
	cfg := configuration.NewConfiguration().Host
	cfg.Transport.Protocol = tr
	return cfg
}

func TestConfiguration_NewHostNetwork(t *testing.T) {

	tests := map[string]struct {
		cfg           configuration.HostNetwork
		expectedError bool
	}{
		// negative
		// "InvalidAddress":   {addressCfg("invalid"), true},
		// "InvalidTransport": {transportCfg("invalid"), true},

		// positive
		"DefaultConfiguration": {configuration.NewConfiguration().Host, false},
		/*
			"UseStun":              {stunCfg(true), false},
			"NotUseStun":           {stunCfg(false), false},
		*/
		// todo: bootstrap
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			network, err := NewHostNetwork(test.cfg)
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

	panicConf := configuration.NewHostNetwork()
	panicConf.Transport.BehindNAT = false
	panicConf.Transport.Address = "0.0.0.0:0"

	assert.PanicsWithValue(t, "hostnetwork.NewHostNetwork: \n Couldn't start at 0.0.0.0", func() { NewHostNetwork(panicConf) })
}
