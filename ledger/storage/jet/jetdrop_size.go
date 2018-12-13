/*
 *    Copyright 2018 Insolar
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

package jet

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

func encode(data interface{}) []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	enc.MustEncode(data)
	return buf.Bytes()
}

// DropSize contains info about size of drop
type DropSize struct {
	JetID     core.RecordID
	PulseNo   core.PulseNumber
	DropSize  uint64
	Signature []byte
}

func (ds *DropSize) serializeDropSize() []byte {
	result := make([]byte, 0, 64)

	buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(buff, ds.DropSize)
	result = append(result, buff...)

	buff = make([]byte, 4)
	binary.LittleEndian.PutUint32(buff, uint32(ds.PulseNo))
	result = append(result, buff...)

	result = append(result, ds.JetID.Bytes()...)

	return result
}

// WriteHashData writes DropSize data to provided writer. This data is used to calculate DropSize's hash.
func (ds *DropSize) WriteHashData(w io.Writer) (int, error) {
	return w.Write(ds.serializeDropSize())
}

// DropSizeHistory is chain of drop sizes
type DropSizeHistory []DropSize

// DeserializeJetDropSizeHistory deserializes DropSizeHistory
func DeserializeJetDropSizeHistory(ctx context.Context, buff []byte) (DropSizeHistory, error) {
	inslogger.FromContext(ctx).Debugf("DeserializeJetDropSizeHistory starts ... ( buff len: %d)", len(buff))
	dec := codec.NewDecoder(bytes.NewReader(buff), &codec.CborHandle{})
	var dropSizes = DropSizeHistory{}

	err := dec.Decode(&dropSizes)
	if err != nil {
		return nil, errors.Wrapf(err, "[ DeserializeJetDropSizeHistory ] Can't decode DropSizeHistory")
	}

	return dropSizes, nil
}

// Bytes serializes DropSizeHistory
func (dropSizeHistory DropSizeHistory) Bytes(ctx context.Context) []byte {
	inslogger.FromContext(ctx).Debug("DropSizeHistory.Bytes starts ...")
	return encode(dropSizeHistory)
}
