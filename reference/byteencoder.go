package reference

import (
	"bytes"
	"encoding/base64"
	"io"
	"strings"

	base58 "github.com/jbenet/go-base58"
	"github.com/pkg/errors"
)

type ByteEncodeFunc func(source io.ByteReader, builder *strings.Builder) error

func byteEncodeBase58(source io.ByteReader, builder *strings.Builder) error {
	buff := bytes.Buffer{}
	for b, err := source.ReadByte(); err == nil; b, err = source.ReadByte() {
		err := buff.WriteByte(b)
		if err != nil {
			return errors.Wrap(err, "failed to write base58 encoded data to string builder")
		}
	}
	_, err := builder.Write([]byte(base58.Encode(buff.Bytes())))
	return err
}

func byteEncodeBase64(source io.ByteReader, builder *strings.Builder) error {
	buff := bytes.Buffer{}
	for b, err := source.ReadByte(); err == nil; b, err = source.ReadByte() {
		buff.WriteByte(b)
	}
	encoder := base64.NewEncoder(base64.RawURLEncoding, builder)
	_, err := encoder.Write(buff.Bytes())
	if err != nil {
		return errors.Wrap(err, "failed to write base64 encoded data to string builder")
	}
	err = encoder.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close string builder")
	}
	return nil
}
