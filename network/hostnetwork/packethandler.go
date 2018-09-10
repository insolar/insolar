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
	"crypto/rand"
	"log"
	"time"

	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/relay"
	"github.com/insolar/insolar/network/hostnetwork/routing"
	"github.com/insolar/insolar/network/hostnetwork/store"
	"github.com/pkg/errors"
)

// DispatchPacketType checks message type.
func DispatchPacketType(dht *DHT, ctx Context, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	switch msg.Type {
	case packet.TypeFindHost:
		return processFindHost(dht, ctx, msg, packetBuilder)
	case packet.TypeFindValue:
		return processFindValue(dht, ctx, msg, packetBuilder)
	case packet.TypeStore:
		return processStore(dht, ctx, msg)
	case packet.TypePing:
		return processPing(msg, packetBuilder)
	case packet.TypeRPC:
		return processRPC(dht, ctx, msg, packetBuilder)
	case packet.TypeRelay:
		return processRelay(dht, ctx, msg, packetBuilder)
	case packet.TypeCheckOrigin:
		return processCheckOriginRequest(dht, msg, packetBuilder)
	case packet.TypeAuth:
		return processAuthentication(dht, ctx, msg, packetBuilder)
	case packet.TypeObtainIP:
		return processObtainIPRequest(msg, packetBuilder)
	case packet.TypeRelayOwnership:
		return processRelayOwnership(dht, ctx, msg, packetBuilder)
	case packet.TypeKnownOuterHosts:
		return processKnownOuterHosts(dht, msg, packetBuilder)
	case packet.TypeCheckNodePriv:
		return processCheckNodePriv(dht, msg, packetBuilder)
	default:
		return nil, errors.New("unknown request type")
	}
}

func processKnownOuterHosts(dht *DHT, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	data := msg.Data.(*packet.RequestKnownOuterHosts)

	ID := dht.Subnet.HighKnownHosts.ID
	hosts := dht.Subnet.HighKnownHosts.OuterHosts
	if data.OuterHosts > hosts {
		ID = data.ID
		hosts = data.OuterHosts
	}
	response := &packet.ResponseKnownOuterHosts{
		ID:         ID,
		OuterHosts: hosts,
	}

	return packetBuilder.Response(response).Build(), nil
}

func processCheckNodePriv(dht *DHT, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	data := msg.Data.(*packet.RequestCheckNodePriv)
	var response packet.ResponseCheckNodePriv

	if dht.ConfirmNodeRole(data.RoleKey) {
		response.State = packet.Confirmed
	} else {
		response.State = packet.Declined
	}

	return packetBuilder.Response(response).Build(), nil
}

func processRelayOwnership(dht *DHT, ctx Context, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	data := msg.Data.(*packet.RequestRelayOwnership)

	if data.Ready {
		dht.Subnet.PossibleProxyIDs = append(dht.Subnet.PossibleProxyIDs, msg.Sender.ID.KeyString())
	} else {
		for i, j := range dht.Subnet.PossibleProxyIDs {
			if j == msg.Sender.ID.KeyString() {
				dht.Subnet.PossibleProxyIDs = append(dht.Subnet.PossibleProxyIDs[:i], dht.Subnet.PossibleProxyIDs[i+1:]...)
				err := dht.AuthenticationRequest(ctx, "begin", msg.Sender.ID.KeyString())
				if err != nil {
					return nil, err
				}
				err = dht.RelayRequest(ctx, "start", msg.Sender.ID.KeyString())
				if err != nil {
					return nil, err
				}
				break
			}
		}
	}
	response := &packet.ResponseRelayOwnership{Accepted: true}
	return packetBuilder.Response(response).Build(), nil
}

func processFindHost(dht *DHT, ctx Context, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	ht := dht.HtFromCtx(ctx)
	data := msg.Data.(*packet.RequestDataFindHost)
	dht.AddHost(ctx, routing.NewRouteHost(msg.Sender))
	closest := ht.GetClosestContacts(routing.MaxContactsInBucket, data.Target, []*host.Host{msg.Sender})
	response := &packet.ResponseDataFindHost{
		Closest: closest.Hosts(),
	}
	return packetBuilder.Response(response).Build(), nil
}

