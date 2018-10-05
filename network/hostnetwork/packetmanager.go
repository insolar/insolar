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
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/hostnetwork/hosthandler"
	"github.com/insolar/insolar/network/hostnetwork/packet"
)

// ParseIncomingPacket detects a packet type.
func ParseIncomingPacket(hostHandler hosthandler.HostHandler, ctx hosthandler.Context, msg *packet.Packet, packetBuilder packet.Builder) (*packet.Packet, error) {
	return DispatchPacketType(hostHandler, ctx, msg, packetBuilder)
}

// BuildContext builds a context for packet.
func BuildContext(cb ContextBuilder, msg *packet.Packet) hosthandler.Context {
	var ctx hosthandler.Context
	var err error
	if msg.Receiver.ID.Bytes() == nil {
		ctx, err = cb.SetDefaultHost().Build()
	} else {
		ctx, err = cb.SetHostByID(msg.Receiver.ID).Build()
	}
	if err != nil {
		// TODO: Do something sane with error!
		log.Error(err) // don't return this error cuz don't know what to do with
	}
	return ctx
}
