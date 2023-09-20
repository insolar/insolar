package configuration

// AvailabilityChecker holds configuration for checking is network available for process API calls
type AvailabilityChecker struct {
	Enabled        bool
	KeeperURL      string
	RequestTimeout uint
	CheckPeriod    uint
}

func NewAvailabilityChecker() AvailabilityChecker {
	return AvailabilityChecker{
		Enabled: true,
		// TODO: set local keeperd address when its done
		// TODO: launch it in functests
		KeeperURL:      "",
		RequestTimeout: 15,
		CheckPeriod:    5,
	}
}
