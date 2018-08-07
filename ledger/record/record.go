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

package record

// ProjectionType is a "view filter" for record.
// E.g. we can read whole object or just it's hash.
type ProjectionType uint32

// Memory is actual contracts' state, variables etc.
type Memory []byte

// Record is base interface for all records.
type Record interface {
	Domain() ID
}
