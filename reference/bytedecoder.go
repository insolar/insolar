// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package reference

import (
	"encoding/base64"
	"io"

	base58 "github.com/jbenet/go-base58"
	"github.com/pkg/errors"
)

type ByteDecoderFactory interface {
	GetByteDecoder(encodingName string) ByteDecodeFunc
	DefaultDecoder() ByteDecodeFunc
	LegacyDecoder() ByteDecodeFunc
}

type byteDecoderFactory struct {
	decoders map[string]ByteDecodeFunc
}

func byteDecodeBase58(s string, target io.ByteWriter) (stringRead int, err error) {
	bytes := base58.Decode(s)
	if len(s) > 0 && len(bytes) == 0 {
		return 0, errors.New("input string contains bad charachters")
	}

	for _, b := range bytes {
		err := target.WriteByte(b)
		if err != nil {
			return 0, errors.Wrap(err, "failed to write byte")
		}
	}
	return len(s), nil
}

func byteDecodeBase64(s string, target io.ByteWriter) (stringRead int, err error) {
	bytes, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return 0, err
	}

	for _, b := range bytes {
		err := target.WriteByte(b)
		if err != nil {
			return 0, errors.Wrap(err, "failed to write byte")
		}
	}
	return len(s), nil
}

func NewByteDecoderFactory() ByteDecoderFactory {
	return &byteDecoderFactory{
		decoders: map[string]ByteDecodeFunc{
			"base58": byteDecodeBase58,
			"base64": byteDecodeBase64,
		},
	}
}

func (f *byteDecoderFactory) GetByteDecoder(encodingName string) ByteDecodeFunc {
	return f.decoders[encodingName]
}

func (f *byteDecoderFactory) DefaultDecoder() ByteDecodeFunc {
	return f.decoders["base64"]
}

func (f *byteDecoderFactory) LegacyDecoder() ByteDecodeFunc {
	return f.decoders["base58"]
}
