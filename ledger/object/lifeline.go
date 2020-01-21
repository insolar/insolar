// Copyright 2020 Insolar Network Ltd.
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

package object

import (
	"github.com/insolar/insolar/insolar/record"
)

// EncodeLifeline converts lifeline index into binary format.
func EncodeLifeline(index record.Lifeline) []byte {
	res, err := index.Marshal()
	if err != nil {
		panic(err)
	}

	return res
}

// MustDecodeLifeline converts byte array into lifeline index struct.
func MustDecodeLifeline(buff []byte) (index record.Lifeline) {
	idx, err := DecodeLifeline(buff)
	if err != nil {
		panic(err)
	}

	return idx
}

// DecodeLifeline converts byte array into lifeline index struct.
func DecodeLifeline(buff []byte) (record.Lifeline, error) {
	lfl := record.Lifeline{}
	err := lfl.Unmarshal(buff)
	return lfl, err
}

// CloneLifeline returns copy of argument idx value.
func CloneLifeline(idx record.Lifeline) record.Lifeline {
	if idx.LatestState != nil {
		tmp := *idx.LatestState
		idx.LatestState = &tmp
	}

	if idx.LatestRequest != nil {
		r := *idx.LatestRequest
		idx.LatestRequest = &r
	}

	if idx.EarliestOpenRequest != nil {
		tmp := *idx.EarliestOpenRequest
		idx.EarliestOpenRequest = &tmp
	}

	return idx
}
