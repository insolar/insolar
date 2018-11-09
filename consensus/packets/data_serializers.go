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
	"io/ioutil"

	"github.com/pkg/errors"
)

var defaultByteOrder = binary.BigEndian

// ----------------------------------PHASE 1--------------------------------

// routInfoMasks auxiliar constants
const (
	// take low bit
	hasRoutingMask = 0x1

	packetTypeMask   = 0x7f
	packetTypeOffset = 1
)

func (ph *PacketHeader) parseRouteInfo(routInfo uint8) {
	ph.PacketT = PacketType(routInfo&packetTypeMask) >> packetTypeOffset
	ph.HasRouting = (routInfo & hasRoutingMask) == 1
}

func (ph *PacketHeader) compactRouteInfo() uint8 {
	var result uint8
	result |= uint8(ph.PacketT) << packetTypeOffset

	if ph.HasRouting {
		result |= hasRoutingMask
	}

	return result
}

// PulseAndCustomFlags auxiliar constants
const (
	// take bit before high bit
	f00Mask  = 0x40000000
	f00Shift = 30

	// take high bit
	f01Mask   = 0x80000000
	f01Shift  = 31
	pulseMask = 0x3fffffff
)

func (ph *PacketHeader) parsePulseAndCustomFlags(pulseAndCustomFlags uint32) {
	ph.f01 = (pulseAndCustomFlags >> f01Shift) == 1
	ph.f00 = ((pulseAndCustomFlags & f00Mask) >> f00Shift) == 1
	ph.Pulse = pulseAndCustomFlags & pulseMask
}

func (ph *PacketHeader) compactPulseAndCustomFlags() uint32 {
	var result uint32
	if ph.f01 {
		result |= f01Mask
	}
	if ph.f00 {
		result |= f00Mask
	}
	result |= ph.Pulse & pulseMask

	return result
}

// claims auxiliar constants
const (
	claimTypeShift      = 10
	claimHeaderTypeMask = 0xfc00

//	claimHeaderLengthMask = 0x3ff
)

func extractClaimTypeFromHeader(claimHeader uint16) uint8 {
	return uint8((claimHeader & claimHeaderTypeMask) >> claimTypeShift)
}

// func extractClaimLengthFromHeader(claimHeader uint16) uint16 {
// 	return claimHeader & claimHeaderLengthMask
// }

func makeClaimHeader(claim ReferendumClaim) uint16 {
	// TODO: we don't need length
	var result = claim.Length()
	result |= (uint16(claim.Type()) << claimTypeShift)

	return result
}

func (p1p *Phase1Packet) parseReferendumClaim(data []byte) error {
	claimsSize := len(data)
	claimsBufReader := bytes.NewReader(data)
	for claimsSize > 0 {
		startSize := claimsBufReader.Len()
		var claimHeader uint16
		err := binary.Read(claimsBufReader, defaultByteOrder, &claimHeader)
		if err != nil {
			return errors.Wrap(err, "[ PacketHeader.parseReferendumClaim ] Can't read claimHeader")
		}

		claimType := ClaimType(extractClaimTypeFromHeader(claimHeader))
		// TODO: Do we need claimLength?
		//claimLength := extractClaimLengthFromHeader(claimHeader)
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
		}
		err = refClaim.Deserialize(claimsBufReader)
		if err != nil {
			return errors.Wrap(err, "[ PacketHeader.parseReferendumClaim ] Can't deserialize claim")
		}
		p1p.claims = append(p1p.claims, refClaim)

		claimsSize -= startSize - claimsBufReader.Len()
	}

	if claimsSize != 0 {
		return errors.New("[ PacketHeader.parseReferendumClaim ] Problem with claims struct")
	}

	return nil
}

func (p1p *Phase1Packet) compactReferendumClaim() ([]byte, error) {
	result := allocateBuffer(2048)
	for _, claim := range p1p.claims {
		claimHeader := makeClaimHeader(claim)
		err := binary.Write(result, defaultByteOrder, claimHeader)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("[ PacketHeader.compactReferendumClaim ] "+
				"Can't write claim header. Type: %d. Length: %d", claim.Type(), claim.Length()))
		}

		rawClaim, err := claim.Serialize()
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("[ PacketHeader.compactReferendumClaim ] "+
				"Can't serialize claim. Type: %d. Length: %d", claim.Type(), claim.Length()))
		}
		_, err = result.Write(rawClaim)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("[ PacketHeader.compactReferendumClaim ] "+
				"Can't append proofNodePulseRaw."+"Type: %d. Length: %d", claim.Type(), claim.Length()))
		}
	}

	return result.Bytes(), nil
}

