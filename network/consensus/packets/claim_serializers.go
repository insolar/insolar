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
	"fmt"
	"io"

	"github.com/insolar/insolar/insolar"
	"github.com/pkg/errors"
)

func (njc *NodeJoinClaim) deserializeRaw(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &njc.ShortNodeID)
	if err != nil {
		return errors.Wrap(err, "[ NodeJoinClaim.deserializeRaw ] Can't read NodeID")
	}

	err = binary.Read(data, defaultByteOrder, &njc.RelayNodeID)
	if err != nil {
		return errors.Wrap(err, "[ NodeJoinClaim.deserializeRaw ] Can't read RelayNodeID")
	}

	err = binary.Read(data, defaultByteOrder, &njc.ProtocolVersionAndFlags)
	if err != nil {
		return errors.Wrap(err, "[ NodeJoinClaim.deserializeRaw ] Can't read ProtocolVersionAndFlags")
	}

	err = binary.Read(data, defaultByteOrder, &njc.JoinsAfter)
	if err != nil {
		return errors.Wrap(err, "[ NodeJoinClaim.deserializeRaw ] Can't read JoinsAfter")
	}

	err = binary.Read(data, defaultByteOrder, &njc.NodeRoleRecID)
	if err != nil {
		return errors.Wrap(err, "[ NodeJoinClaim.deserializeRaw ] Can't read NodeRoleRecID")
	}

	err = binary.Read(data, defaultByteOrder, &njc.NodeAddress)
	if err != nil {
		return errors.Wrap(err, "[ NodeJoinClaim.deserializeRaw ] Can't read NodeAddress")
	}

	err = binary.Read(data, defaultByteOrder, &njc.NodeRef)
	if err != nil {
		return errors.Wrap(err, "[ NodeJoinClaim.deserializeRaw ] Can't read NodeRef")
	}

	err = binary.Read(data, defaultByteOrder, &njc.NodePK)
	if err != nil {
		return errors.Wrap(err, "[ NodeJoinClaim.deserializeRaw ] Can't read NodePK")
	}
	return nil
}

// Deserialize implements interface method
func (njc *NodeJoinClaim) Deserialize(data io.Reader) error {
	err := njc.deserializeRaw(data)
	if err != nil {
		return err
	}
	err = binary.Read(data, defaultByteOrder, &njc.Signature)
	if err != nil {
		return errors.Wrap(err, "[ NodeJoinClaim.Deserialize ] Can't read Signature")
	}
	return nil
}

// Serialize implements interface method
func (njc *NodeJoinClaim) Serialize() ([]byte, error) {
	result := allocateBuffer(1024)

	rawData, err := njc.SerializeRaw()
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.Serialize ] Failed to serialize a claim without header")
	}

	err = binary.Write(result, defaultByteOrder, rawData)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.Serialize ] Failed to write a data without header")
	}

	err = binary.Write(result, defaultByteOrder, njc.Signature[:])
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.Serialize ] Can't write Signature")
	}

	return result.Bytes(), nil
}

func (njc *NodeJoinClaim) SerializeRaw() ([]byte, error) {
	result := allocateBuffer(1024)

	err := binary.Write(result, defaultByteOrder, njc.ShortNodeID)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.SerializeRaw ] Can't write NodeID")
	}

	err = binary.Write(result, defaultByteOrder, njc.RelayNodeID)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.SerializeRaw ] Can't write RelayNodeID")
	}

	err = binary.Write(result, defaultByteOrder, njc.ProtocolVersionAndFlags)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.SerializeRaw ] Can't write ProtocolVersionAndFlags")
	}

	err = binary.Write(result, defaultByteOrder, njc.JoinsAfter)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.SerializeRaw ] Can't write JoinsAfter")
	}

	err = binary.Write(result, defaultByteOrder, njc.NodeRoleRecID)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.SerializeRaw ] Can't write NodeRoleRecID")
	}

	err = binary.Write(result, defaultByteOrder, njc.NodeAddress)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.SerializeRaw ] Can't write NodeAddress")
	}

	err = binary.Write(result, defaultByteOrder, njc.NodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.SerializeRaw ] Can't write NodeRef")
	}

	err = binary.Write(result, defaultByteOrder, njc.NodePK)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.SerializeRaw ] Can't write NodePK")
	}

	return result.Bytes(), nil
}

