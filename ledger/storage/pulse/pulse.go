/*
 *    Copyright 2019 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package pulse

import (
	"context"

	"github.com/insolar/insolar/core"
)

// Accessor provides methods for accessing pulses.
type Accessor interface {
	ForPulseNumber(context.Context, core.PulseNumber) (core.Pulse, error)
	Latest(ctx context.Context) (core.Pulse, error)
}

// Shifter provides method for removing pulses from storage.
type Shifter interface {
	Shift(ctx context.Context) (pulse core.Pulse, err error)
}

// Appender provides method for appending pulses to storage.
type Appender interface {
	Append(ctx context.Context, pulse core.Pulse) error
}

// Calculator performs calculations for pulses.
type Calculator interface {
	Forwards(ctx context.Context, pn core.PulseNumber, steps int) (core.Pulse, error)
	Backwards(ctx context.Context, pn core.PulseNumber, steps int) (core.Pulse, error)
}