func (p1p *Phase1Packet) Deserialize(data io.Reader) error {
	err := p1p.packetHeader.Deserialize(data)
	if err != nil {
		return errors.Wrap(err, "[ Phase1Packet.Deserialize ] Can't deserialize packetHeader")
	}

	err = p1p.pulseData.Deserialize(data)
	if err != nil {
		return errors.Wrap(err, "[ Phase1Packet.Deserialize ] Can't deserialize pulseData")
	}

	err = p1p.proofNodePulse.Deserialize(data)
	if err != nil {
		return errors.Wrap(err, "[ Phase1Packet.Deserialize ] Can't deserialize proofNodePulse")
	}

	if p1p.hasSection2() {
		claimsBuf, err := ioutil.ReadAll(data)
		if err != nil {
			return errors.Wrap(err, "[ Phase1Packet.Deserialize ] Can't read Section 2")
		}
		claimsSize := len(claimsBuf) - 8

		err = p1p.parseReferendumClaim(claimsBuf[:claimsSize])
		if err != nil {
			return errors.Wrap(err, "[ Phase1Packet.Deserialize ] Can't parseReferendumClaim")
		}

		data = bytes.NewReader(claimsBuf[claimsSize:])
	}

	err = binary.Read(data, defaultByteOrder, &p1p.signature)
	if err != nil {
		return errors.Wrap(err, "[ Phase1Packet.Deserialize ] Can't read signature")
	}

	return nil
}

func (p1p *Phase1Packet) Serialize() ([]byte, error) {
	result := allocateBuffer(2048)

	// serializing of  PacketHeader
	packetHeaderRaw, err := p1p.packetHeader.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase1Packet.Serialize ] Can't serialize packetHeader")
	}
	_, err = result.Write(packetHeaderRaw)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase1Packet.Serialize ] Can't append packetHeader")
	}

	// serializing of  PulseData
	pulseDataRaw, err := p1p.pulseData.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase1Packet.Serialize ] Can't serialize pulseDataRaw")
	}
	_, err = result.Write(pulseDataRaw)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase1Packet.Serialize ] Can't append pulseDataRaw")
	}

	// serializing of ProofNodePulse
	proofNodePulseRaw, err := p1p.proofNodePulse.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase1Packet.Serialize ] Can't serialize proofNodePulseRaw")
	}
	_, err = result.Write(proofNodePulseRaw)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase1Packet.Serialize ] Can't append proofNodePulseRaw")
	}

	// serializing of ReferendumClaim
	claimRaw, err := p1p.compactReferendumClaim()
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase1Packet.Serialize ] Can't append claimRaw")
	}
	_, err = result.Write(claimRaw)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase1Packet.Serialize ] Can't append claimRaw")
	}

	// serializing of signature
	err = binary.Write(result, defaultByteOrder, p1p.signature)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase1Packet.Serialize ] Can't write signature")
	}

	return result.Bytes(), nil

}

// Deserialize implements interface method
func (ph *PacketHeader) Deserialize(data io.Reader) error {
	var routInfo uint8
	err := binary.Read(data, defaultByteOrder, &routInfo)
	if err != nil {
		return errors.Wrap(err, "[ PacketHeader.Deserialize ] Can't read routInfo")
	}
	ph.parseRouteInfo(routInfo)

	var pulseAndCustomFlags uint32
	err = binary.Read(data, defaultByteOrder, &pulseAndCustomFlags)
	if err != nil {
		return errors.Wrap(err, "[ PacketHeader.Deserialize ] Can't read pulseAndCustomFlags")
	}
	ph.parsePulseAndCustomFlags(pulseAndCustomFlags)

	err = binary.Read(data, defaultByteOrder, &ph.OriginNodeID)
	if err != nil {
		return errors.Wrap(err, "[ PacketHeader.Deserialize ] Can't read OriginNodeID")
	}

	err = binary.Read(data, defaultByteOrder, &ph.TargetNodeID)
	if err != nil {
		return errors.Wrap(err, "[ PacketHeader.Deserialize ] Can't read TargetNodeID")
	}

	return nil
}

