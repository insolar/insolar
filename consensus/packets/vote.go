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
	"bytes"
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

const NodeListHashLength = 32

type VoteType uint8

const (
	TypeStateFraudNodeSupplementaryVote = VoteType(iota + 1)
	TypeNodeListSupplementaryVote
	TypeMissingNodeSupplementaryVote
	TypeMissingNode
	TypeMissingNodeClaim
)

type ReferendumVote interface {
	Serializer
	Type() VoteType
}

type StateFraudNodeSupplementaryVote struct {
	Node1PulseProof NodePulseProof
	Node2PulseProof NodePulseProof
	PulseData       PulseData // optional
}

type NodeListSupplementaryVote struct {
	NodeListCount uint16
	NodeListHash  [32]byte
}

type MissingNodeClaim struct {
	NodeIndex uint16
	claimSize uint16

	Claim ReferendumClaim
}

func (mn *MissingNodeClaim) Type() VoteType {
	return TypeMissingNodeClaim
}

type MissingNodeSupplementaryVote struct {
	NodeIndex uint16

	NodePulseProof       NodePulseProof
	GlobuleHashSignature GlobuleHashSignature
	// TODO: make it signed
	NodeClaimUnsigned NodeJoinClaim
}

type MissingNode struct {
	NodeIndex uint16
}

func (mn *MissingNode) Type() VoteType {
	return TypeMissingNode
}

func (mn *MissingNode) Serialize() ([]byte, error) {
	var result bytes.Buffer
	err := binary.Write(&result, defaultByteOrder, mn.NodeIndex)
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNode.Serialize ] failed to write ti a buffer")
	}
	return result.Bytes(), nil
}

func (mn *MissingNode) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &mn.NodeIndex)
	if err != nil {
		return errors.Wrap(err, "[ MissingNode.Deserialize ] failed to read a node index")
	}
	return nil
}

func (mn *MissingNodeClaim) Serialize() ([]byte, error) {
	var result bytes.Buffer
	err := binary.Write(&result, defaultByteOrder, mn.NodeIndex)
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeClaim.Serialize ] Can't write NodeIndex")
	}
	serializedClaim, err := mn.Claim.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeClaim.Serialize ] Can't serialize claim")
	}
	err = binary.Write(&result, defaultByteOrder, uint16(len(serializedClaim)))
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeClaim.Serialize ] Can't write claimSize")
	}
	err = binary.Write(&result, defaultByteOrder, serializedClaim)
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeClaim.Serialize ] Can't write Claim")
	}
	return result.Bytes(), nil
}

func (mn *MissingNodeClaim) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &mn.NodeIndex)
	if err != nil {
		return errors.Wrap(err, "[ MissingNodeClaim.Deserialize ] Can't read NodeIndex")
	}
	err = binary.Read(data, defaultByteOrder, &mn.claimSize)
	if err != nil {
		return errors.Wrap(err, "[ MissingNodeClaim.Deserialize ] Can't read claimSize")
	}
	claimData := make([]byte, mn.claimSize)
	err = binary.Read(data, defaultByteOrder, claimData[:])
	if err != nil {
		return errors.Wrap(err, "[ MissingNodeClaim.Deserialize ] Can't read claim data")
	}
	claims, err := parseReferendumClaim(claimData[:])
	if err != nil {
		return errors.Wrap(err, "[ MissingNodeClaim.Deserialize ] Can't parse claim from claim data")
	}
	mn.Claim = claims[0]
	return nil
}

func (v *NodeListSupplementaryVote) Type() VoteType {
	return TypeNodeListSupplementaryVote
}

// Deserialize implements interface method
func (v *NodeListSupplementaryVote) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &v.NodeListCount)
	if err != nil {
		return errors.Wrap(err, "[ NodeListSupplementaryVote.Deserialize ] Can't read NodeListCount")
	}

	err = binary.Read(data, defaultByteOrder, &v.NodeListHash)
	if err != nil {
		return errors.Wrap(err, "[ NodeListSupplementaryVote.Deserialize ] Can't read NodeListHash")
	}

	return nil
}

// Serialize implements interface method
func (v *NodeListSupplementaryVote) Serialize() ([]byte, error) {
	result := allocateBuffer(34)
	err := binary.Write(result, defaultByteOrder, v.NodeListCount)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeListSupplementaryVote.Serialize ] Can't write NodeListCount")
	}

	err = binary.Write(result, defaultByteOrder, v.NodeListHash)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeListSupplementaryVote.Serialize ] Can't write NodeListHash")
	}

	return result.Bytes(), nil
}

func (v *StateFraudNodeSupplementaryVote) Type() VoteType {
	return TypeStateFraudNodeSupplementaryVote
}

// Deserialize implements interface method
func (v *StateFraudNodeSupplementaryVote) Deserialize(data io.Reader) error {
	err := v.Node1PulseProof.Deserialize(data)
	if err != nil {
		return errors.Wrap(err, "[ StateFraudNodeSupplementaryVote.Deserialize ] Can't read Node1PulseProof")
	}

	err = v.Node2PulseProof.Deserialize(data)
	if err != nil {
		return errors.Wrap(err, "[ StateFraudNodeSupplementaryVote.Deserialize ] Can't read Node2PulseProof")
	}

	err = v.PulseData.Deserialize(data)
	if err != nil {
		return errors.Wrap(err, "[ StateFraudNodeSupplementaryVote.Deserialize ] Can't read PulseData")
	}

	return nil
}

