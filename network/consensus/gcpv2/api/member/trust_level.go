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

package member

type TrustLevel int8

const (
	FraudByBlacklist TrustLevel = -5 // in the blacklist
	FraudByNetwork   TrustLevel = -4 // >2/3 of network have indicated fraud
	FraudByNeighbors TrustLevel = -3 // >50% of neighborhood have indicated fraud
	FraudBySome      TrustLevel = -2 // some nodes have indicated fraud
	// unused                   = -1

	UnknownTrust TrustLevel = 0 // initial state

	TrustBySelf      TrustLevel = 1 // node has provided a liveness proof or NSH
	TrustBySome      TrustLevel = 2 // some nodes have indicated trust (same NSH)
	TrustByNeighbors TrustLevel = 3 // >50% of neighborhood have indicated trust
	TrustByNetwork   TrustLevel = 4 // >2/3 of network have indicated trust
	TrustByMandate   TrustLevel = 5 // on- or off-network node with a temporary mandate, e.g. pulsar or discovery
	TrustByCouncil   TrustLevel = 6 // on- or off-network node with a permanent mandate

	LocalSelfTrust  = TrustByNeighbors // MUST be not less than TrustByNeighbors
	FraudByThisNode = FraudByNeighbors // fraud is detected by this node
)

func (v TrustLevel) abs() int8 {
	if v >= 0 {
		return int8(v)
	}
	return int8(-v)
}

func (v TrustLevel) WrapRange(hi TrustLevel) uint16 {
	return uint16(v) | uint16(hi)<<8
}

func UnwrapTrustRange(wrapped uint16) (lo, hi TrustLevel) {
	return TrustLevel(wrapped), TrustLevel(wrapped >> 8)
}

// Updates only to better/worse levels. Negative level of the same magnitude prevails.
func (v *TrustLevel) Update(newLevel TrustLevel) (modified bool) {
	if newLevel == UnknownTrust || newLevel == *v {
		return false
	}
	if newLevel > UnknownTrust {
		if newLevel.abs() <= v.abs() {
			return false
		}
	} else { // negative prevails hence update on |newLevel| == |v|
		if newLevel.abs() < v.abs() {
			return false
		}
	}
	*v = newLevel
	return true
}

func (v *TrustLevel) UpdateKeepNegative(newLevel TrustLevel) (modified bool) {
	if newLevel > UnknownTrust && *v < UnknownTrust {
		return false
	}
	return v.Update(newLevel)
}

func (v *TrustLevel) IsNegative() bool {
	return *v < UnknownTrust
}
