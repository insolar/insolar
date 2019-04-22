//
// Copyright 2019 Insolar Technologies GmbH
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
//

package jet

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

// Fetcher can be used to get actual jets. It involves fetching jet from other nodes via network and updating local
// jet tree.
type Fetcher interface {
	Fetch(ctx context.Context, target insolar.ID, pulse insolar.PulseNumber) (*insolar.ID, error)
	Release(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber)
}

// Used to queue fetching routines.
type seqEntry struct {
	ch   chan struct{}
	once sync.Once
}

// Used as an id for fetching routines. Each jet is updated individually and independently, but routines with the same
// jets are queued.
type seqKey struct {
	pulse insolar.PulseNumber
	jet   insolar.JetID
}

// Used to pass fetching result over channels.
type fetchResult struct {
	jet *insolar.ID
	err error
}

type fetcher struct {
	Nodes       node.Accessor
	JetStorage  Storage
	MessageBus  insolar.MessageBus
	coordinator Coordinator

	seqMutex  sync.Mutex
	sequencer map[seqKey]*seqEntry
}

// NewFetcher creates new fetcher instance.
func NewFetcher(
	ans node.Accessor,
	js Storage,
	mb insolar.MessageBus,
	jc Coordinator,
) Fetcher {
	return &fetcher{
		Nodes:       ans,
		JetStorage:  js,
		MessageBus:  mb,
		coordinator: jc,
		sequencer:   map[seqKey]*seqEntry{},
	}
}

// Fetch coordinates jet fetching routines. It is safe to call concurrently on the same instance.
//
// Multiple routines enter the fetching section and grouped by jet id and pulse. All groups are executed independently.
// Routines within one group executed sequentially. Each routine goes through steps:
// 1. Look in the local tree. If actual jet is found - return.
// 2. Enter the queue.
// 3. Fetch actual jet over network.
// 4. Update local tree.
// 5. Exit the queue.
func (tu *fetcher) Fetch(
	ctx context.Context, target insolar.ID, pulse insolar.PulseNumber,
) (*insolar.ID, error) {
	ctx, span := instracer.StartSpan(ctx, "jet_tree_updater.fetch_jet")
	defer span.End()

	// Look in the local tree. Return if the actual jet found.
	jetID, actual := tu.JetStorage.ForID(ctx, pulse, target)
	if actual {
		return (*insolar.ID)(&jetID), nil
	}

	// Not actual in our tree, asking neighbors for jet.
	span.Annotate(nil, "tree in DB is not actual")
	key := seqKey{pulse, jetID}

	// Indicates that this routine is the first in the queue and should do the fetching.
	// Other routines wait in the queue.
	executing := false

	tu.seqMutex.Lock()
	if _, ok := tu.sequencer[key]; !ok {
		// Key is not found in the queue. We are the first.
		tu.sequencer[key] = &seqEntry{ch: make(chan struct{})}
		executing = true
	}
	entry := tu.sequencer[key]
	tu.seqMutex.Unlock()

	span.Annotate(nil, "got sequencer entry")

	if !executing {
		// We are not the first, waiting in the queue.
		<-entry.ch

		// Tree was updated in another thread, rechecking.
		span.Annotate(nil, "somebody else updated actuality")
		return tu.Fetch(ctx, target, pulse)
	}

	defer func() {
		// Prevents closing of a closed channel.
		entry.once.Do(func() {
			close(entry.ch)
		})

		// Exiting the queue.
		tu.seqMutex.Lock()
		delete(tu.sequencer, key)
		tu.seqMutex.Unlock()
	}()

	// Fetching jet via network.
	resJet, err := tu.fetch(ctx, target, pulse)
	if err != nil {
		return nil, err
	}

	// Updating local tree.
	tu.JetStorage.Update(ctx, pulse, true, insolar.JetID(*resJet))

	return resJet, nil
}

