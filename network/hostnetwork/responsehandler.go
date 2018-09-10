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
	"bytes"
	"log"
	"time"

	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/relay"
	"github.com/pkg/errors"
)

func handleRelayOwnership(dht *DHT, response *packet.ResponseRelayOwnership, target string) {
	if response.Accepted {
		dht.Subnet.PossibleRelayIDs = append(dht.Subnet.PossibleRelayIDs, target)
	}
}

func handleKnownOuterHosts(dht *DHT, ctx Context, response *packet.ResponseKnownOuterHosts, targetID string) error {
	var err error
	if response.OuterHosts > dht.Subnet.HighKnownHosts.OuterHosts { // update data
		dht.Subnet.HighKnownHosts.OuterHosts = response.OuterHosts
		dht.Subnet.HighKnownHosts.ID = response.ID
	}
	if (response.OuterHosts > dht.Subnet.HighKnownHosts.SelfKnownOuterHosts) &&
		(dht.proxy.ProxyHostsCount() == 0) {
		err = AuthenticationRequest(dht, ctx, "begin", targetID)
		if err != nil {
			return err
		}
		err = RelayRequest(dht, ctx, "start", targetID)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleRelayResponse(dht *DHT, ctx Context, response *packet.ResponseRelay, targetID string) error {
	var err error
	switch response.State {
	case relay.Stopped:
		// stop use this address as relay
		dht.proxy.RemoveProxyHost(targetID)
		err = nil
	case relay.Started:
		// start use this address as relay
		dht.proxy.AddProxyHost(targetID)
		err = nil
	case relay.NoAuth:
		err = errors.New("unable to execute relay because this host not authenticated")
	case relay.Unknown:
		err = errors.New("unknown relay command")
	case relay.Error:
		err = errors.New("relay request error")
	default:
		// unknown state/failed to change state
		err = errors.New("unknown response state")
	}
	return err
}

func handleCheckOriginResponse(dht *DHT, response *packet.ResponseCheckOrigin, targetID string) {
	if bytes.Equal(response.AuthUniqueKey, dht.Auth.SentKeys[targetID]) {
		delete(dht.Auth.SentKeys, targetID)
		dht.Auth.AuthenticatedHosts[targetID] = true
	}
}

func handleCheckNodePrivResponse(dht *DHT, response *packet.ResponseCheckNodePriv, toleKey string) error {
	switch response.State {
	case packet.Error:
		return errors.New(response.Error)
	case packet.Confirmed:
		return nil
	case packet.Declined:
		// TODO: set default unconfirmed role
		// dht.nodesMap[dht.nodesIDMap[roleKey]].SetNodeRole("Unconfirmed default role")
	}
	return nil
}

func handleAuthResponse(dht *DHT, response *packet.ResponseAuth, target string) error {
	var err error
	if (len(response.AuthUniqueKey) != 0) && response.Success {
		dht.Auth.Mut.Lock()
		defer dht.Auth.Mut.Unlock()
		dht.Auth.ReceivedKeys[target] = response.AuthUniqueKey
		err = nil
	} else {
		if response.Success && (len(response.AuthUniqueKey) == 0) { // revoke success
			return err
		}
		if !response.Success {
			err = errors.New("authentication unsuccessful")
		} else if len(response.AuthUniqueKey) == 0 {
			err = errors.New("wrong auth unique key received")
		}
	}
	return err
}

func handleObtainIPResponse(dht *DHT, response *packet.ResponseObtainIP, targetID string) error {
	if response.IP != "" {
		dht.Subnet.SubnetIDs[response.IP] = append(dht.Subnet.SubnetIDs[response.IP], targetID)
	} else {
		return errors.New("received empty IP")
	}
	return nil
}

// RelayRequest sends relay request to target.
func RelayRequest(dht *DHT, ctx Context, command, targetID string) error {
	var typedCommand packet.CommandType
	targetHost, exist, err := dht.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("target for relay request not found")
		return err
	}

	switch command {
	case "start":
		typedCommand = packet.StartRelay
	case "stop":
		typedCommand = packet.StopRelay
	default:
		err = errors.New("unknown command")
		return err
	}
	request := packet.NewRelayPacket(typedCommand, dht.HtFromCtx(ctx).Origin, targetHost)
	future, err := dht.transport.SendRequest(request)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	select {
	case rsp := <-future.Result():
		if rsp == nil {
			err = errors.New("chanel closed unexpectedly")
			return err
		}

		response := rsp.Data.(*packet.ResponseRelay)
		err = handleRelayResponse(dht, ctx, response, targetID)
		if err != nil {
			return err
		}

	case <-time.After(dht.options.PacketTimeout):
		future.Cancel()
		err = errors.New("timeout")
		return err
	}

	return nil
}

func sendRelayedRequest(dht *DHT, request *packet.Packet, ctx Context) {
	_, err := dht.transport.SendRequest(request)
	if err != nil {
		log.Println(err)
	}
}
