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

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

type DropSizeData struct {
	JetID    core.RecordID
	PulseNo  core.PulseNumber
	DropSize uint64
}

func encode(data interface{}) []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	enc.MustEncode(data)
	return buf.Bytes()
}

func (dsd *DropSizeData) Bytes(ctx context.Context) []byte {
	inslogger.FromContext(ctx).Debug("DropSizeData.Bytes starts ...")
	return encode(dsd)
}

type DropSize struct {
	SizeData  DropSizeData
	Signature []byte
}

const MaxLenJetDropSizeList = 10

type DropSizeList []DropSize

func DeserializeJetDropSizeList(ctx context.Context, buff []byte) (DropSizeList, error) {
	inslogger.FromContext(ctx).Debugf("DeserializeJetDropSizeList starts ... ( buff len: %d)", len(buff))
	dec := codec.NewDecoder(bytes.NewReader(buff), &codec.CborHandle{})
	var dropSizes = DropSizeList{}

	err := dec.Decode(&dropSizes)
	if err != nil {
		return nil, errors.Wrapf(err, "[ DeserializeJetDropSizeList ] Can't decode DropSizeList")
	}

	return dropSizes, nil
}

func (dropSizeList DropSizeList) Bytes(ctx context.Context) []byte {
	inslogger.FromContext(ctx).Debug("DropSizeList.Bytes starts ...")
	return encode(dropSizeList)
}
