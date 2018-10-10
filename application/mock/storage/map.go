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

package storage

import (
	"fmt"

	"github.com/satori/go.uuid"
)

// MapStorage is Storage interface implementation with map as place for store.
type MapStorage struct {
	storage map[string]interface{}
	keys    []string
}

// NewMapStorage creates new MapStorage instance with empty storage.
func NewMapStorage() *MapStorage {
	mapStorage := MapStorage{
		storage: make(map[string]interface{}),
		keys:    []string{},
	}
	return &mapStorage
}

// Set store object into storage.
func (m *MapStorage) Set(obj interface{}) (string, error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	record := newUUID.String()
	_, exist := m.storage[record]
	if exist {
		return "", fmt.Errorf("object with record %s already exist", record)
	}
	m.storage[record] = obj
	m.keys = append(m.keys, record)
	return record, nil
}

// Get restore object from storage.
func (m *MapStorage) Get(record string) (interface{}, error) {
	obj, exist := m.storage[record]
	if !exist {
		return nil, fmt.Errorf("object with record %s does not exist", record)
	}
	return obj, nil
}

// GetKeys gives list of all keys in storage.
func (m *MapStorage) GetKeys() []string {
	return m.keys
}
