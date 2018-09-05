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

import "github.com/insolar/insolar/core"

// Types for RPC requests and responses between goplugin and goinsider.
// Calls from goplugin to goinsider go "downwards" and names are
// prefixed with "Down". Reverse calls go "upwards", so "Up" prefix

// DownCallMethodReq is a set of arguments for CallMethod RPC in the runner
type DownCallMethodReq struct { // todo it may use foundation.Context
	Reference core.RecordRef
	Data      []byte
	Method    string
	Arguments core.Arguments
}

// DownCallMethodResp is response from CallMethod RPC in the runner
type DownCallMethodResp struct {
	Data []byte
	Ret  core.Arguments
	Err  error
}

// DownCallConstructorReq is a set of arguments for CallConstructor RPC
// in the runner
type DownCallConstructorReq struct {
	Reference core.RecordRef
	Name      string
	Arguments core.Arguments
}

// DownCallConstructorResp is response from CallConstructor RPC in the runner
type DownCallConstructorResp struct {
	Ret core.Arguments
	Err error
}

// UpGetCodeReq is a set of arguments for GetCode RPC in goplugin
type UpGetCodeReq struct {
	Reference core.RecordRef
}

// UpGetCodeResp is response from GetCode RPC in goplugin
type UpGetCodeResp struct {
	Code []byte
}

// UpRouteReq is a set of arguments for Route RPC in goplugin
type UpRouteReq struct {
	Reference core.RecordRef
	Method    string
	Arguments core.Arguments
}

// UpRouteResp is response from Route RPC in goplugin
type UpRouteResp struct {
	Result core.Arguments
	Err    error
}

// UpRouteConstructorReq is a set of arguments for RouteConstructor RPC in goplugin
type UpRouteConstructorReq struct {
	Reference   core.RecordRef
	Constructor string
	Arguments   core.Arguments
}

// UpRouteConstructorResp is response from RouteConstructor RPC in goplugin
type UpRouteConstructorResp struct {
	Data []byte
	Err  error
}

// Object is an inner representation of storage object for transfering it over API
type Object struct {
	MachineType core.MachineType
	Reference   core.RecordRef
	Data        []byte
}
