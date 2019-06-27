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

package adapters

import (
	"context"
	"crypto/ecdsa"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
)

var defaultRoundTimings = common2.RoundTimings{
	StartPhase0At: 100 * time.Millisecond, // Not scaled

	StartPhase1RetryAt: 200 * time.Millisecond, // 0 for no retries
	EndOfPhase1:        250 * time.Millisecond,
	EndOfPhase2:        400 * time.Millisecond,
	EndOfPhase3:        500 * time.Millisecond,

	BeforeInPhase2ChasingDelay: 0 * time.Millisecond,
	BeforeInPhase3ChasingDelay: 0 * time.Millisecond,
}

func NewLocalNodeConfiguration(ctx context.Context, keyStore insolar.KeyStore) *LocalNodeConfiguration {
	privateKey, err := keyStore.GetPrivateKey("")
	if err != nil {
		panic(err)
	}

	ecdsaPrivateKey := privateKey.(*ecdsa.PrivateKey)

	return &LocalNodeConfiguration{
		ctx:            ctx,
		timings:        defaultRoundTimings,
		secretKeyStore: NewECDSASecretKeyStore(ecdsaPrivateKey),
	}
}

type LocalNodeConfiguration struct {
	ctx            context.Context
	timings        common2.RoundTimings
	secretKeyStore common.SecretKeyStore
}

func (c *LocalNodeConfiguration) GetParentContext() context.Context {
	return c.ctx
}

func (c *LocalNodeConfiguration) GetConsensusTimings(nextPulseDelta uint16, isJoiner bool) common2.RoundTimings {
	if nextPulseDelta == 1 {
		return c.timings
	}
	m := time.Duration(nextPulseDelta) // this is NOT a duration, but a multiplier
	t := c.timings

	t.StartPhase0At *= 1 // don't scale!
	t.StartPhase1RetryAt *= m
	t.EndOfPhase1 *= m
	t.EndOfPhase2 *= m
	t.EndOfPhase3 *= m
	t.BeforeInPhase2ChasingDelay *= m
	t.BeforeInPhase3ChasingDelay *= m

	return t
}

func (c *LocalNodeConfiguration) GetSecretKeyStore() common.SecretKeyStore {
	return c.secretKeyStore
}