func (nac *NodeAnnounceClaim) SerializeRaw() ([]byte, error) {
	nodeJoinPart, err := nac.NodeJoinClaim.SerializeRaw()
	if err != nil {
		return nil, err
	}
	result := allocateBuffer(1024)
	err = binary.Write(result, defaultByteOrder, nodeJoinPart)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeAnnounceClaim.Serialize ] Can't write NodeJoinClaim part")
	}
	err = binary.Write(result, defaultByteOrder, nac.NodeAnnouncerIndex)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeAnnounceClaim.Serialize ] Can't write NodeAnnouncerIndex")
	}
	err = binary.Write(result, defaultByteOrder, nac.NodeJoinerIndex)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeAnnounceClaim.Serialize ] Can't write NodeJoinerIndex")
	}
	err = binary.Write(result, defaultByteOrder, nac.NodeCount)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeAnnounceClaim.Serialize ] Can't write NodeCount")
	}
	err = binary.Write(result, defaultByteOrder, nac.CloudHash)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeAnnounceClaim.Serialize ] Can't write CloudHash")
	}
	return result.Bytes(), nil
}

// Serialize implements interface method
func (nac *NodeAnnounceClaim) Serialize() ([]byte, error) {
	result := allocateBuffer(1024)

	rawData, err := nac.SerializeRaw()
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeAnnounceClaim.Serialize ] Failed to serialize a claim without header")
	}

	err = binary.Write(result, defaultByteOrder, rawData)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeAnnounceClaim.Serialize ] Failed to write a data without header")
	}

	err = binary.Write(result, defaultByteOrder, nac.Signature[:])
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeAnnounceClaim.Serialize ] Can't write Signature")
	}

	return result.Bytes(), nil
}

// Deserialize implements interface method
func (nac *NodeAnnounceClaim) Deserialize(data io.Reader) error {
	err := nac.deserializeRaw(data)
	if err != nil {
		return err
	}
	err = binary.Read(data, defaultByteOrder, &nac.NodeAnnouncerIndex)
	if err != nil {
		return errors.Wrap(err, "[ NodeAnnounceClaim.Deserialize ] Can't read NodeAnnouncerIndex")
	}
	err = binary.Read(data, defaultByteOrder, &nac.NodeJoinerIndex)
	if err != nil {
		return errors.Wrap(err, "[ NodeAnnounceClaim.Deserialize ] Can't read NodeJoinerIndex")
	}
	err = binary.Read(data, defaultByteOrder, &nac.NodeCount)
	if err != nil {
		return errors.Wrap(err, "[ NodeAnnounceClaim.Deserialize ] Can't read NodeCount")
	}
	err = binary.Read(data, defaultByteOrder, &nac.CloudHash)
	if err != nil {
		return errors.Wrap(err, "[ NodeAnnounceClaim.Deserialize ] Can't read CloudHash")
	}
	err = binary.Read(data, defaultByteOrder, &nac.Signature)
	if err != nil {
		return errors.Wrap(err, "[ NodeAnnounceClaim.Deserialize ] Can't read Signature")
	}
	return nil
}

func (nac *NodeAnnounceClaim) Update(nodeJoinerID insolar.Reference, crypto insolar.Signer) error {
	index, err := nac.BitSetMapper.RefToIndex(nodeJoinerID)
	if err != nil {
		return errors.Wrap(err, "[ NodeAnnounceClaim.Update ] failed to map joiner node ID to bitset index")
	}
	nac.NodeJoinerIndex = uint16(index)
	data, err := nac.SerializeRaw()
	if err != nil {
		return errors.Wrap(err, "[ NodeAnnounceClaim.Update ] failed to serialize raw announce claim")
	}
	signature, err := crypto.Sign(data)
	if err != nil {
		return errors.Wrap(err, "[ NodeAnnounceClaim.Update ] failed to sign announce claim")
	}
	sign := signature.Bytes()
	copy(nac.Signature[:], sign[:SignatureLength])
	return nil
}

