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
	"time"

	"github.com/insolar/insolar/network/hostnetwork/hosthandler"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/transport"
	"github.com/pkg/errors"
)

// RelayRequest sends relay request to target.
func RelayRequest(hostHandler hosthandler.HostHandler, command, targetID string) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	if err != nil {
		return err
	}
	var typedCommand packet.CommandType
	targetHost, exist, err := hostHandler.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("RelayRequest: target for relay request not found")
		return err
	}

	switch command {
	case "start":
		typedCommand = packet.StartRelay
	case "stop":
		typedCommand = packet.StopRelay
	default:
		err = errors.New("RelayRequest: unknown command")
		return err
	}
	builder := packet.NewBuilder()
	request := builder.Type(packet.TypeRelay).
		Sender(hostHandler.HtFromCtx(ctx).Origin).
		Receiver(targetHost).
		Request(&packet.RequestRelay{Command: typedCommand}).
		Build()
	future, err := hostHandler.SendRequest(request)

	if err != nil {
		return err
	}

	return checkResponse(hostHandler, future, targetID, request)
}

// CheckOriginRequest send a request to check target host originality
func CheckOriginRequest(hostHandler hosthandler.HostHandler, targetID string) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	targetHost, exist, err := hostHandler.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("CheckOriginRequest: target for relay request not found")
		return err
	}

	builder := packet.NewBuilder()
	request := builder.Type(packet.TypeCheckOrigin).
		Sender(hostHandler.HtFromCtx(ctx).Origin).
		Receiver(targetHost).
		Request(&packet.RequestCheckOrigin{}).
		Build()
	future, err := hostHandler.SendRequest(request)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return checkResponse(hostHandler, future, targetID, request)
}

// AuthenticationRequest sends an authentication request.
func AuthenticationRequest(hostHandler hosthandler.HostHandler, command, targetID string) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	targetHost, exist, err := hostHandler.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("AuthenticationRequest: target for auth request not found")
		return err
	}

	origin := hostHandler.HtFromCtx(ctx).Origin
	var authCommand packet.CommandType
	switch command {
	case "begin":
		authCommand = packet.BeginAuth
	case "revoke":
		authCommand = packet.RevokeAuth
	default:
		err = errors.New("AuthenticationRequest: unknown command")
		return err
	}
	builder := packet.NewBuilder()
	request := builder.Type(packet.TypeAuth).
		Sender(origin).
		Receiver(targetHost).
		Request(&packet.RequestAuth{Command: authCommand}).
		Build()
	future, err := hostHandler.SendRequest(request)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return checkResponse(hostHandler, future, targetID, request)
}

// ObtainIPRequest is request to self IP obtaining.
func ObtainIPRequest(hostHandler hosthandler.HostHandler, targetID string) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	targetHost, exist, err := hostHandler.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("ObtainIPRequest: target for relay request not found")
		return err
	}

	origin := hostHandler.HtFromCtx(ctx).Origin
	request := packet.NewObtainIPPacket(origin, targetHost)

	future, err := hostHandler.SendRequest(request)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return checkResponse(hostHandler, future, targetID, request)
}

// RelayOwnershipRequest sends a relay ownership request.
func RelayOwnershipRequest(hostHandler hosthandler.HostHandler, targetID string) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	if err != nil {
		return err
	}
	targetHost, exist, err := hostHandler.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("relayOwnershipRequest: target for relay request not found")
		return err
	}

	request := packet.NewRelayOwnershipPacket(hostHandler.HtFromCtx(ctx).Origin, targetHost, true)
	future, err := hostHandler.SendRequest(request)

	if err != nil {
		return err
	}

	return checkResponse(hostHandler, future, targetID, request)
}

func checkNodePrivRequest(hostHandler hosthandler.HostHandler, targetID string, roleKey string) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	if err != nil {
		return err
	}
	targetHost, exist, err := hostHandler.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("checkNodePrivRequest: target for check node privileges request not found")
		return err
	}

	origin := hostHandler.HtFromCtx(ctx).Origin
	builder := packet.NewBuilder()
	request := builder.Type(packet.TypeCheckNodePriv).Sender(origin).Receiver(targetHost).Request(&packet.RequestCheckNodePriv{RoleKey: "test string"}).Build()
	future, err := hostHandler.SendRequest(request)

	if err != nil {
		return err
	}

	return checkResponse(hostHandler, future, targetID, request)
}

func knownOuterHostsRequest(hostHandler hosthandler.HostHandler, targetID string, hosts int) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	if err != nil {
		return err
	}
	targetHost, exist, err := hostHandler.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("knownOuterHostsRequest: target for relay request not found")
		return err
	}

	request := packet.NewKnownOuterHostsPacket(hostHandler.HtFromCtx(ctx).Origin, targetHost, hosts)
	future, err := hostHandler.SendRequest(request)

	if err != nil {
		return err
	}

	return checkResponse(hostHandler, future, targetID, request)
}

func checkResponse(hostHandler hosthandler.HostHandler, future transport.Future, targetID string, request *packet.Packet) error {
	var err error
	select {
	case rsp := <-future.Result():
		if rsp == nil {
			return err
		}

		switch request.Type {
		case packet.TypeKnownOuterHosts:
			response := rsp.Data.(*packet.ResponseKnownOuterHosts)
			err = handleKnownOuterHosts(hostHandler, response, targetID)
		case packet.TypeCheckOrigin:
			response := rsp.Data.(*packet.ResponseCheckOrigin)
			handleCheckOriginResponse(hostHandler, response, targetID)
		case packet.TypeAuth:
			response := rsp.Data.(*packet.ResponseAuth)
			err = handleAuthResponse(hostHandler, response, targetID)
		case packet.TypeObtainIP:
			response := rsp.Data.(*packet.ResponseObtainIP)
			err = handleObtainIPResponse(hostHandler, response, targetID)
		case packet.TypeRelayOwnership:
			response := rsp.Data.(*packet.ResponseRelayOwnership)
			handleRelayOwnership(hostHandler, response, targetID)
		case packet.TypeCheckNodePriv:
			response := rsp.Data.(*packet.ResponseCheckNodePriv)
			err = handleCheckNodePrivResponse(hostHandler, response)
		case packet.TypeRelay:
			response := rsp.Data.(*packet.ResponseRelay)
			err = handleRelayResponse(hostHandler, response, targetID)
		case packet.TypeCascadeSend:
			response := rsp.Data.(*packet.ResponseCascadeSend)
			if !response.Success {
				err = errors.New(response.Error)
			}
		}

		if err != nil {
			return err
		}

	case <-time.After(hostHandler.GetPacketTimeout()):
		future.Cancel()
		err = errors.New("knownOuterHostsRequest: timeout")
		return err
	}
	return nil
}
