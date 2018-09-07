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

package packet

import (
	"crypto/rand"
	"log"
	"time"

	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/relay"
	"github.com/insolar/insolar/network/hostnetwork/routing"
	"github.com/insolar/insolar/network/hostnetwork/store"
)

// DispatchPacketType checks message type.
func DispatchPacketType(msg *Packet, origin *host.Host) packetType {
	packetBuilder := NewBuilder().Sender(origin).Receiver(msg.Sender).Type(msg.Type)
	switch msg.Type {
	case TypeFindHost:
		processFindHost(msg, packetBuilder)
	case TypeFindValue:
		processFindValue(msg, packetBuilder)
	case TypeStore:
		processStore(msg, packetBuilder)
	case TypePing:
		processPing(msg, packetBuilder)
	case TypeRPC:
		processRPC(msg, packetBuilder)
	case TypeRelay:
		processRelay(msg, packetBuilder)
	case TypeCheckOrigin:
		processCheckOriginRequest(msg, packetBuilder)
	case TypeAuth:
		processAuthentication(msg, packetBuilder)
	case TypeObtainIP:
		processObtainIPRequest(msg, packetBuilder)
	case TypeRelayOwnership:
		processRelayOwnership(msg, packetBuilder)
	case TypeKnownOuterHosts:
		processKnownOuterHosts(msg, packetBuilder)
	case TypeCheckNodePriv:
		processCheckNodePriv(msg, packetBuilder)
	default:
		log.Println("unknown request type")
	}
}

func processKnownOuterHosts(msg *Packet, packetBuilder Builder) {
	data := msg.Data.(*RequestKnownOuterHosts)

	ID := dht.subnet.HighKnownHosts.ID
	hosts := dht.subnet.HighKnownHosts.OuterHosts
	if data.OuterHosts > hosts {
		ID = data.ID
		hosts = data.OuterHosts
	}
	response := &ResponseKnownOuterHosts{
		ID:         ID,
		OuterHosts: hosts,
	}

	err := dht.transport.SendResponse(msg.RequestID, packetBuilder.Response(response).Build())
	if err != nil {
		log.Println("Failed to send response:", err.Error())
	}
}

func processCheckNodePriv(msg *Packet, packetBuilder Builder) {
	data := msg.Data.(*RequestCheckNodePriv)
	var response ResponseCheckNodePriv

	if dht.confirmNodeRole(data.RoleKey) {
		response.State = Confirmed
	} else {
		response.State = Declined
	}

	err := dht.transport.SendResponse(msg.RequestID, packetBuilder.Response(response).Build())
	if err != nil {
		log.Println("Failed to send response:", err.Error())
	}
}

func processRelayOwnership(msg *Packet, packetBuilder Builder) {
	data := msg.Data.(*RequestRelayOwnership)

	if data.Ready {
		dht.subnet.PossibleProxyIDs = append(dht.subnet.PossibleProxyIDs, msg.Sender.ID.KeyString())
	} else {
		for i, j := range dht.subnet.PossibleProxyIDs {
			if j == msg.Sender.ID.KeyString() {
				dht.subnet.PossibleProxyIDs = append(dht.subnet.PossibleProxyIDs[:i], dht.subnet.PossibleProxyIDs[i+1:]...)
				err := dht.AuthenticationRequest(ctx, "begin", msg.Sender.ID.KeyString())
				if err != nil {
					log.Println("error to send auth request: ", err)
				}
				err = dht.RelayRequest(ctx, "start", msg.Sender.ID.KeyString())
				if err != nil {
					log.Println("error to send relay request: ", err)
				}
				break
			}
		}
	}
	response := &ResponseRelayOwnership{Accepted: true}

	err := dht.transport.SendResponse(msg.RequestID, packetBuilder.Response(response).Build())
	if err != nil {
		log.Println("Failed to send response:", err.Error())
	}
}

func processFindHost(msg *Packet, packetBuilder Builder) {
	ht := dht.htFromCtx(ctx)
	data := msg.Data.(*RequestDataFindHost)
	dht.addHost(ctx, routing.NewRouteHost(msg.Sender))
	closest := ht.GetClosestContacts(routing.MaxContactsInBucket, data.Target, []*host.Host{msg.Sender})
	response := &ResponseDataFindHost{
		Closest: closest.Hosts(),
	}
	err := dht.transport.SendResponse(msg.RequestID, packetBuilder.Response(response).Build())
	if err != nil {
		log.Println("Failed to send response:", err.Error())
	}
}

func processFindValue(msg *Packet, packetBuilder Builder) {
	ht := dht.htFromCtx(ctx)
	data := msg.Data.(*RequestDataFindValue)
	dht.addHost(ctx, routing.NewRouteHost(msg.Sender))
	value, exists := dht.store.Retrieve(data.Target)
	response := &ResponseDataFindValue{}
	if exists {
		response.Value = value
	} else {
		closest := ht.GetClosestContacts(routing.MaxContactsInBucket, data.Target, []*host.Host{msg.Sender})
		response.Closest = closest.Hosts()
	}
	err := dht.transport.SendResponse(msg.RequestID, packetBuilder.Response(response).Build())
	if err != nil {
		log.Println("Failed to send response:", err.Error())
	}
}

