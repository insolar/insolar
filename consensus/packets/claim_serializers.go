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
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// Deserialize implements interface method
func (nb *NodeBroadcast) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &nb.EmergencyLevel)
	if err != nil {
		return errors.Wrap(err, "[ NodeBroadcast.Deserialize ] Can't read EmergencyLevel")
	}

	return nil
}

// Serialize implements interface method
func (nb *NodeBroadcast) Serialize() ([]byte, error) {
	result := allocateBuffer(64)
	err := binary.Write(result, defaultByteOrder, nb.EmergencyLevel)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeBroadcast.Serialize ] Can't write EmergencyLevel")
	}

	return result.Bytes(), nil
}

// Deserialize implements interface method
func (cpa *CapabilityPoolingAndActivation) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &cpa.PollingFlags)
	if err != nil {
		return errors.Wrap(err, "[ NodeBroadcast.Deserialize ] Can't read PollingFlags")
	}

	err = binary.Read(data, defaultByteOrder, &cpa.CapabilityType)
	if err != nil {
		return errors.Wrap(err, "[ CapabilityPoolingAndActivation.Deserialize ] Can't read CapabilityType")
	}

	err = binary.Read(data, defaultByteOrder, &cpa.CapabilityRef)
	if err != nil {
		return errors.Wrap(err, "[ CapabilityPoolingAndActivation.Deserialize ] Can't read CapabilityRef")
	}

	return nil
}

// Serialize implements interface method
func (cpa *CapabilityPoolingAndActivation) Serialize() ([]byte, error) {
	result := allocateBuffer(128)
	err := binary.Write(result, defaultByteOrder, cpa.PollingFlags)
	if err != nil {
		return nil, errors.Wrap(err, "[ CapabilityPoolingAndActivation.Serialize ] Can't write PollingFlags")
	}

	err = binary.Write(result, defaultByteOrder, cpa.CapabilityType)
	if err != nil {
		return nil, errors.Wrap(err, "[ CapabilityPoolingAndActivation.Serialize ] Can't write CapabilityType")
	}

	err = binary.Write(result, defaultByteOrder, cpa.CapabilityRef)
	if err != nil {
		return nil, errors.Wrap(err, "[ CapabilityPoolingAndActivation.Serialize ] Can't write CapabilityRef")
	}

	return result.Bytes(), nil
}

// Deserialize implements interface method
func (nvb *NodeViolationBlame) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &nvb.BlameNodeID)
	if err != nil {
		return errors.Wrap(err, "[ NodeViolationBlame.Deserialize ] Can't read BlameNodeID")
	}

	err = binary.Read(data, defaultByteOrder, &nvb.TypeViolation)
	if err != nil {
		return errors.Wrap(err, "[ NodeViolationBlame.Deserialize ] Can't read TypeViolation")
	}

	return nil
}

// Serialize implements interface method
func (nvb *NodeViolationBlame) Serialize() ([]byte, error) {
	result := allocateBuffer(64)
	err := binary.Write(result, defaultByteOrder, nvb.BlameNodeID)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeViolationBlame.Serialize ] Can't write BlameNodeID")
	}

	err = binary.Write(result, defaultByteOrder, nvb.TypeViolation)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeViolationBlame.Serialize ] Can't write TypeViolation")
	}

	return result.Bytes(), nil
}

