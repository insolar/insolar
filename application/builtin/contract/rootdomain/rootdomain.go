// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
	trimmedPublicKey, err := foundation.ExtractCanonicalPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("extracting canonical pk failed, current value %v", publicKey)
	}
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
	trimmedPublicKey, err := foundation.ExtractCanonicalPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("extracting canonical pk failed, current value %v", publicKey)
	}
	shardIndex := foundation.GetShardIndex(trimmedPublicKey, len(rd.PublicKeyShards))
	if shardIndex >= len(rd.PublicKeyShards) {
		return fmt.Errorf("incorrect public key shard index")
	}
	pks := pkshard.GetObject(rd.PublicKeyShards[shardIndex])
	err = pks.SetRef(trimmedPublicKey, memberRef.String())
	if err != nil {
		return errors.Wrap(err, "failed to set reference in public key shard")
	}
	return nil
}
