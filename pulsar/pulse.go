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

package pulsar

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
)

// NewPulse creates a new pulse with using of custom GeneratedEntropy Generator
func NewPulse(numberDelta uint32, previousPulseNumber insolar.PulseNumber, entropyGenerator entropygenerator.EntropyGenerator) *insolar.Pulse {
	previousPulseNumber += insolar.PulseNumber(numberDelta)
	return &insolar.Pulse{
		PulseNumber:     previousPulseNumber,
		NextPulseNumber: previousPulseNumber + insolar.PulseNumber(numberDelta),
		Entropy:         entropyGenerator.GenerateEntropy(),
	}
}
