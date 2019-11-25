///
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
///

package dropbag

import (
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/pulse"
)

type JetPulse interface {
	GetOnlinePopulation() census.OnlinePopulation
	GetPulseData() pulse.Range
}

type JetSectionId uint16

const (
	DefaultSection JetSectionId = iota // MapSection
	ControlSection                     // drop/dropbag lifecycle
	DustSection                        // transient, general, stays for some time (e.g. log)
	GasSection                         // transient, requests, stays until processed
)

type JetDrop interface {
	PulseNumber() pulse.Number
	GetGlobulaPulse() JetPulse

	FindEntryByKey(longbits.ByteString) JetDropEntry
	//	GetSectionDirectory(JetSectionId) JetDropSection
	GetSection(JetSectionId) JetDropSection
}

type JetSectionType uint8

const (
	DirectorySection JetSectionType = 1 << iota //
	TransientSection
	CustomCryptographySection
	//HeavyPayload
)

type JetSectionDeclaration interface {
	HasDirectory() bool
	//IsSorted
	HasPayload() bool
}

type JetSectionDirectory interface {
	FindByKey(longbits.ByteString) JetDropEntry
	EnumKeys()
}

type JetDropSection interface {
	EnumEntries()
}

type JetDropEntry interface {
	Key() longbits.ByteString
	Section() JetDropSection
	IsAvailable() bool
	Data() []byte
	// ProjectionCache()
}

type KeySet interface {
	// inclusive or exclusive

}
