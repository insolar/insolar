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
	"errors"
	"fmt"
	"io"
	"strings"
)

type ByteDecodeFunc func(s string, target io.ByteWriter) (stringRead int, err error)
type IdentityDecoder func(base *Global, name string) *Global

type ByteDecoderFactory interface {
	GetByteDecoder(encodingName string) ByteDecodeFunc
}

type DecoderOptions uint8

const (
	AllowLegacy DecoderOptions = 1 << iota
	AllowRecords
	IgnoreParity
)

type ReferenceDecoder struct {
	byteDecoderFactory ByteDecoderFactory
	legacyDecoder      ByteDecodeFunc
	defaultDecoder     ByteDecodeFunc

	nameDecoder IdentityDecoder
	options     DecoderOptions
}

func (v ReferenceDecoder) Decode(ref string) (Global, error) {
	schemaPos := strings.IndexRune(ref, ':')
	if schemaPos >= 0 {
		decoder, err := v.parseSchema(ref[:schemaPos], ref)
		if err != nil {
			return Global{}, err
		}
		return v.parseReference(ref[schemaPos+1:], decoder)
	}

	// try to parse the legacy format
	if v.options&AllowLegacy != 0 && len(ref) >= 2*len(LegacyDomainName)+1 {
		domainPos := strings.IndexRune(ref, '.')
		if domainPos > len(LegacyDomainName) && ref[domainPos+1:] == LegacyDomainName {
			return v.parseLegacyAddress(ref, domainPos)
		}
	}

	return v.parseReference(ref, v.defaultDecoder)
}

func (v ReferenceDecoder) parseLegacyAddress(ref string, domainPos int) (Global, error) {
	var result Global

	w := result.addressLocal.asWriter()
	_, err := v.legacyDecoder(ref[:domainPos], w)

	switch {
	case err != nil:
		break
	case !w.isFull():
		err = errors.New("insufficient length")
	case result.addressLocal.getScope() != 0: // there is no scope for legacy
		err = errors.New("invalid scope")
	case !result.tryConvertToSelf():
		err = errors.New("invalid self reference")
	default:
		return result, nil
	}
	return result, fmt.Errorf("unable to parse legacy reference, %s: ref=%s", err.Error(), ref)
}

func (v ReferenceDecoder) parseSchema(scheme, refFull string) (ByteDecodeFunc, error) {
	parts := strings.Split(scheme, "+")
	switch len(parts) {
	case 1:
		if parts[0] == SchemaV1 {
			return v.defaultDecoder, nil
		}
	case 2:
		switch {
		case parts[0] == SchemaV1:
			parts[0] = parts[1]
		case parts[1] == SchemaV1:
			//
		default:
			return nil, fmt.Errorf("unsupported scheme: scheme=%s, ref=%s", scheme, refFull)
		}
	default:
		return nil, fmt.Errorf("invalid scheme: scheme=%s", scheme)
	}
	decoder := v.byteDecoderFactory.GetByteDecoder(parts[0])
	if decoder == nil {
		return nil, fmt.Errorf("unknown decoder: scheme=%s, decoder=%s, ref=%s", scheme, parts[0], refFull)
	}
	return decoder, nil
}

func (v ReferenceDecoder) parseAuthority(ref string) (authority string, remaining string) {
	if len(ref) < 3 || ref[:2] != "//" {
		return "", ref
	}
	ref = ref[2:]

	pos := strings.IndexRune(ref, '/')
	if pos < 0 {
		return ref, ""
	}

	return ref[:pos], ref[pos+1:]
}

func (v ReferenceDecoder) parseReference(refFull string, byteDecoder ByteDecodeFunc) (Global, error) {
	_, ref := v.parseAuthority(refFull)
	if len(ref) == 0 {
		return Global{}, fmt.Errorf("empty reference body: ref=%s", refFull)
	}

	parityPos := strings.IndexRune(ref, '/')
	var parity []byte
	switch {
	case parityPos == 0:
		return Global{}, fmt.Errorf("empty reference body: ref=%s", refFull)
	case parityPos > 0:
		buf := bytes.NewBuffer(make([]byte, 0, pulseAndScopeSize))
		_, err := byteDecoder(ref[parityPos+1:], buf)
		if err != nil {
			return Global{}, fmt.Errorf("unable to decode parity: ref=%s, err=%v", refFull, err)
		}
		ref = ref[:parityPos]
		if v.options&IgnoreParity == 0 {
			parity = buf.Bytes()
		}
	}

	var result Global
	err := v.parseAddress(ref, byteDecoder, &result)

	if err == nil && parity != nil {
		err = result.CheckParity(parity)
	}
	if err == nil {
		return result, nil
	}

	return result, fmt.Errorf("invalid reference, %s: ref=%s", err.Error(), refFull)
}

func (v ReferenceDecoder) parseAddress(ref string, byteDecoder ByteDecodeFunc, result *Global) error {

	domainPos := strings.IndexRune(ref, '.')
	switch {
	case domainPos == 0:
		return errors.New("empty reference body")
	case domainPos > 0:
		domainName := ref[domainPos+1:]
		localAddrName := ref[:domainPos]
		switch domainName {
		case "":
			return errors.New("empty domain name")
		case RecordDomainName:
			if v.options&AllowRecords == 0 {
				return errors.New("record reference is not allowed")
			}
			return v.parseBinaryAddress(localAddrName, byteDecoder, &result.addressLocal)
		case LegacyDomainName:
			return errors.New("legacy domain name")
		default:
			if v.nameDecoder == nil {
				return errors.New("aliases are not allowed")
			}
			resolveBase := v.nameDecoder(nil, domainName)
			if resolveBase == nil {
				return errors.New("unknown domain alias")
			}
			return v.parseAddressWithBase(localAddrName, resolveBase, byteDecoder, result)
		}
	default:
		return v.parseAddressWithBase(ref, &Global{}, byteDecoder, result)
	}
}

func (v ReferenceDecoder) parseAddressWithBase(name string, base *Global, byteDecoder ByteDecodeFunc, result *Global) error {
	err := v.parseBinaryAddress(name, byteDecoder, &result.addressLocal)

	switch {
	case err != nil:
		if err != aliasedReferenceError {
			return err
		}
		if v.nameDecoder == nil {
			return errors.New("aliases are not allowed")
		}
		resolved := v.nameDecoder(base, name)
		if resolved == nil {
			return errors.New("unknown object alias")
		}
		*result = *resolved
	case base.IsEmpty():
		if result.IsEmpty() {
			return nil
		}
		if !result.tryConvertToSelf() {
			return errors.New("invalid self reference")
		}
	case result.addressLocal.IsEmpty():
		return errors.New("empty local part of reference")
	default:
		if !result.tryApplyBase(base) {
			return errors.New("scope mismatch between base and local parts of address")
		}
	}
	return nil
}

var aliasedReferenceError = errors.New("record reference alias")

func (v ReferenceDecoder) parseBinaryAddress(name string, byteDecoder ByteDecodeFunc, result *Local) error {

	switch name[0] {
	case '0':
		if len(name) == 1 {
			return nil
		}
		_, err := byteDecoder(name[1:], result.asWriter())
		if err != nil {
			return err
		}
	case '1':
		w := result.asWriter()
		_, err := byteDecoder(name[1:], w)
		if err != nil {
			return err
		}
		if w.isFull() {
			return errors.New("insufficient address length")
		}
	case '2', '3', '4', '5', '6', '7', '8', '9':
		return errors.New("unsupported address prefix")
	default:
		return aliasedReferenceError
	}
	return nil
}
