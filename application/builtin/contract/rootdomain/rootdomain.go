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

package rootdomain

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/application/builtin/proxy/pkshard"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

// RootDomain is smart contract representing entrance point to system.
type RootDomain struct {
	foundation.BaseContract
	PublicKeyShards []insolar.Reference
}

// GetMemberByPublicKey gets member reference by public key.
// ins:immutable
func (rd *RootDomain) GetMemberByPublicKey(publicKey string) (*insolar.Reference, error) {
	trimmedPublicKey := foundation.TrimPublicKey(publicKey)
	i := foundation.GetShardIndex(trimmedPublicKey, len(rd.PublicKeyShards))
	if i >= len(rd.PublicKeyShards) {
		return nil, fmt.Errorf("incorrect shard index")
	}
	s := pkshard.GetObject(rd.PublicKeyShards[i])
	refStr, err := s.GetRef(trimmedPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get reference in shard")
	}
	ref, err := insolar.NewObjectReferenceFromString(refStr)
	if err != nil {
		return nil, errors.Wrap(err, "bad member reference for this public key")
	}

	return ref, nil
}

// AddNewMemberToPublicKeyMap adds new member to PublicKeyMap.
// ins:immutable
func (rd *RootDomain) AddNewMemberToPublicKeyMap(publicKey string, memberRef insolar.Reference) error {
	trimmedPublicKey := foundation.TrimPublicKey(publicKey)
	shardIndex := foundation.GetShardIndex(trimmedPublicKey, len(rd.PublicKeyShards))
	if shardIndex >= len(rd.PublicKeyShards) {
		return fmt.Errorf("incorrect public key shard index")
	}
	pks := pkshard.GetObject(rd.PublicKeyShards[shardIndex])
	err := pks.SetRef(trimmedPublicKey, memberRef.String())
	if err != nil {
		return errors.Wrap(err, "failed to set reference in public key shard")
	}
	return nil
}
