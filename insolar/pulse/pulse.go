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

package pulse

import (
	"context"

	"github.com/insolar/insolar/insolar"
)

//go:generate minimock -i github.com/insolar/insolar/insolar/pulse.Accessor -o ./ -s _mock.go

// Accessor provides methods for accessing pulses.
type Accessor interface {
	ForPulseNumber(context.Context, insolar.PulseNumber) (insolar.Pulse, error)
	Latest(ctx context.Context) (insolar.Pulse, error)
}

//go:generate minimock -i github.com/insolar/insolar/insolar/pulse.Shifter -o ./ -s _mock.go

// Shifter provides method for removing pulses from storage.
type Shifter interface {
	Shift(ctx context.Context, pn insolar.PulseNumber) (err error)
}

//go:generate minimock -i github.com/insolar/insolar/insolar/pulse.Appender -o ./ -s _mock.go

// Appender provides method for appending pulses to storage.
type Appender interface {
	Append(ctx context.Context, pulse insolar.Pulse) error
}

//go:generate minimock -i github.com/insolar/insolar/insolar/pulse.Calculator -o ./ -s _mock.go

// Calculator performs calculations for pulses.
type Calculator interface {
	Forwards(ctx context.Context, pn insolar.PulseNumber, steps int) (insolar.Pulse, error)
	Backwards(ctx context.Context, pn insolar.PulseNumber, steps int) (insolar.Pulse, error)
}
