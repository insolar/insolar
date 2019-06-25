//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package replica

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/heavy/replica.Transport -o ./ -s _mock.go

// TODO: write docs
type Transport interface {
	Send(receiver, method string, data []byte) ([]byte, error)
	Register(method string, handle Handle)
	Me() string
}

type Handle func(data []byte) ([]byte, error)

func NewInternalTransport(net network.HostNetwork, handlers map[string]Handle) Transport {
	return &internalTransport{net: net, handlers: handlers}
}

type internalTransport struct {
	net         network.HostNetwork
	NodeNetwork insolar.NodeNetwork
	handlers    map[string]Handle
}

func (t *internalTransport) Send(receiver, method string, data []byte) ([]byte, error) {
	receiverHost, err := t.hostByAddress(receiver)
	if err != nil || receiverHost == nil {
		return []byte{}, errors.Wrapf(err, "failed to create host by receiver address")
	}
	req := &packet.RPCRequest{
		Method: method,
		Data:   data,
	}
	future, err := t.net.SendRequestToHost(context.Background(), types.Replication, req, receiverHost)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to send request to host")
	}
	packet, err := future.WaitResponse(10 * time.Second)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get result from host")
	}
	if packet.GetResponse() == nil || packet.GetResponse().GetRPC() == nil {
		inslogger.FromContext(context.Background()).Warnf("Error getting RPC response from node %s: "+
			"got invalid response protobuf message: %s", receiver, packet)
	}
	resp := packet.GetResponse().GetRPC()
	if resp.Result == nil {
		return nil, errors.New("RPC call returned error: " + resp.Error)
	}
	return resp.Result, nil
}

func (t *internalTransport) Register(method string, handle Handle) {
	t.handlers[method] = handle
}

func (t *internalTransport) Me() string {
	return t.net.PublicAddress()
}

func (t *internalTransport) hostByAddress(receiver string) (*host.Host, error) {
	for _, node := range t.NodeNetwork.GetWorkingNodes() {
		if node.Address() == receiver {
			host, err := host.NewHostN(receiver, node.ID())
			if err != nil {
				return nil, errors.Wrapf(err, "failed to create new host by address")
			}
			return host, nil
		}
	}
}
