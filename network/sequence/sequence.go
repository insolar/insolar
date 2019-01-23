/*
 *    Copyright 2019 Insolar
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

package sequence

import (
	"github.com/insolar/insolar/network/utils"
)

type Sequence uint64

type Generator interface {
	Generate() Sequence
}

type generatorImpl struct {
	sequence *uint64
}

func NewGeneratorImpl() Generator {
	return &generatorImpl{
		sequence: new(uint64),
	}
}

func (sg *generatorImpl) Generate() Sequence {
	return Sequence(utils.AtomicLoadAndIncrementUint64(sg.sequence))
}
