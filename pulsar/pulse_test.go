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

package pulsar

import (
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/pulsar/pulsartestutils"
	"github.com/stretchr/testify/require"
)

func TestNewPulse(t *testing.T) {
	generator := &pulsartestutils.MockEntropyGenerator{}
	previousPulse := core.PulseNumber(876)
	expectedPulse := previousPulse + core.PulseNumber(configuration.NewPulsar().NumberDelta)

	result := NewPulse(configuration.NewPulsar().NumberDelta, previousPulse, generator)

	require.Equal(t, result.Entropy[:], pulsartestutils.MockEntropy[:])
	require.Equal(t, result.PulseNumber, core.PulseNumber(expectedPulse))
}
