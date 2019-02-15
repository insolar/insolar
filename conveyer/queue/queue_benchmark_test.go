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

package queue

import (
	"sync"
	"testing"
)

func makeTestQueue() IQueue {
	//return NewLockFreeQueue()
	return NewMutexQueue()
}

func BenchmarkPut(b *testing.B) {
	queue := makeTestQueue()
	for n := 0; n < b.N; n++ {
		queue.SinkPush(n)
	}
}

func BenchmarkPutRemove(b *testing.B) {
	queue := makeTestQueue()
	for n := 0; n < b.N; n++ {
		queue.SinkPush(n)
		queue.RemoveAll()
	}
}

func benchParallel(i int, b *testing.B) {
	queue := makeTestQueue()

	parallelGet := i
	parallelPut := i

	wg := sync.WaitGroup{}
	wg.Add(parallelPut + parallelGet)

	for i := 0; i < parallelGet; i++ {
		go func(wg *sync.WaitGroup, q IQueue) {
			for n := 0; n < b.N; n++ {
				q.SinkPush(n)
			}
			wg.Done()
		}(&wg, queue)
	}

	for i := 0; i < parallelPut; i++ {
		go func(wg *sync.WaitGroup, q IQueue) {
			for n := 0; n < b.N; n++ {
				q.RemoveAll()
			}
			wg.Done()
		}(&wg, queue)
	}

	wg.Wait()
}

func BenchmarkParallel1(b *testing.B) { benchParallel(1, b) }

func BenchmarkParallel2(b *testing.B)  { benchParallel(2, b) }
func BenchmarkParallel3(b *testing.B)  { benchParallel(3, b) }
func BenchmarkParallel10(b *testing.B) { benchParallel(10, b) }
func BenchmarkParallel20(b *testing.B) { benchParallel(20, b) }
func BenchmarkParallel40(b *testing.B) { benchParallel(40, b) }
func BenchmarkParallel80(b *testing.B) { benchParallel(80, b) }
