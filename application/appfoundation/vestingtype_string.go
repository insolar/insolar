// Code generated by "stringer -type=VestingType"; DO NOT EDIT.

package appfoundation

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DefaultVesting-0]
	_ = x[Vesting1-1]
	_ = x[Vesting2-2]
	_ = x[Vesting3-3]
	_ = x[Vesting4-4]
}

const _VestingType_name = "DefaultVestingVesting1Vesting2Vesting3Vesting4"

var _VestingType_index = [...]uint8{0, 14, 22, 30, 38, 46}

func (i VestingType) String() string {
	if i < 0 || i >= VestingType(len(_VestingType_index)-1) {
		return "VestingType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _VestingType_name[_VestingType_index[i]:_VestingType_index[i+1]]
}
