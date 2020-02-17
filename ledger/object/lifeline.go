// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
