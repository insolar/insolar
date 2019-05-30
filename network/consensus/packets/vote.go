//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package packets

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

type NodeListHash [32]byte

type VoteType uint8

const (
	TypeMissingNodeReqVote    = VoteType(1)
	TypeMissingNodeRespVote   = VoteType(2)
	TypeMissingNodeClaimsVote = VoteType(3)
)

type ReferendumVote interface {
	Serializer
	Type() VoteType
	Clone() ReferendumVote
}

type MissingNodeClaimsVote struct {
	NodeIndex uint16
	claimSize uint16

	Claim ReferendumClaim
}

func (mn *MissingNodeClaimsVote) Clone() ReferendumVote {
	clone := *mn
	clone.Claim = mn.Claim.Clone()
	return &clone
}

func (mn *MissingNodeClaimsVote) Type() VoteType {
	return TypeMissingNodeClaimsVote
}

type MissingNodeRespVote struct {
	NodeIndex uint16

	NodePulseProof NodePulseProof
	// TODO: make it signed
	NodeClaimUnsigned NodeJoinClaim
}

func (v *MissingNodeRespVote) Clone() ReferendumVote {
	clone := *v
	return &clone
}

type MissingNodeReqVote struct {
	NodeIndex uint16
}

func (mn *MissingNodeReqVote) Clone() ReferendumVote {
	clone := *mn
	return &clone
}

func (mn *MissingNodeReqVote) Type() VoteType {
	return TypeMissingNodeReqVote
}

func (mn *MissingNodeReqVote) Serialize() ([]byte, error) {
	var result bytes.Buffer
	err := binary.Write(&result, defaultByteOrder, mn.NodeIndex)
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeReqVote.Serialize ] failed to write ti a buffer")
	}
	return result.Bytes(), nil
}

func (mn *MissingNodeReqVote) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &mn.NodeIndex)
	if err != nil {
		return errors.Wrap(err, "[ MissingNodeReqVote.Deserialize ] failed to read a node index")
	}
	return nil
}

func (mn *MissingNodeClaimsVote) Serialize() ([]byte, error) {
	var result bytes.Buffer
	err := binary.Write(&result, defaultByteOrder, mn.NodeIndex)
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeClaimsVote.Serialize ] Can't write NodeIndex")
	}
	serializedClaim, err := mn.Claim.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeClaimsVote.Serialize ] Can't serialize claim")
	}
	err = binary.Write(&result, defaultByteOrder, uint16(len(serializedClaim)))
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeClaimsVote.Serialize ] Can't write claimSize")
	}
	err = binary.Write(&result, defaultByteOrder, serializedClaim)
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeClaimsVote.Serialize ] Can't write Claim")
	}
	return result.Bytes(), nil
}

func (mn *MissingNodeClaimsVote) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &mn.NodeIndex)
	if err != nil {
		return errors.Wrap(err, "[ MissingNodeClaimsVote.Deserialize ] Can't read NodeIndex")
	}
	err = binary.Read(data, defaultByteOrder, &mn.claimSize)
	if err != nil {
		return errors.Wrap(err, "[ MissingNodeClaimsVote.Deserialize ] Can't read claimSize")
	}
	claimData := make([]byte, mn.claimSize)
	err = binary.Read(data, defaultByteOrder, claimData)
	if err != nil {
		return errors.Wrap(err, "[ MissingNodeClaimsVote.Deserialize ] Can't read claim data")
	}
	claims, err := parseReferendumClaim(claimData)
	if err != nil {
		return errors.Wrap(err, "[ MissingNodeClaimsVote.Deserialize ] Can't parse claim from claim data")
	}
	mn.Claim = claims[0]
	return nil
}

func (v *MissingNodeRespVote) Type() VoteType {
	return TypeMissingNodeRespVote
}

// Deserialize implements interface method
func (v *MissingNodeRespVote) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &v.NodeIndex)
	if err != nil {
		return errors.Wrap(err, "[ MissingNodeRespVote.Deserialize ] Can't read NodeIndex")
	}
	err = v.NodePulseProof.Deserialize(data)
	if err != nil {
		return errors.Wrap(err, "[ MissingNodeRespVote.Deserialize ] Can't read NodePulseProof")
	}
	err = v.NodeClaimUnsigned.deserializeRaw(data)
	if err != nil {
		return errors.Wrap(err, "[ MissingNodeRespVote.Deserialize ] Can't read NodeClaimUnsigned")
	}

	return nil
}

// Serialize implements interface method
func (v *MissingNodeRespVote) Serialize() ([]byte, error) {
	result := allocateBuffer(1024)

	err := binary.Write(result, defaultByteOrder, v.NodeIndex)
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeRespVote.Serialize ] Can't write NodeIndex")
	}

	nodePulseProofRaw, err := v.NodePulseProof.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeRespVote.Serialize ] Can't serialize NodePulseProof")
	}
	_, err = result.Write(nodePulseProofRaw)
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeRespVote.Serialize ] Can't append NodePulseProof")
	}

	joinClaim, err := v.NodeClaimUnsigned.SerializeRaw()
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeRespVote.Serialize ] Can't serialize join claim")
	}
	_, err = result.Write(joinClaim)
	if err != nil {
		return nil, errors.Wrap(err, "[ MissingNodeRespVote.Serialize ] Can't write join claim")
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
		case TypeMissingNodeRespVote:
			refVote = &MissingNodeRespVote{}
		case TypeMissingNodeReqVote:
			refVote = &MissingNodeReqVote{}
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