func allocateBuffer(n int) *bytes.Buffer {
	buf := make([]byte, 0, n)
	result := bytes.NewBuffer(buf)
	return result
}

// Serialize implements interface method
func (ph *PacketHeader) Serialize() ([]byte, error) {
	result := allocateBuffer(64)
	routeInfo := ph.compactRouteInfo()
	err := binary.Write(result, defaultByteOrder, routeInfo)
	if err != nil {
		return nil, errors.Wrap(err, "[ PacketHeader.Serialize ] Can't write routeInfo")
	}

	pulseAndCustomFlags := ph.compactPulseAndCustomFlags()
	err = binary.Write(result, defaultByteOrder, pulseAndCustomFlags)
	if err != nil {
		return nil, errors.Wrap(err, "[ PacketHeader.Serialize ] Can't write pulseAndCustomFlags")
	}

	err = binary.Write(result, defaultByteOrder, ph.OriginNodeID)
	if err != nil {
		return nil, errors.Wrap(err, "[ PacketHeader.Serialize ] Can't write OriginNodeID")
	}

	err = binary.Write(result, defaultByteOrder, ph.TargetNodeID)
	if err != nil {
		return nil, errors.Wrap(err, "[ PacketHeader.Serialize ] Can't write TargetNodeID")
	}

	return result.Bytes(), nil
}

// Deserialize implements interface method
func (pde *PulseDataExt) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &pde.NextPulseDelta)
	if err != nil {
		return errors.Wrap(err, "[ PulseDataExt.Deserialize ] Can't read NextPulseDelta")
	}

	err = binary.Read(data, defaultByteOrder, &pde.PrevPulseDelta)
	if err != nil {
		return errors.Wrap(err, "[ PulseDataExt.Deserialize ] Can't read PrevPulseDelta")
	}

	err = binary.Read(data, defaultByteOrder, &pde.OriginID)
	if err != nil {
		return errors.Wrap(err, "[ PulseDataExt.Deserialize ] Can't read OriginID")
	}

	err = binary.Read(data, defaultByteOrder, &pde.EpochPulseNo)
	if err != nil {
		return errors.Wrap(err, "[ PulseDataExt.Deserialize ] Can't read EpochPulseNo")
	}

	err = binary.Read(data, defaultByteOrder, &pde.PulseTimestamp)
	if err != nil {
		return errors.Wrap(err, "[ PulseDataExt.Deserialize ] Can't read PulseTimestamp")
	}

	err = binary.Read(data, defaultByteOrder, &pde.Entropy)
	if err != nil {
		return errors.Wrap(err, "[ PulseDataExt.Deserialize ] Can't read Entropy")
	}

	return nil
}

// Serialize implements interface method
func (pde *PulseDataExt) Serialize() ([]byte, error) {
	result := allocateBuffer(256)
	err := binary.Write(result, defaultByteOrder, pde.NextPulseDelta)
	if err != nil {
		return nil, errors.Wrap(err, "[ PulseDataExt.Serialize ] Can't write NextPulseDelta")
	}

	err = binary.Write(result, defaultByteOrder, pde.PrevPulseDelta)
	if err != nil {
		return nil, errors.Wrap(err, "[ PulseDataExt.Serialize ] Can't write PrevPulseDelta")
	}

	err = binary.Write(result, defaultByteOrder, pde.OriginID)
	if err != nil {
		return nil, errors.Wrap(err, "[ PulseDataExt.Serialize ] Can't write OriginID")
	}

	err = binary.Write(result, defaultByteOrder, pde.EpochPulseNo)
	if err != nil {
		return nil, errors.Wrap(err, "[ PulseDataExt.Serialize ] Can't write EpochPulseNo")
	}

	err = binary.Write(result, defaultByteOrder, pde.PulseTimestamp)
	if err != nil {
		return nil, errors.Wrap(err, "[ PulseDataExt.Serialize ] Can't write PulseTimestamp")
	}

	err = binary.Write(result, defaultByteOrder, pde.Entropy)
	if err != nil {
		return nil, errors.Wrap(err, "[ PulseDataExt.Serialize ] Can't write Entropy")
	}

	return result.Bytes(), nil
}

