//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package args

import "sort"

const PrimeCount = 6550
const MaxPrime = 65581

func Prime(i int) int {
	switch n := len(primes); {
	case i >= n:
		return int(primesOverUint16[i-n])
	default:
		return int(primes[i])
	}
}

func Primes() PrimeList {
	return allPrimes
}

func OddPrimes() PrimeList {
	return oddPrimes
}

func IsolatedOddPrimes() PrimeList {
	return isolatedOddPrimes
}

var allPrimes = PrimeList{primes, primesOverUint16[0]}
var oddPrimes = PrimeList{primes[1:], primesOverUint16[0]}
var isolatedOddPrimes = PrimeList{isolatedPrimes[1:], isolatedPrimeOverUint16[0]}

type PrimeList struct {
	primes   []uint16
	maxPrime uint
}

func (m PrimeList) search(v int) uint {
	return uint(sort.Search(len(m.primes), func(i int) bool { return int(m.primes[i]) >= v }))
}

func (m PrimeList) Len() int {
	return len(m.primes) + 1
}

func (m PrimeList) Prime(i int) int {
	if len(m.primes) == i {
		return int(m.maxPrime)
	}
	return int(m.primes[i])
}

func (m PrimeList) Ceiling(v int) int {
	return m.Prime(int(m.search(v)))
}

func (m PrimeList) Floor(v int) int {
	switch n := m.search(v); {
	case n == 0:
		return int(m.primes[0])
	case n == uint(len(m.primes)):
		if v >= int(m.maxPrime) {
			return int(m.maxPrime)
		}
		return int(m.primes[len(m.primes)-1])
	case int(m.primes[n]) == v:
		return v
	default:
		return int(m.primes[n-1])
	}
}

func (m PrimeList) Nearest(v int) int {
	switch n := m.search(v); {
	case n == 0:
		return int(m.primes[0])
	case n == uint(len(m.primes)):
		return int(m.maxPrime)
	case (v - int(m.primes[n-1])) <= (int(m.primes[n]) - v):
		return int(m.primes[n-1])
	default:
		return int(m.primes[n])
	}
}

func (m PrimeList) Max() int {
	return int(m.maxPrime)
}

func (m PrimeList) Min() int {
	return int(m.primes[0])
}
