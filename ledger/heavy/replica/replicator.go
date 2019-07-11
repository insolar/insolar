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

package replica

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
)

type Replicator struct {
	DB             store.DB                    `inject:""`
	CryptoService  insolar.CryptographyService `inject:""`
	ServiceNetwork insolar.Network             `inject:""`
	Transport      Transport                   `inject:""`
	config         configuration.Configuration
	jetKeeper      JetKeeper
	target         Target
}

func NewReplicator(
	cfg configuration.Configuration,
	jetKeeper JetKeeper,
) *Replicator {
	return &Replicator{config: cfg, jetKeeper: jetKeeper}
}

func (r *Replicator) Init(ctx context.Context) error {
	// TODO: inject sequencer
	replicaConfig := r.config.Ledger.Replica
	role := RoleBy(replicaConfig.Role)
	switch role {
	case Root:
		parent := NewParent(r.DB, r.jetKeeper, r.CryptoService)
		registerParent(parent, r.Transport)
	case Replica:
		remote := NewRemoteParent(r.Transport, replicaConfig.ParentAddress)
		r.target = NewTarget(r.DB, replicaConfig, remote, r.CryptoService)
		registerTarget(r.target, r.Transport)

		parent := NewParent(r.DB, r.jetKeeper, r.CryptoService)
		registerParent(parent, r.Transport)
	case Observer:
		parent := NewRemoteParent(r.Transport, replicaConfig.ParentAddress)
		r.target = NewTarget(r.DB, replicaConfig, parent, r.CryptoService)
		registerTarget(r.target, r.Transport)
	}
	return nil
}

func (r *Replicator) Start(ctx context.Context) error {
	replicaConfig := r.config.Ledger.Replica
	role := RoleBy(replicaConfig.Role)
	inslogger.FromContext(ctx).Warnf("Starting replicator config", replicaConfig)
	switch role {
	case Replica, Observer:
		if cmp, ok := r.target.(component.Starter); ok {
			return cmp.Start(ctx)
		}
	}
	return nil
}

func registerParent(parent Parent, transport Transport) {
	transport.Register("replica.Subscribe", func(data []byte) ([]byte, error) {
		sub := Subscription{}
		err := insolar.Deserialize(data, &sub)
		if err != nil {
			return []byte{}, errors.Wrapf(err, "failed to deserialize subscription data")
		}
		target := NewRemoteTarget(transport, sub.Target)
		err = parent.Subscribe(target, sub.At)
		if err != nil {
			return []byte{}, errors.Wrapf(err, "failed to call parent.Subscribe")
		}
		return []byte{}, nil
	})
	transport.Register("replica.Pull", func(data []byte) ([]byte, error) {
		pr := PullRequest{}
		err := insolar.Deserialize(data, &pr)
		if err != nil {
			return []byte{}, errors.Wrapf(err, "failed to deserialize pull request data")
		}
		packet, err := parent.Pull(pr.Scope, pr.From, pr.Limit)
		if err != nil {
			return []byte{}, errors.Wrapf(err, "failed to call parent.Pull")
		}
		return packet, nil
	})
}

func registerTarget(target Target, transport Transport) {
	transport.Register("replica.Notify", func(data []byte) ([]byte, error) {
		err := target.Notify()
		if err != nil {
			return []byte{}, errors.Wrapf(err, "failed to call target.Notify")
		}
		return []byte{}, nil
	})
}
