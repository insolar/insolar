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
	"crypto/rand"
	"time"

	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/hosthandler"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/relay"
	"github.com/insolar/insolar/network/hostnetwork/routing"
	"github.com/insolar/insolar/network/hostnetwork/store"
	"github.com/pkg/errors"
)

// DispatchPacketType checks message type.
func DispatchPacketType(
	hostHandler hosthandler.HostHandler,
	ctx hosthandler.Context,
	msg *packet.Packet,
	packetBuilder packet.Builder,
) (*packet.Packet, error) { // nolint: gocyclo
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
	case packet.TypeAuthentication:
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
	case packet.TypePulse:
		return processPulse(hostHandler, ctx, msg, packetBuilder)
	case packet.TypeGetRandomHosts:
		return processGetRandomHosts(hostHandler, ctx, msg, packetBuilder)
	case packet.TypeCheckSignedNonce:
		return processCheckSignedNonce(hostHandler, ctx, msg, packetBuilder)
	case packet.TypeCheckPublicKey:
		return processCheckPublicKey(hostHandler, ctx, msg, packetBuilder)
	case packet.TypeActiveNodes:
		return processActiveNodes(hostHandler, packetBuilder)
	case packet.TypeDisconnect:
		return processDisconnect(hostHandler, packetBuilder)
	default:
		return nil, errors.New("unknown request type")
	}
}

func processDisconnect(hostHandler hosthandler.HostHandler, packetBuilder packet.Builder) (*packet.Packet, error) {
	// TODO: disconnect from active list
	return packetBuilder.Response(&packet.ResponseDisconnect{Disconnected: true, Error: nil}).Build(), nil
}

func processCheckPublicKey(
	hostHandler hosthandler.HostHandler,
	ctx hosthandler.Context,
	msg *packet.Packet,
	packetBuilder packet.Builder) (*packet.Packet, error) {
	// TODO: do real check key.
	exist := true
	nonce, err := time.Now().MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal nonce")
	}
	return packetBuilder.Response(&packet.ResponseCheckPublicKey{Nonce: nonce, Exist: exist}).Build(), nil
}

func processCheckSignedNonce(
	hostHandler hosthandler.HostHandler,
	ctx hosthandler.Context,
	msg *packet.Packet,
	packetBuilder packet.Builder) (*packet.Packet, error) {
	// TODO: do real check sign.
	parsed := true
	return packetBuilder.Response(&packet.ResponseCheckSignedNonce{Success: parsed}).Build(), nil
}

func processGetRandomHosts(
	hostHandler hosthandler.HostHandler,
	ctx hosthandler.Context,
	msg *packet.Packet,
	packetBuilder packet.Builder) (*packet.Packet, error) {

	data := msg.Data.(*packet.RequestGetRandomHosts)
	ht := hostHandler.HtFromCtx(ctx)
	if data.HostsNumber <= 0 {
		return packetBuilder.Response(&packet.ResponseGetRandomHosts{
			Hosts: nil, Error: "hosts number should be more than zero"}).Build(), nil
	}
	hosts := ht.GetHosts(data.HostsNumber)
	// TODO: handle scenario when we get less hosts than requested
	return packetBuilder.Response(&packet.ResponseGetRandomHosts{Hosts: hosts, Error: ""}).Build(), nil
}

func processPulse(hostHandler hosthandler.HostHandler, ctx hosthandler.Context, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	data := msg.Data.(*packet.RequestPulse)
	pm := hostHandler.GetNetworkCommonFacade().GetPulseManager()
	if pm == nil {
		return nil, errors.New("PulseManager is not initialized")
	}
	currentPulse, err := pm.Current()
	if err != nil {
		return nil, errors.Wrap(err, "could not get current pulse")
	}
	log.Debugf("got new pulse number: %d", currentPulse.PulseNumber)
	if (data.Pulse.PulseNumber > currentPulse.PulseNumber) &&
		(data.Pulse.PulseNumber >= currentPulse.NextPulseNumber) {
		err = pm.Set(data.Pulse)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to set pulse")
		}
		log.Debugf("set new current pulse number: %d", currentPulse.PulseNumber)
		ht := hostHandler.HtFromCtx(ctx)
		hosts := ht.GetMulticastHosts()
		go ResendPulseToKnownHosts(hostHandler, hosts, data)
	}
	return packetBuilder.Response(&packet.ResponsePulse{Success: true, Error: ""}).Build(), nil
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
			return nil, errors.Wrap(err, "AuthenticationRequest failed")
		}
		err = RelayRequest(hostHandler, "start", msg.Sender.ID.String())
		if err != nil {
			return nil, errors.Wrap(err, "RelayRequest failed")
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
		return nil, errors.Wrap(err, "Failed to store data")
	}
	return nil, nil
}

func processPing(msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	log.Debugln("recv ping message from " + msg.Sender.Address.String())
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
		log.Debug("relay request from unknown host rejected")
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
	data := msg.Data.(*packet.RequestAuthentication)
	switch data.Command {
	case packet.BeginAuthentication:
		if hostHandler.HostIsAuthenticated(msg.Sender.ID.String()) {
			// TODO: whats next?
			response := &packet.ResponseAuthentication{
				Success:       false,
				AuthUniqueKey: nil,
			}

			return packetBuilder.Response(response).Build(), nil
		}
		key := make([]byte, 512)
		_, err := rand.Read(key) // crypto/rand
		if err != nil {
			return nil, errors.Wrap(err, "Failed to generate random key")
		}
		hostHandler.AddAuthSentKey(msg.Sender.ID.String(), key)
		response := &packet.ResponseAuthentication{
			Success:       true,
			AuthUniqueKey: key,
		}

		// TODO process verification msg.Sender host
		// confirmed
		err = CheckOriginRequest(hostHandler, msg.Sender.ID.String())
		if err != nil {
			return nil, errors.Wrap(err, "CheckOriginRequest failed")
		}

		return packetBuilder.Response(response).Build(), nil
	case packet.RevokeAuthentication:
		hostHandler.RemoveAuthHost(msg.Sender.ID.String())
		response := &packet.ResponseAuthentication{
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
	err = hostHandler.GetNetworkCommonFacade().GetCascade().SendToNextLayer(data.Data, data.RPC.Method, data.RPC.Args)
	if err != nil {
		log.Debug("failed to send message to next cascade layer")
	}

	return packetBuilder.Response(response).Build(), err
}

func processActiveNodes(hostHandler hosthandler.HostHandler, packetBuilder packet.Builder) (*packet.Packet, error) {
	response := &packet.ResponseActiveNodes{ActiveNodes: hostHandler.GetActiveNodesList()}
	return packetBuilder.Response(response).Build(), nil
}