// Deserialize implements interface method
func (pd *PulseData) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &pd.PulseNumber)
	if err != nil {
		return errors.Wrap(err, "[ PulseData.Deserialize ] Can't read PulseNumer")
	}

	pd.Data = &PulseDataExt{}

	err = pd.Data.Deserialize(data)
	if err != nil {
		return errors.Wrap(err, "[ PulseData.Deserialize ] Can't read PulseDataExt")
	}

	return nil
}

// Serialize implements interface method
func (pd *PulseData) Serialize() ([]byte, error) {
	result := allocateBuffer(64)
	err := binary.Write(result, defaultByteOrder, pd.PulseNumber)
	if err != nil {
		return nil, errors.Wrap(err, "[ PulseData.Serialize ] Can't write PulseNumer")
	}

	pulseDataExtRaw, err := pd.Data.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ PulseData.Serialize ] Can't write PulseDataExt")
	}

	_, err = result.Write(pulseDataExtRaw)
	if err != nil {
		return nil, errors.Wrap(err, "[ PulseData.Serialize ] Can't append PulseDataExt")
	}

	return result.Bytes(), nil
}

// Deserialize implements interface method
func (npp *NodePulseProof) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &npp.NodeStateHash)
	if err != nil {
		return errors.Wrap(err, "[ NodePulseProof.Deserialize ] Can't read NodeStateHash")
	}

	err = binary.Read(data, defaultByteOrder, &npp.NodeSignature)
	if err != nil {
		return errors.Wrap(err, "[ NodePulseProof.Deserialize ] Can't read NodeSignature")
	}

	return nil
}

// Serialize implements interface method
func (npp *NodePulseProof) Serialize() ([]byte, error) {
	result := allocateBuffer(128)
	err := binary.Write(result, defaultByteOrder, npp.NodeStateHash)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodePulseProof.Serialize ] Can't write NodeStateHash")
	}

	err = binary.Write(result, defaultByteOrder, npp.NodeSignature)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodePulseProof.Serialize ] Can't write NodeSignature")
	}

	return result.Bytes(), nil
}

// Deserialize implements interface method
func (nb *NodeBroadcast) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &nb.EmergencyLevel)
	if err != nil {
		return errors.Wrap(err, "[ NodeBroadcast.Deserialize ] Can't read EmergencyLevel")
	}

	err = binary.Read(data, defaultByteOrder, &nb.length)
	if err != nil {
		return errors.Wrap(err, "[ NodeBroadcast.Deserialize ] Can't read length")
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

	err = binary.Write(result, defaultByteOrder, nb.length)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeBroadcast.Serialize ] Can't write length")
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

	err = binary.Read(data, defaultByteOrder, &cpa.length)
	if err != nil {
		return errors.Wrap(err, "[ CapabilityPoolingAndActivation.Deserialize ] Can't read length")
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

	err = binary.Write(result, defaultByteOrder, cpa.length)
	if err != nil {
		return nil, errors.Wrap(err, "[ CapabilityPoolingAndActivation.Serialize ] Can't write length")
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

	err = binary.Read(data, defaultByteOrder, &nvb.claimType)
	if err != nil {
		return errors.Wrap(err, "[ NodeViolationBlame.Deserialize ] Can't read claimType")
	}

	err = binary.Read(data, defaultByteOrder, &nvb.length)
	if err != nil {
		return errors.Wrap(err, "[ NodeViolationBlame.Deserialize ] Can't read length")
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

	err = binary.Write(result, defaultByteOrder, nvb.claimType)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeViolationBlame.Serialize ] Can't write claimType")
	}

	err = binary.Write(result, defaultByteOrder, nvb.length)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeViolationBlame.Serialize ] Can't write length")
	}

	return result.Bytes(), nil
}

// Deserialize implements interface method
func (njc *NodeJoinClaim) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &njc.NodeID)
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

	// err = binary.Read(data, defaultByteOrder, &njc.length)
	// if err != nil {
	// 	return errors.Wrap(err, "[ NodeJoinClaim.Deserialize ] Can't read length")
	// }

	return nil
}

