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
	"errors"
	"fmt"
	"io"
	"strings"
)

type ByteDecodeFunc func(s string, target io.ByteWriter) (stringRead int, err error)

type IdentityDecoder func(base Holder, name string) *Global

type DecoderOptions uint8

const (
	AllowLegacy DecoderOptions = 1 << iota
	AllowRecords
	IgnoreParity
)

func DefaultDecoder() GlobalDecoder {
	return NewDefaultDecoder(AllowLegacy | AllowRecords)
}

type GlobalDecoder interface {
	Decode(ref string) (Global, error)
	DecodeHolder(ref string) (Holder, error)
}

type decoder struct {
	byteDecoderFactory ByteDecoderFactory
	legacyDecoder      ByteDecodeFunc
	defaultDecoder     ByteDecodeFunc

	nameDecoder IdentityDecoder
	options     DecoderOptions
}

func NewDecoder(options DecoderOptions, factory ByteDecoderFactory) GlobalDecoder {
	return &decoder{
		byteDecoderFactory: factory,
		legacyDecoder:      factory.LegacyDecoder(),
		defaultDecoder:     factory.DefaultDecoder(),

		options: options,
	}
}

func NewDefaultDecoder(options DecoderOptions) GlobalDecoder {
	return NewDecoder(options, NewByteDecoderFactory())
}

func (v decoder) Decode(ref string) (result Global, err error) {
	result, err = v.decode(ref)
	if err != nil {
		return Global{}, err
	}
	return result, nil
}

func (v decoder) DecodeHolder(ref string) (Holder, error) {
	result, err := v.decode(ref)
	if err != nil {
		return nil, err
	}
	return result.tryConvertToCompact(), nil
}

func (v decoder) decode(ref string) (result Global, err error) {
	schemaPos := strings.IndexRune(ref, ':')
	if schemaPos >= 0 {
		decoder, err := v.parseSchema(ref[:schemaPos], ref)
		if err != nil {
			return result, err
		}
		return v.parseReference(ref[schemaPos+1:], decoder)
	}

	// try to parse the legacy format
	if v.options&AllowLegacy != 0 && len(ref) >= 2*len(LegacyDomainName)+1 {
		domainPos := strings.IndexRune(ref, '.')
		if domainPos >= len(LegacyDomainName) && ref[domainPos+1:] == LegacyDomainName {
			result.addressLocal, err = v.parseLegacyAddress(ref, domainPos)
			if err == nil {
				result.convertToSelf()
			}
			return result, err
		}
	}

	return v.parseReference(ref, v.defaultDecoder)
}

func (v decoder) parseLegacyAddress(ref string, domainPos int) (resultLocal Local, _ error) {
	w := resultLocal.asWriter()
	_, err := v.legacyDecoder(ref[:domainPos], w)

	switch {
	case err != nil:
		break
	case !w.isFull():
		err = errors.New("insufficient length")
	case resultLocal.getScope() != 0: // there is no scope for legacy
		err = errors.New("invalid scope")
	case !resultLocal.canConvertToSelf():
		err = errors.New("invalid self reference")
	default:
		return resultLocal, nil
	}
	return resultLocal, fmt.Errorf("unable to parse legacy reference: ref=%v err=%v", ref, err)
}

func (v decoder) parseSchema(schema, refFull string) (ByteDecodeFunc, error) {
	parts := strings.Split(schema, "+")
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
			return nil, fmt.Errorf("unsupported schema: schema=%s, ref=%v", schema, refFull)
		}
	default:
		return nil, fmt.Errorf("invalid schema: schema=%s", schema)
	}
	decoder := v.byteDecoderFactory.GetByteDecoder(parts[0])
	if decoder == nil {
		return nil, fmt.Errorf("unknown decoder: schema=%s, decoder=%s, ref=%v", schema, parts[0], refFull)
	}
	return decoder, nil
}

func (v decoder) parseAuthority(ref string) (authority string, remaining string) {
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

func (v decoder) parseReference(refFull string, byteDecoder ByteDecodeFunc) (result Global, _ error) {
	if err := func() error {
		_, ref := v.parseAuthority(refFull)
		if len(ref) == 0 {
			return fmt.Errorf("empty reference body: ref=%v", refFull)
		}

		parityPos := strings.IndexRune(ref, '/')
		var parity []byte
		switch {
		case parityPos == 0:
			return fmt.Errorf("empty reference body: ref=%v", refFull)
		case parityPos > 0:
			encodedParity := ref[parityPos+1:]
			if encodedParity[0] != '2' {
				return fmt.Errorf("invalid parity prefix: ref=%v, parity=%s", refFull, encodedParity)
			}
			buf := bytes.NewBuffer(make([]byte, 0, LocalBinaryPulseAndScopeSize))
			_, err := byteDecoder(encodedParity[1:], buf)
			if err != nil {
				return fmt.Errorf("unable to decode parity: ref=%v, err=%v", refFull, err)
			}
			ref = ref[:parityPos]
			if v.options&IgnoreParity == 0 {
				parity = buf.Bytes()
			}
		}

		err := v.parseAddress(ref, byteDecoder, &result)

		if err == nil && parity != nil {
			err = CheckParity(&result, parity)
		}
		if err != nil {
			return fmt.Errorf("invalid reference: ref=%v err=%v", refFull, err)
		}
		return nil
	}(); err != nil {
		return Global{}, err
	}
	return result, nil
}

func (v decoder) parseAddress(ref string, byteDecoder ByteDecodeFunc, result *Global) error {

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
			resolveBase := &Global{}

			err := v.parseBinaryAddress(domainName, byteDecoder, &resolveBase.addressLocal)
			switch {
			case err == nil:
				if !resolveBase.tryConvertToSelf() {
					return errors.New("invalid reference base")
				}
			case err == errAliasedReference:
				if v.nameDecoder == nil {
					return errors.New("aliases are not allowed")
				}
				resolveBase = v.nameDecoder(nil, domainName)
				if resolveBase == nil {
					return errors.New("unknown domain alias")
				}
			default:
				return err
			}
			return v.parseAddressWithBase(localAddrName, resolveBase, byteDecoder, result)
		}
	default:
		return v.parseAddressWithBase(ref, &Global{}, byteDecoder, result)
	}
}

func (v decoder) parseAddressWithBase(name string, base *Global, byteDecoder ByteDecodeFunc, result *Global) error {
	err := v.parseBinaryAddress(name, byteDecoder, &result.addressLocal)

	switch {
	case err != nil:
		if err != errAliasedReference {
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

var errAliasedReference = errors.New("record reference alias")

func (v decoder) parseBinaryAddress(name string, byteDecoder ByteDecodeFunc, result *Local) error {

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
		if !w.isFull() {
			return errors.New("insufficient address length")
		}
	case '2', '3', '4', '5', '6', '7', '8', '9':
		return errors.New("unsupported address prefix")
	default:
		return errAliasedReference
	}
	return nil
}
