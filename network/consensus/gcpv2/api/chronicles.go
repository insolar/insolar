// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package api

import (
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/pulse"
)

type ConsensusChronicles interface {
	GetProfileFactory(ksf cryptkit.KeyStoreFactory) profiles.Factory

	GetActiveCensus() census.Active
	GetExpectedCensus() census.Expected
	GetLatestCensus() (lastCensus census.Operational, expectedCensus bool)
	GetRecentCensus(pn pulse.Number) census.Operational
	// FindArchivedCensus(pn common.PulseNumber) Archived
}