func processStore(msg *Packet, packetBuilder Builder) {
	data := msg.Data.(*RequestDataStore)
	dht.addHost(ctx, routing.NewRouteHost(msg.Sender))
	key := store.NewKey(data.Data)
	expiration := dht.getExpirationTime(ctx, key)
	replication := time.Now().Add(dht.options.ReplicateTime)
	err := dht.store.Store(key, data.Data, replication, expiration, false)
	if err != nil {
		log.Println("Failed to store data:", err.Error())
	}
}

func processPing(msg *Packet, packetBuilder Builder) {
	log.Println("recv ping message from " + msg.Sender.Address.String())
	err := dht.transport.SendResponse(msg.RequestID, packetBuilder.Response(nil).Build())
	if err != nil {
		log.Println("Failed to send response:", err.Error())
	}
}

func processRPC(msg *Packet, packetBuilder Builder) {
	data := msg.Data.(*RequestDataRPC)
	dht.addHost(ctx, routing.NewRouteHost(msg.Sender))
	result, err := dht.rpc.Invoke(msg.Sender, data.Method, data.Args)
	response := &ResponseDataRPC{
		Success: true,
		Result:  result,
		Error:   "",
	}
	if err != nil {
		response.Success = false
		response.Error = err.Error()
	}
	err = dht.transport.SendResponse(msg.RequestID, packetBuilder.Response(response).Build())
	if err != nil {
		log.Println("Failed to send response:", err.Error())
	}
}

// Precess relay request.
func processRelay(msg *Packet, packetBuilder Builder) {
	var err error
	if !dht.auth.authenticatedHosts[msg.Sender.ID.KeyString()] {
		log.Print("relay request from unknown host rejected")
		response := &ResponseRelay{
			State: relay.NoAuth,
		}

		err = dht.transport.SendResponse(msg.RequestID, packetBuilder.Response(response).Build())
	} else {
		data := msg.Data.(*RequestRelay)
		dht.addHost(ctx, routing.NewRouteHost(msg.Sender))

		var state relay.State

		switch data.Command {
		case StartRelay:
			err = dht.relay.AddClient(msg.Sender)
			state = relay.Started
		case StopRelay:
			err = dht.relay.RemoveClient(msg.Sender)
			state = relay.Stopped
		default:
			state = relay.Unknown
		}

		if err != nil {
			state = relay.Error
		}

		response := &ResponseRelay{
			State: state,
		}

		err = dht.transport.SendResponse(msg.RequestID, packetBuilder.Response(response).Build())
	}
	if err != nil {
		log.Println("Failed to send response:", err.Error())
	}
}

func processAuthentication(msg *Packet, packetBuilder Builder) {
	data := msg.Data.(*RequestAuth)
	switch data.Command {
	case BeginAuth:
		if dht.auth.authenticatedHosts[msg.Sender.ID.KeyString()] {
			// TODO: whats next?
			response := &ResponseAuth{
				Success:       false,
				AuthUniqueKey: nil,
			}

			err := dht.transport.SendResponse(msg.RequestID, packetBuilder.Response(response).Build())
			if err != nil {
				log.Println("Failed to send response:", err)
			}
			break
		}
		key := make([]byte, 512)
		_, err := rand.Read(key) // crypto/rand
		if err != nil {
			log.Println("failed to create auth key. ", err)
			return
		}
		dht.auth.SentKeys[msg.Sender.ID.KeyString()] = key
		response := &ResponseAuth{
			Success:       true,
			AuthUniqueKey: key,
		}

		err = dht.transport.SendResponse(msg.RequestID, packetBuilder.Response(response).Build())
		if err != nil {
			log.Println("Failed to send response:", err)
		}
		// TODO process verification msg.Sender host
		// confirmed
		err = dht.CheckOriginRequest(ctx, msg.Sender.ID.KeyString())
		if err != nil {
			log.Println("error: ", err)
		}
	case RevokeAuth:
		delete(dht.auth.authenticatedHosts, msg.Sender.ID.KeyString())
		response := &ResponseAuth{
			Success:       true,
			AuthUniqueKey: nil,
		}

		err := dht.transport.SendResponse(msg.RequestID, packetBuilder.Response(response).Build())
		if err != nil {
			log.Println("Failed to send response:", err)
		}
	default:
		log.Println("unknown auth command")
	}
}

func processCheckOriginRequest(msg *Packet, packetBuilder Builder) {
	dht.auth.mut.Lock()
	defer dht.auth.mut.Unlock()
	if key, ok := dht.auth.ReceivedKeys[msg.Sender.ID.KeyString()]; ok {
		response := &ResponseCheckOrigin{AuthUniqueKey: key}
		err := dht.transport.SendResponse(msg.RequestID, packetBuilder.Response(response).Build())
		if err != nil {
			log.Println("Failed to send check origin response:", err)
		}
	} else {
		log.Println("CheckOrigin request from unregistered host")
	}
}

func processObtainIPRequest(msg *Packet, packetBuilder Builder) {
	response := &ResponseObtainIP{IP: msg.RemoteAddress}
	err := dht.transport.SendResponse(msg.RequestID, packetBuilder.Response(response).Build())
	if err != nil {
		log.Println("Failed to send obtain IP response:", err)
	}
}
