/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted (subject to the limitations in the disclaimer below) provided that
 * the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *  * Neither the name of Insolar Technologies nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
 * BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
 * CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING,
 * BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package packets

import (
	"github.com/insolar/insolar/core"
)

//go:generate stringer -type=PacketType
type PacketType uint8

const (
	Phase1 = PacketType(iota + 1)
	Phase2
	Phase3
)

const HashLength = 64
const SignatureLength = 66
const PublicKeyLength = 66
const ReferenceLength = 64

// ------------------------------PACKET HEADER------------------------------

type PacketHeader struct {
	PacketT    PacketType
	HasRouting bool
	//-----------------
	f01   bool
	f00   bool
	Pulse uint32
	//-----------------
	OriginNodeID core.ShortNodeID
	TargetNodeID core.ShortNodeID
}

// ------------------------------PULSE DATA------------------------------

// PulseDataExt is a pulse data extension.
type PulseDataExt struct {
	NextPulseDelta uint16
	PrevPulseDelta uint16
	OriginID       [16]byte
	EpochPulseNo   uint32
	PulseTimestamp uint32
	Entropy        core.Entropy
}

// PulseData is a pulse data.
type PulseData struct {
	PulseNumber uint32
	Data        PulseDataExt
}

type NodePulseProof struct {
	NodeStateHash [HashLength]byte
	NodeSignature [SignatureLength]byte
}

func (npp *NodePulseProof) StateHash() []byte {
	return npp.NodeStateHash[:]
}

func (npp *NodePulseProof) Signature() []byte {
	return npp.NodeSignature[:]
}
