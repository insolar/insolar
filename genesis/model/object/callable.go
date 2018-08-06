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

package object

// Callable allows itself to be called by its reference.
type Callable interface {
	Object
	GetReference() Reference
	SetReference(reference Reference)
}

// BaseCallable is a base implementation of Callable.
type BaseCallable struct {
	BaseObject
	reference Reference
}

// GetReference return reference.
func (bc *BaseCallable) GetReference() Reference {
	return bc.reference
}

// SetReference sets reference.
func (bc *BaseCallable) SetReference(reference Reference) {
	bc.reference = reference
}
