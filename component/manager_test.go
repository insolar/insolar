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

package component

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Interface1 interface {
	Method1()
}

type Interface2 interface {
	Method2()
}

type Component1 struct {
	field1     string
	Interface2 Interface2 `inject:""`
	asd        int
	started    bool
}

func (cm *Component1) Start(ctx context.Context) error {
	cm.Method1()
	cm.Interface2.Method2()
	return nil
}

func (cm *Component1) Method1() {
	fmt.Println("Component1.Method1 called")
}

type Component2 struct {
	field2     string
	Interface1 Interface1 `inject:""`
	dsa        string
	started    bool
}

func (cm *Component2) Start(ctx context.Context) error {
	cm.Interface1.Method1()
	cm.Method2()
	return nil
}

func (cm *Component2) Stop(ctx context.Context) error {
	return nil
}

func (cm *Component2) Method2() {
	fmt.Println("Component2.Method2 called")
}

func TestComponentManager_Register(t *testing.T) {
	cm := Manager{}
	cm.Register(&Component1{}, &Component2{})

	assert.NoError(t, cm.Start(nil))
	assert.NoError(t, cm.Stop(nil))
}
