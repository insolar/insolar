// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package tests

import (
	"context"
	"time"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
)

const defaultNshGenerationDelay = time.Millisecond * 0

var RoundTimingsFor1s = api.RoundTimings{
	StartPhase0At: 100 * time.Millisecond, // Not scaled

	StartPhase1RetryAt: 00 * time.Millisecond, // 0 for no retries
	EndOfPhase1:        250 * time.Millisecond,
	EndOfPhase2:        400 * time.Millisecond,
	EndOfPhase3:        500 * time.Millisecond,
	EndOfConsensus:     600 * time.Millisecond,

	BeforeInPhase2ChasingDelay: 0 * time.Millisecond,
	BeforeInPhase3ChasingDelay: 0 * time.Millisecond,
}

var EphemeralRoundTimingsFor1s = api.RoundTimings{
	StartPhase0At: 100 * time.Millisecond, // Not scaled

	StartPhase1RetryAt: 0 * time.Millisecond, // 0 for no retries
	EndOfPhase1:        200 * time.Millisecond,
	EndOfPhase2:        600 * time.Millisecond,
	EndOfPhase3:        800 * time.Millisecond,
	EndOfConsensus:     900 * time.Millisecond,

	BeforeInPhase2ChasingDelay: 0 * time.Millisecond,
	BeforeInPhase3ChasingDelay: 0 * time.Millisecond,
}

func NewEmuLocalConfig(ctx context.Context) api.LocalNodeConfiguration {
	r := emuLocalConfig{timings: RoundTimingsFor1s, ephemeralTimings: EphemeralRoundTimingsFor1s, ctx: ctx}
	return &r
}

type emuLocalConfig struct {
	timings          api.RoundTimings
	ephemeralTimings api.RoundTimings
	ctx              context.Context
}

func (r *emuLocalConfig) GetNodeCountHint() int {
	return 10
}

func (r *emuLocalConfig) GetParentContext() context.Context {
	return r.ctx
}

func (r *emuLocalConfig) PublicKeyStore() {
}

func (r *emuLocalConfig) AsPublicKeyStore() cryptkit.PublicKeyStore {
	return r
}

func (r *emuLocalConfig) PrivateKeyStore() {
}

func (r *emuLocalConfig) getConsensusTimings(t api.RoundTimings, nextPulseDelta uint16) api.RoundTimings {
	if nextPulseDelta == 1 {
		return t
	}
	m := time.Duration(nextPulseDelta) // this is NOT a duration, but a multiplier

	t.StartPhase0At *= 1 // don't scale!
	t.StartPhase1RetryAt *= m
	t.EndOfPhase1 *= m
	t.EndOfPhase2 *= m
	t.EndOfPhase3 *= m
	t.EndOfConsensus *= m
	t.BeforeInPhase2ChasingDelay *= m
	t.BeforeInPhase3ChasingDelay *= m

	return t
}

func (r *emuLocalConfig) GetConsensusTimings(nextPulseDelta uint16) api.RoundTimings {
	return r.getConsensusTimings(r.timings, nextPulseDelta)
}

func (r *emuLocalConfig) GetEphemeralTimings(nextPulseDelta uint16) api.RoundTimings {
	return r.getConsensusTimings(r.ephemeralTimings, nextPulseDelta)
}

func (r *emuLocalConfig) GetSecretKeyStore() cryptkit.SecretKeyStore {
	return r
}
