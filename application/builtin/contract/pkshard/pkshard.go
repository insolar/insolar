// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pkshard

import (
	"github.com/pkg/errors"

	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

// PKShard - shard contract for public keys.
type PKShard struct {
	foundation.BaseContract
	Map foundation.StableMap
}

// New creates new member.
func New(members foundation.StableMap) (*PKShard, error) {
	return &PKShard{
		Map: members,
	}, nil
}

// GetRef gets ref by key.
// ins:immutable
func (s *PKShard) GetRef(key string) (string, error) {
	if ref, ok := s.Map[key]; !ok {
		return "", errors.New("failed to find reference by key")
	} else {
		return ref, nil
	}
}

// SetRef sets reference with public key as a key.
func (s *PKShard) SetRef(key string, ref string) error {
	if _, ok := s.Map[key]; ok {
		return errors.New("can't set reference because this key already exists")
	}
	s.Map[key] = ref
	return nil
}
