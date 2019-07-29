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

type Dependencies struct {
	FetchJet     func(*FetchJet)
	CheckJet     func(*CheckJet)
	WaitHot      func(*WaitHot)
	WaitHotWM    func(*WaitHotWM)
	EnsureIndex  func(*EnsureIndex)
	SendObject   func(*SendObject)
	GetCode      func(*GetCode)
	GetRequest   func(*GetRequest)
	SetRequest   func(*SetRequest)
	SetResult    func(*SetResult)
	GetPendings  func(*GetPendings)
	GetJet       func(*GetJet)
	HotObjects   func(*HotObjects)
	PassState    func(*PassState)
	CalculateID  func(*CalculateID)
	SetCode      func(*SetCode)
	SendRequests func(*SendRequests)
	HasPendings  func(*HasPendings)
}

// NewDependenciesMock returns all dependencies for handlers.
// It's all empty.
// Use it ONLY for tests.
func NewDependenciesMock() *Dependencies {
	return &Dependencies{
		FetchJet:     func(*FetchJet) {},
		CheckJet:     func(*CheckJet) {},
		WaitHot:      func(*WaitHot) {},
		WaitHotWM:    func(*WaitHotWM) {},
		EnsureIndex:  func(*EnsureIndex) {},
		SendObject:   func(*SendObject) {},
		GetCode:      func(*GetCode) {},
		SetRequest:   func(*SetRequest) {},
		SetResult:    func(*SetResult) {},
		GetPendings:  func(*GetPendings) {},
		GetJet:       func(*GetJet) {},
		HotObjects:   func(*HotObjects) {},
		PassState:    func(*PassState) {},
		CalculateID:  func(*CalculateID) {},
		SetCode:      func(*SetCode) {},
		SendRequests: func(*SendRequests) {},
		HasPendings:  func(*HasPendings) {},
	}
}
