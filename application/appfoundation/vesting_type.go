package appfoundation

// VestingType type of vesting process
type VestingType int

//go:generate stringer -type=VestingType
const (
	DefaultVesting VestingType = iota
	Vesting1
	Vesting2
	Vesting3
	Vesting4
)
