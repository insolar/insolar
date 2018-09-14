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
	"github.com/insolar/insolar/network/hostnetwork/hosthandler"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/relay"
	"github.com/insolar/insolar/network/hostnetwork/routing"
	"github.com/insolar/insolar/network/hostnetwork/store"
	"github.com/pkg/errors"
)

// DispatchPacketType checks message type.
func DispatchPacketType(hostHandler hosthandler.HostHandler, ctx hosthandler.Context, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	switch msg.Type {
	case packet.TypeFindHost:
		return processFindHost(hostHandler, ctx, msg, packetBuilder)
	case packet.TypeFindValue:
		return processFindValue(hostHandler, ctx, msg, packetBuilder)
	case packet.TypeStore:
		return processStore(hostHandler, ctx, msg)
	case packet.TypePing:
		return processPing(msg, packetBuilder)
	case packet.TypeRPC:
		return processRPC(hostHandler, ctx, msg, packetBuilder)
	case packet.TypeRelay:
		return processRelay(hostHandler, ctx, msg, packetBuilder)
	case packet.TypeCheckOrigin:
		return processCheckOriginRequest(hostHandler, msg, packetBuilder)
	case packet.TypeAuth:
		return processAuthentication(hostHandler, msg, packetBuilder)
	case packet.TypeObtainIP:
		return processObtainIPRequest(msg, packetBuilder)
	case packet.TypeRelayOwnership:
		return processRelayOwnership(hostHandler, msg, packetBuilder)
	case packet.TypeKnownOuterHosts:
		return processKnownOuterHosts(hostHandler, msg, packetBuilder)
	case packet.TypeCheckNodePriv:
		return processCheckNodePriv(hostHandler, msg, packetBuilder)
	case packet.TypeCascadeSend:
		return processCascadeSend(hostHandler, ctx, msg, packetBuilder)
	default:
		return nil, errors.New("unknown request type")
	}
}

