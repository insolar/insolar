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

package hash

import (
	"hash"
)

type hashWrapper struct {
	hash    hash.Hash
	sumFunc func([]byte) []byte
}

func (h *hashWrapper) Write(p []byte) (n int, err error) {
	return h.hash.Write(p)
}

func (h *hashWrapper) Sum(b []byte) []byte {
	return h.hash.Sum(b)
}

func (h *hashWrapper) Reset() {
	h.hash.Reset()
}

func (h *hashWrapper) Size() int {
	return h.hash.Size()
}

func (h *hashWrapper) BlockSize() int {
	return h.hash.BlockSize()
}

