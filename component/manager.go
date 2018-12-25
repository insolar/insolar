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
	"reflect"

	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

// Manager provide methods to manage components lifecycle
type Manager struct {
	components []interface{}
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
			if _, ok := fieldMeta.Tag.Lookup("inject"); ok && component.Field(i).IsNil() {
				log.Debugf("ComponentManager: Component %s need inject: ", componentType.String(), fieldMeta.Name)
				m.mustInject(component, fieldMeta)
			}
		}
	}
}

func (m *Manager) mustInject(component reflect.Value, fieldMeta reflect.StructField) {
	for _, componentMeta := range m.components {
		componentType := reflect.ValueOf(componentMeta).Type()

		if componentType.Implements(fieldMeta.Type) {
			field := component.FieldByName(fieldMeta.Name)
			field.Set(reflect.ValueOf(componentMeta))

			log.Debugf(
				"ComponentManager: Inject interface %s with %s: ",
				field.Type().String(),
				componentType.String(),
			)
			return
		}
	}

	panic(fmt.Sprintf(
		"Component %s injects not existing component with interface %s to field %s",
		component.Type().String(),
		fieldMeta.Type.String(),
		fieldMeta.Name,
	))
}

// Start invokes Start method of all components which implements Starter interface
func (m *Manager) Start(ctx context.Context) error {
	for _, c := range m.components {
		name := reflect.TypeOf(c).Elem().String()
		if s, ok := c.(Starter); ok {
			log.Debugln("ComponentManager: Start component: ", name)
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
		name := reflect.TypeOf(c).Elem().String()
		s, ok := c.(Initer)
		if !ok {
			log.Debugf("ComponentManager: Component %s has no Init method", name)
			continue
		}
		log.Debugln("ComponentManager: Init component: ", name)
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
		name := reflect.TypeOf(m.components[i]).Elem().String()
		if s, ok := m.components[i].(Stopper); ok {
			log.Debugln("ComponentManager: Stop component: ", name)

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
