// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package artifacts

import (
	"context"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
)

type PulseAccessorLRU struct {
	pulses pulse.Accessor
	client Client
	cache  *lru.Cache
}

func NewPulseAccessorLRU(pulses pulse.Accessor, client Client, size int) *PulseAccessorLRU {
	cache, err := lru.New(size)
	if err != nil {
		panic("failed to init pulse cache")
	}

	return &PulseAccessorLRU{
		pulses: pulses,
		client: client,
		cache:  cache,
	}
}

func (p *PulseAccessorLRU) ForPulseNumber(ctx context.Context, pn insolar.PulseNumber) (insolar.Pulse, error) {
	var (
		foundPulse insolar.Pulse
		err        error
	)

	val, ok := p.cache.Get(pn)
	if ok {
		return val.(insolar.Pulse), nil
	}

	foundPulse, err = p.pulses.ForPulseNumber(ctx, pn)
	if err == nil {
		p.cache.Add(pn, foundPulse)
		return foundPulse, nil
	}

	foundPulse, err = p.client.GetPulse(ctx, pn)
	if err == nil {
		p.cache.Add(pn, foundPulse)
		return foundPulse, nil
	}

	return insolar.Pulse{}, errors.Wrapf(err, "couldn't get pulse for pulse number: %s", pn)
}

func (p *PulseAccessorLRU) Latest(ctx context.Context) (insolar.Pulse, error) {
	return p.pulses.Latest(ctx)
}
