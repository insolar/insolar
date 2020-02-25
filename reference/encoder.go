// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package reference

import (
	"bytes"
	"strings"

	"github.com/pkg/errors"
)

type IdentityEncoder func(ref *Global) (domain, object string)

type EncoderOptions uint8

const (
	Parity EncoderOptions = 1 << iota
	EncodingSchema
	FormatSchema

	NilRef   = "<nil>" // non-parsable
	SchemaV1 = "insolar"
)

type Encoder interface {
	Encode(ref *Global) (string, error)
	EncodeToBuilder(ref *Global, b *strings.Builder) error
	EncodeRecord(rec *Local) (string, error)
}

func DefaultEncoder() Encoder {
	return NewBase64Encoder(FormatSchema)
}

func Base58Encoder() Encoder {
	return NewBase58Encoder(FormatSchema)
}

type encoder struct {
	nameEncoder     IdentityEncoder
	byteEncoder     ByteEncodeFunc
	byteEncoderName string
	authorityName   string
	options         EncoderOptions
}

func NewBase58Encoder(opts EncoderOptions) Encoder {
	return &encoder{
		nameEncoder:     nil,
		byteEncoder:     byteEncodeBase58,
		byteEncoderName: "base58",
		authorityName:   "",
		options:         opts,
	}
}

func NewBase64Encoder(opts EncoderOptions) Encoder {
	return &encoder{
		nameEncoder:     nil,
		byteEncoder:     byteEncodeBase64,
		byteEncoderName: "base64",
		authorityName:   "",
		options:         opts | FormatSchema,
	}
}

func (v encoder) Encode(ref *Global) (string, error) {
	b := strings.Builder{}
	err := v.EncodeToBuilder(ref, &b)
	return b.String(), err
}

func (v encoder) EncodeToBuilder(ref *Global, b *strings.Builder) error {
	if ref == nil {
		b.WriteString(NilRef)
		return nil
	}

	v.appendPrefix(b)

	if ref.IsEmpty() {
		b.WriteString("0")
	}
	if ref.IsRecordScope() {
		return v.encodeRecord(&ref.addressLocal, b)
	}

	var domainName, objectName string

	if v.nameEncoder != nil {
		domainName, objectName = v.nameEncoder(ref)
	}

	if objectName != "" {
		if IsReservedName(objectName) || !IsValidObjectName(objectName) {
			return errors.Errorf("illegal object name from IdentityEncoder: ref=%v, domain='%s', object='%s'", ref, domainName, objectName)
		}
		b.WriteString(objectName)
	} else {
		err := v.encodeBinary(&ref.addressLocal, b)
		if err != nil {
			return err
		}
	}

	switch {
	case domainName != "":
		if IsReservedName(domainName) || !IsValidDomainName(domainName) {
			return errors.Errorf("illegal domain name from IdentityEncoder: ref=%v, domain='%s', object='%s'", ref, domainName, objectName)
		}
		b.WriteByte('.')
		b.WriteString(domainName)
	case ref.IsSelfScope():
		// nothing
	default:
		b.WriteByte('.')
		err := v.encodeBinary(&ref.addressBase, b)
		if err != nil {
			return err
		}
	}

	if v.options&Parity != 0 {
		parity := ref.GetParity()
		if len(parity) > 0 {
			b.WriteByte('/')
			err := v.byteEncoder(bytes.NewReader(parity), b)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (v encoder) appendPrefix(b *strings.Builder) {

	switch v.options & (EncodingSchema | FormatSchema) {
	case EncodingSchema | FormatSchema:
		b.WriteString(SchemaV1 + "+" + v.byteEncoderName + ":")
	case FormatSchema:
		b.WriteString(SchemaV1 + ":")
	case EncodingSchema:
		b.WriteString(v.byteEncoderName + ":")
	}

	if len(v.authorityName) > 0 {
		b.WriteString("//")
		b.WriteString(v.authorityName)
		b.WriteByte('/')
	}
}

func (v encoder) encodeBinary(rec *Local, b *strings.Builder) error {
	if rec.IsEmpty() {
		b.WriteByte('0')
		return nil
	}
	pn := rec.GetPulseNumber()
	switch {
	case pn.IsTimePulse():
		b.WriteByte('1')
		// full encode
		err := v.byteEncoder(rec.AsReader(), b)
		if err != nil {
			return err
		}

	case pn.IsSpecial():
		b.WriteString("0")

		limit := len(rec.hash) - 1
		for ; limit >= 0 && rec.hash[limit] == 0; limit-- {
		}
		limit += 1 + LocalBinaryPulseAndScopeSize

		err := v.byteEncoder(rec.asReader(uint8(limit)), b)
		if err != nil {
			return err
		}
	default:
		panic("unexpected")
	}
	return nil
}

func (v encoder) encodeRecord(rec *Local, b *strings.Builder) error {
	if rec.IsEmpty() {
		b.WriteString("0." + RecordDomainName)
		return nil
	}
	if rec.getScope() != 0 {
		panic("illegal value")
	}
	err := v.encodeBinary(rec, b)
	if err != nil {
		return err
	}
	b.WriteString("." + RecordDomainName)

	return nil
}

func (v encoder) EncodeRecord(rec *Local) (string, error) {
	if rec == nil {
		return NilRef, nil
	}

	b := strings.Builder{}
	v.appendPrefix(&b)
	err := v.encodeRecord(rec, &b)
	return b.String(), err
}
