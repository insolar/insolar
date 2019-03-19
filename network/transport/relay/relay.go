/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted (subject to the limitations in the disclaimer below) provided that
 * the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *  * Neither the name of Insolar Technologies nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
 * BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
 * CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING,
 * BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package relay

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/transport/host"

	"errors"
)

// State is alias for relaying state
type State int

const (
	// Unknown unknown relay state.
	Unknown = State(iota + 1)
	// Started is relay type means relaying started.
	Started
	// Stopped is relay type means relaying stopped.
	Stopped
	// Error is relay type means error state change.
	Error
	// NoAuth - this error returns if host tries to send relay request but not authenticated.
	NoAuth
)

// Relay Interface for relaying
type Relay interface {
	// AddClient add client to relay list.
	AddClient(host *host.Host) error
	// RemoveClient removes client from relay list.
	RemoveClient(host *host.Host) error
	// ClientsCount - clients count.
	ClientsCount() int
	// NeedToRelay returns true if origin host is proxy for target host.
	NeedToRelay(targetAddress string) bool
}

type relay struct {
	clients []*host.Host
}

// NewRelay constructs relay list.
func NewRelay() Relay {
	return &relay{
		clients: make([]*host.Host, 0),
	}
}

// AddClient add client to relay list.
func (r *relay) AddClient(host *host.Host) error {
	if _, n := r.findClient(host.NodeID); n != nil {
		return errors.New("client exists already")
	}
	r.clients = append(r.clients, host)
	return nil
}

// RemoveClient removes client from relay list.
func (r *relay) RemoveClient(host *host.Host) error {
	idx, n := r.findClient(host.NodeID)
	if n == nil {
		return errors.New("client not found")
	}
	r.clients = append(r.clients[:idx], r.clients[idx+1:]...)
	return nil
}

// ClientsCount - returns clients count.
func (r *relay) ClientsCount() int {
	return len(r.clients)
}

// NeedToRelay returns true if origin host is proxy for target host.
func (r *relay) NeedToRelay(targetAddress string) bool {
	for i := 0; i < r.ClientsCount(); i++ {
		if r.clients[i].Address.String() == targetAddress {
			return true
		}
	}
	return false
}

func (r *relay) findClient(id core.RecordRef) (int, *host.Host) {
	for idx, hostIterator := range r.clients {
		if hostIterator.NodeID.Equal(id) {
			return idx, hostIterator
		}
	}
	return -1, nil
}
