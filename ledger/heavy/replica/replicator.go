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
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/replica/intergrity"
	"github.com/insolar/insolar/ledger/heavy/sequence"
	"github.com/insolar/insolar/platformpolicy"
)

type Replicator struct {
	Sequencer     sequence.Sequencer          `inject:""`
	CryptoService insolar.CryptographyService `inject:""`
	Transport     Transport                   `inject:""`
	config        configuration.Configuration
	pulses        pulse.Accessor
	jetKeeper     JetKeeper
	target        Target
	cmps          component.Manager
}

func NewReplicator(
	cfg configuration.Configuration,
	pulses pulse.Accessor,
	jetKeeper JetKeeper,
) *Replicator {
	cmps := component.Manager{}
	return &Replicator{config: cfg, pulses: pulses, jetKeeper: jetKeeper, cmps: cmps}
}

func (r *Replicator) Init(ctx context.Context) error {
	replicaConfig := r.config.Ledger.Replica
	role := RoleBy(replicaConfig.Role)
	switch role {
	case Root:
		parent := NewParent()
		r.registerParent(parent)
	case Replica:
		remoteParent := NewRemoteParent(r.Transport, replicaConfig.ParentAddress)
		r.target = NewTarget(replicaConfig, remoteParent)
		r.registerTarget(r.target)
		validator := makeValidator(replicaConfig, r.CryptoService)
		r.cmps.Register(validator)

		parent := NewParent()
		r.registerParent(parent)
	case Observer:
		remoteParent := NewRemoteParent(r.Transport, replicaConfig.ParentAddress)
		r.target = NewTarget(replicaConfig, remoteParent)
		r.registerTarget(r.target)
		validator := makeValidator(replicaConfig, r.CryptoService)
		r.cmps.Register(validator)
	}
	provider := makeProvider(r.CryptoService)
	r.cmps.Inject(r.Sequencer, r.pulses, r.jetKeeper, provider)
	return nil
}

func (r *Replicator) Start(ctx context.Context) error {
	replicaConfig := r.config.Ledger.Replica
	role := RoleBy(replicaConfig.Role)
	switch role {
	case Replica, Observer:
		if cmp, ok := r.target.(component.Starter); ok {
			return cmp.Start(ctx)
		}
	}
	return nil
}

func (r *Replicator) registerParent(parent Parent) {
	r.Transport.Register("replica.Subscribe", func(data []byte) ([]byte, error) {
		ctx := context.Background()
		sub := Subscription{}
		err := insolar.Deserialize(data, &sub)
		if err != nil {
			return []byte{}, errors.Wrapf(err, "failed to deserialize subscription data")
		}
		target := NewRemoteTarget(r.Transport, sub.Target)
		err = parent.Subscribe(ctx, target, sub.At)
		if err != nil {
			return []byte{}, errors.Wrapf(err, "failed to call parent.Subscribe")
		}
		return []byte{}, nil
	})
	r.Transport.Register("replica.Pull", func(data []byte) ([]byte, error) {
		ctx := context.Background()
		pr := PullRequest{}
		err := insolar.Deserialize(data, &pr)
		if err != nil {
			return []byte{}, errors.Wrapf(err, "failed to deserialize pull request data")
		}
		packet, total, err := parent.Pull(ctx, pr.Page)
		if err != nil {
			return []byte{}, errors.Wrapf(err, "failed to call parent.Pull")
		}
		packet, err = insolar.Serialize(&PullReply{Data: packet, Total: total})
		if err != nil {
			return []byte{}, errors.Wrapf(err, "failed to serialize PullReply")
		}
		return packet, nil
	})
	r.cmps.Register(parent)
}

func (r *Replicator) registerTarget(target Target) {
	r.Transport.Register("replica.Notify", func(data []byte) ([]byte, error) {
		ctx := context.Background()
		n := Notification{}
		err := insolar.Deserialize(data, &n)
		if err != nil {
			return []byte{}, errors.Wrapf(err, "failed to deserialize notification data")
		}
		err = target.Notify(ctx, n.Pulse)
		if err != nil {
			return []byte{}, errors.Wrapf(err, "failed to call target.Notify")
		}
		return []byte{}, nil
	})
	r.cmps.Register(target)
}

func makeProvider(cryptoService insolar.CryptographyService) intergrity.Provider {
	return intergrity.NewProvider(cryptoService)
}

func makeValidator(cfg configuration.Replica, cryptoService insolar.CryptographyService) intergrity.Validator {
	logger := inslogger.FromContext(context.Background())
	kp := platformpolicy.NewKeyProcessor()
	pubKey, err := kp.ImportPublicKeyPEM([]byte(cfg.ParentPubKey))
	if err != nil {
		logger.Error(errors.Wrap(err, "failed to import a public key from PEM"))
		return nil
	}
	return intergrity.NewValidator(cryptoService, pubKey)
}
