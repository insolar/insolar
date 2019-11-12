//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package jetid

import "testing"

func TestPrefixTree_Print(t *testing.T) {
	pt := PrefixTree{}
	PrintTable()

	Split(0, 0)
	PrintTable()

	Merge(0, 1)
	PrintTable()

	Split(0, 0)
	PrintTable()

	Split(0, 1)
	PrintTable()

	Split(0, 2)
	PrintTable()

	Split(1, 1)
	PrintTable()

	Split(3, 2)
	PrintTable()

	Merge(0, 3)
	PrintTable()

	Merge(0, 2)
	PrintTable()

	Merge(3, 3)
	PrintTable()

	Merge(1, 2)
	PrintTable()

	Merge(0, 1)
	PrintTable()
}

func TestPrefixTree_SplitMax0(t *testing.T) {
	pt := PrefixTree{}
	Split(0, 0)
	Split(0, 1)
	Split(0, 2)
	Split(0, 3)
	Split(0, 4)
	Split(0, 5)
	Split(0, 6)
	Split(0, 7)
	Split(0, 8)
	Split(0, 9)
	Split(0, 10)
	Split(0, 11)
	Split(0, 12)
	Split(0, 13)
	Split(0, 14)
	Split(0, 15)
	PrintTable()
	Merge(0, 16)
	PrintTable()
}

func TestPrefixTree_SplitMax1(t *testing.T) {
	pt := PrefixTree{}
	Split(0, 0)
	Split(1, 1)
	Split(3, 2)
	Split(7, 3)
	Split(15, 4)
	Split(31, 5)
	Split(63, 6)
	Split(127, 7)
	Split(255, 8)
	Split(511, 9)
	Split(1023, 10)
	Split(2047, 11)
	Split(4095, 12)
	Split(8191, 13)
	Split(16383, 14)
	Split(32767, 15)
	PrintTable()
	Merge(32767, 16)
	PrintTable()
}
