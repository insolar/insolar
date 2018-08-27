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
	"log"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/relay"
	"github.com/insolar/insolar/network/hostnetwork/rpc"
	"github.com/insolar/insolar/network/hostnetwork/store"
	"github.com/insolar/insolar/network/hostnetwork/transport"
)

/*
//todo: interface for HostNetwork
type HostNetwork interface {
	RPC
}
*/

// NewHostNetwork creates and returns DHT network.
func NewHostNetwork(cfg configuration.HostNetwork) (*DHT, error) {

	proxy := relay.NewProxy()

	tp, err := transport.NewTransport(cfg.Transport, proxy)
	if err != nil {
		return nil, err
	}

	originAddress, err := host.NewAddress(tp.PublicAddress())
	if err != nil {
		return nil, err
	}

	origin, err := host.NewOrigin(nil, originAddress)
	if err != nil {
		return nil, err
	}

	options := &Options{BootstrapHosts: getBootstrapHosts(cfg.BootstrapHosts)}

	network, err := NewDHT(
		store.NewMemoryStoreFactory().Create(),
		origin,
		tp,
		rpc.NewRPCFactory(nil).Create(),
		options,
		proxy)
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


