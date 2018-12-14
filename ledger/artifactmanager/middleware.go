/*
 *    Copyright 2018 Insolar
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

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/pkg/errors"
)

type middleware struct {
	db             *storage.DB
	jetCoordinator core.JetCoordinator
}

type jetKey struct{}

func contextWithJet(ctx context.Context, jetID core.RecordID) context.Context {
	return context.WithValue(ctx, jetKey{}, jetID)
}

func jetFromContext(ctx context.Context) core.RecordID {
	val := ctx.Value(jetKey{})
	jet, ok := val.(core.RecordID)
	if !ok {
		panic("failed to extract jet from context")
	}

	return jet
}

func (m *middleware) checkJet(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		msg := parcel.Message()
		target := msg.DefaultTarget().Record()
		if target == nil {
			return nil, errors.New("unexpected message")
		}
		tree, err := m.db.GetJetTree(ctx, target.Pulse())
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch jet tree")
		}
		jetID := tree.Find(*target)
		if err != nil {
			return nil, err
		}

		isMine, err := m.jetCoordinator.AmI(ctx, core.DynamicRoleLightExecutor, target, target.Pulse())
		if err != nil {
			return nil, err
		}
		if !isMine {
			return &reply.JetMiss{JetID: *jetID}, nil
		}

		return handler(contextWithJet(ctx, *jetID), parcel)
	}
}

func (m *middleware) saveParcel(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		jetID := jetFromContext(ctx)
		pulse, err := m.db.GetLatestPulse(ctx)
		if err != nil {
			return nil, err
		}
		err = m.db.SetMessage(ctx, jetID, pulse.Pulse.PulseNumber, parcel)
		if err != nil {
			return nil, err
		}

		return handler(ctx, parcel)
	}
}
