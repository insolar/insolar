// Copyright 2020 Insolar Technologies GmbH
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

package configuration

// Ledger holds configuration for ledger.
type LedgerLight struct {
	// JetSplit holds jet split configuration.
	JetSplit JetSplit

	// CleanerDelay holds value of pulses, that should happen before end of LightChainLimit and start
	// of LME's data cleaning
	CleanerDelay int

	// MaxNotificationsPerPulse holds the limit for abandoned requests notifications limit
	MaxNotificationsPerPulse uint

	// FilamentCacheLimit holds the limit for cache items for an object
	FilamentCacheLimit int
}

// JetSplit holds configuration for jet split.
type JetSplit struct {
	// RecordsCountThreshold is a drop threshold in records to perform split for jet.
	ThresholdRecordsCount int
	// ThresholdOverflowCount is a how many times in row ThresholdRecordsCount should be surpassed.
	ThresholdOverflowCount int
	// DepthLimit limits jet tree depth (maximum possible jets = 2^DepthLimit)
	DepthLimit uint8
}

// NewLedger creates new default Ledger configuration.
func NewLedgerLight() LedgerLight {
	return LedgerLight{
		JetSplit: JetSplit{
			// TODO: find best default values
			ThresholdRecordsCount:  100,
			ThresholdOverflowCount: 3,
			DepthLimit:             5, // limit to 32 jets
		},

		CleanerDelay:             3,    // 3 pulses
		MaxNotificationsPerPulse: 100,  // 100 objects
		FilamentCacheLimit:       3000, // 3000 records for every object
	}
}
