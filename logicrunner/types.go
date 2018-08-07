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

package logicrunner

// Reference is a contract address
type Reference string

// Object is an inner representation of storage object for transfwering it over API
type Object struct {
	MachineType MachineType
	Reference   Reference
	Data        []byte
}

// Argument is a dedicated type for arguments, that represented as bynary cbored blob
type Argument []byte
