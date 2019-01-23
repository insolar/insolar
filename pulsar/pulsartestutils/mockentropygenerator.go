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

// Package pulsartestutil - test utils for pulsar package
package pulsartestutils

import (
	"github.com/insolar/insolar/core"
)

// MockEntropy for pulsar's tests
var MockEntropy = [64]byte{1, 2, 3, 4, 5, 6, 7, 8}

// MockEntropyGenerator implements EntropyGenerator and is being used for tests
type MockEntropyGenerator struct {
}

// GenerateEntropy returns mocked entropy
func (MockEntropyGenerator) GenerateEntropy() core.Entropy {
	return MockEntropy
}
