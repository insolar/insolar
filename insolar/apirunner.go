// Copyright 2020 Insolar Network Ltd.
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

package insolar

import (
	"context"
)

// APIRunner
type APIRunner interface {
	IsAPIRunner() bool
}

//go:generate minimock -i github.com/insolar/insolar/insolar.AvailabilityChecker -o ../testutils -s _mock.go -g

// AvailabilityChecker component checks if insolar network can't process any new requests
type AvailabilityChecker interface {
	IsAvailable(context.Context) bool
}
