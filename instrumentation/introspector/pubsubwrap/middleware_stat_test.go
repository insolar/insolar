// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pubsubwrap

import (
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiddlewareStat(t *testing.T) {
	mi := NewMessageStatByType()
	expected := []struct {
		payloadType payload.Type
		payload     payload.Payload
		count       int
	}{
		{
			payload.TypeGetObject,
			&payload.GetObject{Polymorph: uint32(payload.TypeGetObject)},
			10,
		},
		{
			payload.TypeGetCode,
			&payload.GetCode{Polymorph: uint32(payload.TypeGetCode)},
			4,
		},
	}
	for _, ex := range expected {
		b, err := ex.payload.Marshal()
		require.NoError(t, err, "payload should be marshaled w/o errors")
		var meta payload.Meta
		meta.Payload = b
		metaBytes, err := meta.Marshal()
		require.NoError(t, err, "meta should be marshaled w/o errors")
		msg := &message.Message{
			Payload: metaBytes,
		}
		for i := 0; i < ex.count; i++ {
			mi.Filter(msg)
		}
	}

	stat, err := mi.GetMessagesStat(nil, nil)
	require.NoError(t, err, "GetMessagesStat should not failed")
	require.Equal(t, len(expected), len(stat.Counters), "expects statistic for the same types count")

	for _, ex := range expected {
		typ := ex.payloadType.String()
		assert.Equalf(t, int64(ex.count), mi.stat[typ], "check %v stat", typ)
	}
}
