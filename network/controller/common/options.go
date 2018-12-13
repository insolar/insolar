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

package common

import (
	"time"
)

// Options contains configuration options for the local host.
type Options struct {
	// The maximum time to wait for a response to ping request.
	PingTimeout time.Duration

	// The maximum time to wait for a response to any packet.
	PacketTimeout time.Duration

	// InfiniteBootstrap bool

	// Bootstrap reconnect timeout
	BootstrapTimeout time.Duration

	// Min bootstrap retry timeout
	MinTimeout time.Duration

	// Max bootstrap retry timeout
	MaxTimeout time.Duration

	// Multiplier for boostrap retry time
	TimeoutMult time.Duration

	// True - infinity tries to bootstrap
	InfinityBootstrap bool

	// HandshakeSession TTL
	HandshakeSessionTTL time.Duration
}
