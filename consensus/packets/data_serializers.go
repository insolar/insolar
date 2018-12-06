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

func (p1p *Phase1Packet) DeserializeWithoutHeader(data io.Reader, header *PacketHeader) error {
	if header == nil {
		return errors.New("[ Phase1Packet.DeserializeWithoutHeader ] Can't deserialize pulseData")
	}
	if header.PacketT != Phase1 {
		return errors.New("[ Phase1Packet.DeserializeWithoutHeader ] Wrong packet type")
	}

	p1p.packetHeader = *header

	err := p1p.pulseData.Deserialize(data)
	if err != nil {
		return errors.Wrap(err, "[ Phase1Packet.DeserializeWithoutHeader ] Can't deserialize pulseData")
	}

	err = p1p.proofNodePulse.Deserialize(data)
	if err != nil {
		return errors.Wrap(err, "[ Phase1Packet.DeserializeWithoutHeader ] Can't deserialize proofNodePulse")
	}

	if p1p.hasSection2() {
		claimsBuf, err := ioutil.ReadAll(data)
		if err != nil {
			return errors.Wrap(err, "[ Phase1Packet.DeserializeWithoutHeader ] Can't read Section 2")
		}
		claimsSize := len(claimsBuf) - SignatureLength

		p1p.claims, err = parseReferendumClaim(claimsBuf[:claimsSize])
		if err != nil {
			return errors.Wrap(err, "[ Phase1Packet.DeserializeWithoutHeader ] Can't parseReferendumClaim")
		}

		data = bytes.NewReader(claimsBuf[claimsSize:])
	}

	err = binary.Read(data, defaultByteOrder, &p1p.Signature)
	if err != nil {
		return errors.Wrap(err, "[ Phase1Packet.DeserializeWithoutHeader ] Can't read signature")
	}

	return nil
}

func (p1p *Phase1Packet) Deserialize(data io.Reader) error {
	err := p1p.packetHeader.Deserialize(data)
	if err != nil {
		return errors.Wrap(err, "[ Phase1Packet.Deserialize ] Can't deserialize packetHeader")
	}

	err = p1p.DeserializeWithoutHeader(data, &p1p.packetHeader)
	if err != nil {
		return errors.Wrap(err, "[ Phase1Packet.Deserialize ] Can't deserialize body")
	}

	return nil
}

func (p1p *Phase1Packet) Serialize() ([]byte, error) {
	result := allocateBuffer(phase1PacketMaxSize)

	if !p1p.hasSection2() && len(p1p.claims) > 0 {
		return nil, errors.New("invalid Phase1Packet")
	}

	raw, err := p1p.RawBytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get raw bytes")
	}
	result.Write(raw)

	return result.Bytes(), nil
}

func (p1p *Phase1Packet) RawBytes() ([]byte, error) {
	result := allocateBuffer(2048)

	// serializing of  packetHeader
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
	claimRaw, err := serializeClaims(p1p.claims)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase1Packet.Serialize ] Can't append claimRaw")
	}
	_, err = result.Write(claimRaw)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase1Packet.Serialize ] Can't append claimRaw")
	}

	// serializing of signature
	err = binary.Write(result, defaultByteOrder, p1p.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase1Packet.Serialize ] Can't write signature")
	}

	return result.Bytes(), nil

}

func allocateBuffer(n int) *bytes.Buffer {
	buf := make([]byte, 0, n)
	result := bytes.NewBuffer(buf)
	return result
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

// ----------------------------------PHASE 2--------------------------------

func (p2p *Phase2Packet) DeserializeWithoutHeader(data io.Reader, header *PacketHeader) error {
	if header == nil {
		return errors.New("[ Phase2Packet.DeserializeWithoutHeader ] Can't deserialize pulseData")
	}
	if header.PacketT != Phase2 {
		return errors.New("[ Phase2Packet.DeserializeWithoutHeader ] Wrong packet type")
	}

	p2p.packetHeader = *header

	err := binary.Read(data, defaultByteOrder, &p2p.globuleHashSignature)
	if err != nil {
		return errors.Wrap(err, "[ Phase2Packet.Deserialize ] Can't read globuleHashSignature")
	}

	// err = p2p.bitSet.Deserialize(data)
	// if err != nil {
	// 	return errors.Wrap(err, "[ Phase2Packet.Deserialize ] Can't deserialize bitSet")
	// }

	err = binary.Read(data, defaultByteOrder, &p2p.SignatureHeaderSection1)
	if err != nil {
		return errors.Wrap(err, "[ Phase2Packet.Deserialize ] Can't read SignatureHeaderSection1")
	}

	// TODO: add reading Referendum vote

	err = binary.Read(data, defaultByteOrder, &p2p.SignatureHeaderSection2)
	if err != nil {
		return errors.Wrap(err, "[ Phase2Packet.Deserialize ] Can't read SignatureHeaderSection2")
	}

	return nil
}

func (p2p *Phase2Packet) Deserialize(data io.Reader) error {
	err := p2p.packetHeader.Deserialize(data)
	if err != nil {
		return errors.Wrap(err, "[ Phase2Packet.Deserialize ] Can't deserialize packetHeader")
	}

	err = p2p.DeserializeWithoutHeader(data, &p2p.packetHeader)
	if err != nil {
		return errors.Wrap(err, "[ Phase2Packet.Deserialize ] Can't deserialize body")
	}

	return nil

}

func (p2p *Phase2Packet) Serialize() ([]byte, error) {
	result := allocateBuffer(2048)

	raw1, err := p2p.RawFirstPart()
	if err != nil {
		return nil, errors.Wrap(err, "failed to serialize")
	}

	result.Write(raw1)

	err = binary.Write(result, defaultByteOrder, p2p.SignatureHeaderSection1)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase2Packet.Serialize ] Can't write SignatureHeaderSection1")
	}

	err = binary.Write(result, defaultByteOrder, p2p.SignatureHeaderSection2)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase2Packet.Serialize ] Can't write SignatureHeaderSection2")
	}

	// bitSetRaw, err := p2p.bitSet.Serialize()
	// if err != nil {
	// 	return nil, errors.Wrap(err, "[ Phase2Packet.Serialize ] Can't serialize bitSet")
	// }
	// _, err = result.Write(bitSetRaw)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "[ Phase2Packet.Serialize ] Can't append bitSet")
	// }

	return result.Bytes(), nil
}

