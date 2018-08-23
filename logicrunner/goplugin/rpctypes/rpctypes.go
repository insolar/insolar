/*
 *    Copyright 2018 INS Ecosystem
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

package rpctypes

import "github.com/insolar/insolar/logicrunner"

// Types for RPC requests and responses between goplugin and goinsider.
// Calls from goplugin to goinsider go "downwards" and names are
// prefixed with "Down". Reverse calls go "upwards", so "Up" prefix

// DownCallReq is a set of arguments for Call RPC in the runner
type DownCallReq struct { // todo it may use foundation.Context
	Reference logicrunner.Reference
	Data      []byte
	Method    string
	Arguments logicrunner.Arguments
}

// DownCallResp is response from Call RPC in the runner
type DownCallResp struct {
	Data []byte
	Ret  logicrunner.Arguments
	Err  error
}
