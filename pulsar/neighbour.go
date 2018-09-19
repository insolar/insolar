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

package pulsar

import (
	"crypto/ecdsa"
	"net"
	"net/rpc"
	"sync"

	"github.com/insolar/insolar/log"

	"github.com/insolar/insolar/configuration"
)

type RpcConnection struct {
	sync.Mutex
	*rpc.Client
}

type Neighbour struct {
	ConnectionType    configuration.ConnectionType
	ConnectionAddress string
	OutgoingClient    *RpcConnection
	PublicKey         *ecdsa.PublicKey
}

func (neighbour *Neighbour) CheckAndRefreshConnection(rpcErr error) {
	var err error
	var conn net.Conn
	if rpcErr == rpc.ErrShutdown {
		log.Info("Restarting RPC Connection due to error")
		neighbour.OutgoingClient.Lock()
		conn, err = net.Dial(neighbour.ConnectionType.String(), neighbour.ConnectionAddress)
		neighbour.OutgoingClient.Client = rpc.NewClient(conn)
		neighbour.OutgoingClient.Unlock()
	}

	if err != nil {
		log.Error("Unable to initialize connection to RPC")
	}
}
