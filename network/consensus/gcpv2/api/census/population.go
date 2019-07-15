package census

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
)

type OfflinePopulation interface {
	FindRegisteredProfile(identity endpoints.Inbound) profiles.Host
	// FindPulsarProfile(pulsarId PulsarId) PulsarProfile
}

type OnlinePopulation interface {
	FindProfile(nodeID insolar.ShortNodeID) profiles.ActiveNode
	GetCount() int
	GetProfiles() []profiles.ActiveNode
	GetLocalProfile() profiles.LocalNode
}

type EvictedPopulation interface {
	FindProfile(nodeID insolar.ShortNodeID) profiles.EvictedNode
	GetCount() int
	GetProfiles() []profiles.EvictedNode
}

type PopulationBuilder interface {
	GetCount() int
	AddJoinerProfile(intro profiles.NodeIntroProfile) profiles.Updatable
	RemoveProfile(nodeID insolar.ShortNodeID)
	GetUnorderedProfiles() []profiles.Updatable
	FindProfile(nodeID insolar.ShortNodeID) profiles.Updatable
	GetLocalProfile() profiles.Updatable
	RemoveOthers()
}
