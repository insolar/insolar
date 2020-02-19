// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
