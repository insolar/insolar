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

// IAdapterResponse gives access to reponse of iadapter
//go:generate minimock -i github.com/insolar/insolar/conveyor/interfaces/iadapter.IAdapterResponse -o ./ -s _mock.go
type IAdapterResponse interface {
	GetAdapterID() uint32
	GetElementID() uint32
	GetHandlerID() uint32
	GetRespPayload() interface{}
}

// IAdapterNestedEvent gives access to nested event of iadapter
//go:generate minimock -i github.com/insolar/insolar/conveyor/interfaces/iadapter.IAdapterNestedEvent -o ./ -s _mock.go
type IAdapterNestedEvent interface {
	GetAdapterID() uint32
	GetParentElementID() uint32
	GetHandlerID() uint32
	GetEventPayload() interface{}
}
