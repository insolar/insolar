// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package serialization

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"

	"github.com/pkg/errors"
)

const (
	packetMaxSize = 2048
	headerSize    = 16
	signatureSize = 64
)

type serializeSetter interface {
	setPayloadLength(uint16)
	setSignature(signature cryptkit.SignatureHolder)
}

type deserializeGetter interface {
	getPayloadLength() uint16
}

type packetContext struct {
	context.Context
	PacketHeaderAccessor

	header *Header

	fieldContext          FieldContext
	neighbourNodeID       insolar.ShortNodeID
	announcedJoinerNodeID insolar.ShortNodeID
}

func newPacketContext(ctx context.Context, header *Header) packetContext {
	ctx, _ = inslogger.WithFields(ctx, map[string]interface{}{
		"packet_flags":   fmt.Sprintf("%08b", header.PacketFlags),
		"payload_length": header.getPayloadLength(),
	})

	return packetContext{
		Context:              ctx,
		PacketHeaderAccessor: header,
		header:               header,
	}
}

func (pc *packetContext) InContext(ctx FieldContext) bool {
	return pc.fieldContext == ctx
}

func (pc *packetContext) SetInContext(ctx FieldContext) {
	pc.fieldContext = ctx
}

func (pc *packetContext) GetNeighbourNodeID() insolar.ShortNodeID {
	if pc.neighbourNodeID.IsAbsent() {
		panic("illegal value")
	}

	return pc.neighbourNodeID
}

func (pc *packetContext) SetNeighbourNodeID(nodeID insolar.ShortNodeID) {
	pc.neighbourNodeID = nodeID
}

func (pc *packetContext) GetAnnouncedJoinerNodeID() insolar.ShortNodeID {
	return pc.announcedJoinerNodeID
}

func (pc *packetContext) SetAnnouncedJoinerNodeID(nodeID insolar.ShortNodeID) {
	if nodeID.IsAbsent() {
		panic("illegal value")
	}

	pc.announcedJoinerNodeID = nodeID
}

type trackableWriter struct {
	io.Writer
	totalWrite int64
}

func newTrackableWriter(writer io.Writer) *trackableWriter {
	return &trackableWriter{Writer: writer}
}

func (w *trackableWriter) Write(p []byte) (int, error) {
	n, err := w.Writer.Write(p)
	w.totalWrite += int64(n)
	return n, err
}

type trackableReader struct {
	io.Reader
	totalRead int64
}

func newTrackableReader(reader io.Reader) *trackableReader {
	return &trackableReader{Reader: reader}
}

func (r *trackableReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.totalRead += int64(n)
	return n, err
}

type serializeContext struct {
	packetContext
	PacketHeaderModifier

	writer   *trackableWriter
	digester cryptkit.DataDigester
	signer   cryptkit.DigestSigner
	setter   serializeSetter

	buf1         [packetMaxSize]byte
	buf2         [packetMaxSize]byte
	bodyBuffer   *bytes.Buffer
	bodyTracker  *trackableWriter
	packetBuffer *bytes.Buffer
}

func newSerializeContext(ctx packetContext, writer *trackableWriter, digester cryptkit.DataDigester, signer cryptkit.DigestSigner, callback serializeSetter) *serializeContext {
	serializeCtx := &serializeContext{
		packetContext:        ctx,
		PacketHeaderModifier: ctx.header,

		writer:   writer,
		digester: digester,
		signer:   signer,
		setter:   callback,
	}

	serializeCtx.bodyBuffer = bytes.NewBuffer(serializeCtx.buf1[0:0:packetMaxSize])
	serializeCtx.bodyTracker = newTrackableWriter(serializeCtx.bodyBuffer)
	serializeCtx.packetBuffer = bytes.NewBuffer(serializeCtx.buf2[0:0:packetMaxSize])

	return serializeCtx
}

func (ctx *serializeContext) Write(p []byte) (int, error) {
	// Uncomment on debug. Too verbose
	// inslogger.FromContext(ctx).WithSkipFrameCount(3).Debugf("Serializing bytes %d", len(p))

	return ctx.bodyTracker.Write(p)
}

func (ctx *serializeContext) Finalize() (int64, error) {
	var (
		totalWrite int64
		err        error
	)

	payloadLength := ctx.bodyTracker.totalWrite + headerSize + signatureSize

	if payloadLength > int64(math.MaxUint16) { // Will overflow
		return totalWrite, errors.New("payload too big")
	}
	ctx.setter.setPayloadLength(uint16(payloadLength))

	// TODO: set receiver id = 0
	if err := ctx.header.SerializeTo(ctx, ctx.packetBuffer); err != nil {
		return totalWrite, ErrMalformedHeader(err)
	}

	if _, err := ctx.bodyBuffer.WriteTo(ctx.packetBuffer); err != nil {
		return totalWrite, ErrMalformedPacketBody(err)
	}

	readerForSignature := bytes.NewReader(ctx.packetBuffer.Bytes())
	digest := ctx.digester.GetDigestOf(readerForSignature)
	signedDigest := digest.SignWith(ctx.signer)
	signature := signedDigest.GetSignatureHolder()
	ctx.setter.setSignature(signature)

	if _, err := signature.WriteTo(ctx.packetBuffer); err != nil {
		return totalWrite, ErrMalformedPacketSignature(err)
	}

	if totalWrite, err = ctx.packetBuffer.WriteTo(ctx.writer); totalWrite != payloadLength {
		return totalWrite, ErrPayloadLengthMismatch(payloadLength, totalWrite)
	}

	return totalWrite, err
}

type deserializeContext struct {
	packetContext

	reader *trackableReader
	getter deserializeGetter
}

func newDeserializeContext(ctx packetContext, reader *trackableReader, callback deserializeGetter) *deserializeContext {
	deserializeCtx := &deserializeContext{
		packetContext: ctx,

		reader: reader,
		getter: callback,
	}
	return deserializeCtx
}

func (ctx *deserializeContext) Read(p []byte) (int, error) {
	return ctx.reader.Read(p)
}

func (ctx *deserializeContext) Finalize() (int64, error) {
	if payloadLength := int64(ctx.getter.getPayloadLength()); payloadLength != ctx.reader.totalRead {
		return ctx.reader.totalRead, ErrPayloadLengthMismatch(payloadLength, ctx.reader.totalRead)
	}

	return ctx.reader.totalRead, nil
}