func processKnownOuterHosts(hostHandler hosthandler.HostHandler, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	data := msg.Data.(*packet.RequestKnownOuterHosts)

	ID := hostHandler.GetHighKnownHostID()
	hosts := hostHandler.GetOuterHostsCount()
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

func processCheckNodePriv(hostHandler hosthandler.HostHandler, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	data := msg.Data.(*packet.RequestCheckNodePriv)
	var response packet.ResponseCheckNodePriv

	if hostHandler.ConfirmNodeRole(data.RoleKey) {
		response.State = packet.Confirmed
	} else {
		response.State = packet.Declined
	}

	return packetBuilder.Response(response).Build(), nil
}

func processRelayOwnership(hostHandler hosthandler.HostHandler, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	data := msg.Data.(*packet.RequestRelayOwnership)

	if data.Ready {
		hostHandler.AddPossibleProxyID(msg.Sender.ID.String())
	} else {
		hostHandler.RemovePossibleProxyID(msg.Sender.ID.String())
		err := AuthenticationRequest(hostHandler, "begin", msg.Sender.ID.String())
		if err != nil {
			return nil, err
		}
		err = RelayRequest(hostHandler, "start", msg.Sender.ID.String())
		if err != nil {
			return nil, err
		}
	}
	response := &packet.ResponseRelayOwnership{Accepted: true}
	return packetBuilder.Response(response).Build(), nil
}

func processFindHost(hostHandler hosthandler.HostHandler, ctx hosthandler.Context, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	ht := hostHandler.HtFromCtx(ctx)
	data := msg.Data.(*packet.RequestDataFindHost)
	hostHandler.AddHost(ctx, routing.NewRouteHost(msg.Sender))
	closest := ht.GetClosestContacts(routing.MaxContactsInBucket, data.Target, []*host.Host{msg.Sender})
	response := &packet.ResponseDataFindHost{
		Closest: closest.Hosts(),
	}
	return packetBuilder.Response(response).Build(), nil
}

func processFindValue(hostHandler hosthandler.HostHandler, ctx hosthandler.Context, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	ht := hostHandler.HtFromCtx(ctx)
	data := msg.Data.(*packet.RequestDataFindValue)
	hostHandler.AddHost(ctx, routing.NewRouteHost(msg.Sender))
	value, exists := hostHandler.StoreRetrieve(data.Target)
	response := &packet.ResponseDataFindValue{}
	if exists {
		response.Value = value
	} else {
		closest := ht.GetClosestContacts(routing.MaxContactsInBucket, data.Target, []*host.Host{msg.Sender})
		response.Closest = closest.Hosts()
	}
	return packetBuilder.Response(response).Build(), nil
}

func processStore(hostHandler hosthandler.HostHandler, ctx hosthandler.Context, msg *packet.Packet) (*packet.Packet, error) {
	data := msg.Data.(*packet.RequestDataStore)
	hostHandler.AddHost(ctx, routing.NewRouteHost(msg.Sender))
	key := store.NewKey(data.Data)
	expiration := hostHandler.GetExpirationTime(ctx, key)
	replication := time.Now().Add(hostHandler.GetReplicationTime())
	err := hostHandler.Store(key, data.Data, replication, expiration, false)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func processPing(msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	log.Println("recv ping message from " + msg.Sender.Address.String())
	return packetBuilder.Response(nil).Build(), nil
}

func processRPC(hostHandler hosthandler.HostHandler, ctx hosthandler.Context, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	data := msg.Data.(*packet.RequestDataRPC)
	hostHandler.AddHost(ctx, routing.NewRouteHost(msg.Sender))
	result, err := hostHandler.InvokeRPC(msg.Sender, data.Method, data.Args)
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
func processRelay(hostHandler hosthandler.HostHandler, ctx hosthandler.Context, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	var err error
	var state relay.State
	var packet1 *packet.Packet
	if !hostHandler.HostIsAuthenticated(msg.Sender.ID.String()) {
		log.Print("relay request from unknown host rejected")
		response := &packet.ResponseRelay{
			State: relay.NoAuth,
		}

		packet1, err = packetBuilder.Response(response).Build(), nil
	} else {
		data := msg.Data.(*packet.RequestRelay)
		hostHandler.AddHost(ctx, routing.NewRouteHost(msg.Sender))

		switch data.Command {
		case packet.StartRelay:
			err = hostHandler.AddRelayClient(msg.Sender)
			state = relay.Started
		case packet.StopRelay:
			err = hostHandler.RemoveRelayClient(msg.Sender)
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
		packet1, err = packetBuilder.Response(response).Build(), nil
	}
	return packet1, err
}

func processAuthentication(hostHandler hosthandler.HostHandler, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	data := msg.Data.(*packet.RequestAuth)
	switch data.Command {
	case packet.BeginAuth:
		if hostHandler.HostIsAuthenticated(msg.Sender.ID.String()) {
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
		hostHandler.AddAuthSentKey(msg.Sender.ID.String(), key)
		response := &packet.ResponseAuth{
			Success:       true,
			AuthUniqueKey: key,
		}

		// TODO process verification msg.Sender host
		// confirmed
		err = CheckOriginRequest(hostHandler, msg.Sender.ID.String())
		if err != nil {
			return nil, err
		}

		return packetBuilder.Response(response).Build(), nil
	case packet.RevokeAuth:
		hostHandler.RemoveAuthHost(msg.Sender.ID.String())
		response := &packet.ResponseAuth{
			Success:       true,
			AuthUniqueKey: nil,
		}

		return packetBuilder.Response(response).Build(), nil
	}
	return nil, errors.New("unknown auth command")
}

func processCheckOriginRequest(hostHandler hosthandler.HostHandler, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	if key, ok := hostHandler.KeyIsReceived(msg.Sender.ID.String()); ok {
		response := &packet.ResponseCheckOrigin{AuthUniqueKey: key}
		return packetBuilder.Response(response).Build(), nil
	}
	return nil, errors.New("CheckOrigin request from unregistered host")
}

func processObtainIPRequest(msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	response := &packet.ResponseObtainIP{IP: msg.RemoteAddress}
	return packetBuilder.Response(response).Build(), nil
}

func processCascadeSend(hostHandler hosthandler.HostHandler, ctx hosthandler.Context, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	data := msg.Data.(*packet.RequestCascadeSend)

	hostHandler.AddHost(ctx, routing.NewRouteHost(msg.Sender))
	_, err := hostHandler.InvokeRPC(msg.Sender, data.RPC.Method, data.RPC.Args)
	response := &packet.ResponseCascadeSend{
		Success: true,
		Error:   "",
	}
	if err != nil {
		response.Success = false
		response.Error = err.Error()
	}
	hostHandler.GetNetworkCommonFacade().GetCascade().SendToNextLayer(data.Data, data.RPC.Method, data.RPC.Args)

	return packetBuilder.Response(response).Build(), err
}
