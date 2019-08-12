///
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
///

package ref

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type ByteEncodeFunc func(source io.ByteReader, b *strings.Builder)
type IdentityEncoder func(ref *Global) (domain, object string)

const NilRef = "<nil>" //non-parsable

type EncoderOptions uint8

const (
	Parity EncoderOptions = 1 << iota
	EncodingSchema
	FormatSchema
)

const SchemaV1 = "insolarv1"

type ReferenceEncoder struct {
	nameEncoder     IdentityEncoder
	byteEncoder     ByteEncodeFunc
	byteEncoderName string
	authorityName   string
	options         EncoderOptions
}

func (v ReferenceEncoder) Encode(ref *Global) string {
	b := strings.Builder{}
	v.EncodeToBuilder(ref, &b)
	return b.String()
}

func (v ReferenceEncoder) EncodeToBuilder(ref *Global, b *strings.Builder) {
	if ref == nil {
		b.WriteString(NilRef)
		return
	}

	v.appendPrefix(b)

	if ref.IsEmpty() {
		b.WriteString("0")
	}
	if ref.IsRecordScope() {
		v.encodeRecord(&ref.addressLocal, b)
		return
	}

	var domainName, objectName string

	if v.nameEncoder != nil {
		domainName, objectName = v.nameEncoder(ref)
	}

	if objectName != "" {
		if IsReservedName(objectName) || !IsValidObjectName(objectName) {
			panic(fmt.Errorf("illegal object name from IdentityEncoder: ref=%v, domain='%s', object='%s'", ref, domainName, objectName))
		}
		b.WriteString(objectName)
	} else {
		v.encodeBinary(&ref.addressLocal, b)
	}

	switch {
	case domainName != "":
		if IsReservedName(domainName) || !IsValidDomainName(domainName) {
			panic(fmt.Errorf("illegal domain name from IdentityEncoder: ref=%v, domain='%s', object='%s'", ref, domainName, objectName))
		}
		b.WriteByte('.')
		b.WriteString(domainName)
	case ref.IsSelfScope():
		// nothing
	default:
		b.WriteByte('.')
		v.encodeBinary(&ref.addressBase, b)
	}

	if v.options&Parity != 0 {
		parity := ref.GetParity()
		if len(parity) > 0 {
			b.WriteByte('/')
			v.byteEncoder(bytes.NewBuffer(parity), b)
		}
	}
}

func (v ReferenceEncoder) appendPrefix(b *strings.Builder) {

	if v.options&(EncodingSchema|FormatSchema) != 0 {
		b.WriteString(v.byteEncoderName)
		if v.options&FormatSchema != 0 {
			b.WriteString("+" + SchemaV1)
		}
		b.WriteByte(':')
	}

	if len(v.authorityName) > 0 {
		b.WriteString("//")
		b.WriteString(v.authorityName)
		b.WriteByte('/')
	}
}

func (v ReferenceEncoder) encodeBinary(rec *Local, b *strings.Builder) {
	if rec.IsEmpty() {
		b.WriteByte('0')
		return
	}
	pn := rec.GetPulseNumber()
	switch {
	case pn.IsTimePulse():
		b.WriteByte('1')
		//full encode
		v.byteEncoder(rec.AsReader(), b)

	case pn.IsSpecial():
		b.WriteString("0")

		limit := len(rec.hash) - 1
		for ; limit >= 0 && rec.hash[limit] == 0; limit-- {
		}
		limit += 1 + pulseAndScopeSize

		v.byteEncoder(rec.asReader(uint8(limit)), b)
	default:
		panic("unexpected")
	}
}

func (v ReferenceEncoder) encodeRecord(rec *Local, b *strings.Builder) {
	if rec.IsEmpty() {
		b.WriteString("0." + RecordDomainName)
		return
	}
	if rec.getScope() != 0 {
		panic("illegal value")
	}
	v.encodeBinary(rec, b)
	b.WriteString("." + RecordDomainName)
}

func (v ReferenceEncoder) EncodeRecord(rec *Local) string {
	if rec == nil {
		return NilRef
	}
	b := strings.Builder{}
	v.appendPrefix(&b)
	v.encodeRecord(rec, &b)
	return b.String()
}
