//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package core

import (
	"context"
	"time"

	"github.com/insolar/insolar/network/consensus/common/chaser"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
)

type PollingWorker struct {
	ctx context.Context

	polls   []api.MaintenancePollFunc
	pollCmd chan api.MaintenancePollFunc
}

func (p *PollingWorker) Start(ctx context.Context, pollingInterval time.Duration) {
	if p.ctx != nil {
		panic("illegal state")
	}
	p.ctx = ctx
	p.pollCmd = make(chan api.MaintenancePollFunc, 10)

	go p.pollingWorker(pollingInterval)
}

func (p *PollingWorker) AddPoll(fn api.MaintenancePollFunc) {
	p.pollCmd <- fn
}

func (p *PollingWorker) pollingWorker(pollingInterval time.Duration) {
	pollingTimer := chaser.NewChasingTimer(pollingInterval)

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-pollingTimer.Channel():
			if p.scanPolls() {
				pollingTimer.RestartChase()
			}
		case add := <-p.pollCmd:
			if add == nil {
				continue
			}
			p.polls = append(p.polls, add)
			if len(p.polls) == 1 {
				pollingTimer.RestartChase()
			}
		}
	}
}

func (p *PollingWorker) scanPolls() bool {
	j := 0
	for i, poll := range p.polls {
		if !poll(p.ctx) {
			p.polls[i] = nil
			continue
		}
		if i != j {
			p.polls[i] = nil
			p.polls[j] = poll
		}
		j++
	}
	p.polls = p.polls[:j]
	return j > 0
}
