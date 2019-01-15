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

package component

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SubInterface1 interface {
	Method3()
}

type SubInterface2 interface {
	Method4()
}

type SubComponent1 struct {
	field1        string
	SubInterface2 SubInterface2 `inject:"subcomponent"`
	Interface1    Interface1    `inject:""`
	asd           int
	started       bool
}

func (s *SubComponent1) Start(ctx context.Context) error {
	s.Method3()
	s.SubInterface2.Method4()
	s.Interface1.Method1()
	return nil
}

func (cm *SubComponent1) Method3() {
	fmt.Println("SubComponent1.Method3 called")
}

type SubComponent2 struct {
	field2        string
	SubInterface1 SubInterface1 `inject:"subcomponent"`
	dsa           string
	started       bool
}

func (s *SubComponent2) Start(ctx context.Context) error {
	s.SubInterface1.Method3()
	s.Method4()
	return nil
}

func (cm *SubComponent2) Stop(ctx context.Context) error {
	return nil
}

func (cm *SubComponent2) Method4() {
	fmt.Println("SubComponent2.Method4 called")
}

type BigComponent struct {
	// subcomponents
	SubInterface1 SubInterface1 `inject:"subcomponent"`
	SubInterface2 SubInterface2 `inject:"subcomponent"`

	// components
	Interface1 Interface1 `inject:""`

	cm *Manager
}

func (b *BigComponent) Init(ctx context.Context) error {
	b.cm.Inject(&SubComponent1{}, &SubComponent2{}, b)
	return nil
}

func (b *BigComponent) Start(ctx context.Context) error {
	b.cm.Start(ctx)
	b.SubInterface1.Method3()
	b.SubInterface2.Method4()
	return nil
}

func NewBigComponent(componentManager *Manager) *BigComponent {
	return &BigComponent{cm: NewManager(componentManager)}
}

func TestManager_Subcomponents(t *testing.T) {

	rootCm := NewManager(nil)

	big := NewBigComponent(rootCm)
	c1 := &Component1{}
	c2 := &Component2{}

	rootCm.Inject(big, c1, c2)
	err := rootCm.Init(context.Background())
	assert.NoError(t, err)

	assert.NotNil(t, c1.Interface2)
	assert.NotNil(t, c2.Interface1)
	assert.NotNil(t, big.Interface1)

	assert.NotNil(t, big.SubInterface1)
	assert.NotNil(t, big.SubInterface2)

	assert.NotNil(t, big.SubInterface1.(*SubComponent1).Interface1)
	assert.NotNil(t, big.SubInterface1.(*SubComponent1).SubInterface2)

	require.NoError(t, rootCm.Start(nil))
	require.NoError(t, rootCm.Stop(nil))
}
