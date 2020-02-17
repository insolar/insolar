// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package coreapi

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
			pollingTimer.ClearExpired()

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
