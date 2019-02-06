package artifactmanager

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/storage"
)

type jetTreeUpdater struct {
	ActiveNodesStorage storage.NodeStorage
	JetStorage         storage.JetStorage
	MessageBus         core.MessageBus
	JetCoordinator     core.JetCoordinator

	seqMutex  sync.Mutex
	sequencer map[string]*struct {
		sync.Mutex
		done bool
	}
}

func newJetTreeUpdater(
	ans storage.NodeStorage,
	js storage.JetStorage, mb core.MessageBus, jc core.JetCoordinator,
) *jetTreeUpdater {
	return &jetTreeUpdater{
		ActiveNodesStorage: ans,
		JetStorage:         js,
		MessageBus:         mb,
		JetCoordinator:     jc,
		sequencer: map[string]*struct {
			sync.Mutex
			done bool
		}{},
	}
}

func (jtu *jetTreeUpdater) fetchJet(
	ctx context.Context, target core.RecordID, pulse core.PulseNumber,
) (*core.RecordID, error) {
	ctx, span := instracer.StartSpan(ctx, "jet_tree_updater.fetch_jet")
	defer span.End()

	// Look in the local tree. Return if the actual jet found.
	tree, err := jtu.JetStorage.GetJetTree(ctx, pulse)
	if err != nil {
		return nil, err
	}

	jetID, actual := tree.Find(target)
	if actual {
		inslogger.FromContext(ctx).Debugf(
			"we believe object %s is in JET %s", target.String(), jetID.DebugString(),
		)
		return jetID, nil
	}

	span.Annotate(nil, "tree in DB is not actual")
	inslogger.FromContext(ctx).Debugf(
		"jet %s is not actual in our tree, asking neighbors for jet of object %s",
		jetID.DebugString(), target.String(),
	)

	key := fmt.Sprintf("%d:%s", pulse, jetID.String())

	jtu.seqMutex.Lock()
	if _, ok := jtu.sequencer[key]; !ok {
		jtu.sequencer[key] = &struct {
			sync.Mutex
			done bool
		}{}
	}
	mu := jtu.sequencer[key]
	jtu.seqMutex.Unlock()

	span.Annotate(nil, "got sequencer entry")

	mu.Lock()
	if mu.done {
		mu.Unlock()
		span.Annotate(nil, "somebody else updated actuality")
		inslogger.FromContext(ctx).Debugf(
			"somebody else updated actuality of jet %s, rechecking our DB",
			jetID.DebugString(),
		)
		return jtu.fetchJet(ctx, target, pulse)
	}
	defer func() {
		inslogger.FromContext(ctx).Debugf("done fetching jet, cleaning")

		mu.done = true
		mu.Unlock()

		jtu.seqMutex.Lock()
		inslogger.FromContext(ctx).Debugf("deleting sequencer for jet %s", jetID.DebugString())
		delete(jtu.sequencer, key)
		jtu.seqMutex.Unlock()
	}()

	resJet, err := jtu.fetchActualJetFromOtherNodes(ctx, target, pulse)
	if err != nil {
		return nil, err
	}

	err = jtu.JetStorage.UpdateJetTree(ctx, pulse, true, *resJet)
	if err != nil {
		inslogger.FromContext(ctx).Error(
			errors.Wrapf(err, "couldn't actualize jet %s", resJet.DebugString()),
		)
	}

	return resJet, nil
}

func (jtu *jetTreeUpdater) fetchActualJetFromOtherNodes(
	ctx context.Context, target core.RecordID, pulse core.PulseNumber,
) (*core.RecordID, error) {
	ctx, span := instracer.StartSpan(ctx, "jet_tree_updater.fetch_jet_from_other_nodes")
	defer span.End()

	nodes, err := jtu.otherNodesForPulse(ctx, pulse)
	if err != nil {
		return nil, err
	}

	num := len(nodes)

	wg := sync.WaitGroup{}
	wg.Add(num)

	replies := make([]*reply.Jet, num)
	for i, node := range nodes {
		go func(i int, node core.Node) {
			ctx, span := instracer.StartSpan(ctx, "jet_tree_updater.one_node_get_jet")
			defer span.End()

			defer wg.Done()

			nodeID := node.ID()
			rep, err := jtu.MessageBus.Send(
				ctx,
				&message.GetJet{Object: target, Pulse: pulse},
				&core.MessageSendOptions{Receiver: &nodeID},
			)
			if err != nil {
				inslogger.FromContext(ctx).Error(
					errors.Wrap(err, "couldn't get jet"),
				)
				return
			}

			r, ok := rep.(*reply.Jet)
			if !ok {
				inslogger.FromContext(ctx).Errorf("middleware.fetchActualJetFromOtherNodes: unexpected reply: %#v\n", rep)
				return
			}

			inslogger.FromContext(ctx).Debugf(
				"Got jet %s from %s node, actual is %s",
				r.ID.DebugString(), node.ID().String(), r.Actual,
			)
			replies[i] = r
		}(i, node)
	}
	wg.Wait()

	seen := make(map[core.RecordID]struct{})
	res := make([]*core.RecordID, 0)
	for _, r := range replies {
		if r == nil {
			continue
		}
		if !r.Actual {
			continue
		}
		if _, ok := seen[r.ID]; ok {
			continue
		}

		seen[r.ID] = struct{}{}
		res = append(res, &r.ID)
	}

	if len(res) == 1 {
		inslogger.FromContext(ctx).Debugf(
			"got jet %s as actual for object %s on pulse %d",
			res[0].DebugString(), target.String(), pulse,
		)
		return res[0], nil
	} else if len(res) == 0 {
		inslogger.FromContext(ctx).Errorf(
			"all lights for pulse %d have no actual jet for object %s",
			pulse, target.String(),
		)
		return nil, errors.New("impossible situation")
	} else {
		inslogger.FromContext(ctx).Errorf(
			"lights said different actual jet for object %s",
			target.String(),
		)
		return nil, errors.New("nodes returned more than one unique jet")
	}
}

func (jtu *jetTreeUpdater) otherNodesForPulse(
	ctx context.Context, pulse core.PulseNumber,
) ([]core.Node, error) {
	ctx, span := instracer.StartSpan(ctx, "jet_tree_updater.other_nodes_for_pulse")
	defer span.End()

	nodes, err := jtu.ActiveNodesStorage.GetActiveNodesByRole(pulse, core.StaticRoleLightMaterial)
	if err != nil {
		return nil, err
	}

	me := jtu.JetCoordinator.Me()
	for i := range nodes {
		if nodes[i].ID() == me {
			nodes = append(nodes[:i], nodes[i+1:]...)
			break
		}
	}

	num := len(nodes)
	if num == 0 {
		inslogger.FromContext(ctx).Error(
			"This shouldn't happen. We're solo active light material",
		)

		return nil, errors.New("impossible situation")
	}

	return nodes, nil
}
