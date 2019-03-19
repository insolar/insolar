/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package object

import (
	"bytes"
	"encoding/binary"

	"github.com/insolar/insolar"
	"github.com/insolar/insolar/core"
	"github.com/ugorji/go/codec"
)

// SerializeType returns binary representation of provided type.
func SerializeType(id TypeID) []byte {
	buf := make([]byte, TypeIDSize)
	binary.BigEndian.PutUint32(buf, uint32(id))
	return buf
}

// DeserializeType returns type from provided binary representation.
func DeserializeType(buf []byte) TypeID {
	return TypeID(binary.BigEndian.Uint32(buf))
}

// SerializeRecord returns binary representation of provided record.
func SerializeRecord(rec VirtualRecord) []byte {
	typeBytes := SerializeType(TypeFromRecord(rec))
	buff := bytes.NewBuffer(typeBytes)
	enc := codec.NewEncoder(buff, &codec.CborHandle{})
	enc.MustEncode(rec)
	return buff.Bytes()
}

// DeserializeRecord returns record decoded from bytes.
func DeserializeRecord(buf []byte) VirtualRecord {
	t := DeserializeType(buf[:TypeIDSize])
	dec := codec.NewDecoderBytes(buf[TypeIDSize:], &codec.CborHandle{})
	rec := RecordFromType(t)
	dec.MustDecode(&rec)
	return rec
}

// CalculateIDForBlob calculate id for blob with using current pulse number
func CalculateIDForBlob(scheme core.PlatformCryptographyScheme, pulseNumber insolar.PulseNumber, blob []byte) *insolar.ID {
	hasher := scheme.IntegrityHasher()
	_, err := hasher.Write(blob)
	if err != nil {
		panic(err)
	}
	return core.NewRecordID(pulseNumber, hasher.Sum(nil))
}
