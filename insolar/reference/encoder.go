//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package reference

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

type Encoder struct {
	nameEncoder     IdentityEncoder
	byteEncoder     ByteEncodeFunc
	byteEncoderName string
	authorityName   string
	options         EncoderOptions
}

func (v Encoder) Encode(ref *Global) string {
	b := strings.Builder{}
	v.EncodeToBuilder(ref, &b)
	return b.String()
}

func (v Encoder) EncodeToBuilder(ref *Global, b *strings.Builder) {
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
			v.byteEncoder(bytes.NewReader(parity), b)
		}
	}
}

func (v Encoder) appendPrefix(b *strings.Builder) {

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

func (v Encoder) encodeBinary(rec *Local, b *strings.Builder) {
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

func (v Encoder) encodeRecord(rec *Local, b *strings.Builder) {
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

func (v Encoder) EncodeRecord(rec *Local) string {
	if rec == nil {
		return NilRef
	}
	b := strings.Builder{}
	v.appendPrefix(&b)
	v.encodeRecord(rec, &b)
	return b.String()
}
