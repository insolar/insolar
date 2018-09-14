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

package servicenetwork

import (
	"strings"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/messagerouter/message"
	"github.com/prometheus/common/log"
	"github.com/stretchr/testify/assert"
)

type componentManager struct {
	components     core.Components
	interfaceNames []string
}

func (cm *componentManager) register(interfaceName string, component core.Component) {
	cm.interfaceNames = append(cm.interfaceNames, interfaceName)
	cm.components[interfaceName] = component
}

func (cm *componentManager) linkAll() {
	for _, name := range cm.interfaceNames {
		_ = cm.components[name].Start(cm.components)
	}
}

func (cm *componentManager) stopAll() {
	for i := len(cm.interfaceNames) - 1; i >= 0; i-- {
		name := cm.interfaceNames[i]
		log.Infoln("Stop component: ", name)
		err := cm.components[name].Stop()
		if err != nil {
			log.Errorln("failed to stop component ", name, " : ", err.Error())
		}
	}
}

func TestNewServiceNetwork(t *testing.T) {
	cfg := configuration.NewConfiguration()
	_, err := NewServiceNetwork(cfg.Host, cfg.Node)
	assert.NoError(t, err)
}

func TestServiceNetwork_GetAddress(t *testing.T) {
	cfg := configuration.NewConfiguration()
	network, err := NewServiceNetwork(cfg.Host, cfg.Node)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(network.GetAddress(), strings.Split(cfg.Host.Transport.Address, ":")[0]))
}

func TestServiceNetwork_GetHostNetwork(t *testing.T) {
	cfg := configuration.NewConfiguration()
	network, err := NewServiceNetwork(cfg.Host, cfg.Node)
	assert.NoError(t, err)
	host, _ := network.GetHostNetwork()
	assert.NotNil(t, host)
}

func TestServiceNetwork_SendMessage(t *testing.T) {
	cfg := configuration.NewConfiguration()
	network, err := NewServiceNetwork(cfg.Host, cfg.Node)
	assert.NoError(t, err)

	msg := &message.CallMethodMessage{
		ObjectRef: core.String2Ref("test"),
		Method:    "test",
		Arguments: []byte("test"),
	}

	network.SendMessage(core.String2Ref("test"), "test", msg)
}

func TestServiceNetwork_Start(t *testing.T) {
	cfg := configuration.NewConfiguration()
	network, err := NewServiceNetwork(cfg.Host, cfg.Node)
	assert.NoError(t, err)
	cm := componentManager{components: make(core.Components), interfaceNames: make([]string, 0)}
	cm.register("core.Network", network)
	cm.linkAll()
	cm.stopAll()
}
