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

// Package logicrunner - infrastructure for executing smartcontracts
package logicrunner

// MachineType is a type of virtual machine
type MachineType int

// Real constants of MachineType
const (
	MachineTypeBuiltin MachineType = iota
	MachineTypeGoPlugin
)

// LogicRunner is a general interface of contract executor
type LogicRunner interface {
	Start()
	Stop()
	Exec(object Object, method string, args Arguments) (ret Arguments, err error)
}
