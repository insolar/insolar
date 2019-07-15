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

package serialization

import (
	"io"

	"github.com/pkg/errors"
)

const (
	claimHeaderLengthMask = 0x1F
	claimHeaderTypeShift  = 10
)

type ClaimType uint8

const (
	claimTypeEmpty   = ClaimType(0)
	claimTypeGeneric = ClaimType(1)
)

type ClaimHeader struct {
	TypeAndLength uint16 `insolar-transport:"header;[0-9]=length;[10-15]=header:ClaimType;group=Claims"` // [00-09] ByteLength [10-15] ClaimClass
	// actual payload
}

type GenericClaim struct {
	// ByteSize>=2
	ClaimHeader
	Payload []byte
}

type EmptyClaim struct {
	// ByteSize=2
	ClaimHeader `insolar-transport:"delimiter;ClaimType=0;length=header"`
}

type ClaimList struct {
	// ByteSize>=2
	Claims      []GenericClaim
	EndOfClaims EmptyClaim // ByteSize=2 - indicates end of claims
}

func newClaimHeader(t ClaimType, length int) ClaimHeader {
	var h ClaimHeader
	h.TypeAndLength = (h.TypeAndLength | uint16(t)<<claimHeaderTypeShift) | uint16(length)
	return h
}

func NewClaimList() ClaimList {
	return ClaimList{
		Claims:      make([]GenericClaim, 0),
		EndOfClaims: EmptyClaim{newClaimHeader(claimTypeEmpty, 0)},
	}
}

func NewGenericClaim(payload []byte) GenericClaim {
	return GenericClaim{
		ClaimHeader: newClaimHeader(claimTypeGeneric, len(payload)),
		Payload:     payload,
	}
}

func (cl *ClaimList) Push(claim GenericClaim) {
	cl.Claims = append(cl.Claims, claim)
}

func (ch *ClaimHeader) ClaimType() ClaimType {
	return ClaimType(ch.TypeAndLength >> claimHeaderTypeShift)
}

// Length returns claim length without header
func (ch *ClaimHeader) Length() int {
	return int(ch.TypeAndLength & claimHeaderLengthMask)
}

func (cl *ClaimList) SerializeTo(ctx SerializeContext, writer io.Writer) error {
	for _, c := range cl.Claims {
		err := c.SerializeTo(ctx, writer)
		if err != nil {
			return err
		}
	}
	return write(writer, cl.EndOfClaims)
}

func (cl *ClaimList) DeserializeFrom(ctx DeserializeContext, reader io.Reader) error {
	cl.EndOfClaims = EmptyClaim{newClaimHeader(claimTypeEmpty, 0)}
	for {
		header := ClaimHeader{}
		err := header.DeserializeFrom(ctx, reader)
		if err != nil {
			return ErrMalformedHeader(err)
		}
		switch header.ClaimType() {
		case claimTypeEmpty:
			return nil
		case claimTypeGeneric:
			if cl.Claims == nil {
				cl.Claims = make([]GenericClaim, 0)
			}

			claim := GenericClaim{header, make([]byte, header.Length())}
			limitReader := io.LimitReader(reader, int64(header.Length()))
			n, err := limitReader.Read(claim.Payload)
			if err != nil || n != header.Length() {
				return ErrPayloadLengthMismatch(int64(header.Length()), int64(n))
			}
			cl.Claims = append(cl.Claims, claim)
			return nil
		default:
			return ErrMalformedHeader(errors.New("unknown claim type"))
		}
	}
}

func (ch *ClaimHeader) SerializeTo(ctx SerializeContext, writer io.Writer) error {
	return write(writer, ch)
}

func (ch *ClaimHeader) DeserializeFrom(ctx DeserializeContext, reader io.Reader) error {
	return read(reader, ch)
}

func (c *GenericClaim) SerializeTo(ctx SerializeContext, writer io.Writer) error {
	err := c.ClaimHeader.SerializeTo(ctx, writer)
	if err != nil {
		return err
	}
	return write(writer, c.Payload)
}
