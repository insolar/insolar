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

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
)

type BootstrapController struct {
	options   *Options
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
	bootstrapCount := len(bc.options.BootstrapHosts)
	bootstrapHosts := make(chan *host.Host, bootstrapCount)
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

	counter := 0
	result := make([]*host.Host, 0)
Loop:
	for {
		select {
		case bootstrapHost := <-bootstrapHosts:
			result = append(result, bootstrapHost)
			counter++
			if counter == bootstrapCount {
				break Loop
			}
		case <-time.After(bc.options.BootstrapTimeout):
			log.Warnf("Bootstrap timeout, successful bootstraps: %d/%d", counter, bootstrapCount)
			break Loop
		}
	}

	if counter == 0 {
		return errors.New("Failed to bootstrap to any of discovery nodes")
	}
	bc.bootstrapHosts = result
	return nil
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

func (bc *BootstrapController) Start(components core.Components) {
	bc.transport.RegisterPacketHandler(types.Bootstrap, bc.processBootstrap)
}

func NewBootstrapController(options *Options, transport hostnetwork.InternalTransport) *BootstrapController {
	return &BootstrapController{options: options, transport: transport, pinger: NewPinger(transport)}
}