func processFindValue(dht *DHT, ctx Context, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	ht := dht.HtFromCtx(ctx)
	data := msg.Data.(*packet.RequestDataFindValue)
	dht.AddHost(ctx, routing.NewRouteHost(msg.Sender))
	value, exists := dht.Store.Retrieve(data.Target)
	response := &packet.ResponseDataFindValue{}
	if exists {
		response.Value = value
	} else {
		closest := ht.GetClosestContacts(routing.MaxContactsInBucket, data.Target, []*host.Host{msg.Sender})
		response.Closest = closest.Hosts()
	}
	return packetBuilder.Response(response).Build(), nil
}

func processStore(dht *DHT, ctx Context, msg *packet.Packet) (*packet.Packet, error) {
	data := msg.Data.(*packet.RequestDataStore)
	dht.AddHost(ctx, routing.NewRouteHost(msg.Sender))
	key := store.NewKey(data.Data)
	expiration := dht.GetExpirationTime(ctx, key)
	replication := time.Now().Add(dht.GetReplicationTime())
	err := dht.Store.Store(key, data.Data, replication, expiration, false)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func processPing(msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	log.Println("recv ping message from " + msg.Sender.Address.String())
	return packetBuilder.Response(nil).Build(), nil
}

func processRPC(dht *DHT, ctx Context, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	data := msg.Data.(*packet.RequestDataRPC)
	dht.AddHost(ctx, routing.NewRouteHost(msg.Sender))
	result, err := dht.InvokeRPC(msg.Sender, data.Method, data.Args)
	response := &packet.ResponseDataRPC{
		Success: true,
		Result:  result,
		Error:   "",
	}
	if err != nil {
		response.Success = false
		response.Error = err.Error()
	}
	return packetBuilder.Response(response).Build(), nil
}

// Precess relay request.
func processRelay(dht *DHT, ctx Context, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	var err error
	if !dht.Auth.AuthenticatedHosts[msg.Sender.ID.KeyString()] {
		log.Print("relay request from unknown host rejected")
		response := &packet.ResponseRelay{
			State: relay.NoAuth,
		}

		return packetBuilder.Response(response).Build(), nil
	} else {
		data := msg.Data.(*packet.RequestRelay)
		dht.AddHost(ctx, routing.NewRouteHost(msg.Sender))

		var state relay.State

		switch data.Command {
		case packet.StartRelay:
			err = dht.Relay.AddClient(msg.Sender)
			state = relay.Started
		case packet.StopRelay:
			err = dht.Relay.RemoveClient(msg.Sender)
			state = relay.Stopped
		default:
			state = relay.Unknown
		}

		if err != nil {
			state = relay.Error
		}

		response := &packet.ResponseRelay{
			State: state,
		}

		return packetBuilder.Response(response).Build(), nil
	}
}

func processAuthentication(dht *DHT, ctx Context, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	data := msg.Data.(*packet.RequestAuth)
	switch data.Command {
	case packet.BeginAuth:
		if dht.Auth.AuthenticatedHosts[msg.Sender.ID.KeyString()] {
			// TODO: whats next?
			response := &packet.ResponseAuth{
				Success:       false,
				AuthUniqueKey: nil,
			}

			return packetBuilder.Response(response).Build(), nil
		}
		key := make([]byte, 512)
		_, err := rand.Read(key) // crypto/rand
		if err != nil {
			return nil, err
		}
		dht.Auth.SentKeys[msg.Sender.ID.KeyString()] = key
		response := &packet.ResponseAuth{
			Success:       true,
			AuthUniqueKey: key,
		}

		// TODO process verification msg.Sender host
		// confirmed
		err = dht.CheckOriginRequest(ctx, msg.Sender.ID.KeyString())
		if err != nil {
			return nil, err
		}

		return packetBuilder.Response(response).Build(), nil
	case packet.RevokeAuth:
		delete(dht.Auth.AuthenticatedHosts, msg.Sender.ID.KeyString())
		response := &packet.ResponseAuth{
			Success:       true,
			AuthUniqueKey: nil,
		}

		return packetBuilder.Response(response).Build(), nil
	}
	return nil, errors.New("unknown auth command")
}

func processCheckOriginRequest(dht *DHT, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	dht.Auth.Mut.Lock()
	defer dht.Auth.Mut.Unlock()
	if key, ok := dht.Auth.ReceivedKeys[msg.Sender.ID.KeyString()]; ok {
		response := &packet.ResponseCheckOrigin{AuthUniqueKey: key}
		return packetBuilder.Response(response).Build(), nil
	} else {
		return nil, errors.New("CheckOrigin request from unregistered host")
	}
}

func processObtainIPRequest(msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	response := &packet.ResponseObtainIP{IP: msg.RemoteAddress}
	return packetBuilder.Response(response).Build(), nil
}