// Serialize implements interface method
func (v *StateFraudNodeSupplementaryVote) Serialize() ([]byte, error) {
	result := allocateBuffer(1024)

	node1PulseProofRaw, err := v.Node1PulseProof.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ StateFraudNodeSupplementaryVote.Serialize ] Can't serialize Node1PulseProof")
	}

	_, err = result.Write(node1PulseProofRaw)
	if err != nil {
		return nil, errors.Wrap(err, "[ StateFraudNodeSupplementaryVote.Serialize ] Can't append Node1PulseProof")
	}

	node2PulseProofRaw, err := v.Node2PulseProof.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ StateFraudNodeSupplementaryVote.Serialize ] Can't serialize Node2PulseProof")
	}

	_, err = result.Write(node2PulseProofRaw)
	if err != nil {
		return nil, errors.Wrap(err, "[ StateFraudNodeSupplementaryVote.Serialize ] Can't append Node2PulseProof")
	}

	// serializing of  PulseData
	pulseDataRaw, err := v.PulseData.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ StateFraudNodeSupplementaryVote.Serialize ] Can't serialize PulseData")
	}
	_, err = result.Write(pulseDataRaw)
	if err != nil {
		return nil, errors.Wrap(err, "[ StateFraudNodeSupplementaryVote.Serialize ] Can't append PulseData")
	}

	return result.Bytes(), nil
}

func (v *MissingNodeSupplementaryVote) Type() VoteType {
	return TypeMissingNodeSupplementaryVote
}

// Deserialize implements interface method
func (v *MissingNodeSupplementaryVote) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &v.NodeIndex)
	if err != nil {
		return errors.Wrap(err, "[ MissingNodeSupplementaryVote.Deserialize ] Can't read NodeIndex")
	}
	err = v.NodePulseProof.Deserialize(data)
	if err != nil {
		return errors.Wrap(err, "[ MissingNodeSupplementaryVote.Deserialize ] Can't read NodePulseProof")
	}
	err = binary.Read(data, defaultByteOrder, &v.GlobuleHashSignature)
	if err != nil {
		return errors.Wrap(err, "[ MissingNodeSupplementaryVote.Deserialize ] Can't read GlobuleHashSignature")
	}
	err = v.NodeClaimUnsigned.deserializeRaw(data)
	if err != nil {
		return errors.Wrap(err, "[ MissingNodeSupplementaryVote.Deserialize ] Can't read NodeClaimUnsigned")
	}

	return nil
}

// Serialize implements interface method
func (v *MissingNodeSupplementaryVote) Serialize() ([]byte, error) {
	result := allocateBuffer(1024)

	err := binary.Write(result, defaultByteOrder, v.NodeIndex)
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeSupplementaryVote.Serialize ] Can't write NodeIndex")
	}

	nodePulseProofRaw, err := v.NodePulseProof.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeSupplementaryVote.Serialize ] Can't serialize NodePulseProof")
	}
	_, err = result.Write(nodePulseProofRaw)
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeSupplementaryVote.Serialize ] Can't append NodePulseProof")
	}

	_, err = result.Write(v.GlobuleHashSignature[:])
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeSupplementaryVote.Serialize ] Can't write GlobuleHashSignature")
	}

	joinClaim, err := v.NodeClaimUnsigned.SerializeRaw()
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeSupplementaryVote.Serialize ] Can't serialize join claim")
	}
	_, err = result.Write(joinClaim)
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeSupplementaryVote.Serialize ] Can't write join claim")
	}

	return result.Bytes(), nil
}

func parseReferendumVotes(data []byte) ([]ReferendumVote, error) {
	votesSize := len(data)
	votesReader := bytes.NewReader(data)
	result := make([]ReferendumVote, 0)

	// get claim header
	for votesSize > 0 {
		startSize := votesReader.Len()
		var voteHeader uint16
		err := binary.Read(votesReader, defaultByteOrder, &voteHeader)
		if err != nil {
			return nil, errors.Wrap(err, "[ PacketHeader.parseReferendumVotes ] Can't read voteHeader")
		}

		voteType := VoteType(extractTypeFromHeader(voteHeader))
		// TODO: Do we need voteLength?
		// voteLength := extractVoteLengthFromHeader(voteHeader)
		var refVote ReferendumVote

		switch voteType {
		case TypeStateFraudNodeSupplementaryVote:
			refVote = &StateFraudNodeSupplementaryVote{}
		case TypeNodeListSupplementaryVote:
			refVote = &NodeListSupplementaryVote{}
		case TypeMissingNodeSupplementaryVote:
			refVote = &MissingNodeSupplementaryVote{}
		case TypeMissingNode:
			refVote = &MissingNode{}
		default:
			return nil, errors.Wrap(err, "[ PacketHeader.parseReferendumVotes ] Unsupported vote type.")
		}
		err = refVote.Deserialize(votesReader)
		if err != nil {
			return nil, errors.Wrap(err, "[ PacketHeader.parseReferendumVotes ] Can't deserialize vote")
		}
		result = append(result, refVote)

		votesSize -= startSize - votesReader.Len()
	}

	if votesSize != 0 {
		return nil, errors.New("[ PacketHeader.parseReferendumVotes ] Problem with vote struct")
	}

	return result, nil
}
