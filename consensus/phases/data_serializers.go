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

package phases

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

type Serializer interface {
	Serialize() ([]byte, error)
	Deserialize(data io.Reader) error
}

// Deserialize implements interface method
func (ph *PacketHeader) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &ph.Routing)
	if err != nil {
		return errors.Wrap(err, "[ PacketHeader.Deserialize ] Can't read Routing")
	}

	err = binary.Read(data, defaultByteOrder, &ph.Pulse)
	if err != nil {
		return errors.Wrap(err, "[ PacketHeader.Deserialize ] Can't read Pulse")
	}

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
	result := new(bytes.Buffer)
	err := binary.Write(result, defaultByteOrder, ph.Routing)
	if err != nil {
		return nil, errors.Wrap(err, "[ PacketHeader.Serialize ] Can't write Routing")
	}

	err = binary.Write(result, defaultByteOrder, ph.Pulse)
	if err != nil {
		return nil, errors.Wrap(err, "[ PacketHeader.Serialize ] Can't write Pulse")
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
	result := new(bytes.Buffer)
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
	err := binary.Read(data, defaultByteOrder, &pd.PulseNumer)
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
	result := new(bytes.Buffer)
	err := binary.Write(result, defaultByteOrder, pd.PulseNumer)
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
	result := new(bytes.Buffer)
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

	err = binary.Read(data, defaultByteOrder, &nb.claimType)
	if err != nil {
		return errors.Wrap(err, "[ NodeBroadcast.Deserialize ] Can't read claimType")
	}

	err = binary.Read(data, defaultByteOrder, &nb.length)
	if err != nil {
		return errors.Wrap(err, "[ NodeBroadcast.Deserialize ] Can't read length")
	}

	return nil
}

// Serialize implements interface method
func (nb *NodeBroadcast) Serialize() ([]byte, error) {
	result := new(bytes.Buffer)
	err := binary.Write(result, defaultByteOrder, nb.EmergencyLevel)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeBroadcast.Serialize ] Can't write EmergencyLevel")
	}

	err = binary.Write(result, defaultByteOrder, nb.claimType)
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeBroadcast.Serialize ] Can't write claimType")
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

	err = binary.Read(data, defaultByteOrder, &cpa.claimType)
	if err != nil {
		return errors.Wrap(err, "[ CapabilityPoolingAndActivation.Deserialize ] Can't read claimType")
	}

	err = binary.Read(data, defaultByteOrder, &cpa.length)
	if err != nil {
		return errors.Wrap(err, "[ CapabilityPoolingAndActivation.Deserialize ] Can't read length")
	}

	return nil
}

// Serialize implements interface method
func (cpa *CapabilityPoolingAndActivation) Serialize() ([]byte, error) {
	result := new(bytes.Buffer)
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

	err = binary.Write(result, defaultByteOrder, cpa.claimType)
	if err != nil {
		return nil, errors.Wrap(err, "[ CapabilityPoolingAndActivation.Serialize ] Can't write claimType")
	}

	err = binary.Write(result, defaultByteOrder, cpa.length)
	if err != nil {
		return nil, errors.Wrap(err, "[ CapabilityPoolingAndActivation.Serialize ] Can't write length")
	}

	return result.Bytes(), nil
}