func (p2p *Phase2Packet) RawFirstPart() ([]byte, error) {
	result := allocateBuffer(2048)

	packetHeaderRaw, err := p2p.packetHeader.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase2Packet.Serialize ] Can't serialize PacketHeader")
	}
	_, err = result.Write(packetHeaderRaw)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase2Packet.Serialize ] Can't append PacketHeader")
	}

	err = binary.Write(result, defaultByteOrder, p2p.globuleHashSignature)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase2Packet.Serialize ] Can't write globuleHashSignature")
	}

	return result.Bytes(), nil
}

func (p2p *Phase2Packet) RawSecondPart() ([]byte, error) {
	// TODO: add serialising Referendum vote
	return nil, nil
}

// ----------------------------------PHASE 3--------------------------------

func (p3p *Phase3Packet) Serialize() ([]byte, error) {
	rawBytes, err := p3p.RawBytes()
	if err != nil {
		return nil, errors.Wrap(err, "[ Serialize ] failed to get a raw bytes")
	}

	var data bytes.Buffer
	_, err = data.Write(rawBytes)
	if err != nil {
		return nil, errors.Wrap(err, "[ Serialize ] failed to write a raw bytes to buffer")
	}

	_, err = data.Write(p3p.SignatureHeaderSection1[:])
	if err != nil {
		return nil, errors.Wrap(err, "[ Serialize ] failed to write a signature to buffer")
	}

	return data.Bytes(), nil
}

func (p3p *Phase3Packet) RawBytes() ([]byte, error) {
	header, err := p3p.packetHeader.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ RawBytes ] failed to serialize p3p header")
	}

	bitset, err := p3p.deviantBitSet.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ RawBytes ] failed to serialize bitset")
	}

	var data bytes.Buffer

	_, err = data.Write(header)
	if err != nil {
		return nil, errors.Wrap(err, "[ RawBytes ] failed to write a header to buffer")
	}

	_, err = data.Write(bitset)
	if err != nil {
		return nil, errors.Wrap(err, "[ RawBytes ] failed to write a bitset to buffer")
	}

	_, err = data.Write(p3p.globuleHashSignature[:])
	if err != nil {
		return nil, errors.Wrap(err, "[ RawBytes ] failed to write a bitset to buffer")
	}

	return data.Bytes(), nil
}

func (p3p *Phase3Packet) Deserialize(data io.Reader) error {
	err := p3p.packetHeader.Deserialize(data)
	if err != nil {
		return errors.Wrap(err, "[ Deserialize ] failed to deserialize p3p header")
	}

	err = p3p.DeserializeWithoutHeader(data, &p3p.packetHeader)
	if err != nil {
		return errors.Wrap(err, "[ Deserialize ] failed to deserialize p3p data")
	}
	return nil
}

func (p3p *Phase3Packet) DeserializeWithoutHeader(data io.Reader, header *PacketHeader) error {
	bitset, err := DeserializeBitSet(data)
	if err != nil {
		return errors.Wrap(err, "[ DeserializeWithoutHeader ] failed to deserialize a bitset")
	}
	p3p.deviantBitSet = bitset

	err = binary.Read(data, defaultByteOrder, &p3p.globuleHashSignature)
	if err != nil {
		return errors.Wrap(err, "[ DeserializeWithoutHeader ] failed to deserialize p3p globule hash")
	}

	err = binary.Read(data, defaultByteOrder, &p3p.SignatureHeaderSection1)
	if err != nil {
		return errors.Wrap(err, "[ DeserializeWithoutHeader ] failed to deserialize p3p signature")
	}

	return nil
}
