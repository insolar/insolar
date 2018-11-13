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
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
)

type builder struct {
	sender *host.Host
	t      types.PacketType
	data   interface{}
}

func (b *builder) Type(packetType types.PacketType) network.RequestBuilder {
	b.t = packetType
	return b
}

func (b *builder) Data(data interface{}) network.RequestBuilder {
	b.data = data
	return b
}

func (b *builder) GetSender() core.RecordRef {
	return b.sender.NodeID
}

func (b *builder) GetSenderHost() *host.Host {
	return b.sender
}

func (b *builder) GetType() types.PacketType {
	return b.t
}

func (b *builder) GetData() interface{} {
	return b.data
}

func (b *builder) Build() network.Request {
	return b
}
