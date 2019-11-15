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
