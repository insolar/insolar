//
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
//

package proc

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/bus"
)

type Dependencies struct {
	FetchJet            func(*FetchJet)
	CheckJet            func(*CheckJet)
	WaitHot             func(*WaitHot)
	WaitHotWM           func(*WaitHotWM)
	GetIndex            func(*EnsureIndex)
	EnsureIndex         func(*EnsureIndexWM)
	SendObject          func(*SendObject)
	GetCode             func(*GetCode)
	GetRequest          func(*GetRequest)
	SetBlob             func(*SetBlob)
	SetRequest          func(*SetRequest)
	SetResult           func(*SetResult)
	ActivateObject      func(*ActivateObject)
	DeactivateObject    func(*DeactivateObject)
	UpdateObject        func(*UpdateObject)
	RegisterChild       func(*RegisterChild)
	GetPendingRequests  func(*GetPendingRequests)
	GetPendingRequestID func(*GetPendingRequestID)
	GetJet              func(*GetJet)
	GetChildren         func(*GetChildren)
	HotData             func(*HotData)
	PassState           func(*PassState)
	CalculateID         func(*CalculateID)
	SetCode             func(*SetCode)
	SendRequests        func(*SendRequests)
	GetDelegate         func(*GetDelegate)
}

type ReturnReply struct {
	ReplyTo chan<- bus.Reply
	Err     error
	Reply   insolar.Reply
}

func (p *ReturnReply) Proceed(ctx context.Context) error {
	select {
	case p.ReplyTo <- bus.Reply{Reply: p.Reply, Err: p.Err}:
	case <-ctx.Done():
	}
	return nil
}

// NewDependenciesMock returns all dependencies for handlers.
// It's all empty.
// Use it ONLY for tests.
func NewDependenciesMock() *Dependencies {
	return &Dependencies{
		FetchJet:            func(*FetchJet) {},
		CheckJet:            func(*CheckJet) {},
		WaitHot:             func(*WaitHot) {},
		WaitHotWM:           func(*WaitHotWM) {},
		GetIndex:            func(*EnsureIndex) {},
		EnsureIndex:         func(*EnsureIndexWM) {},
		SendObject:          func(*SendObject) {},
		GetCode:             func(*GetCode) {},
		GetRequest:          func(*GetRequest) {},
		SetBlob:             func(*SetBlob) {},
		SetRequest:          func(*SetRequest) {},
		SetResult:           func(*SetResult) {},
		ActivateObject:      func(*ActivateObject) {},
		DeactivateObject:    func(*DeactivateObject) {},
		UpdateObject:        func(*UpdateObject) {},
		RegisterChild:       func(*RegisterChild) {},
		GetPendingRequests:  func(*GetPendingRequests) {},
		GetPendingRequestID: func(*GetPendingRequestID) {},
		GetJet:              func(*GetJet) {},
		GetChildren:         func(*GetChildren) {},
		HotData:             func(*HotData) {},
		PassState:           func(*PassState) {},
		CalculateID:         func(*CalculateID) {},
		SetCode:             func(*SetCode) {},
		GetDelegate:         func(*GetDelegate) {},
	}
}
