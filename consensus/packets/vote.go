/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
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
	TypeNodeJoinSupplementaryVote = VoteType(iota + 1)
	TypeStateFraudNodeSupplementaryVote
	TypeNodeListSupplementaryVote
	TypeMissingNodeSupplementaryVote
	TypeMissingNode
)

type ReferendumVote interface {
	Serializer
	Type() VoteType
}

type NodeJoinSupplementaryVote struct {
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

type MissingNodeSupplementaryVote struct {
	NodePulseProof NodePulseProof
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

func (v *NodeJoinSupplementaryVote) Type() VoteType {
	return TypeNodeJoinSupplementaryVote
}

// Deserialize implements interface method
func (v *NodeJoinSupplementaryVote) Deserialize(data io.Reader) error {
	return nil
}

// Serialize implements interface method
func (v *NodeJoinSupplementaryVote) Serialize() ([]byte, error) {
	return nil, nil
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
	result := allocateBuffer(packetMaxSize)

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
	err := v.NodePulseProof.Deserialize(data)
	if err != nil {
		return errors.Wrap(err, "[ MissingNodeSupplementaryVote.Deserialize ] Can't read NodePulseProof")
	}

	return nil
}

// Serialize implements interface method
func (v *MissingNodeSupplementaryVote) Serialize() ([]byte, error) {
	result := allocateBuffer(1024)

	nodePulseProofRaw, err := v.NodePulseProof.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeSupplementaryVote.Serialize ] Can't serialize NodePulseProof")
	}

	_, err = result.Write(nodePulseProofRaw)
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeSupplementaryVote.Serialize ] Can't append NodePulseProof")
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
		case TypeNodeJoinSupplementaryVote:
			refVote = &NodeJoinSupplementaryVote{}
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
