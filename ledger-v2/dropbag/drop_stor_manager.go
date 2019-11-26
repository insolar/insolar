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
	"github.com/insolar/insolar/ledger-v2/jetid"
	"github.com/insolar/insolar/ledger-v2/keyset"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/pulse"
)

type DropStorageManager interface {
	// will actually open CompositeDropStorage that may have multiple DropStorage(s) incorporated
	OpenStorage(jetId jetid.ShortJetId, pn pulse.Number) (DropStorage, error)
	OpenPulseStorage(pn pulse.Number) (CompositeDropPerPulseData, error)

	// storages are lazy-closed, but can be marked for closing
	UnusedStorage(CompositeDropStorage)

	// returns nil when the required storage is not open
	GetOpenedStorage(jetId jetid.ShortJetId, pn pulse.Number) DropStorage

	BuildStorage(pr pulse.Range) CompositeDropStorageBuilder
}

type CompositeDropStorageBuilder interface {
	CompositeDropStorageBuilder()
}

type CompositeDropPerPulseData interface {
	CoveringRange() pulse.Range // latest PulseData + earliest pulse number - will not include intermediate PulseData
	PulseDataCount() int
	GetPerPulseData(int) DropPerPulseData
	FindPerPulseData(pulse.Number) DropPerPulseData
}

type CompositeDropStorage interface {
	CoveringRange() pulse.Range // latest PulseData + earliest pulse number - will not include intermediate PulseData
	PulseData() CompositeDropPerPulseData

	// identified by the latest pulse
	GetDropStorage(jetId jetid.ShortJetId, pn pulse.Number) DropStorage

	// Jets
	// Cabinet -> StorageCabinet

	NoticeUsage()
}

type DropType uint8

const (
	_ DropType = iota
	RegularDropType
	SummaryDropType
	ArchivedSummaryDropType
	ArchivedDropType
)

type DropPerPulseData interface {
	PulseRange() pulse.Range // complete range, with all pulse data included
	JetTree() DropJetTree

	OnlinePopulation() census.OnlinePopulation
	OfflinePopulation() census.OfflinePopulation
}

type DropJetTree interface {
	MinDepth() uint8
	MaxDepth() uint8
	Count() int
	PrefixToJetId(prefix jetid.Prefix) jetid.ShortJetId
	KeyToJetId(keyset.Key) jetid.ShortJetId
}

type DropStorage interface {
	Composite() CompositeDropStorage

	JetId() jetid.FullJetId
	PulseNumber() pulse.Number // the rightmost/latest pulse number of this drop
	PulseData() DropPerPulseData

	DropType() DropType

	// a simple lookup method, looks through all primary directories
	FindByKey(keyset.Key)

	FindSection(DropSectionId) DropSection
	FindDirectory(DropSectionId) DropSectionDirectory
	// === one for all jets === hash ref to StorageCabinet?
	// PulseData

	// FindByKey - across sections?
	// DropSection
}

type DropSectionId uint8

const ()

type DropSectionDirectory interface {
	DropSectionId() DropSectionId
	DropSection() DropSection
	FindByKey(keyset.Key)
}

type DropSection interface {
	DropSectionId() DropSectionId
	// Record listing
}
