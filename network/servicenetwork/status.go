package servicenetwork

import (
	"context"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/version"
)

var startTime time.Time

func (n *ServiceNetwork) GetNetworkStatus() insolar.StatusReply {
	var reply insolar.StatusReply
	reply.NetworkState = n.Gatewayer.Gateway().GetState()

	np, err := n.PulseAccessor.GetLatestPulse(context.Background())
	if err != nil {
		np = *insolar.GenesisPulse
	}
	reply.Pulse = np

	activeNodes := n.NodeKeeper.GetAccessor(np.PulseNumber).GetActiveNodes()
	workingNodes := n.NodeKeeper.GetAccessor(np.PulseNumber).GetWorkingNodes()

	reply.ActiveListSize = len(activeNodes)
	reply.WorkingListSize = len(workingNodes)

	reply.Nodes = activeNodes
	reply.Origin = n.NodeKeeper.GetOrigin()

	reply.Version = version.Version

	reply.Timestamp = time.Now()
	reply.StartTime = startTime

	return reply
}

func init() {
	startTime = time.Now()
}
