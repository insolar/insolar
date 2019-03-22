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

package core

import (
	"context"
)

type LeaveApproved struct{}

// TerminationHandler handles such node events as graceful stop, abort, etc.
//go:generate minimock -i github.com/insolar/insolar/core.TerminationHandler -o ../testutils -s _mock.go
type TerminationHandler interface {
	Leave(context.Context, PulseNumber) chan LeaveApproved
	OnLeaveApproved(context.Context)
	// Abort forces to stop all node components
	Abort()
}
