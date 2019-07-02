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

package logicrunner

type RequestResult struct {
	Result       []byte
	NewMemory    []byte
	Activation   bool
	Deactivation bool
}

func NewRequestResult(res []byte) *RequestResult {
	return &RequestResult{
		Result: res,
	}
}

func (rr *RequestResult) Activate(mem []byte) {
	rr.Activation = true
	rr.NewMemory = mem
}

func (rr *RequestResult) Update(mem []byte) {
	rr.NewMemory = mem
}

func (rr *RequestResult) Deactivate() {
	rr.Deactivation = true
}
