/*
 *    Copyright 2018 Insolar
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

package fakepulsar

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
)

// Fakepulsar needed when the network starts and can't receive a real pulse.

const delta uint32 = 0
const prevPulseNumber = 0

// GetFakePulse creates and returns a fake pulse.
func GetFakePulse() *core.Pulse {
	return pulsar.NewPulse(delta, prevPulseNumber, &entropygenerator.StandardEntropyGenerator{})
}
