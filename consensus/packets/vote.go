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

package packets

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

type VoteType uint8

const (
	TypeNodeJoinSupplementaryVote = VoteType(iota + 1)
	TypeStateFraudNodeSupplementaryVote
	TypeNodeListSupplementaryVote
	TypeMissingNodeSupplementaryVote
)

type ReferendumVote interface {
	Serializer
	Type() VoteType
}

// todo: unused, remove
type NodeJoinSupplementaryVote struct {
	NodeListCount uint16
	NodeListHash  [32]byte
}

func (nlv *NodeJoinSupplementaryVote) Type() VoteType {
	return TypeNodeJoinSupplementaryVote
}

// Deserialize implements interface method
func (v *NodeJoinSupplementaryVote) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, v.NodeListCount)
	if err != nil {
		return errors.Wrap(err, "[ NodeListVote.Deserialize ] Can't read NodeListCount")
	}

	err = binary.Read(data, defaultByteOrder, v.NodeListHash)
	if err != nil {
		return errors.Wrap(err, "[ NodeListVote.Deserialize ] Can't read NodeListHash")
	}

	return nil
}

// Serialize implements interface method
func (nlv *NodeJoinSupplementaryVote) Serialize() ([]byte, error) {
	result := allocateBuffer(64)
	err := binary.Write(result, defaultByteOrder, nlv.NodeListCount)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeListVote.Serialize ] Can't write NodeListCount")
	}

	err = binary.Write(result, defaultByteOrder, nlv.NodeListHash)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeListVote.Serialize ] Can't write NodeListHash")
	}

	return result.Bytes(), nil
}
