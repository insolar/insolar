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
