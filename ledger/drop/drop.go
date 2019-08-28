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

package drop

import (
	"context"

	"github.com/insolar/insolar/insolar"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/drop.Modifier -o ./ -s _mock.go -g

// Modifier provides an interface for modifying jetdrops.
type Modifier interface {
	Set(ctx context.Context, drop Drop) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/drop.Accessor -o ./ -s _mock.go -g

// Accessor provides an interface for accessing jetdrops.
type Accessor interface {
	ForPulse(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) (Drop, error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/drop.Cleaner -o ./ -s _mock.go -g

// Cleaner provides an interface for removing jetdrops from a storage.
type Cleaner interface {
	DeleteForPN(ctx context.Context, pulse insolar.PulseNumber)
}
