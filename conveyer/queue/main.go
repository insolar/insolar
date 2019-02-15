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
	"fmt"
	"sync"
)

func main() {
	queue := NewMutexQueue()

	//parallelGet := 30

	parallelPut := 30
	wg := sync.WaitGroup{}
	//wg.Add(parallelPut*2 + parallelGet)
	wg.Add(parallelPut)

	numIterations := 100

	for i := 0; i < parallelPut; i++ {
		go func(wg *sync.WaitGroup, q IQueue) {
			//fmt.Println("START")
			for i := 0; i < numIterations; i++ {
				q.SinkPush(i)
			}
			wg.Done()
		}(&wg, queue)
	}

	// for i := 0; i < parallelGet; i++ {
	// 	go func(wg *sync.WaitGroup, q IQueue) {
	// 		//fmt.Println("START")
	// 		for i := 0; i < numIterations; i++ {
	// 			q.RemoveAll()
	// 		}
	// 		wg.Done()
	// 	}(&wg, queue)
	// }
	//
	// for i := 0; i < parallelPut; i++ {
	// 	go func(wg *sync.WaitGroup, q IQueue) {
	// 		//fmt.Println("START")
	// 		input := []interface{}{}
	// 		for i := 0; i < numIterations; i++ {
	// 			input = append(input, i)
	// 		}
	//
	// 		queue.SinkPushAll(input)
	// 		wg.Done()
	// 	}(&wg, queue)
	// }

	wg.Wait()

	result := queue.RemoveAll()

	fmt.Printf("Num Elements: %d . Must be: %d\n", len(result), parallelPut*numIterations)

}
