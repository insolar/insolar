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

package messagebus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	queue := NewExpiryQueue(10 * time.Second)
	queue.Push("a")
	queue.Push("b")
	queue.Push("c")
	assert.NotNil(t, queue)
	assert.Equal(t, queue.items.Len(), 3)
	assert.Equal(t, queue.Pop(), "a")
	assert.Equal(t, queue.items.Len(), 2)
	assert.Equal(t, queue.Pop(), "b")
	assert.Equal(t, queue.items.Len(), 1)
	assert.Equal(t, queue.Pop(), "c")
	assert.Equal(t, queue.items.Len(), 0)
	assert.Nil(t, queue.Pop())
	queue.Push("a")
	queue.Push("b")
	queue.Push("c")
	values := []string{}
	for _, v := range queue.PopValues() {
		values = append(values, v.(string))
	}
	assert.Equal(t, values, []string{"a", "b", "c"})
}

func TestExpiryQueue(t *testing.T) {
	queue := NewExpiryQueue(20 * time.Millisecond)
	queue.Push("a")
	queue.Push("b")
	queue.Push("c")
	time.Sleep(40 * time.Millisecond)
	assert.Equal(t, 0, queue.items.Len())

}
