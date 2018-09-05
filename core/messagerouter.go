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

package core

// Arguments is a dedicated type for arguments, that represented as bynary cbored blob
type Arguments []byte

// Message is a routable packet, ATM just a method call
type Message struct {
	Caller      struct{}
	Constructor bool
	Reference   RecordRef
	Method      string
	Arguments   Arguments
}

// Response to a `Message`
type Response struct {
	Data   []byte
	Result []byte
	Error  error
}

// MessageRouter interface
type MessageRouter interface {
	Component
	Route(msg Message) (resp Response, err error)
}
