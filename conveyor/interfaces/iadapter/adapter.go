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

package iadapter

// Response gives access to response of adapter
//go:generate minimock -i github.com/insolar/insolar/conveyor/interfaces/iadapter.Response -o ./ -s _mock.go
type Response interface {
	// GetAdapterID returns adapter id
	GetAdapterID() uint32
	// GetElementID returns element id
	GetElementID() uint32
	// GetHandlerID returns handler id
	GetHandlerID() uint32
	// GetRespPayload returns payload
	GetRespPayload() interface{}
}

// NestedEvent gives access to nested event of adapter
//go:generate minimock -i github.com/insolar/insolar/conveyor/interfaces/iadapter.NestedEvent -o ./ -s _mock.go
type NestedEvent interface {
	// GetAdapterID returns adapter id
	GetAdapterID() uint32
	// GetParentElementID returns parent element id
	GetParentElementID() uint32
	// GetHandlerID returns handler id
	GetHandlerID() uint32
	// GetEventPayload returns event payload
	GetEventPayload() interface{}
}
