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
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/ledger/ledgertestutil"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/stretchr/testify/assert"
)

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

	e := &message.CallMethod{
		ObjectRef: core.NewRefFromBase58("test"),
		Method:    "test",
		Arguments: []byte("test"),
	}

	network.SendMessage(core.NewRefFromBase58("test"), "test", e)
}

func TestServiceNetwork_Start(t *testing.T) {
	cfg := configuration.NewConfiguration()
	network, err := NewServiceNetwork(cfg.Host, cfg.Node)
	assert.NoError(t, err)
	err = network.Start(core.Components{})
	assert.NoError(t, err)

	err = network.Stop()
	assert.NoError(t, err)
}

func mockServiceConfiguration(host string, bootstrapHosts []string, nodeID string) (configuration.HostNetwork, configuration.NodeNetwork) {
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

type mockLedger struct {
	PM core.PulseManager
}

func (l *mockLedger) GetArtifactManager() core.ArtifactManager {
	return nil
}

func (l *mockLedger) GetJetCoordinator() core.JetCoordinator {
	return nil
}

func (l *mockLedger) GetPulseManager() core.PulseManager {
	return l.PM
}

func TestServiceNetwork_SendMessage2(t *testing.T) {
	firstNodeId := "4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj"
	secondNodeId := "53jNWvey7Nzyh4ZaLdJDf3SRgoD4GpWuwHgrgvVVGLbDkk3A7cwStSmBU2X7s4fm6cZtemEyJbce9dM9SwNxbsxf"

	firstNode, _ := NewServiceNetwork(mockServiceConfiguration(
		"127.0.0.1:10000",
		[]string{"127.0.0.1:10001"},
		firstNodeId))
	secondNode, _ := NewServiceNetwork(mockServiceConfiguration(
		"127.0.0.1:10001",
		nil,
		secondNodeId))

	secondNode.Start(core.Components{})
	firstNode.Start(core.Components{})

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

	e := &message.CallMethod{
		ObjectRef: core.NewRefFromBase58("test"),
		Method:    "test",
		Arguments: []byte("test"),
	}

	firstNode.SendMessage(core.NewRefFromBase58(secondNodeId), "test", e)
	success := waitTimeout(&wg, 20*time.Millisecond)

	assert.True(t, success)
}

func TestServiceNetwork_SendCascadeMessage(t *testing.T) {
	firstNodeId := "4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj"
	secondNodeId := "53jNWvey7Nzyh4ZaLdJDf3SRgoD4GpWuwHgrgvVVGLbDkk3A7cwStSmBU2X7s4fm6cZtemEyJbce9dM9SwNxbsxf"

	firstNode, _ := NewServiceNetwork(mockServiceConfiguration(
		"127.0.0.1:10000",
		[]string{"127.0.0.1:10001"},
		firstNodeId))
	secondNode, _ := NewServiceNetwork(mockServiceConfiguration(
		"127.0.0.1:10001",
		nil,
		secondNodeId))

	secondNode.Start(core.Components{})
	firstNode.Start(core.Components{})

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

	e := &message.CallMethod{
		ObjectRef: core.NewRefFromBase58("test"),
		Method:    "test",
		Arguments: []byte("test"),
	}

	c := core.Cascade{
		NodeIds:           []core.RecordRef{core.NewRefFromBase58(secondNodeId)},
		ReplicationFactor: 2,
		Entropy:           core.Entropy{0},
	}

	firstNode.SendCascadeMessage(c, "test", e)
	success := waitTimeout(&wg, 20*time.Millisecond)

	assert.True(t, success)

	err := firstNode.SendCascadeMessage(c, "test", nil)
	assert.NotNil(t, err)
	c.ReplicationFactor = 0
	err = firstNode.SendCascadeMessage(c, "test", e)
	assert.NotNil(t, err)
	c.ReplicationFactor = 2
	c.NodeIds = nil
	err = firstNode.SendCascadeMessage(c, "test", e)
	assert.NotNil(t, err)
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
	wg.Add(len(nodeIds) - 1)
	services := make([]*ServiceNetwork, 0)

	defer func() {
		for _, service := range services {
			service.Stop()
		}
	}()

	// init node and register test function
	initService := func(node string, bHosts []string) (service *ServiceNetwork, host string) {
		host = prefix + strconv.Itoa(port)
		service, _ = NewServiceNetwork(mockServiceConfiguration(host, bHosts, node))
		service.Start(core.Components{})
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
	// first node that will send cascade message to all other nodes
	var firstService *ServiceNetwork
	for i, node := range nodes {
		service, _ := initService(node.String(), bootstrapHosts)
		if i == 0 {
			firstService = service
		}
	}

	e := &message.CallMethod{
		ObjectRef: core.NewRefFromBase58("test"),
		Method:    "test",
		Arguments: []byte("test"),
	}
	// send cascade message to all nodes except the first
	c := core.Cascade{
		NodeIds:           nodeIds[1:],
		ReplicationFactor: 2,
		Entropy:           core.Entropy{0},
	}
	firstService.SendCascadeMessage(c, "test", e)
	success := waitTimeout(&wg, 100*time.Millisecond)

	assert.True(t, success)

	hostHandler, ctx := firstService.GetHostNetwork()
	// routing table should return total of 11 hosts
	assert.Equal(t, 11, len(hostHandler.HtFromCtx(ctx).GetHosts(100)))
	// when we request 4 hosts, routing table should return 4 hosts
	assert.Equal(t, 4, len(hostHandler.HtFromCtx(ctx).GetHosts(4)))
}

func Test_processPulse(t *testing.T) {
	firstNodeId := "4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj"
	secondNodeId := "53jNWvey7Nzyh4ZaLdJDf3SRgoD4GpWuwHgrgvVVGLbDkk3A7cwStSmBU2X7s4fm6cZtemEyJbce9dM9SwNxbsxf"

	firstNode, _ := NewServiceNetwork(mockServiceConfiguration(
		"127.0.0.1:10000",
		[]string{"127.0.0.1:10001"},
		firstNodeId))
	secondNode, _ := NewServiceNetwork(mockServiceConfiguration(
		"127.0.0.1:10001",
		nil,
		secondNodeId))
	firstLedger := &mockLedger{PM: &hostnetwork.MockPulseManager{}}
	mpm := hostnetwork.MockPulseManager{}
	var wg sync.WaitGroup
	wg.Add(1)
	mpm.SetCallback(func(pulse core.Pulse) {
		if pulse.PulseNumber == core.PulseNumber(1) {
			wg.Done()
		}
	})
	secondLedger := &mockLedger{PM: &mpm}

	firstNode.hostNetwork.GetNetworkCommonFacade().SetMessageBus(ledgertestutil.NewMessageBusMock())
	secondNode.hostNetwork.GetNetworkCommonFacade().SetMessageBus(ledgertestutil.NewMessageBusMock())

	secondNode.Start(core.Components{Ledger: secondLedger})
	firstNode.Start(core.Components{Ledger: firstLedger})

	defer func() {
		firstNode.Stop()
		secondNode.Stop()
	}()

	// pulse number is zero in MockPulseManager before receiving any pulses (default)
	firstStoredPulse, _ := firstLedger.GetPulseManager().Current()
	assert.Equal(t, core.PulseNumber(0), firstStoredPulse.PulseNumber)

	hh := firstNode.hostNetwork
	pckt := packet.NewBuilder().Type(packet.TypePulse).Request(
		&packet.RequestPulse{Pulse: core.Pulse{PulseNumber: 1, Entropy: core.Entropy{0}}}).
		Build()
	// imitate receiving pulse from the pulsar
	hostnetwork.DispatchPacketType(hh, hostnetwork.GetDefaultCtx(hh), pckt, packet.NewBuilder())

	// pulse is stored on the first node
	firstStoredPulse, _ = firstLedger.GetPulseManager().Current()
	assert.Equal(t, core.PulseNumber(1), firstStoredPulse.PulseNumber)

	// pulse is passed to the second node and stored there, too
	success := waitTimeout(&wg, time.Millisecond*10)
	assert.True(t, success)
	secondStoredPulse, _ := secondLedger.GetPulseManager().Current()
	assert.Equal(t, core.PulseNumber(1), secondStoredPulse.PulseNumber)
}

func Test_processPulse2(t *testing.T) {
	nodeIds := []core.RecordRef{
		core.NewRefFromBase58("4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj"),
		core.NewRefFromBase58("53jNWvey7Nzyh4ZaLdJDf3SRgoD4GpWuwHgrgvVVGLbDkk3A7cwStSmBU2X7s4fm6cZtemEyJbce9dM9SwNxbsxf"),
		core.NewRefFromBase58("9uE5MEWQB2yfKY8kTgTNovWii88D4anmf7GAiovgcxx6Uc6EBmZ212mpyMa1L22u9TcUUd94i8LvULdsdBoG8ed"),
		core.NewRefFromBase58("4qXdYkfL9U4tL3qRPthdbdajtafR4KArcXjpyQSEgEMtpuin3t8aZYmMzKGRnXHBauytaPQ6bfwZyKZzRPpR6gyX"),
	}

	prefix := "127.0.0.1:"
	port := 11000
	bootstrapNodes := nodeIds[len(nodeIds)-2:]
	bootstrapHosts := make([]string, 0)
	services := make([]*ServiceNetwork, 0)
	ledgers := make([]core.Ledger, 0)

	var wg sync.WaitGroup
	wg.Add(len(nodeIds))

	defer func() {
		for _, service := range services {
			service.Stop()
		}
	}()

	// init node and register test function
	initService := func(node string, bHosts []string) (host string) {
		mpm := hostnetwork.MockPulseManager{}
		mpm.SetCallback(func(pulse core.Pulse) {
			if pulse.PulseNumber == core.PulseNumber(1) {
				wg.Done()
			}
		})
		ledger := &mockLedger{PM: &mpm}
		ledgers = append(ledgers, ledger)

		host = prefix + strconv.Itoa(port)
		service, _ := NewServiceNetwork(mockServiceConfiguration(host, bHosts, node))
		service.hostNetwork.GetNetworkCommonFacade().SetMessageBus(ledgertestutil.NewMessageBusMock())
		service.Start(core.Components{Ledger: ledger})
		port++
		services = append(services, service)
		return
	}

	for _, node := range bootstrapNodes {
		host := initService(node.String(), nil)
		bootstrapHosts = append(bootstrapHosts, host)
	}
	nodes := nodeIds[:len(nodeIds)-2]
	for _, node := range nodes {
		initService(node.String(), bootstrapHosts)
	}

	lastIndex := len(services) - 1

	// pulse number is zero in MockPulseManager before receiving any pulses (default)
	ll := ledgers[lastIndex]
	firstStoredPulse, _ := ll.GetPulseManager().Current()
	assert.Equal(t, core.PulseNumber(0), firstStoredPulse.PulseNumber)

	// time.Sleep(time.Millisecond * 100)

	hh := services[lastIndex].hostNetwork
	pckt := packet.NewBuilder().Type(packet.TypePulse).Request(
		&packet.RequestPulse{Pulse: core.Pulse{PulseNumber: 1, Entropy: core.Entropy{0}}}).
		Build()
	// imitate receiving pulse from the pulsar on the last started service
	hostnetwork.DispatchPacketType(hh, hostnetwork.GetDefaultCtx(hh), pckt, packet.NewBuilder())

	// pulse is stored on the first node
	firstStoredPulse, _ = ll.GetPulseManager().Current()
	assert.Equal(t, core.PulseNumber(1), firstStoredPulse.PulseNumber)

	// pulse is passed to the other 11 nodes and stored there, too
	success := waitTimeout(&wg, time.Millisecond*200)
	assert.True(t, success)

	for _, ldgr := range ledgers {
		pulse, _ := ldgr.GetPulseManager().Current()
		assert.Equal(t, core.PulseNumber(1), pulse.PulseNumber)
	}
}
