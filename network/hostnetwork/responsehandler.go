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
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/hostnetwork/hosthandler"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/relay"
	"github.com/pkg/errors"
)

func handleRelayOwnership(hostHandler hosthandler.HostHandler, response *packet.ResponseRelayOwnership, target string) {
	if response.Accepted {
		hostHandler.AddPossibleRelayID(target)
	}
}

func handleKnownOuterHosts(hostHandler hosthandler.HostHandler, response *packet.ResponseKnownOuterHosts, targetID string) error {
	var err error
	if response.OuterHosts > hostHandler.GetOuterHostsCount() { // update data
		hostHandler.SetOuterHostsCount(response.OuterHosts)
		hostHandler.SetHighKnownHostID(response.ID)
	}
	if (response.OuterHosts > hostHandler.GetSelfKnownOuterHosts()) &&
		(hostHandler.GetProxyHostsCount() == 0) {
		err = AuthenticationRequest(hostHandler, "begin", targetID)
		if err != nil {
			return err
		}
		err = RelayRequest(hostHandler, "start", targetID)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleRelayResponse(hostHandler hosthandler.HostHandler, response *packet.ResponseRelay, targetID string) error {
	var err error
	switch response.State {
	case relay.Stopped:
		// stop use this address as relay
		hostHandler.RemoveProxyHost(targetID)
		err = nil
	case relay.Started:
		// start use this address as relay
		hostHandler.AddProxyHost(targetID)
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

func handleCheckOriginResponse(hostHandler hosthandler.HostHandler, response *packet.ResponseCheckOrigin, targetID string) {
	if hostHandler.EqualAuthSentKey(targetID, response.AuthUniqueKey) {
		hostHandler.RemoveAuthSentKeys(targetID)
		hostHandler.SetAuthStatus(targetID, true)
	}
}

func handleCheckNodePrivResponse(hostHandler hosthandler.HostHandler, response *packet.ResponseCheckNodePriv) error {
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

func handleAuthResponse(hostHandler hosthandler.HostHandler, response *packet.ResponseAuth, target string) error {
	var err error
	if (len(response.AuthUniqueKey) != 0) && response.Success {
		hostHandler.AddReceivedKey(target, response.AuthUniqueKey)
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

func handleObtainIPResponse(hostHandler hosthandler.HostHandler, response *packet.ResponseObtainIP, targetID string) error {
	if response.IP != "" {
		hostHandler.AddSubnetID(response.IP, targetID)
	} else {
		return errors.New("received empty IP")
	}
	return nil
}

func sendRelayedRequest(hostHandler hosthandler.HostHandler, request *packet.Packet, ctx hosthandler.Context) {
	_, err := hostHandler.SendRequest(request)
	if err != nil {
		log.Debugln(err)
	}
}
