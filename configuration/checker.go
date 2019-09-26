//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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
