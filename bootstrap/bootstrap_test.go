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

package bootstrap

import (
	"reflect"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/messagebus"
	"github.com/stretchr/testify/assert"
)

// linkAll - link dependency for all components
func linkAll(cm core.Components) {
	v := reflect.ValueOf(cm)
	for i := 0; i < v.NumField(); i++ {
		componentName := v.Field(i).String()
		log.Infof("Starting component `%s` ...", componentName)
		if v.Field(i).Interface() == nil {
			continue
		}
		err := v.Field(i).Interface().(core.Component).Start(cm)
		if err != nil {
			log.Fatalf("failed to start component %s : %s", componentName, err.Error())
		}

		log.Infof("Component `%s` successfully started", componentName)
	}
}

func TestNewBootstrapper(t *testing.T) {
	t.Skip("mock MessageBus needed")
	mb, err := messagebus.NewMessageBus()
	assert.NotNil(t, mb)
	assert.NoError(t, err)

	l, err := ledger.NewLedger(configuration.NewLedger())
	assert.NoError(t, err)
	assert.NotNil(t, l)

	bootstrapper, err := NewBootstrapper(configuration.NewBootstrap())
	assert.NoError(t, err)
	assert.NotNil(t, bootstrapper)

	components := core.Components{Ledger: l, MessageBus: mb, Bootstrapper: bootstrapper}
	//err = bootstrapper.Start(components)
	//assert.NoError(t, err)
	linkAll(components)

	err = bootstrapper.Stop()
	assert.NoError(t, err)
}