// Serialize implements interface method
func (njc *NodeJoinClaim) Serialize() ([]byte, error) {
	result := allocateBuffer(1024)
	err := binary.Write(result, defaultByteOrder, njc.NodeID)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.Serialize ] Can't write NodeID")
	}

	err = binary.Write(result, defaultByteOrder, njc.RelayNodeID)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.Serialize ] Can't write RelayNodeID")
	}

	err = binary.Write(result, defaultByteOrder, njc.ProtocolVersionAndFlags)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.Serialize ] Can't write ProtocolVersionAndFlags")
	}

	err = binary.Write(result, defaultByteOrder, njc.JoinsAfter)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.Serialize ] Can't write JoinsAfter")
	}

	err = binary.Write(result, defaultByteOrder, njc.NodeRoleRecID)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.Serialize ] Can't write NodeRoleRecID")
	}

	err = binary.Write(result, defaultByteOrder, njc.NodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.Serialize ] Can't write NodeRef")
	}

	err = binary.Write(result, defaultByteOrder, njc.NodePK)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeJoinClaim.Serialize ] Can't write NodePK")
	}

	// err = binary.Write(result, defaultByteOrder, njc.length)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "[ NodeJoinClaim.Serialize ] Can't write length")
	// }

	return result.Bytes(), nil
}

// Deserialize implements interface method
func (nlc *NodeLeaveClaim) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &nlc.length)
	if err != nil {
		return errors.Wrap(err, "[ NodeLeaveClaim.Deserialize ] Can't read length")
	}

	return nil
}

// Serialize implements interface method
func (nlc *NodeLeaveClaim) Serialize() ([]byte, error) {
	result := allocateBuffer(64)
	err := binary.Write(result, defaultByteOrder, nlc.length)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeLeaveClaim.Serialize ] Can't write length")
	}

	return result.Bytes(), nil
}

// ----------------------------------PHASE 2--------------------------------

// Deserialize implements interface method
func (rv *ReferendumVote) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &rv.Type)
	if err != nil {
		return errors.Wrap(err, "[ ReferendumVote.Deserialize ] Can't read Type")
	}

	err = binary.Read(data, defaultByteOrder, &rv.Length)
	if err != nil {
		return errors.Wrap(err, "[ ReferendumVote.Deserialize ] Can't read Length")
	}

	return nil
}

// Serialize implements interface method
func (rv *ReferendumVote) Serialize() ([]byte, error) {
	result := allocateBuffer(64)
	err := binary.Write(result, defaultByteOrder, rv.Type)
	if err != nil {
		return nil, errors.Wrap(err, "[ ReferendumVote.Serialize ] Can't write Type")
	}

	err = binary.Write(result, defaultByteOrder, rv.Length)
	if err != nil {
		return nil, errors.Wrap(err, "[ ReferendumVote.Serialize ] Can't write Length")
	}

	return result.Bytes(), nil
}

// Deserialize implements interface method
func (nlv *NodeListVote) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &nlv.NodeListCount)
	if err != nil {
		return errors.Wrap(err, "[ NodeListVote.Deserialize ] Can't read NodeListCount")
	}

	err = binary.Read(data, defaultByteOrder, &nlv.NodeListHash)
	if err != nil {
		return errors.Wrap(err, "[ NodeListVote.Deserialize ] Can't read NodeListHash")
	}

	return nil
}

