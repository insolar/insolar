/*
 *    Copyright 2019 Insolar Technologies
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
	"reflect"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/log"
)

// Manager provide methods to manage components lifecycle
type Manager struct {
	parent     *Manager
	components []interface{}
}

// NewManager creates new component manager
func NewManager(parent *Manager) *Manager {
	return &Manager{parent: parent}
}

// Register components in Manager and inject required dependencies.
// Register can inject interfaces only, tag public struct fields with `inject:""`.
// If the injectable struct already has a value on the tagged field, the value WILL NOT be overridden.
func (m *Manager) Register(components ...interface{}) {
	m.components = append(m.components, components...)
}

// Inject components in Manager and inject required dependencies
// Inject can inject interfaces only, tag public struct fields with `inject:""`
func (m *Manager) Inject(components ...interface{}) {
	m.Register(components...)

	for _, componentMeta := range m.components {
		component := reflect.ValueOf(componentMeta).Elem()
		componentType := component.Type()
		log.Debugf("ComponentManager: Inject component: %s", componentType.String())

		for i := 0; i < componentType.NumField(); i++ {
			fieldMeta := componentType.Field(i)
			if value, ok := fieldMeta.Tag.Lookup("inject"); ok && component.Field(i).IsNil() {
				if value == "subcomponent" && m.parent == nil {
					continue
				}
				log.Debugf("ComponentManager: Component %s need inject: %s", componentType.String(), fieldMeta.Name)
				m.mustInject(component, fieldMeta)
			}
		}
	}
}

func (m *Manager) mustInject(component reflect.Value, fieldMeta reflect.StructField) {
	found := false
	if m.parent != nil {
		found = injectDependency(component, fieldMeta, m.parent.components)
	}
	found = found || injectDependency(component, fieldMeta, m.components)
	if found {
		return
	}

	panic(fmt.Sprintf(
		"Component %s injects not existing component with interface %s to field %s",
		component.Type().String(),
		fieldMeta.Type.String(),
		fieldMeta.Name,
	))
}

func injectDependency(component reflect.Value, dependencyMeta reflect.StructField, components []interface{}) (injectFound bool) {
	for _, componentMeta := range components {
		componentType := reflect.ValueOf(componentMeta).Type()

		if componentType.Implements(dependencyMeta.Type) {
			field := component.FieldByName(dependencyMeta.Name)
			field.Set(reflect.ValueOf(componentMeta))

			log.Debugf(
				"ComponentManager: Inject interface %s with %s: ",
				field.Type().String(),
				componentType.String(),
			)
			return true
		}
	}
	return false
}

func (m *Manager) isManaged(component interface{}) bool {
	// TODO: refactor this behavior
	if m.parent == nil {
		return true
	}
	for _, c := range m.parent.components {
		if c == component {
			return false
		}
	}
	return true
}

// Start invokes Start method of all components which implements Starter interface
func (m *Manager) Start(ctx context.Context) error {
	for _, c := range m.components {
		if !m.isManaged(c) {
			continue
		}
		name := reflect.TypeOf(c).Elem().String()
		if s, ok := c.(Starter); ok {
			log.Debug("ComponentManager: Start component: ", name)
			err := s.Start(ctx)
			if err != nil {
				return errors.Wrap(err, "Failed to start components.")
			}
			log.Debugf("ComponentManager: Component %s started ", name)
		} else {
			log.Debugf("ComponentManager: Component %s has no Start method", name)
		}
	}
	return nil
}

// Init invokes Init method of all components which implements Initer interface
func (m *Manager) Init(ctx context.Context) error {
	for _, c := range m.components {
		if !m.isManaged(c) {
			continue
		}
		name := reflect.TypeOf(c).Elem().String()
		s, ok := c.(Initer)
		if !ok {
			log.Debugf("ComponentManager: Component %s has no Init method", name)
			continue
		}
		log.Debug("ComponentManager: Init component: ", name)
		err := s.Init(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to init components.")
		}
	}
	return nil
}

// Stop invokes Stop method of all components which implements Starter interface
func (m *Manager) Stop(ctx context.Context) error {

	for i := len(m.components) - 1; i >= 0; i-- {
		if !m.isManaged(m.components[i]) {
			continue
		}
		name := reflect.TypeOf(m.components[i]).Elem().String()
		if s, ok := m.components[i].(Stopper); ok {
			log.Debug("ComponentManager: Stop component: ", name)

			err := s.Stop(ctx)
			if err != nil {
				return errors.Wrap(err, "Failed to stop components.")
			}
		} else {
			log.Debugf("ComponentManager: Component %s has no Stop method", name)
		}
	}
	return nil
}
