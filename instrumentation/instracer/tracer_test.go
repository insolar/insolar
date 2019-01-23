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

package instracer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

func TestTracerBasics(t *testing.T) {
	ctx := inslogger.ContextWithTrace(context.Background(), "tracenotdefined")
	_, err := instracer.RegisterJaeger("server", "nodeRef", "localhost:6831", "", 1)
	assert.NoError(t, err)
	_, span := instracer.StartSpan(ctx, "root")
	assert.NotNil(t, span)
}
