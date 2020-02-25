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

package component

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

// Logger interface provides methods for debug logging
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
}

// DefaultLogger logs to std out
type DefaultLogger struct{}

func (l *DefaultLogger) Debug(v ...interface{}) {
	log.Println(v...)
}

func (l *DefaultLogger) Debugf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// NoLogger skips log messages
type NoLogger struct{}

func (l *NoLogger) Debug(v ...interface{}) {
}

func (l *NoLogger) Debugf(format string, v ...interface{}) {
}

// Manager provide methods to manage components lifecycle
type Manager struct {
	parent     *Manager
	components []interface{}
	logger     Logger
	startStopLock sync.Mutex
	started       bool
}

// NewManager creates new component manager with default logger
func NewManager(parent *Manager) *Manager {
	return &Manager{parent: parent, logger: &DefaultLogger{}}
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

		for i := 0; i < componentType.NumField(); i++ {
			fieldMeta := componentType.Field(i)
			if value, ok := fieldMeta.Tag.Lookup("inject"); ok && component.Field(i).IsNil() {
				if value == "subcomponent" && m.parent == nil {
					continue
				}
				m.mustInject(component, fieldMeta)
			}
		}
	}
}

func (m *Manager) mustInject(component reflect.Value, fieldMeta reflect.StructField) {
	found := false
	if m.parent != nil {
		found = m.injectDependency(component, fieldMeta, m.parent.components)
	}
	found = found || m.injectDependency(component, fieldMeta, m.components)
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

func (m *Manager) injectDependency(component reflect.Value, dependencyMeta reflect.StructField, components []interface{}) (injectFound bool) {
	for _, componentMeta := range components {
		componentType := reflect.ValueOf(componentMeta).Type()

		if componentType.Implements(dependencyMeta.Type) {
			field := component.FieldByName(dependencyMeta.Name)
			field.Set(reflect.ValueOf(componentMeta))

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
	m.startStopLock.Lock()
	defer m.startStopLock.Unlock()

	for _, c := range m.components {
		if !m.isManaged(c) {
			continue
		}
		name := reflect.TypeOf(c).Elem().String()
		if s, ok := c.(Starter); ok {
			m.logger.Debug("ComponentManager: Start component: ", name)
			err := s.Start(ctx)
			if err != nil {
				return errors.Wrap(err, "Failed to start components.")
			}
			m.logger.Debugf("ComponentManager: Component %s started ", name)
		}
	}

	m.started = true
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
			continue
		}
		m.logger.Debug("ComponentManager: Init component: ", name)
		err := s.Init(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to init components.")
		}
	}
	return nil
}

// GracefulStop invokes GracefulStop method of all components which implements Starter interface
func (m *Manager) GracefulStop(ctx context.Context) error {
	for i := len(m.components) - 1; i >= 0; i-- {
		if !m.isManaged(m.components[i]) {
			continue
		}
		name := reflect.TypeOf(m.components[i]).Elem().String()
		if s, ok := m.components[i].(GracefulStopper); ok {
			m.logger.Debug("ComponentManager: GracefulStop component: ", name)

			err := s.GracefulStop(ctx)
			if err != nil {
				return errors.Wrap(err, "Failed to gracefully stop components.")
			}
		}
	}
	return nil
}

// Stop invokes Stop method of all components which implements Starter interface
func (m *Manager) Stop(ctx context.Context) error {
	m.startStopLock.Lock()
	defer m.startStopLock.Unlock()

	if !m.started {
		m.logger.Debug("ComponentManager: components are not started. Skip stopping")
		return nil
	}

	for i := len(m.components) - 1; i >= 0; i-- {
		if !m.isManaged(m.components[i]) {
			continue
		}
		name := reflect.TypeOf(m.components[i]).Elem().String()
		if s, ok := m.components[i].(Stopper); ok {
			m.logger.Debug("ComponentManager: Stop component: ", name)

			err := s.Stop(ctx)
			if err != nil {
				return errors.Wrap(err, "Failed to stop components.")
			}
		}
	}
	return nil
}

// SetLogger sets custom DefaultLogger
func (m *Manager) SetLogger(logger Logger) {
	m.logger = logger
}
