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

package configuration

// JetCoordinator holds configuration for JetCoordinator.
type JetCoordinator struct {
	VirtualExecutor string
	LightExecutor   string
	HeavyExecutor   string

	VirtualValidators []string
	LightValidators   []string
}

// JetCoordinator creates new default JetCoordinator configuration.
func NewJetCoordinator() JetCoordinator {
	return JetCoordinator{
		VirtualExecutor: "",
		LightExecutor:   "",
		HeavyExecutor:   "",

		VirtualValidators: make([]string, 0),
		LightValidators:   make([]string, 0),
	}
}

// Ledger holds configuration for ledger.
type Ledger struct {
	// DataDirectory is a directory where database's files live.
	DataDirectory string
	// TxRetriesOnConflict defines how many retries on transaction conflicts
	// storage update methods should do.
	TxRetriesOnConflict int
	// JetCoordinator defines jet coordinator configuration.
	JetCoordinator JetCoordinator
}

// NewLedger creates new default Ledger configuration.
func NewLedger() Ledger {
	return Ledger{
		DataDirectory:       "./data",
		TxRetriesOnConflict: 3,
		JetCoordinator:      NewJetCoordinator(),
	}
}
