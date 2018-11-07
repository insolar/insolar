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

package controller

import (
	"encoding/gob"
	"time"

	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
)

type BootstrapController struct {
	options   *common.Options
	transport hostnetwork.InternalTransport
	pinger    *Pinger

	bootstrapHosts []*host.Host
}

type BootstrapRequest struct {
	// pass node certificate, i guess
}

type BootstrapResponse struct {
	Code         BootstrapCode
	RedirectHost string
	RejectReason string
}

type BootstrapCode uint8

const (
	BootstrapAccepted = BootstrapCode(iota + 1)
	BootstrapRejected
	BootstrapRedirected
)

func init() {
	gob.Register(&BootstrapRequest{})
	gob.Register(&BootstrapResponse{})
}

func (bc *BootstrapController) GetBootstrapHosts() []*host.Host {
	return bc.bootstrapHosts
}

func (bc *BootstrapController) Bootstrap() error {
	if len(bc.options.BootstrapHosts) == 0 {
		bc.bootstrapHosts = make([]*host.Host, 0)
		return nil
	}
	ch := bc.getBootstrapHostsChannel()
	hosts := bc.waitResultsFromChannel(ch)

	if len(hosts) == 0 {
		return errors.New("Failed to bootstrap to any of discovery nodes")
	}
	bc.bootstrapHosts = hosts
	return nil
}

func (bc *BootstrapController) getBootstrapHostsChannel() <-chan *host.Host {
	bootstrapHosts := make(chan *host.Host, len(bc.options.BootstrapHosts))
	for _, bootstrapAddress := range bc.options.BootstrapHosts {
		go func(address string, ch chan<- *host.Host) {
			log.Infof("Starting bootstrap to address %s", address)
			bootstrapHost, err := bc.bootstrap(address)
			if err != nil {
				log.Errorf("Error bootstrapping to address %s: %s", address, err.Error())
				return
			}
			bootstrapHosts <- bootstrapHost
		}(bootstrapAddress, bootstrapHosts)
	}
	return bootstrapHosts
}

func (bc *BootstrapController) waitResultsFromChannel(ch <-chan *host.Host) []*host.Host {
	result := make([]*host.Host, 0)
	for {
		select {
		case bootstrapHost := <-ch:
			result = append(result, bootstrapHost)
			if len(result) == len(bc.options.BootstrapHosts) {
				return result
			}
		case <-time.After(bc.options.BootstrapTimeout):
			log.Warnf("Bootstrap timeout, successful bootstraps: %d/%d", len(result), len(bc.options.BootstrapHosts))
			return result
		}
	}
	return result
}

func (bc *BootstrapController) bootstrap(address string) (*host.Host, error) {
	// TODO: add infinite bootstrap
	bootstrapHost, err := bc.pinger.Ping(address, bc.options.PingTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to ping address %s", address)
	}
	request := bc.transport.NewRequestBuilder().Type(types.Bootstrap).Data(&BootstrapRequest{}).Build()
	future, err := bc.transport.SendRequestPacket(request, bootstrapHost)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to ping address %s", address)
	}
	response, err := future.GetResponse(bc.options.BootstrapTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get response to bootstrap request")
	}
	data := response.GetData().(*BootstrapResponse)
	if data.Code == BootstrapRejected {
		return nil, errors.New("Rejected: " + data.RejectReason)
	}
	if data.Code == BootstrapRedirected {
		return bc.bootstrap(data.RedirectHost)
	}
	return response.GetSenderHost(), nil
}

func (bc *BootstrapController) processBootstrap(request network.Request) (network.Response, error) {
	// TODO: check certificate and redirect logic
	return bc.transport.BuildResponse(request, &BootstrapResponse{Code: BootstrapAccepted}), nil
}

func (bc *BootstrapController) Start() {
	bc.transport.RegisterPacketHandler(types.Bootstrap, bc.processBootstrap)
}

func NewBootstrapController(options *common.Options, transport hostnetwork.InternalTransport) *BootstrapController {
	return &BootstrapController{options: options, transport: transport, pinger: NewPinger(transport)}
}
