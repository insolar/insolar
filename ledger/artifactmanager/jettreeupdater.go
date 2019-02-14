/*
 *    Copyright 2019 Insolar Technologies
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

package artifactmanager

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

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
		return jetID, nil
	}

	// Not actual in our tree, asking neighbors for jet.
	span.Annotate(nil, "tree in DB is not actual")
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
		// Tree was updated in another thread, rechecking.
		span.Annotate(nil, "somebody else updated actuality")
		return jtu.fetchJet(ctx, target, pulse)
	}
	defer func() {
		mu.done = true
		mu.Unlock()

		jtu.seqMutex.Lock()
		delete(jtu.sequencer, key)
		jtu.seqMutex.Unlock()
	}()

	ch := jtu.fetchActualJetFromOtherNodes(ctx, target, pulse)
	res := <-ch

	resJet, err := res.jet, res.err
	if err != nil {
		return nil, err
	}

	err = jtu.JetStorage.UpdateJetTree(ctx, pulse, true, *resJet)
	if err != nil {
		inslogger.FromContext(ctx).Error(
			errors.Wrapf(err, "failed actualize jet %s", resJet.DebugString()),
		)
	}

	return resJet, nil
}

type result struct {
	jet *core.RecordID
	err error
}

func (jtu *jetTreeUpdater) fetchActualJetFromOtherNodes(
	ctx context.Context, target core.RecordID, pulse core.PulseNumber,
) chan result {
	ctx, span := instracer.StartSpan(ctx, "jet_tree_updater.fetch_jet_from_other_nodes")
	defer span.End()

	ch := make(chan result, 1)

	go func() {
		nodes, err := jtu.otherNodesForPulse(ctx, pulse)
		if err != nil {
			ch <- result{nil, err}
			return
		}

		num := len(nodes)

		wg := sync.WaitGroup{}
		wg.Add(num)

		found := uint32(0)

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

				if !r.Actual {
					return
				}

				if atomic.CompareAndSwapUint32(&found, 0, 1) {
					jet := r.ID
					ch <- result{&jet, nil }
					close(ch)
				}

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
			if _, ok := seen[r.ID]; ok {
				continue
			}

			seen[r.ID] = struct{}{}
			res = append(res, &r.ID)
		}

		if len(res) == 0 {
			inslogger.FromContext(ctx).WithFields(map[string]interface{}{
				"pulse":  pulse,
				"object": target.DebugString(),
			}).Error("all lights for pulse have no actual jet for object")
			ch <- result{nil, errors.New("impossible situation") }
			close(ch)
		} else if len(res) > 1 {
			inslogger.FromContext(ctx).WithFields(map[string]interface{}{
				"pulse":  pulse,
				"object": target.DebugString(),
			}).Error("lights said different actual jet for object")
		}
	}()

	return ch
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
