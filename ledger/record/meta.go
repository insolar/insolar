/*
 *    Copyright 2018 Insolar
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

package record

// GenesisRecord is the first record created on storage. It's used to link root objects to it.
type GenesisRecord struct {
}

// IsDeactivation determines if current state is deactivation.
func (*GenesisRecord) IsDeactivation() bool {
	return false
}

// IsAmend determines if current state is amend.
func (*GenesisRecord) IsAmend() bool {
	return false
}

// GetMemory returns state memory.
func (*GenesisRecord) GetMemory() []byte {
	return nil
}

// GetCode returns state code.
func (*GenesisRecord) GetCode() *Reference {
	return nil
}

// ChildRecord is a child activation record. Its used for children iterating.
type ChildRecord struct {
	PrevChild *ID

	Ref Reference // Reference to the child's head.
}

// Next returns next record.
func (r *ChildRecord) Next() *ID {
	if r == nil {
		return nil
	}

	return r.PrevChild
}
