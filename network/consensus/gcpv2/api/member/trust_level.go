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
