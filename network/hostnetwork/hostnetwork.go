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

package hostnetwork

import (
	"log"
	"strings"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/cascade"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/hosthandler"
	"github.com/insolar/insolar/network/hostnetwork/id"
	"github.com/insolar/insolar/network/hostnetwork/relay"
	"github.com/insolar/insolar/network/hostnetwork/rpc"
	"github.com/insolar/insolar/network/hostnetwork/store"
	"github.com/insolar/insolar/network/hostnetwork/transport"
)

type NodeNetwork interface {
	ResolveHostID(ref core.RecordRef) string
	GetID() core.RecordRef
}

// NewHostNetwork creates and returns DHT network.
func NewHostNetwork(cfg configuration.HostNetwork, nn NodeNetwork, cascade *cascade.Cascade) (*DHT, error) {

	if strings.Contains(cfg.Transport.Address, "0.0.0.0") && !cfg.Transport.BehindNAT {
		log.Fatal("hostnetwork.NewHostNetwork: \n Couldn't start at 0.0.0.0")
	}

	proxy := relay.NewProxy()

	tp, err := transport.NewTransport(cfg.Transport, proxy)
	if err != nil {
		return nil, err
	}

	originAddress, err := host.NewAddress(tp.PublicAddress())
	if err != nil {
		return nil, err
	}

	originID := id.ID(nn.ResolveHostID(nn.GetID()))
	origin, err := host.NewOrigin([]id.ID{originID}, originAddress)
	if err != nil {
		return nil, err
	}

	options := &Options{BootstrapHosts: getBootstrapHosts(cfg.BootstrapHosts)}
	ncf := hosthandler.NewFacade(rpc.NewRPCFactory(nil).Create(), cascade)

	network, err := NewDHT(
		store.NewMemoryStoreFactory().Create(),
		origin,
		tp,
		ncf,
		options,
		proxy,
	)
	if err != nil {
		return nil, err
	}

	return network, nil
}

func getBootstrapHosts(addresses []string) []*host.Host {
	var hosts []*host.Host
	for _, a := range addresses {
		address, err := host.NewAddress(a)
		if err != nil {
			log.Fatalln("Failed to create bootstrap address:", err.Error())
		}
		hosts = append(hosts, host.NewHost(address))
	}
	return hosts
}
