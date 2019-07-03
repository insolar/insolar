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
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/insolar/insolar/network/consensus/common"
	"github.com/pkg/errors"
)

type SerializerTo interface {
	SerializeTo(writer io.Writer, signer common.DataSigner) error
}

var defaultByteOrder = binary.BigEndian

const (
	fieldBufSize  = 2048
	packetBufSize = 2048
)

func serializeTo(writer io.Writer, signer common.DataSigner, data interface{}) error {
	var fieldInternalBuffer [fieldBufSize]byte

	checksumBuffer := &bytes.Buffer{}
	fieldBuf := bytes.NewBuffer(fieldInternalBuffer[0:0:fieldBufSize])

	v := reflect.ValueOf(data)
	vt := v.Type()

	fmt.Println(vt.NumField())

	for i := 0; i < vt.NumField(); i++ {

		fvt := vt.Field(i)
		fv := v.Field(i)

		// Skip nil pointer fields since it's optional
		if fv.Kind() == reflect.Ptr && fv.IsNil() {
			if !isOptional(fvt) {
				return errors.Errorf("Invalid nil field: %s.%s", vt.Name(), fvt.Name)
			}

			continue
		}

		if shouldGenerateSignature(fvt) && signer != nil {
			sd := signer.GetSignOfData(checksumBuffer)
			sigBytes := sd.GetSignature().AsBytes()
			bits := common.NewBits512FromBytes(sigBytes)
			fv.Set(reflect.ValueOf(bits))
		}

		err := writeValue(fv, fieldBuf, signer)
		if err != nil {
			return err
		}

		b := fieldBuf.Bytes()

		_, err = checksumBuffer.Write(b)
		if err != nil {
			return err
		}

		if !shouldIgnoreInSerialization(fvt) {
			_, err := writer.Write(b)
			if err != nil {
				return err
			}
		}

		fieldBuf.Reset()
	}

	return nil
}

func writeValue(fv reflect.Value, writer io.Writer, signer common.DataSigner) error {
	var err error

	val := fv.Interface()
	switch v := val.(type) {
	case SerializerTo:
		err = v.SerializeTo(writer, signer)
	default:
		err = binary.Write(writer, defaultByteOrder, val)
	}

	return err
}

func shouldIgnoreInSerialization(field reflect.StructField) bool {
	tag, ok := field.Tag.Lookup("insolar-transport")
	if !ok {
		return false
	}

	return strings.Contains(tag, "ignore=send")
}

func shouldGenerateSignature(field reflect.StructField) bool {
	tag, ok := field.Tag.Lookup("insolar-transport")
	if !ok {
		return false
	}

	return strings.Contains(tag, "generate=signature")
}

func isOptional(field reflect.StructField) bool {
	tag, ok := field.Tag.Lookup("insolar-transport")
	if !ok {
		return false
	}

	return strings.Contains(tag, "optional=")
}