// Serialize implements interface method
func (nlc *NodeLeaveClaim) Serialize() ([]byte, error) {
	var result bytes.Buffer
	err := binary.Write(&result, defaultByteOrder, nlc.ETA)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeLeaveClaim.Serialize ] failed to write ETA to buffer")
	}
	return result.Bytes(), nil
}

// Deserialize implements interface method
func (nlc *NodeLeaveClaim) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &nlc.ETA)
	if err != nil {
		return errors.Wrap(err, "[ NodeLeaveClaim.Deserialize ] failed to read a ETA")
	}
	return nil
}

// Serialize implements interface method
func (cnc *ChangeNetworkClaim) Serialize() ([]byte, error) {
	var result bytes.Buffer
	err := binary.Write(&result, defaultByteOrder, []byte(cnc.Address))
	if err != nil {
		return nil, errors.Wrap(err, "[ ChangeNetworkClaim.Serialize ] failed to write ETA to buffer")
	}
	return result.Bytes(), nil
}

// Deserialize implements interface method
func (cnc *ChangeNetworkClaim) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &cnc.Address)
	if err != nil {
		return errors.Wrap(err, "[ ChangeNetworkClaim.Deserialize ] failed to read a ETA")
	}
	return nil
}

func serializeClaims(claims []ReferendumClaim) ([]byte, error) {
	result := allocateBuffer(packetMaxSize)
	for _, claim := range claims {
		claimHeader := makeClaimHeader(claim)
		err := binary.Write(result, defaultByteOrder, claimHeader)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("[ serializeClaims ] "+
				"Can't write claim header. Type: %d. Length: %d", claim.Type(), getClaimSize(claim)))
		}

		rawClaim, err := claim.Serialize()
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("[ serializeClaims ] "+
				"Can't serialize claim. Type: %d. Length: %d", claim.Type(), getClaimSize(claim)))
		}
		_, err = result.Write(rawClaim)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("[ serializeClaims ] "+
				"Can't write claim. Type: %d. Length: %d", claim.Type(), getClaimSize(claim)))
		}
	}

	return result.Bytes(), nil
}

func parseReferendumClaim(data []byte) ([]ReferendumClaim, error) {
	claimsSize := len(data)
	claimsBufReader := bytes.NewReader(data)
	result := make([]ReferendumClaim, 0)

	// get claim header
	for claimsSize > 0 {
		startSize := claimsBufReader.Len()
		var claimHeader uint16
		err := binary.Read(claimsBufReader, defaultByteOrder, &claimHeader)
		if err != nil {
			return nil, errors.Wrap(err, "[ PacketHeader.parseReferendumClaim ] Can't read claimHeader")
		}

		claimType := ClaimType(extractTypeFromHeader(claimHeader))
		// TODO: Do we need claimLength?
		// claimLength := extractLengthFromHeader(claimHeader)
		var refClaim ReferendumClaim

		switch claimType {
		case TypeNodeJoinClaim:
			refClaim = &NodeJoinClaim{}
		case TypeNodeLeaveClaim:
			refClaim = &NodeLeaveClaim{}
		case TypeNodeAnnounceClaim:
			refClaim = &NodeAnnounceClaim{}
		case TypeChangeNetworkClaim:
			refClaim = &ChangeNetworkClaim{}
		default:
			return nil, errors.Wrap(err, "[ PacketHeader.parseReferendumClaim ] Unsupported claim type.")
		}
		err = refClaim.Deserialize(claimsBufReader)
		if err != nil {
			return nil, errors.Wrap(err, "[ PacketHeader.parseReferendumClaim ] Can't deserialize claim")
		}
		result = append(result, refClaim)

		claimsSize -= startSize - claimsBufReader.Len()
	}

	if claimsSize != 0 {
		return nil, errors.New("[ PacketHeader.parseReferendumClaim ] Problem with claims struct")
	}

	return result, nil
}
