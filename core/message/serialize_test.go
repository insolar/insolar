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

package message

import (
	"bytes"
	"context"
	"github.com/insolar/insolar/core"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

func TestSerializeSigned(t *testing.T) {
	msg := &SetRecord{
		Record: []byte{0x0A},
	}
	signMsgIn := &Parcel{
		Msg:       msg,
		Signature: nil,
	}

	signMsgOut, err := DeserializeParcel(bytes.NewBuffer(ParcelToBytes(signMsgIn)))
	require.NoError(t, err)

	require.Equal(t, signMsgIn, signMsgOut)
	require.Equal(t, signMsgIn.Message(), signMsgOut.Message())
}

func TestSerializeSignedFail(t *testing.T) {
	msg := &SetRecord{
		Record: []byte{0x0A},
	}

	signMsgIn := &Parcel{
		Msg:       msg,
		Signature: nil,
	}

	signMsgOut, err := Deserialize(bytes.NewBuffer(ParcelToBytes(signMsgIn)))
	require.Error(t, err)
	require.Nil(t, signMsgOut)
}

func TestSerializeSignedWithContext(t *testing.T) {
	msg := &SetRecord{
		Record: []byte{0x0A},
	}
	ctxIn := context.Background()
	traceid := "testtraceid"
	ctxIn = inslogger.ContextWithTrace(context.Background(), traceid)
	ctxIn = instracer.SetBaggage(ctxIn, instracer.Entry{Key: "traceid", Value: traceid})

	signMsgIn := &Parcel{
		Msg:           msg,
		Signature:     nil,
		ServiceData:   ServiceData{
			TraceSpanData: instracer.MustSerialize(ctxIn),
			LogTraceID:    inslogger.TraceID(ctxIn),
			LogLevel:      core.DebugLevel,
		},
	}

	signMsgOut, err := DeserializeParcel(bytes.NewBuffer(ParcelToBytes(signMsgIn)))
	require.NoError(t, err)

	ctxOut := signMsgOut.Context(context.Background())
	require.Equal(t, inslogger.TraceID(ctxIn), traceid)
	require.Equal(t, inslogger.TraceID(ctxIn), inslogger.TraceID(ctxOut))
	require.Equal(t, instracer.GetBaggage(ctxOut), instracer.GetBaggage(ctxIn))
	require.Equal(t, core.NoLevel, inslogger.GetLoggerLevel(ctxIn))
	require.Equal(t, core.DebugLevel, inslogger.GetLoggerLevel(ctxOut))
}