// Deserialize implements interface method
func (njc *NodeJoinClaim) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &njc.ShortNodeID)
	if err != nil {
		return errors.Wrap(err, "[ NodeJoinClaim.Deserialize ] Can't read NodeID")
	}

	err = binary.Read(data, defaultByteOrder, &njc.RelayNodeID)
	if err != nil {
		return errors.Wrap(err, "[ NodeJoinClaim.Deserialize ] Can't read RelayNodeID")
	}

	err = binary.Read(data, defaultByteOrder, &njc.ProtocolVersionAndFlags)
	if err != nil {
		return errors.Wrap(err, "[ NodeJoinClaim.Deserialize ] Can't read ProtocolVersionAndFlags")
	}

	err = binary.Read(data, defaultByteOrder, &njc.JoinsAfter)
	if err != nil {
		return errors.Wrap(err, "[ NodeJoinClaim.Deserialize ] Can't read JoinsAfter")
	}

	err = binary.Read(data, defaultByteOrder, &njc.NodeRoleRecID)
	if err != nil {
		return errors.Wrap(err, "[ NodeJoinClaim.Deserialize ] Can't read NodeRoleRecID")
	}

	err = binary.Read(data, defaultByteOrder, &njc.NodeRef)
	if err != nil {
		return errors.Wrap(err, "[ NodeJoinClaim.Deserialize ] Can't read NodeRef")
	}

	err = binary.Read(data, defaultByteOrder, &njc.NodePK)
	if err != nil {
		return errors.Wrap(err, "[ NodeJoinClaim.Deserialize ] Can't read NodePK")
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

	rawData, err := njc.SerializeWithoutSign()
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

func (njc *NodeJoinClaim) SerializeWithoutSign() ([]byte, error) {
	result := allocateBuffer(1024)

	err := binary.Write(result, defaultByteOrder, njc.ShortNodeID)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.SerializeWithoutSign ] Can't write NodeID")
	}

	err = binary.Write(result, defaultByteOrder, njc.RelayNodeID)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.SerializeWithoutSign ] Can't write RelayNodeID")
	}

	err = binary.Write(result, defaultByteOrder, njc.ProtocolVersionAndFlags)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.SerializeWithoutSign ] Can't write ProtocolVersionAndFlags")
	}

	err = binary.Write(result, defaultByteOrder, njc.JoinsAfter)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.SerializeWithoutSign ] Can't write JoinsAfter")
	}

	err = binary.Write(result, defaultByteOrder, njc.NodeRoleRecID)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.SerializeWithoutSign ] Can't write NodeRoleRecID")
	}

	err = binary.Write(result, defaultByteOrder, njc.NodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.SerializeWithoutSign ] Can't write NodeRef")
	}

	err = binary.Write(result, defaultByteOrder, njc.NodePK)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.SerializeWithoutSign ] Can't write NodePK")
	}

	return result.Bytes(), nil
}

// Serialize implements interface method
func (nac *NodeAnnounceClaim) Serialize() ([]byte, error) {
	nodeJoinPart, err := nac.NodeJoinClaim.Serialize()
	if err != nil {
		return nil, err
	}
	result := allocateBuffer(1024)
	err = binary.Write(result, defaultByteOrder, nodeJoinPart)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeAnnounceClaim.Serialize ] Can't write NodeJoinClaim part")
	}
	err = binary.Write(result, defaultByteOrder, nac.NodeIndex)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeAnnounceClaim.Serialize ] Can't write NodeIndex")
	}
	err = binary.Write(result, defaultByteOrder, nac.NodeCount)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeAnnounceClaim.Serialize ] Can't write NodeCount")
	}
	return result.Bytes(), nil
}

// Deserialize implements interface method
func (nac *NodeAnnounceClaim) Deserialize(data io.Reader) error {
	err := nac.NodeJoinClaim.Deserialize(data)
	if err != nil {
		return err
	}
	err = binary.Read(data, defaultByteOrder, &nac.NodeIndex)
	if err != nil {
		return errors.Wrap(err, "[ NodeAnnounceClaim.Deserialize ] Can't read NodeIndex")
	}
	err = binary.Read(data, defaultByteOrder, &nac.NodeCount)
	if err != nil {
		return errors.Wrap(err, "[ NodeAnnounceClaim.Deserialize ] Can't read NodeCount")
	}
	return nil
}

// Deserialize implements interface method
func (nlc *NodeLeaveClaim) Deserialize(data io.Reader) error {
	return nil
}

// Serialize implements interface method
func (nlc *NodeLeaveClaim) Serialize() ([]byte, error) {
	return nil, nil
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
		case TypeCapabilityPollingAndActivation:
			refClaim = &CapabilityPoolingAndActivation{}
		case TypeNodeViolationBlame:
			refClaim = &NodeViolationBlame{}
		case TypeNodeBroadcast:
			refClaim = &NodeBroadcast{}
		case TypeNodeLeaveClaim:
			refClaim = &NodeLeaveClaim{}
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
