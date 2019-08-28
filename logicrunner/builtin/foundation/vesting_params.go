package foundation

type VestingParams struct {
	Lockup      int64 `json:"lockupInPulses"`
	Vesting     int64 `json:"vestingInPulses"`
	VestingStep int64 `json:"vestingStepInPulses"`
}
