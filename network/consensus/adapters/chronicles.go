package adapters

import (
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/censusimpl"
)

func NewChronicles(pf profiles.Factory) censusimpl.LocalConsensusChronicles {
	return censusimpl.NewLocalChronicles(pf)
}

func NewCensusForJoiner(
	localNode profiles.StaticProfile,
	vc census.VersionedRegistries,
	vf cryptkit.SignatureVerifierFactory,
) *censusimpl.PrimingCensusTemplate {

	return censusimpl.NewPrimingCensusForJoiner(localNode, vc, vf, true)
}

func NewCensus(
	localNode profiles.StaticProfile,
	nodes []profiles.StaticProfile,
	vc census.VersionedRegistries,
	vf cryptkit.SignatureVerifierFactory,
) *censusimpl.PrimingCensusTemplate {

	return censusimpl.NewPrimingCensus(nodes, localNode, vc, vf, true)
}
