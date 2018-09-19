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
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/eventbus/event"
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

	msg := &event.CallMethodMessage{
		ObjectRef: core.NewRefFromBase58("test"),
		Method:    "test",
		Arguments: []byte("test"),
	}

	network.SendMessage(core.NewRefFromBase58("test"), "test", msg)
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

func mockConfiguration(host string, bootstrapHosts []string, nodeID string) (configuration.HostNetwork, configuration.NodeNetwork) {
	transport := configuration.Transport{Protocol: "UTP", Address: host, BehindNAT: false}
	h := configuration.HostNetwork{
		Transport:      transport,
		IsRelay:        false,
		BootstrapHosts: bootstrapHosts,
	}

	n := configuration.NodeNetwork{Node: &configuration.Node{ID: nodeID}}

	return h, n
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return true // completed normally
	case <-time.After(timeout):
		return false // timed out
	}
}

func TestServiceNetwork_SendMessage2(t *testing.T) {
	firstNodeId := "4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj"
	secondNodeId := "53jNWvey7Nzyh4ZaLdJDf3SRgoD4GpWuwHgrgvVVGLbDkk3A7cwStSmBU2X7s4fm6cZtemEyJbce9dM9SwNxbsxf"

	firstNode, _ := NewServiceNetwork(mockConfiguration(
		"127.0.0.1:10000",
		[]string{"127.0.0.1:10001"},
		firstNodeId))
	secondNode, _ := NewServiceNetwork(mockConfiguration(
		"127.0.0.1:10001",
		nil,
		secondNodeId))

	secondNode.Start(nil)
	firstNode.Start(nil)

	defer func() {
		firstNode.Stop()
		secondNode.Stop()
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	secondNode.RemoteProcedureRegister("test", func(args [][]byte) ([]byte, error) {
		wg.Done()
		return nil, nil
	})

	msg := &event.CallMethodMessage{
		ObjectRef: core.NewRefFromBase58("test"),
		Method:    "test",
		Arguments: []byte("test"),
	}

	firstNode.SendMessage(core.NewRefFromBase58(secondNodeId), "test", msg)
	success := waitTimeout(&wg, 20*time.Millisecond)

	assert.True(t, success)
}

func TestServiceNetwork_SendCascadeMessage(t *testing.T) {
	firstNodeId := "4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj"
	secondNodeId := "53jNWvey7Nzyh4ZaLdJDf3SRgoD4GpWuwHgrgvVVGLbDkk3A7cwStSmBU2X7s4fm6cZtemEyJbce9dM9SwNxbsxf"

	firstNode, _ := NewServiceNetwork(mockConfiguration(
		"127.0.0.1:10000",
		[]string{"127.0.0.1:10001"},
		firstNodeId))
	secondNode, _ := NewServiceNetwork(mockConfiguration(
		"127.0.0.1:10001",
		nil,
		secondNodeId))

	secondNode.Start(nil)
	firstNode.Start(nil)

	defer func() {
		firstNode.Stop()
		secondNode.Stop()
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	secondNode.RemoteProcedureRegister("test", func(args [][]byte) ([]byte, error) {
		wg.Done()
		return nil, nil
	})

	msg := &event.CallMethodMessage{
		ObjectRef: core.NewRefFromBase58("test"),
		Method:    "test",
		Arguments: []byte("test"),
	}

	c := core.Cascade{
		NodeIds:           []core.RecordRef{core.NewRefFromBase58(secondNodeId)},
		ReplicationFactor: 2,
		Entropy:           core.Entropy{0},
	}

	firstNode.SendCascadeMessage(c, "test", msg)
	success := waitTimeout(&wg, 20*time.Millisecond)

	assert.True(t, success)
}

func TestServiceNetwork_SendCascadeMessage2(t *testing.T) {
	nodeIds := []core.RecordRef{
		core.NewRefFromBase58("4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj"),
		core.NewRefFromBase58("53jNWvey7Nzyh4ZaLdJDf3SRgoD4GpWuwHgrgvVVGLbDkk3A7cwStSmBU2X7s4fm6cZtemEyJbce9dM9SwNxbsxf"),
		core.NewRefFromBase58("9uE5MEWQB2yfKY8kTgTNovWii88D4anmf7GAiovgcxx6Uc6EBmZ212mpyMa1L22u9TcUUd94i8LvULdsdBoG8ed"),
		core.NewRefFromBase58("4qXdYkfL9U4tL3qRPthdbdajtafR4KArcXjpyQSEgEMtpuin3t8aZYmMzKGRnXHBauytaPQ6bfwZyKZzRPpR6gyX"),
		core.NewRefFromBase58("5q5rnvayXyKszoWofxp4YyK7FnLDwhsqAXKxj6H7B5sdEsNn4HKNFoByph4Aj8rGptdWL54ucwMQrySMJgKavxX1"),
		core.NewRefFromBase58("5tsFDwNLMW4GRHxSbBjjxvKpR99G4CSBLRqZAcpqdSk5SaeVcDL3hCiyjjidCRJ7Lu4VZoANWQJN2AgPvSRgCghn"),
		core.NewRefFromBase58("48UWM6w7YKYCHoP7GHhogLvbravvJ6bs4FGETqXfgdhF9aPxiuwDWwHipeiuNBQvx7zyCN9wFxbuRrDYRoAiw5Fj"),
		core.NewRefFromBase58("5owQeqWyHcobFaJqS2BZU2o2ZRQ33GojXkQK6f8vNLgvNx6xeWRwenJMc53eEsS7MCxrpXvAhtpTaNMPr3rjMHA"),
		core.NewRefFromBase58("xF12WfbkcWrjrPXvauSYpEGhkZT2Zha53xpYh5KQdmGHMywJNNgnemfDN2JfPV45aNQobkdma4dsx1N7Xf5wCJ9"),
		core.NewRefFromBase58("4VgDz9o23wmYXN9mEiLnnsGqCEEARGByx1oys2MXtC6M94K85ZpB9sEJwiGDER61gHkBxkwfJqtg9mAFR7PQcssq"),
		core.NewRefFromBase58("48g7C8QnH2CGMa62sNaL1gVVyygkto8EbMRHv168psCBuFR2FXkpTfwk4ZwpY8awFFXKSnWspYWWQ7sMMk5W7s3T"),
		core.NewRefFromBase58("Lvssptdwq7tatd567LUfx2AgsrWZfo4u9q6FJgJ9BgZK8cVooZv2A7F7rrs1FS5VpnTmXhr6XihXuKWVZ8i5YX9"),
	}

	prefix := "127.0.0.1:"
	port := 10000
	bootstrapNodes := nodeIds[len(nodeIds)-2:]
	bootstrapHosts := make([]string, 0)
	var wg sync.WaitGroup
	wg.Add(11)
	services := make([]*ServiceNetwork, 0)

	defer func() {
		for _, service := range services {
			service.Stop()
		}
	}()

	// init node and register test function
	initService := func(node string, bHosts []string) (service *ServiceNetwork, host string) {
		host = prefix + strconv.Itoa(port)
		service, _ = NewServiceNetwork(mockConfiguration(host, bHosts, node))
		service.Start(nil)
		service.RemoteProcedureRegister("test", func(args [][]byte) ([]byte, error) {
			wg.Done()
			return nil, nil
		})
		port++
		services = append(services, service)
		return
	}

	for _, node := range bootstrapNodes {
		_, host := initService(node.String(), nil)
		bootstrapHosts = append(bootstrapHosts, host)
	}
	nodes := nodeIds[:len(nodeIds)-2]
	// first node that will send cascade event to all other nodes
	var firstService *ServiceNetwork
	for i, node := range nodes {
		service, _ := initService(node.String(), bootstrapHosts)
		if i == 0 {
			firstService = service
		}
	}

	msg := &event.CallMethodMessage{
		ObjectRef: core.NewRefFromBase58("test"),
		Method:    "test",
		Arguments: []byte("test"),
	}
	// send cascade event to all nodes except the first
	c := core.Cascade{
		NodeIds:           nodeIds[1:],
		ReplicationFactor: 2,
		Entropy:           core.Entropy{0},
	}
	firstService.SendCascadeMessage(c, "test", msg)
	success := waitTimeout(&wg, 100*time.Millisecond)

	assert.True(t, success)
}
