package transport

import "math"

type NeighbourhoodSizes struct {
	NeighbourhoodSize           int
	NeighbourhoodTrustThreshold int
	JoinersPerNeighbourhood     int
	JoinersBoost                int
	ExtendingNeighbourhoodLimit uint8
}

func (sizes *NeighbourhoodSizes) VerifySizes() {
	if sizes.NeighbourhoodSize < 4 {
		panic("neighbourSize can not be less than 4")
	}
	if sizes.NeighbourhoodTrustThreshold < 1 || sizes.NeighbourhoodTrustThreshold > math.MaxUint8 {
		panic("neighbourhood trust threshold must be in [1..MaxUint8]")
	}
	// if neighbourSize > math.MaxInt8 { panic("neighbourSize can not be more than 127") }
	if sizes.JoinersPerNeighbourhood < 2 {
		panic("neighbourJoiners can not be less than 2")
	}
	if sizes.JoinersBoost < 0 {
		panic("joinersBoost can not be less than 0")
	}
	if sizes.JoinersBoost+sizes.JoinersPerNeighbourhood > sizes.NeighbourhoodSize-1 {
		panic("joiners + boost are more than neighbourSize - 1")
	}
}