// Release unlocks all the queses on the branch for provided jet. I.e. all the jets that are higher in the tree on the
// current branch get released and "fall through" until they hit provided jet or branch out.
func (tu *fetcher) Release(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) {
	tu.seqMutex.Lock()
	defer tu.seqMutex.Unlock()

	depth := jetID.Depth()
	for {
		key := seqKey{pulse, jetID}
		if v, ok := tu.sequencer[key]; ok {
			// Unlocking jets queue.
			v.once.Do(func() {
				close(v.ch)
			})

			delete(tu.sequencer, key)
		}

		if depth == 0 {
			break
		}
		// Iterating over jet parents (going up the tree).
		jetID = Parent(jetID)
		depth--
	}
}

// Fetching jet over network.
func (tu *fetcher) fetch(
	ctx context.Context, target insolar.ID, pulse insolar.PulseNumber,
) (*insolar.ID, error) {
	ctx, span := instracer.StartSpan(ctx, "jet_tree_updater.fetch_jet_from_other_nodes")
	defer span.End()

	// Fetching result will be written here.
	ch := make(chan fetchResult, 1)

	go func() {
		// Other nodes that might have the actual jet.
		nodes, err := tu.nodesForPulse(ctx, pulse)
		if err != nil {
			ch <- fetchResult{nil, err}
			return
		}

		num := len(nodes)

		wg := sync.WaitGroup{}
		wg.Add(num)

		once := sync.Once{}

		replies := make([]*reply.Jet, num)
		for i, node := range nodes {
			// Asking all the nodes concurrently.
			go func(i int, node insolar.Node) {
				ctx, span := instracer.StartSpan(ctx, "jet_tree_updater.one_node_get_jet")
				defer span.End()

				defer wg.Done()

				nodeID := node.ID
				// Asking the node for jet.
				rep, err := tu.MessageBus.Send(
					ctx,
					&message.GetJet{Object: target, Pulse: pulse},
					&insolar.MessageSendOptions{Receiver: &nodeID},
				)
				if err != nil {
					inslogger.FromContext(ctx).Error(
						errors.Wrap(err, "couldn't get jet"),
					)
					return
				}

				r, ok := rep.(*reply.Jet)
				if !ok {
					inslogger.FromContext(ctx).Errorf("middleware.fetch: unexpected reply: %#v\n", rep)
					return
				}

				if !r.Actual {
					return
				}

				// Only one routine writes the result. The rest will still collect their result for future comparison.
				// We compare all the results to find potential problems.
				once.Do(func() {
					jetID := r.ID
					ch <- fetchResult{&jetID, nil}
					close(ch)
				})

				replies[i] = r
			}(i, node)
		}
		wg.Wait()

		// Collect non-nil replies (only actual).
		res := make(map[insolar.ID]struct{})
		for _, r := range replies {
			if r == nil {
				continue
			}

			res[r.ID] = struct{}{}
		}

		if len(res) == 0 {
			// No one knows the actual jet.
			inslogger.FromContext(ctx).WithFields(map[string]interface{}{
				"pulse":  pulse,
				"object": target.DebugString(),
			}).Error("all lights for pulse have no actual jet for object")
			ch <- fetchResult{nil, errors.New("all lights for pulse have no actual jet for object")}
			close(ch)
		} else if len(res) > 1 {
			// We have multiple different opinions on the actual jet.
			inslogger.FromContext(ctx).WithFields(map[string]interface{}{
				"pulse":  pulse,
				"object": target.DebugString(),
			}).Error("lights said different actual jet for object")
		}
	}()

	res := <-ch
	return res.jet, res.err
}

// All light materials except ourselves.
func (tu *fetcher) nodesForPulse(ctx context.Context, pulse insolar.PulseNumber) ([]insolar.Node, error) {
	ctx, span := instracer.StartSpan(ctx, "jet_tree_updater.other_nodes_for_pulse")
	defer span.End()

	res, err := tu.Nodes.InRole(pulse, insolar.StaticRoleLightMaterial)
	if err != nil {
		return nil, err
	}

	me := tu.coordinator.Me()
	for i := range res {
		if res[i].ID == me {
			res = append(res[:i], res[i+1:]...)
			break
		}
	}

	num := len(res)
	if num == 0 {
		inslogger.FromContext(ctx).Error(
			"This shouldn't happen. We're solo active light material",
		)

		return nil, errors.New("no other light to fetch jet tree data from")
	}

	return res, nil
}
