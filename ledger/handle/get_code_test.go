///
// Copyright 2019 Insolar Technologies GmbH
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
///

package handle

import (
	"context"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/ledger/proc"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestGetCode_Present(t *testing.T) {
	jetID := gen.JetID()
	codeRef := gen.Reference()
	parcel := testutils.NewParcelMock(t)
	parcel.MessageFunc = func() insolar.Message {
		return &message.GetCode{
			Code: codeRef,
		}
	}
	msg := bus.Message{
		Parcel: parcel,
	}
	gc := GetCode{
		dep: &proc.Dependencies{
			FetchJet: func(jet *proc.FetchJet) *proc.FetchJet {
				return jet
			},
			GetCode: func(code *proc.GetCode) *proc.GetCode {
				return code
			},
		},
		Message: msg,
	}
	fl := flow.NewFlowMock(t)
	fl.ProcedureFunc = func(ctx context.Context, pr flow.Procedure) error {
		if fetchJet, ok := pr.(*proc.FetchJet); ok {
			require.Equal(t, "TypeGetCode", fetchJet.Parcel.Message().Type().String())
			fetchJet.Result.Jet = jetID
			return nil
		} else if getCode, ok := pr.(*proc.GetCode); ok {
			require.Equal(t, getCode.JetID, jetID)
			require.Equal(t, getCode.Code, codeRef)
			require.Equal(t, getCode.Message, msg)
			return nil
		}
		t.Fatal("you shouldn't be here")
		return nil
	}
	ctx := context.Background()
	err := gc.Present(ctx, fl)
	require.NoError(t, err)
}

func TestGetCode_Present_MissJet(t *testing.T) {
	// jetID := gen.JetID()
	codeRef := gen.Reference()
	parcel := testutils.NewParcelMock(t)
	parcel.MessageFunc = func() insolar.Message {
		return &message.GetCode{
			Code: codeRef,
		}
	}
	replyCh := make(chan bus.Reply)
	msg := bus.Message{
		Parcel:  parcel,
		ReplyTo: replyCh,
	}
	gc := GetCode{
		dep: &proc.Dependencies{
			FetchJet: func(jet *proc.FetchJet) *proc.FetchJet {
				return jet
			},
			GetCode: func(code *proc.GetCode) *proc.GetCode {
				return code
			},
		},
		Message: msg,
	}
	fl := flow.NewFlowMock(t)
	fl.ProcedureFunc = func(ctx context.Context, pr flow.Procedure) error {
		if fetchJet, ok := pr.(*proc.FetchJet); ok {
			require.Equal(t, "TypeGetCode", fetchJet.Parcel.Message().Type().String())
			fetchJet.Result.Miss = true
			return nil
		} else if reply, ok := pr.(*proc.ReturnReply); ok {
			var expected (chan<- bus.Reply) = replyCh
			require.Equal(t, expected, reply.ReplyTo)
			return nil
		}
		t.Fatal("you shouldn't be here")
		return nil
	}
	ctx := context.Background()
	err := gc.Present(ctx, fl)
	require.EqualError(t, err, "jet miss")
}