// Serialize implements interface method
func (nlv *NodeListVote) Serialize() ([]byte, error) {
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

// DeviantBitSet auxiliar constants
const (
	// take high bit
	compressedSetMask   = 0x80
	compressedSetOffset = 7

	highBitLengthFlagMask   = 0x40
	highBitLengthFlagOffset = 6
	lowBitLengthMask        = 0x3f
)

func (dbs *DeviantBitSet) parsePackedData(packedData uint8) {
	dbs.CompressedSet = (packedData >> compressedSetOffset) == 1
	dbs.HighBitLengthFlag = ((packedData & highBitLengthFlagMask) >> highBitLengthFlagOffset) == 1
	dbs.LowBitLength = packedData & lowBitLengthMask
}

func (dbs *DeviantBitSet) compactPacketData() uint8 {
	var result uint8

	if dbs.CompressedSet {
		result |= compressedSetMask
	}
	if dbs.HighBitLengthFlag {
		result |= highBitLengthFlagMask
	}

	result |= dbs.LowBitLength & lowBitLengthMask

	return result
}

// Deserialize implements interface method
func (dbs *DeviantBitSet) Deserialize(data io.Reader) error {
	var packedData uint8
	err := binary.Read(data, defaultByteOrder, &packedData)
	if err != nil {
		return errors.Wrap(err, "[ DeviantBitSet.Deserialize ] Can't read packedData")
	}
	dbs.parsePackedData(packedData)

	// TODO: these fields are optional
	err = binary.Read(data, defaultByteOrder, &dbs.HighBitLength)
	if err != nil {
		return errors.Wrap(err, "[ DeviantBitSet.Deserialize ] Can't read HighBitLength")
	}

	return nil
	// // TODO: calc correct size
	// dbs.Payload = make([]byte, transport.GetUDPMaxPacketSize())
	// n, err := data.Read(dbs.Payload)
	// if err != nil {
	// 	return errors.Wrap(err, "[ DeviantBitSet.Deserialize ] Can't read Payload")
	// }
	// dbs.Payload = dbs.Payload[:n]
	//
	// return nil
}

// Serialize implements interface method
func (dbs *DeviantBitSet) Serialize() ([]byte, error) {
	result := allocateBuffer(2048)

	packedData := dbs.compactPacketData()
	err := binary.Write(result, defaultByteOrder, packedData)
	if err != nil {
		return nil, errors.Wrap(err, "[ DeviantBitSet.Serialize ] Can't write packedData")
	}

	// TODO: these fields are optional
	err = binary.Write(result, defaultByteOrder, dbs.HighBitLength)
	if err != nil {
		return nil, errors.Wrap(err, "[ DeviantBitSet.Serialize ] Can't write HighBitLength")
	}

	return result.Bytes(), nil
	// _, err = result.Write(dbs.Payload)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "[ DeviantBitSet.Serialize ] Can't write Payload")
	// }
	//
	// return result.Bytes(), nil
}

func (phase2Packet *Phase2Packet) Deserialize(data io.Reader) error {
	err := phase2Packet.packetHeader.Deserialize(data)
	if err != nil {
		return errors.Wrap(err, "[ Phase2Packet.Deserialize ] Can't deserialize packetHeader")
	}

	err = binary.Read(data, defaultByteOrder, &phase2Packet.globuleHashSignature)
	if err != nil {
		return errors.Wrap(err, "[ Phase2Packet.Deserialize ] Can't read globuleHashSignature")
	}

	err = phase2Packet.deviantBitSet.Deserialize(data)
	if err != nil {
		return errors.Wrap(err, "[ Phase2Packet.Deserialize ] Can't deserialize deviantBitSet")
	}

	err = binary.Read(data, defaultByteOrder, &phase2Packet.signatureHeaderSection1)
	if err != nil {
		return errors.Wrap(err, "[ Phase2Packet.Deserialize ] Can't read signatureHeaderSection1")
	}

	// TODO: add reading Referendum vote

	err = binary.Read(data, defaultByteOrder, &phase2Packet.signatureHeaderSection2)
	if err != nil {
		return errors.Wrap(err, "[ Phase2Packet.Deserialize ] Can't read signatureHeaderSection2")
	}

	return nil
}

func (phase2Packet *Phase2Packet) Serialize() ([]byte, error) {
	result := allocateBuffer(2048)

	// serializing of  PacketHeader
	packetHeaderRaw, err := phase2Packet.packetHeader.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase2Packet.Serialize ] Can't serialize packetHeader")
	}
	_, err = result.Write(packetHeaderRaw)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase2Packet.Serialize ] Can't append packetHeader")
	}

	err = binary.Write(result, defaultByteOrder, phase2Packet.globuleHashSignature)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase2Packet.Serialize ] Can't write globuleHashSignature")
	}

	// serializing of deviantBitSet
	deviantBitSetRaw, err := phase2Packet.deviantBitSet.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase2Packet.Serialize ] Can't serialize deviantBitSet")
	}
	_, err = result.Write(deviantBitSetRaw)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase2Packet.Serialize ] Can't append deviantBitSet")
	}

	err = binary.Write(result, defaultByteOrder, phase2Packet.signatureHeaderSection1)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase2Packet.Serialize ] Can't write signatureHeaderSection1")
	}

	// TODO: add serialising Referendum vote

	err = binary.Write(result, defaultByteOrder, phase2Packet.signatureHeaderSection2)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase2Packet.Serialize ] Can't write signatureHeaderSection2")
	}

	return result.Bytes(), nil

}
