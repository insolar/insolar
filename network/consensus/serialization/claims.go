// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
